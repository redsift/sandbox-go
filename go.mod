module github.com/redsift/sandbox-go

go 1.23

require (
	github.com/redsift/go-mangosock v0.2.1
	github.com/redsift/go-sandbox-rpc v0.2.0
	github.com/stretchr/testify v1.2.2
	golang.org/x/mod v0.20.0
)

require (
	github.com/Microsoft/go-winio v0.4.16 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/sys v0.0.0-20210320140829-1e4c9ba3b0c4 // indirect
	golang.org/x/xerrors v0.0.0-20191011141410-1b5146add898 // indirect
	nanomsg.org/go-mangos v1.4.0 // indirect
)

replace server => /run/sandbox/sift/server
