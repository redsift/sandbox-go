module github.com/redsift/sandbox-go

go 1.16

require (
	github.com/Microsoft/go-winio v0.4.16 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-mangos/mangos v1.2.1-0.20171226193431-e5e6478f44cd // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/redsift/go-mangosock v0.2.1
	github.com/redsift/go-sandbox-rpc v0.1.0
	github.com/stretchr/testify v1.2.2
	golang.org/x/mod v0.4.0
	golang.org/x/sys v0.0.0-20210320140829-1e4c9ba3b0c4 // indirect
)

replace server => /run/sandbox/sift/server
