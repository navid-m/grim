package main

import (
	"fmt"
	"grim/rpp"
	"log"
)

func main() {
	project, err := rpp.Load("data/test.rpp")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(rpp.ParseProjectInfo(project))
}
