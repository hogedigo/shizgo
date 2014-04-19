package main

import (
	"github.com/go-martini/martini"
)

func main() {
	m := martini.Classic()
	m.Get("/hello", func() string {
			return "Hello world!"
		})
	m.Run()
}

