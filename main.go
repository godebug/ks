package main

import (
	"github.com/godebug/ks/history"
	"log"
)

func main() {
	h, err := history.NewHistory("n.csv")
	if err != nil {
		log.Fatal("Load history: " + err.Error())
	}
	print(h)
}
