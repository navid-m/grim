package rpp

import (
	"strconv"
	"strings"
)

// ProjectInfo contains general information about the project
type ProjectInfo struct {
	ProjectName  string
	Tempo        float64
	SampleRate   int
	Tracks       int
	IsFullFormat bool
}

// ParseProjectInfo extracts general project information from the Element tree
func ParseProjectInfo(element *Element) ProjectInfo {
	info := ProjectInfo{
		ProjectName:  "Unknown",
		Tempo:        120.0, // default
		SampleRate:   44100, // default
		Tracks:       0,
		IsFullFormat: false,
	}

	// Traverse the element tree and extract information
	for _, child := range element.Children {
		switch child.Tag {
		case "REAPER_PROJECT":
			for _, attr := range child.Attrib {
				// Example of extracting some info from REAPER_PROJECT tag
				if strings.HasPrefix(attr, "TEMPO") {
					info.Tempo = parseTempo(attr)
				}
			}
		case "TRACK":
			info.Tracks++
		}
	}

	return info
}

// parseTempo extracts the tempo from the attribute string (simple example)
func parseTempo(attr string) float64 {
	// Example: "TEMPO 120.000"
	parts := strings.Split(attr, " ")
	if len(parts) == 2 {
		if tempo, err := strconv.ParseFloat(parts[1], 64); err == nil {
			return tempo
		}
	}
	return 120.0 // default tempo
}
