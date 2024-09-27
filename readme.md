# Grim

Library for parsing REAPER `.rpp` project files (from the Reaper DAW).
Like [this](https://github.com/Perlence/rpp) but as a useful library.

## Features

-  Extract information like project settings, tracks, items, and FX chains from `.rpp` files.
-  Parse REAPER project info (name, tempo, sample rate, etc).
-  Extract track and item details.
-  Access FX chains and plugin information.

## Installation

```bash
go get github.com/navid-m/grim
```

## Usage example

```go
package main

import (
    "fmt"
    "os"
    "github.com/navid-m/grim/rpp"
)

func main() {
    file, err := os.Open("some_project.rpp")
    if err != nil {
        panic(err)
    }
    defer file.Close()

    element, err := rpp.LoadFromReader(file)
    if err != nil {
        panic(err)
    }

    project := rpp.ParseProjectInfo(element)
    fmt.Println("Project Name:", project.ProjectName)
    fmt.Println("Tempo:", project.Tempo)
    fmt.Println("Tracks:", len(project.Tracks.TrackList))
}
```

## Running Tests

To run the tests:

```bash
go test ./...
```

## License

This is licensed under the GPLv3 License.
