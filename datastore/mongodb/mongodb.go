package mongodb

import (
	"context"
	"github.com/marcoEgger/genki/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
	"time"
)

type MongoDB struct {
	db  *mongo.Client
	uri string
}

// New will connect to the MongoDB server using the given URI
//
//goland:noinspection GoUnusedExportedFunction
func New(uri string) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db, err := mongo.Connect(ctx, options.Client().SetMonitor(otelmongo.NewMonitor()).ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	logger.Info("connected to mongodb")

	return &MongoDB{
		db:  db,
		uri: uri,
	}, nil
}

// Close is just a proxy for convenient access to db.Close()
func (m MongoDB) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	logger.Info("disconnecting from mongodb")
	return m.db.Disconnect(ctx)
}

// DB is just a proxy for convenient access to the underlying mongodb client implementation
// This method is used a lot, therefore it's name is abbreviated.
func (m MongoDB) DB() *mongo.Client {
	return m.db
}
