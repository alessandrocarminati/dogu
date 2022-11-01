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
Host:"https://localtunnel.me"
Target:"localhost"
Port:8080
