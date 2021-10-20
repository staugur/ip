## ip

My IP Service

### Deploy

A web app, the default listen address: `0.0.0.0:7000`,
the generated executable file name is `mip`,
the help option is `mip -h/--help`.

### binary

Go to [releases](https://github.com/staugur/ip/releases) and download package, run it.

### docker

```bash
$ docker run -d --name ip --net=host staugur/ip
```

### API

The default use `/` is the prefix,
and it can be defined as other(`prefix` option is used at startup),
such as `/v1`.

#### /myip

Show client ip

#### /addr

Show client ip, isp, area, country, string.

### /rest

The returned result is the same as `/addr`, but the type is json.
