package models

type GreetingRepository interface {
	FindGreetingByName(name string) (*Greeting, error)
}