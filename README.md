## ip

My IP Service: the IP address database is derived from [lionsoul/ip2region](https://github.com/lionsoul/ip2region),
and is for reference only!

### Deploy

A web app, the default listen address: `0.0.0.0:7000`,
the generated executable file name is `mip`,
the help option is `mip -h/--help`:
```bash
$ ./mip -h
Usage of mip:
  -db string
        the ip2region.xdb filepath (default "data/ip2region.xdb")
  -host string
        http listen host (default "0.0.0.0")
  -port uint
        http listen port (default 7000)
  -prefix string
        route prefix
  -v    show version and exit
```

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
