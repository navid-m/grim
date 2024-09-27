package rpp

import (
	"io"
	"log"
	"os"
	"strings"
)

// Load and parse an RPP file from a given path
func Load(filePath string) (*Element, error) {
	if !strings.Contains(strings.ToLower(filePath), ".rpp") {
		log.Fatalf("Invalid file passed. Pass a valid RPP file.")
	}

	file, err := os.Open(filePath)

	if err != nil {
		return nil, err
	}

	defer file.Close()
	toret, err := LoadFromReader(file)
	toret.RootFileName = filePath

	return toret, err
}

// Parses an RPP file from a reader
func LoadFromReader(r io.Reader) (*Element, error) {
	content, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	parser := NewParser(string(content))
	return parser.Parse()
}

// Write the parsed Element tree to a writer
func Dump(element *Element, writer io.Writer) error {
	_, err := writer.Write([]byte(elementToString(element, 0)))
	return err
}

// Converts element to string
func elementToString(e *Element, indent int) string {
	result := strings.Repeat(" ", indent*2) + "<" + e.Tag + "\n"
	for _, attr := range e.Attrib {
		result += strings.Repeat(" ", (indent+1)*2) + attr + "\n"
	}
	for _, child := range e.Children {
		result += elementToString(child, indent+1)
	}
	result += strings.Repeat(" ", indent*2) + ">\n"
	return result
}
