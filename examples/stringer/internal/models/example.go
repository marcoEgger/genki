package models

import (
	"errors"
	"fmt"
)

var (
	ErrGreetingEmptyName = errors.New("greeting.name cannot be empty")
	ErrGetHelloTemplate  = errors.New("cannot get hello template for name")
)

type Greeting struct {
	Template string
	Name     string
	Rendered string
}

func (g *Greeting) Render() string {
	return fmt.Sprintf(g.Template, g.Name)
}

func (g *Greeting) Validate() error {
	if g.Name == "" {
		return ErrGreetingEmptyName
	}

	return nil
}
