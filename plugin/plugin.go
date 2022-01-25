package plugin

import (
	"fmt"

	"github.com/team4yf/yf-fpm-server-go/ctx"
	"github.com/team4yf/yf-fpm-server-go/fpm"
)

type WSNamespace struct {
	Name string
	Path string
}
type WSOptions struct {
	Enable    bool
	Namespace []WSNamespace
}

type Body struct {
	Sender    string   `json:"sender"`
	Receiver  []string `json:"receiver"`
	Namespace []string `json:"namespace"`
	Boardcast bool     `json:"boardcast"`
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
			hub := newHub()
			go hub.run()
			fpmApp.BindHandler("/ws/{channel}", func(c *ctx.Ctx, _ *fpm.Fpm) {
				channel := c.Param("channel")
				fmt.Println("Channel:", channel)
				serveWs(hub, c.GetResponse(), c.GetRequest())
			})
			fpmApp.AddBizModule("ws", &fpm.BizModule{
				"send": func(param *fpm.BizParam) (data interface{}, err error) {
					body := Body{}
					if err = param.Convert(&body); err != nil {
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
					for _, client := range hub.Clients {
						client.Send <- []byte(body.Payload)
					}
					data = 1
					return
				},
				"clients": func(param *fpm.BizParam) (data interface{}, err error) {
					clients := []string{}
					for id := range hub.Clients {
						clients = append(clients, id)
					}
					data = clients
					return
				},
			})
		},
	})
}
