GOOS?="linux"
GOARCH?="amd64"

all:	main.go config.go
	GOARCH=${GOARCH} GOOS=${GOOS} go build
386_var:
	$(eval GOARCH=386)
arm_var:
	$(eval GOARCH=arm)
arm64_var:
	$(eval GOARCH=arm64)
mips_var:
	$(eval GOARCH=mips)
mipsle_var:
	$(eval GOARCH=mipsle)

386: 386_var all
arm: arm_var all
arm64: arm64_var all
mips: mips_var all
mipsle: mipsle_var all
