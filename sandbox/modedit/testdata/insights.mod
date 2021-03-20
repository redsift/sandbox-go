module server

go 1.15

require (
	github.com/blevesearch/bleve v1.0.14
	github.com/blevesearch/blevex v1.0.0
	github.com/cenkalti/backoff/v4 v4.1.0
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/redsift/go-sandbox-rpc v0.0.0-20200511193908-a6e5d4529bf7
	github.com/redsift/go-stats v0.2.0
	github.com/stretchr/testify v1.6.1
	golang.org/x/sys v0.0.0-20201223074533-0d417f636930 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
	github.com/redsift/go-sandbox-rpc v0.1.0 // indirect
)

replace (
	github.com/blevesearch/bleve v1.0.14 => github.com/redsift/bleve v0.5.1-0.20201222154652-c8a1b8b0f852
	github.com/tecbot/gorocksdb => github.com/redsift/gorocksdb v0.0.0-20180109115255-d1d69065a9b9
)
