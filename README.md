# fpm-go-plugin-ws

## Config

```yaml
ws:
  enable: true
  namespace:
    foo:
      name: "foo"
    bar:
      name: "bar"
```

## Events

- #ws/receive

  ```json
	{
		"namespace": "",
		"message": "",
		"clientID": "",
	}
  ```
- #ws/errors
- #ws/connect
- #ws/close