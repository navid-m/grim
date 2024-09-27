package tests

import (
	"testing"

	"github.com/navid-m/grim/rpp"
)

func TestParseProjectInfo(t *testing.T) {
	element := &rpp.Element{
		RootFileName: "test_project.rpp",
		Tag:          "REAPER_PROJECT",
		Attrib:       []string{"TEMPO 120.00", "SAMPLERATE 48000"},
		Children: []*rpp.Element{
			{
				Tag:    "TRACK 12345678-1234-5678-1234-567812345678",
				Attrib: []string{},
				Children: []*rpp.Element{
					{
						Tag: "FXCHAIN",
						Children: []*rpp.Element{
							{Tag: "VST", Attrib: []string{"PRESETNAME Default", "WNDRECT 0 0 100 100"}},
						},
					},
				},
			},
			{
				Tag:    "ITEM",
				Attrib: []string{"POSITION 1.00", "LENGTH 2.00", "NAME TestItem", "GUID abc123"},
				Children: []*rpp.Element{
					{Tag: "SOURCE", Attrib: []string{"FILE test.wav"}},
				},
			},
		},
	}

	project := rpp.ParseProjectInfo(element)

	if project.ProjectName != "test_project" {
		t.Errorf("Expected ProjectName to be 'test_project', got %s", project.ProjectName)
	}

	if project.Tempo != 120.00 {
		t.Errorf("Expected Tempo to be 120.00, got %f", project.Tempo)
	}

	if project.SampleRate != 48000 {
		t.Errorf("Expected SampleRate to be 48000, got %d", project.SampleRate)
	}

	if len(project.Tracks.TrackList) != 1 {
		t.Errorf("Expected 1 track, got %d", len(project.Tracks.TrackList))
	}

	if len(project.Items) != 1 {
		t.Errorf("Expected 1 item, got %d", len(project.Items))
	}
}

// Test string method for ProjectInfo
func TestProjectInfoString(t *testing.T) {
	projectInfo := rpp.ProjectInfo{
		ProjectName:      "Test Project",
		OriginalPlatform: "Windows (win64)",
		Tempo:            128.0,
		LoopEnabled:      true,
		SampleRate:       44100,
		Tracks: rpp.Tracks{
			TrackList: []rpp.Track{
				{Number: 1, GUID: "1234-5678"},
			},
		},
		Items: []rpp.Item{
			{Name: "Test Item", Position: 1.0, Length: 2.0, Guid: "abcd-efg"},
		},
		FXChains: []rpp.FXChain{
			{Instances: []rpp.FXInstance{
				{PresetName: "Default", FxId: "vst-123", Vst: "Some VST"},
			}},
		},
	}

	expectedOutput := `Project Name: Test Project
Original Platform: Windows (win64)
Tempo: 128.00
Loop enabled: true
Sample Rate: 44100
Tracks: 1
FX Chains: 1
FX Chain #1:
	- Preset: Default, FXID: vst-123, VST: Some VST, WndRect: [0 0 0 0], Show: 0, LastSel: 0, Docked: 0, Bypass: [0 0 0]
Items:
Name: Test Item, Position: 1.00, Length: 2.00, Playrate: 0.00, Source: , GUID: abcd-efg`

	if projectInfo.String() != expectedOutput {
		t.Errorf("Expected:\n%s\nGot:\n%s", expectedOutput, projectInfo.String())
	}
}
