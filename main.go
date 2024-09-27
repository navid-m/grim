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
	// fmt.Printf("Project Name: %s\n", projectInfo.ProjectName)
	// fmt.Printf("Tempo: %.2f\n", projectInfo.Tempo)
	// fmt.Printf("Sample Rate: %d\n", projectInfo.SampleRate)
	// fmt.Printf("Tracks: %d\n", projectInfo.Tracks)

	//rpp.Dump(project, os.Stdout)
}
