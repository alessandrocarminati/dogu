# Dogu (道具)
Dogu is a small utility that can be used st control via http command a machine behind a firewall.
It uses local tunnel to expose a http service to the internet, and implements a small http server that executes a few commands.

# build
Makefile provides a few targest:
Default target compiles for x86_64,
```
$ make
GOARCH="amd64" GOOS="linux" go build
```
Additional target are supported: `386`, `arm`, `arm64` and `mips`.
All these targets build for linux.
If another go supported target is needed, GOLANG env variables can be overridden on the make command line.
```
$ make GOOS=windows
GOARCH="amd64" GOOS=windows go build
$ make mipsle GOMIPS=softfloat 
GOARCH=mipsle GOOS="linux" go build
```

# usage
Dogu supports a few command line switches:
```
$ ./dogu
Missing needed arg
App Name: nav
Descr: kernel symbol navigator
	-h	<v>	Specifies localtunnel host
	-r	<v>	Specifies request domain
	-p	<v>	Specifies http service local port
	-j	<v>	Specifies config file
```
Only the `-r` switch is mandatory, and it is needed to specify the base URL for the service:

Default values:

| field  | Value                  |
|--------|------------------------|
| Host   | https://localtunnel.me |
| Target | localhost              |
| Port   | 8080                   |

# http commands
All commands are sent using the `GET` verb easing its usage from commandline.

Targeted usage is something like:

```
$ wget -O - -q "https://example.loca.lt/hello"
Service is alive.
```

| function   | args | arg strings | description                                                                 |
|------------|------|-------------|-----------------------------------------------------------------------------|
| hello      | 0    |             | Sends back a hello string                                                   |
| cmd_fore   | 1    | cmd         | Executes a command in foreground                                            |
| cmd_back   | 1    | cmd         | Executes a command in background                                            |
| cmd_backc  | 0    |             | Returns the output for the last background command                          |
| upd_script | 2    | name, b64pl | Uploads a script/executable sets execution flag, payload is meant in base64 |
| getlog     | 0    |             | Return the list of the commands received                                    |

