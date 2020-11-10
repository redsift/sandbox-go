module github.com/redsift/sandbox-go

go 1.15

require (
	nanomsg.org/go-mangos v1.4.0
	github.com/redsift/go-sandbox-rpc v0.0.0-20190108170927-e56484a1d427
)

replace server => /run/sandbox/sift/server
