## ip

My IP Service

### Deploy

A web app, the default listen address: `0.0.0.0:7000`,
the generated executable file name is `ip`,
the help option is `ip -h/--help`.

### binary

go to [releases](releases) and download package, run it.

### docker

```bash
$ docker run -d --name ip --net=host staugur/ip
```

### API

默认使用根为前缀，可以定义为其他（启动时使用prefix选项），比如 /v1

The default use `/` is the prefix,
and it can be defined as other(`prefix` option is used at startup),
such as `/v1`.

#### /myip

Show client ip

#### /addr

Show client ip, isp, area, country, string.

### /rest

The returned result is the same as `/addr`, but the type is json.
