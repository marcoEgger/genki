module github.com/marcoEgger/genki/examples/stringer

go 1.13

replace github.com/marcoEgger/genki => ../../

require (
	github.com/gogo/protobuf v1.3.2
	github.com/golang/protobuf v1.5.2
	github.com/marcoEgger/genki v0.0.0-20191103174705-04f7563b417a
	github.com/maxbrunsfeld/counterfeiter/v6 v6.2.2 // indirect
	github.com/mokiat/gostub v1.3.0 // indirect
	github.com/prometheus/client_golang v1.12.1
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/streadway/amqp v1.0.0
	google.golang.org/grpc v1.44.0
	gopkg.in/urfave/cli.v1 v1.20.0 // indirect
)
