module github.com/redsift/sandbox-go

go 1.15

require (
	github.com/redsift/go-mangosock v0.1.2
	github.com/redsift/go-sandbox-rpc v0.1.0
	golang.org/x/mod v0.4.0
	server v1.0.0
)

replace server => /run/sandbox/sift/server

replace github.com/blevesearch/bleve v1.0.14 => github.com/redsift/bleve v0.5.1-0.20201222154652-c8a1b8b0f852

replace github.com/tecbot/gorocksdb => github.com/redsift/gorocksdb v0.0.0-20180109115255-d1d69065a9b9
