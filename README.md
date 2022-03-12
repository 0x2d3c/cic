### Chat In Cmd
- Build
  - `go build cic.go`
- Usage of cic
```shell
  -f string
        used config file path (default "cfg.json")
  -m string
        running mode (default "s")
```
- Server

```shell
./cic -m s -f cfg.json
```

- Client

```shell
./cic -m c -f cfg.json
```