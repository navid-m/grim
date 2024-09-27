package rpp

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/rodaine/table"
)

// The FX chain itself
type FXChain struct {
	Instances []FXInstance
}

// An FX chain item
type FXInstance struct {
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
	Number   int
	GUID     string
	FXChains []FXChain
}

// Bunch of tracks
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

// Stringer implementation for FXChain
func (f FXInstance) String() string {
	return fmt.Sprintf("\t- Preset: %s, FXID: %s, VST: %s, WndRect: %v, Show: %d, LastSel: %d, Docked: %d, Bypass: %v",
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

// Stringer implementation for FXChain
func (t FXChain) String() string {
	res := ""
	for _, x := range t.Instances {
		res += x.String()
	}
	return res
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

// Stringer implementation for ProjectInfo
func (p ProjectInfo) String() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Project Name: %s\n", p.ProjectName))
	sb.WriteString(fmt.Sprintf("Original Platform: %s\n", p.OriginalPlatform))
	sb.WriteString(fmt.Sprintf("Tempo: %.2f\n", p.Tempo))
	sb.WriteString(fmt.Sprintf("Loop enabled: %s\n", strconv.FormatBool(p.LoopEnabled)))
	sb.WriteString(fmt.Sprintf("Sample Rate: %d\n", p.SampleRate))
	sb.WriteString(fmt.Sprintf("Tracks: %d\n", len(p.Tracks.TrackList)))
	sb.WriteString(fmt.Sprintf("FX Chains: %d\n", len(p.FXChains)))

	if len(p.FXChains) != 0 {
		for c, fx := range p.FXChains {
			sb.WriteString("FX Chain #" + strconv.Itoa(c+1) + ":\n")
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

// Tabler implementation for ProjectInfo
func (p ProjectInfo) AsTable() table.Table {
	tbl := table.New("Trait", "Value")

	headerFmt := color.New(color.FgHiBlack, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgGreen).SprintfFunc()
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	tbl.AddRow("Project Name", p.ProjectName)
	tbl.AddRow("Original Platform", p.OriginalPlatform)
	tbl.AddRow("Tempo", p.Tempo)
	tbl.AddRow("Loop enabled", strconv.FormatBool(p.LoopEnabled))
	tbl.AddRow("Sample Rate", p.SampleRate)
	tbl.AddRow("Tracks", len(p.Tracks.TrackList))
	tbl.AddRow("FX Chains", len(p.FXChains))

	for currentIter, chain := range p.FXChains {
		oneIndexedCurrent := currentIter + 1

		for instanceIter, instance := range chain.Instances {
			if strings.TrimSpace(instance.Vst) == "" {
				continue
			}
			tbl.AddRow(fmt.Sprintf("FXC %d I %d VST Name", oneIndexedCurrent, instanceIter), instance.Vst)
			if strings.TrimSpace(instance.PresetName) != "" {
				tbl.AddRow(fmt.Sprintf("FXC %d I %d Preset Name: %s", oneIndexedCurrent, instanceIter, instance.PresetName))
				tbl.AddRow(fmt.Sprintf("FXC %d I %d ID: %s", oneIndexedCurrent, instanceIter, instance.FxId))
				tbl.AddRow(fmt.Sprintf("FXC %d I %d Bypass", oneIndexedCurrent, instanceIter), instance.Bypass)
				tbl.AddRow(fmt.Sprintf("FXC %d I %d Docked", oneIndexedCurrent, instanceIter), instance.Docked)
				tbl.AddRow(fmt.Sprintf("FXC %d I %d Last Selection", oneIndexedCurrent, instanceIter), instance.LastSel)
			}
		}
	}

	if len(p.Items) != 0 {
		for _, item := range p.Items {
			tbl.AddRow("Items", item.String())
		}
	}

	return tbl
}
