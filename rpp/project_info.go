package rpp

import (
	"fmt"
	"strconv"
	"strings"
)

// An FX chain item
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

// A track item
type Item struct {
	Position float64
	Length   float64
	Loop     int
	Name     string
	Source   Source
	Guid     string
	Playrate float64
}

// (FLAC, WAV, etc.)
type Source struct {
	Type string
	File string
}

// Some track in the playlist
type Track struct {
	Number int
	GUID   string
}

type Tracks struct {
	TrackList []Track
}

// General project information
type ProjectInfo struct {
	ProjectName      string
	OriginalPlatform string
	Tempo            float64
	LoopEnabled      bool
	SampleRate       int
	Tracks           Tracks
	Items            []Item
	FXChains         []FXChain
}

func cleanProjectName(s string) string {
	initName := ""
	if strings.LastIndex(s, "/") == -1 {
		initName = s
	} else {
		initName = s[strings.LastIndex(s, "/")+1:]
	}
	return strings.Replace(initName, ".rpp", "", -1)
}

// Get project information from the Element tree
func ParseProjectInfo(element *Element) ProjectInfo {
	info := ProjectInfo{
		ProjectName:      cleanProjectName(element.RootFileName),
		OriginalPlatform: "Unknown",
		Tempo:            120.0,
		LoopEnabled:      false,
		SampleRate:       44100,
		Tracks:           Tracks{},
		Items:            []Item{},
		FXChains:         []FXChain{},
	}

	if strings.Contains(element.Tag, "win64") {
		info.OriginalPlatform = "Windows (win64)"
	} else if strings.Contains(element.Tag, "darwin") {
		info.OriginalPlatform = "Mac OS X (darwin)"
	} else {
		info.OriginalPlatform = "Linux"
	}

	for _, attr := range element.Attrib {
		if strings.HasPrefix(attr, "TEMPO") {
			info.Tempo = parseTempo(attr)
		} else if strings.HasPrefix(attr, "SAMPLERATE") {
			info.SampleRate = parseSampleRate(attr)
		}
	}

	trackCount := 0

	for _, child := range element.Children {
		switch {
		case child.Tag == "ITEM":
			item := parseItem(child)
			info.Items = append(info.Items, item)
		case child.Tag == "FXCHAIN":
			fxChain := parseFXChain(child)
			info.FXChains = append(info.FXChains, fxChain)
		case strings.HasPrefix(child.Tag, "TRACK"):
			trackCount++
			info.Tracks.TrackList = append(
				info.Tracks.TrackList, Track{
					Number: trackCount,
					GUID: strings.Replace(
						child.Tag,
						"TRACK ",
						"",
						-1,
					),
				},
			)
		}
	}
	return info
}

// Parse tempo from the tempo attribute
func parseTempo(attr string) float64 {
	parts := strings.Fields(attr)
	if len(parts) > 1 {
		if tempo, err := strconv.ParseFloat(parts[1], 64); err == nil {
			return tempo
		}
	}
	return 120.0
}

// Parse sample rate from the SAMPLERATE attribute
func parseSampleRate(attr string) int {
	parts := strings.Fields(attr)
	if len(parts) > 1 {
		if rate, err := strconv.Atoi(parts[1]); err == nil {
			return rate
		}
	}
	return 44100
}

// Parse items within the project
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
	for _, child := range element.Children {
		if child.Tag == "SOURCE" {
			item.Source = parseSource(child)
		}
	}
	return item
}

// Parse sources for items
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

// Parse FX chain information
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

// Parse float from an attribute string (e.g., "POSITION 1.234")
func parseFloat(attr string) float64 {
	parts := strings.Fields(attr)
	if len(parts) > 1 {
		if value, err := strconv.ParseFloat(parts[1], 64); err == nil {
			return value
		}
	}
	return 0.0
}

// Parse GUID from an attribute string
func parseGUID(attr string) string {
	return strings.TrimSpace(attr)
}

// Parse NAME from an attribute string
func parseName(attr string) string {
	parts := strings.Split(attr, " ")
	if len(parts) > 1 {
		return strings.Join(parts[1:], " ")
	}
	return ""
}

// Parse FILE path from an attribute string
func parseFile(attr string) string {
	parts := strings.Split(attr, " ")
	if len(parts) > 1 {
		return parts[1]
	}
	return ""
}

// Parse window rectangle for FX chain
func parseWndRect(attr string) [4]int {
	parts := strings.Fields(attr)
	if len(parts) >= 5 {
		x1, _ := strconv.Atoi(parts[1])
		y1, _ := strconv.Atoi(parts[2])
		x2, _ := strconv.Atoi(parts[3])
		y2, _ := strconv.Atoi(parts[4])
		return [4]int{x1, y1, x2, y2}
	}
	return [4]int{0, 0, 0, 0}
}

// Parse FXID from an attribute string
func parseFxId(attr string) string {
	return strings.TrimSpace(attr)
}

// Parse preset name from an attribute string
func parsePresetName(attr string) string {
	parts := strings.Split(attr, " ")
	if len(parts) > 1 {
		return strings.Join(parts[1:], " ")
	}
	return ""
}

// Stringer implementation for ProjectInfo
func (p ProjectInfo) String() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Project Name: %s\n", p.ProjectName))
	sb.WriteString(fmt.Sprintf("Original Platform: %s\n", p.OriginalPlatform))
	sb.WriteString(fmt.Sprintf("Tempo: %.2f\n", p.Tempo))
	sb.WriteString(fmt.Sprintf("Loop enabled: %s\n", strconv.FormatBool(p.LoopEnabled)))
	sb.WriteString(fmt.Sprintf("Sample Rate: %d\n", p.SampleRate))
	sb.WriteString(fmt.Sprintf("Tracks: %d\n", len(p.Tracks.TrackList)))

	if len(p.FXChains) != 0 {
		sb.WriteString("FX Chains:\n")
		for _, fx := range p.FXChains {
			sb.WriteString(fx.String())
		}
	}

	if len(p.Items) != 0 {
		sb.WriteString("Items:\n")
		for _, item := range p.Items {
			sb.WriteString(item.String())
		}

	}

	return strings.TrimSuffix(sb.String(), "\n")
}

// Stringer implementation for FXChain
func (f FXChain) String() string {
	return fmt.Sprintf("Preset: %s, FXID: %s, VST: %s, WndRect: %v, Show: %d, LastSel: %d, Docked: %d, Bypass: %v",
		f.PresetName, f.FxId, f.Vst, f.WndRect, f.Show, f.LastSel, f.Docked, f.Bypass)
}

// Stringer implementation for Item
func (i Item) String() string {
	return fmt.Sprintf("Name: %s, Position: %.2f, Length: %.2f, Playrate: %.2f, Source: %s, GUID: %s",
		i.Name, i.Position, i.Length, i.Playrate, i.Source.File, i.Guid)
}

// Stringer implementation for Track
func (t Track) String() string {
	return fmt.Sprintf("No: %d, GUID: %s", t.Number, t.GUID)
}

// Stringer implementation for Tracks
func (t Tracks) String() string {
	sep := strings.Repeat("-", 10)
	res := "\n" + sep + "\n"
	for _, x := range t.TrackList {
		res += fmt.Sprintln("Track number: " + strconv.Itoa(x.Number))
		res += fmt.Sprintln("GUID: " + x.GUID)
		res += sep + "\n"
	}
	return res
}

// Stringer implementation for Source
func (s Source) String() string {
	return fmt.Sprintf("Type: %s, File: %s", s.Type, s.File)
}
