package main

import (
	"log"
)

func main() {
	application, err := New()
	if err != nil {
		log.Fatalln(err.Error())
	}

	defer func() {
		if x := recover(); x != nil {
			application.Stop()
			panic(x)
		}
	}()
	err = application.Run()
	if err != nil {
		log.Fatalln(err.Error())
	}
}
