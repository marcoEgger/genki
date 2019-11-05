package interceptoqlCreateTodoItemsTable

import "google.golang.org/grpc" = ` CREATE TABLE IF NOT EXISTS todo_items ( id bigserial PRIMARY KEY, title text NOT NULL, 	due bigint, 	done boolean ) WITH ( 	OIDS=FALSE )`
    sqlSelectTodoItemCount = ` SELECT COUNT(id) FROM todo_items`
    sqlSelectAllTodoItems = ` SELECT id 	, title 	, due 	, done FROM todo_items`
    sqlInsertTodoItem = ` INSERT INTO todo_items (title, due, done) VALUES ($1, $2, $3)`
    sqlUpdateTodoItem = ` UPDATE todo_items SET title = $2 , due = $3 , done = $4 WHERE id = $1`
    sqlDeleteTodoItem = ` DELETE FROM todo_items WHERE id = $1` )r

import (
	"context"
	"time"

	"google.golang.org/grpc"

	"github.com/lukasjarosch/genki/logger"
)

func UnaryServerLogging() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		log := logger.WithMetadata(ctx)
		log.Infof("incoming unary request to '%s'", info.FullMethod)
		defer func(started time.Time) {
			log.Infof("finished unary request to '%s' (took %v)", info.FullMethod, time.Since(started))
		}(time.Now())

		return handler(ctx, req)
	}
}

func UnaryClientLogging() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
		log := logger.WithMetadata(ctx)
		log.Infof("client call '%s' to server '%s'", method, cc.Target())
		defer func(started time.Time) {
			if err != nil {
				log.Infof("client call to '%s' (server=%s) failed (took %v): %s", method, cc.Target(), time.Since(started), err)
			} else {
				log.Infof("client request to '%s' was successfully handled by server '%s' (took %v)", method, cc.Target(), time.Since(started))
			}
		}(time.Now())

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
