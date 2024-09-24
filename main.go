package main

import (
	"fmt"
	"grim/rpp"
	"log"
	"os"
)

func main() {
	project, err := rpp.Load("example.rpp")
	if err != nil {
		log.Fatal(err)
	}

	projectInfo := rpp.ParseProjectInfo(project)

	fmt.Printf("Project Name: %s\n", projectInfo.ProjectName)
	fmt.Printf("Tempo: %.2f\n", projectInfo.Tempo)
	fmt.Printf("Sample Rate: %d\n", projectInfo.SampleRate)
	fmt.Printf("Tracks: %d\n", projectInfo.Tracks)

	rpp.Dump(project, os.Stdout)
}
