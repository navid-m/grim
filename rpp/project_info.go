package rpp

import (
	"fmt"
	"strconv"
	"strings"
)

// Sub-Struct for FX chain items
type FXChain struct {
	WndRect    [4]int
	Show       int
	LastSel    int
	Docked     int
	Bypass     [3]int
	PresetName string
	FxId       string
	Vst        string
}

// Sub-Struct for Items
type Item struct {
	Position float64
	Length   float64
	Loop     int
	Name     string
	Source   Source
	Guid     string
	Playrate float64
}

// Sub-Struct for Source (FLAC, WAV, etc.)
type Source struct {
	Type string
	File string
}

// ProjectInfo contains general project information
type ProjectInfo struct {
	ProjectName      string
	OriginalPlatform string
	Tempo            float64
	LoopEnabled      bool
	SampleRate       int
	Tracks           int
	Items            []Item
	FXChains         []FXChain
}

// Extract general project information from the Element tree
func ParseProjectInfo(element *Element) ProjectInfo {
	info := ProjectInfo{
		ProjectName:      element.RootFileName,
		OriginalPlatform: "Unknown",
		Tempo:            120.0, // Default tempo value
		LoopEnabled:      false,
		SampleRate:       44100, // Default sample rate
		Tracks:           0,
		Items:            []Item{},
		FXChains:         []FXChain{},
	}

	// Determine the platform from the tag
	if strings.Contains(element.Tag, "win64") {
		info.OriginalPlatform = "Windows (win64)"
	} else if strings.Contains(element.Tag, "darwin") {
		info.OriginalPlatform = "Mac OS X (darwin)"
	} else {
		info.OriginalPlatform = "Linux"
	}

	// Traverse the element tree and extract information
	for _, attr := range element.Attrib {
		if strings.HasPrefix(attr, "TEMPO") {
			info.Tempo = parseTempo(attr) // Parse tempo from the TEMPO attribute
		}
	}

	// Traverse children to gather additional information
	for _, child := range element.Children {
		switch child.Tag {
		case "ITEM":
			item := parseItem(child)
			info.Items = append(info.Items, item)
		case "FXCHAIN":
			fxChain := parseFXChain(child)
			info.FXChains = append(info.FXChains, fxChain)
		}
	}

	return info
}

func parseTempo(attr string) float64 {
	parts := strings.Fields(attr)
	if len(parts) > 1 {
		// The second part should be the tempo value (BPM)
		if tempo, err := strconv.ParseFloat(parts[1], 64); err == nil {
			return tempo
		}
	}
	// Return default tempo if parsing fails
	return 120.0
}

func parseItem(element *Element) Item {
	item := Item{}
	for _, attr := range element.Attrib {
		switch {
		case strings.HasPrefix(attr, "POSITION"):
			item.Position = parseFloat(attr)
		case strings.HasPrefix(attr, "LENGTH"):
			item.Length = parseFloat(attr)
		case strings.HasPrefix(attr, "NAME"):
			item.Name = parseName(attr)
		case strings.HasPrefix(attr, "PLAYRATE"):
			item.Playrate = parseFloat(attr)
		case strings.HasPrefix(attr, "GUID"):
			item.Guid = parseGUID(attr)
		}
	}
	// Parse nested <SOURCE> element
	for _, child := range element.Children {
		if child.Tag == "SOURCE" {
			item.Source = parseSource(child)
		}
	}
	return item
}

// parseSource extracts the source information
func parseSource(element *Element) Source {
	source := Source{}
	for _, attr := range element.Attrib {
		if strings.HasPrefix(attr, "FILE") {
			source.File = parseFile(attr)
		}
	}
	source.Type = element.Tag
	return source
}

// parseFXChain extracts information from the <FXCHAIN> element
func parseFXChain(element *Element) FXChain {
	chain := FXChain{}
	for _, attr := range element.Attrib {
		switch {
		case strings.HasPrefix(attr, "WNDRECT"):
			chain.WndRect = parseWndRect(attr)
		case strings.HasPrefix(attr, "PRESETNAME"):
			chain.PresetName = parsePresetName(attr)
		case strings.HasPrefix(attr, "FXID"):
			chain.FxId = parseFxId(attr)
		}
	}
	return chain
}

// Example parse functions
func parseFloat(attr string) float64 {
	return 0.0
}

func parseGUID(attr string) string {
	// Extract GUID from the attribute
	return attr
}

func parseName(attr string) string {
	// Extract NAME
	return attr
}

func parseFile(attr string) string {
	// Extract FILE path
	return attr
}

func parseWndRect(attr string) [4]int {
	// Example parse for WNDRECT (returns dummy values)
	return [4]int{0, 0, 0, 0}
}

func parseFxId(attr string) string {
	// Extract FXID
	return attr
}

func parsePresetName(attr string) string {
	// Extract preset name
	return attr
}

// Stringer implementation for pretty output
func (p ProjectInfo) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Project Name: %s\n", p.ProjectName))
	sb.WriteString(fmt.Sprintf("Original Platform: %s\n", p.OriginalPlatform))
	sb.WriteString(fmt.Sprintf("Tempo: %.2f\n", p.Tempo))
	sb.WriteString(fmt.Sprintf("Loop enabled: %s\n", strconv.FormatBool(p.LoopEnabled)))
	sb.WriteString(fmt.Sprintf("Sample Rate: %d\n", p.SampleRate))
	sb.WriteString(fmt.Sprintf("Tracks: %d\n", p.Tracks))

	if len(p.FXChains) != 0 {
		sb.WriteString("FX Chains:\n")
		for _, fx := range p.FXChains {
			sb.WriteString(fmt.Sprintf("  Preset: %s, FXID: %s\n", fx.PresetName, fx.FxId))
		}
	}

	if len(p.Items) != 0 {
		sb.WriteString("Items:\n")
		for _, item := range p.Items {
			sb.WriteString(fmt.Sprintf("  Name: %s, Position: %.2f, Length: %.2f, Playrate: %.2f, Source: %s\n",
				item.Name, item.Position, item.Length, item.Playrate, item.Source.File))
		}

	}

	return sb.String()
}
