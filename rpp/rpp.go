package rpp

import (
	"io"
	"os"
	"strings"
)

// Load loads and parses an RPP file from the given path
func Load(filePath string) (*Element, error) {
	file, err := os.Open(filePath)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	toret, err := LoadFromReader(file)
	toret.RootFileName = filePath

	return toret, err
}

// LoadFromReader parses an RPP file from an io.Reader
func LoadFromReader(r io.Reader) (*Element, error) {
	content, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	parser := NewParser(string(content))
	return parser.Parse()
}

// Dump writes the parsed Element tree to a writer in a human-readable format
func Dump(element *Element, writer io.Writer) error {
	_, err := writer.Write([]byte(elementToString(element, 0)))
	return err
}

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
