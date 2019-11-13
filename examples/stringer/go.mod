module github.com/lukasjarosch/genki/examples/stringer

go 1.13

replace github.com/lukasjarosch/genki => ../../

require (
	github.com/gogo/protobuf v1.2.1
	github.com/golang/protobuf v1.3.2
	github.com/lukasjarosch/genki v0.0.0-20191103174705-04f7563b417a
	github.com/maxbrunsfeld/counterfeiter/v6 v6.2.2 // indirect
	github.com/mokiat/gostub v1.3.0 // indirect
	github.com/prometheus/client_golang v1.2.1
	github.com/streadway/amqp v0.0.0-20190827072141-edfb9018d271
	google.golang.org/grpc v1.24.0
	gopkg.in/urfave/cli.v1 v1.20.0 // indirect
)
