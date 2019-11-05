package datastore

import (
	"database/sql"

	"github.com/lukasjarosch/genki/examples/stringer/internal/models"
	"github.com/lukasjarosch/genki/logger"
)

type DbGreeting struct {
	*models.Greeting
}
func (greeting *DbGreeting) scan(rows *sql.Rows)  {
	err := rows.Scan(
		&greeting.Name,
		&greeting.Template,
	)
	if err != nil {
	    logger.Warnf("failed to scan greeting from mysql database: %s", err)
	}
}

type DbGreetings []*DbGreeting
func (greetings *DbGreetings) scan(rows *sql.Rows)  {
	for rows.Next() {
		greeting := &DbGreeting{&models.Greeting{}}
		greeting.scan(rows)
		*greetings = append(*greetings, greeting)
	}
}

type mySqlRepository struct {
	db *sql.DB
}

func NewMySQL() *mySqlRepository {
	// TODO connect
	return &mySqlRepository{}
}

func (repo *mySqlRepository) FindGreetingByName(name string) (*models.Greeting, error)  {
	greeting := &DbGreeting{}
	rows, err := repo.db.Query(sqlFindGreeting, name)
	if err != nil {
	    return nil, err
	}
	greeting.scan(rows)
	return greeting.Greeting, nil
}
