module github.com/lukasjarosch/genki/examples/stringer

go 1.13

replace github.com/lukasjarosch/genki => ../../

require (
	github.com/golang/protobuf v1.3.2
	github.com/lukasjarosch/enki v0.0.0-20191025210149-3fe35c746369 // indirect
	github.com/lukasjarosch/genki v0.0.0-00010101000000-000000000000
	github.com/prometheus/client_golang v1.2.1 // indirect
	github.com/spf13/pflag v1.0.5
	google.golang.org/grpc v1.24.0
)
