package datastore

import "github.com/lukasjarosch/genki/examples/stringer/internal/models"

type inMem struct {

}

func NewInMem() *inMem {
	return &inMem{}
}

func (db *inMem) FindGreetingByName(name string) (*models.Greeting, error)  {
	return &models.Greeting{
		Template: "Whoop whopp, %s",
		Name:     name,
	}, nil
}
