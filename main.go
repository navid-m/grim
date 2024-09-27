package main

import (
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
	rpp.ParseProjectInfo(project).AsTable().Print()
}
