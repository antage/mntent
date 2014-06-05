# mntent

`mntent` is mtab/fstab parser for Go.

## Requirements

`mntent` requires Go 1.1 or above.

## Installation

```
go get github.com/antage/mntent
```

## Usage

```go
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/antage/mntent"
)

func main() {
	entries, err := mntent.Parse("/etc/fstab")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't open /etc/fstab: %s\n", err)
		os.Exit(1)
	}

	for _, ent := range entries {
		fmt.Printf("FS name: %s\n", ent.Name)
		fmt.Printf("\tDir: %s\n", ent.Directory)
		fmt.Printf("\tTypes: %v\n", ent.Types)
		fmt.Printf("\tOptions: %s\n", strings.Join(ent.Options, ","))
		fmt.Printf("\tDump frequency: %d\n", ent.DumpFrequency)
		fmt.Printf("\tFSCK pass number: %d\n", ent.PassNumber)
	}
}
```
