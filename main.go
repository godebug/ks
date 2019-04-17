package main

import (
	"github.com/godebug/ks/answer"
	"log"
)

func main() {
	var a answer.Answer
	err := a.Answer("n.csv")
	if err != nil {
		log.Fatal(err.Error())
	}
	err = a.Serve(":8000")
	if err != nil {
		log.Fatal(err.Error())
	}
}
