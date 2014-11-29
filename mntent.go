package mntent

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Entry struct {
	Name          string
	Directory     string
	Types         []string
	Options       []string
	DumpFrequency int
	PassNumber    int
}

func Parse(filename string) ([]*Entry, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buf := bufio.NewReader(file)

	entries := make([]*Entry, 0, 4)

	for {
		line, rserr := buf.ReadString('\n')
		if rserr != nil && rserr != io.EOF {
			return nil, rserr
		}
		entry, err := parseLine(line)
		if err != nil {
			return nil, err
		}
		if entry != nil {
			entries = append(entries, entry)
		}

		if rserr == io.EOF {
			break
		}
	}

	return entries, nil
}

func unescape(s string) string {
	s = strings.Replace(s, "\\011", "\t", -1)
	s = strings.Replace(s, "\\012", "\n", -1)
	s = strings.Replace(s, "\\040", " ", -1)
	s = strings.Replace(s, "\\134", "\\", -1)
	return s
}

var splitRegExp = regexp.MustCompile("\\s+")

func parseLine(untrimmedLine string) (*Entry, error) {
	line := strings.TrimSpace(untrimmedLine)
	if len(line) == 0 {
		return nil, nil
	}
	if strings.HasPrefix(line, "#") {
		return nil, nil
	}

	fields := splitRegExp.Split(line, -1)
	if len(fields) != 6 {
		return nil, errors.New(fmt.Sprintf("Each line must consist 6 fields but got %d", len(fields)))
	}

	entry := &Entry{}
	entry.Name = unescape(fields[0])
	entry.Directory = unescape(fields[1])
	entry.Types = strings.Split(unescape(fields[2]), ",")
	entry.Options = strings.Split(unescape(fields[3]), ",")

	num, err := strconv.ParseUint(fields[4], 10, 31)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Can't parse dump frequency field: %s", err))
	}
	entry.DumpFrequency = int(num)

	num, err = strconv.ParseUint(fields[5], 10, 31)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Can't parse pass number field: %s", err))
	}
	entry.PassNumber = int(num)

	return entry, nil
}
