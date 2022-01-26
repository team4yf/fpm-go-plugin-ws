package plugin

import (
	"errors"

	"github.com/team4yf/fpm-go-pkg/utils"
	"github.com/team4yf/yf-fpm-server-go/ctx"
	"github.com/team4yf/yf-fpm-server-go/fpm"
)

type WSNamespace struct {
	Name       string `json:"name"`
	Path       string `json:"path"`
	Validation string `json:"validation"`
	Extra      string `json:"extra"`
}
type WSOptions struct {
	Enable    bool                    `json:"enable"`
	Namespace map[string]*WSNamespace `json:"namespace"`
}

type Body struct {
	Sender    string   `json:"sender"`
	Receiver  []string `json:"receiver"`
	Namespace string   `json:"namespace"`
	Payload   string   `json:"payload"`
}

func init() {

	fpm.RegisterByPlugin(&fpm.Plugin{
		Name: "fpm-plugin-ws",
		V:    "0.0.1",
		Handler: func(fpmApp *fpm.Fpm) {
			if !fpmApp.HasConfig("ws") {
				return
			}
			options := WSOptions{}
			if err := fpmApp.FetchConfig("ws", &options); err != nil {
				fpmApp.Logger.Errorf("Failed to fetch websocket config: %v", err)
				return
			}
			if !options.Enable {
				return
			}
			hubs := map[string]*Hub{}
			for _, v := range options.Namespace {
				hub := NewHub(v.Name, v)
				hubs[v.Name] = hub
				go hub.Run()
			}
			fpmApp.BindHandler("/ws/{ns}", func(c *ctx.Ctx, _ *fpm.Fpm) {
				ns := c.Param("ns")
				h, ok := hubs[ns]
				if !ok {
					c.Fail(errors.New("NAME_SPACE_NOT_FOUND"))
					return
				}
				if h.Options.Validation == "jwt" {
					token := c.Query("token")
					if ok, _ := utils.CheckToken(token); !ok {
						c.Fail(errors.New("TOKEN_NO_VALID"))
						return
					}
				}
				if h.Options.Validation == "key" {
					token := c.Query("token")
					if token != h.Options.Extra {
						c.Fail(errors.New("TOKEN_NO_VALID"))
						return
					}
				}
				serveWs(h, c)
			})
			fpmApp.AddBizModule("ws", &fpm.BizModule{
				"send": func(param *fpm.BizParam) (data interface{}, err error) {
					body := Body{}
					if err = param.Convert(&body); err != nil {
						return
					}
					hub, ok := hubs[body.Namespace]
					if !ok {
						err = errors.New("NAME_SPACE_NOT_FOUND")
						return
					}
					for _, id := range body.Receiver {
						hub.Clients[id].Send <- []byte(body.Payload)
					}
					data = 1
					return
				},
				"broadcast": func(param *fpm.BizParam) (data interface{}, err error) {
					body := Body{}
					if err = param.Convert(&body); err != nil {
						return
					}
					hub, ok := hubs[body.Namespace]
					if !ok {
						err = errors.New("NAME_SPACE_NOT_FOUND")
						return
					}
					for _, client := range hub.Clients {
						client.Send <- []byte(body.Payload)
					}
					data = 1
					return
				},
				"clients": func(param *fpm.BizParam) (data interface{}, err error) {
					clients := map[string][]string{}
					for _, hub := range hubs {
						ids := []string{}
						for id := range hub.Clients {
							ids = append(ids, id)
						}
						clients[hub.Namespace] = ids
					}
					data = clients
					return
				},
			})
		},
	})
}
