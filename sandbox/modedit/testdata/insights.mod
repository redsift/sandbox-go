module server

go 1.15

require (
	github.com/redsift/ccc v0.1.0 // indirect
)

replace (
	github.com/blevesearch/bleve v1.0.14 => github.com/redsift/bleve v0.5.1-0.20201222154652-c8a1b8b0f852
	github.com/tecbot/gorocksdb => github.com/redsift/gorocksdb v0.0.0-20180109115255-d1d69065a9b9
)
