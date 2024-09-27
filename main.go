package main

import (
	"fmt"
	"grim/rpp"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Missing project path argument")
	}
	project, err := rpp.Load(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(rpp.ParseProjectInfo(project))
}
