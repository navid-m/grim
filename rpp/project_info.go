package rpp

import (
	"strconv"
	"strings"
)

const defaultTempo float64 = 100
const defaultSampleRate int = 44100
const defaultPlatform string = "Unknown"

// Get project information from the Element tree
func ParseProjectInfo(element *Element) ProjectInfo {
	info := ProjectInfo{
		ProjectName:      cleanProjectName(element.RootFileName),
		OriginalPlatform: defaultPlatform,
		Tempo:            defaultTempo,
		LoopEnabled:      false,
		SampleRate:       defaultSampleRate,
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
		case strings.HasPrefix(child.Tag, "TRACK"):
			trackCount++
			track := Track{
				Number: trackCount,
				GUID:   strings.Replace(child.Tag, "TRACK ", "", -1),
			}

			for _, trackChild := range child.Children {
				if trackChild.Tag == "FXCHAIN" {
					fxChain := parseFXChain(trackChild)
					track.FXChains = append(track.FXChains, fxChain)
					info.FXChains = append(info.FXChains, fxChain)
				}
			}

			info.Tracks.TrackList = append(info.Tracks.TrackList, track)
		}
	}
	return info
}

func parseFXChain(element *Element) FXChain {
	chain := FXChain{}
	for _, child := range element.Children {
		if strings.HasPrefix(child.Tag, "VST") || strings.HasPrefix(child.Tag, "JS") {
			instance := parseFXInstance(child)
			chain.Instances = append(chain.Instances, instance)
		}
	}
	return chain
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

// Parse sample rate from the sample rate attribute
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

// Parse FX instance information
func parseFXInstance(element *Element) FXInstance {
	instance := FXInstance{}

	if pros := parseVst(element.Tag); len(pros) > 1 {
		instance.Vst = pros[1]
		if strings.TrimSpace(instance.Vst) == "" {
			instance.Vst = pros[0]
		}
	}

	for _, attr := range element.Attrib {
		switch {
		case strings.HasPrefix(attr, "WNDRECT"):
			instance.WndRect = parseWndRect(attr)
		case strings.HasPrefix(attr, "PRESETNAME"):
			instance.PresetName = parsePresetName(attr)
		case strings.HasPrefix(attr, "FXID"):
			instance.FxId = parseFxId(attr)
		case strings.HasPrefix(attr, "BYPASS"):
			instance.Bypass = parseBypass(attr)
		}
	}

	return instance
}

func parseVst(tag string) []string {
	start := strings.Index(tag, "\"")
	end := strings.LastIndex(tag, "\"")

	if start != -1 && end != -1 && end > start {
		substring := tag[start : end+1]

		parts := strings.FieldsFunc(substring, func(r rune) bool {
			return r == '"'
		})

		var result []string
		for _, part := range parts {
			if part != "" {
				result = append(result, strings.TrimSpace(part))
			}
		}

		return result
	}
	return nil
}

// Parse the bypass value
func parseBypass(attr string) [3]int {
	parts := strings.Fields(attr)
	var bypass [3]int
	for i := 1; i <= 3 && i < len(parts); i++ {
		bypass[i-1], _ = strconv.Atoi(parts[i])
	}
	return bypass
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

// Parse filepath from an attribute string
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

// Remove non-pure-name things from project name
func cleanProjectName(s string) string {
	initName := ""
	if strings.LastIndex(s, "/") == -1 {
		initName = s
	} else {
		initName = s[strings.LastIndex(s, "/")+1:]
	}
	return strings.Replace(initName, ".rpp", "", -1)
}
