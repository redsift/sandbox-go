module github.com/redsift/sandbox-go

go 1.27.4

require (
	github.com/redsift/go-mangosock v0.2.1
	github.com/redsift/go-sandbox-rpc v0.4.0
	github.com/stretchr/testify v1.12.1
	golang.org/x/mod v0.43.0
)

require (
	github.com/Microsoft/go-winio v0.6.4 // indirect
	go.yaml.in/yaml/v3 v3.0.6 // indirect
	golang.org/x/sys v0.50.0 // indirect
	nanomsg.org/go-mangos v1.4.0 // indirect
)

replace server => /run/sandbox/sift/server
