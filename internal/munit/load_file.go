package munit

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// LoadFile loads units from a file
func LoadFile(filename string) (units []Unit, err error) {
	var f *os.File
	if f, err = os.Open(filename); err != nil {
		err = fmt.Errorf("failed to open unit file %s: %w", filename, err)
		return
	}
	defer f.Close()

	dec := yaml.NewDecoder(f)
	docNum := 0
	for {
		var unit Unit
		docNum++
		if err = dec.Decode(&unit); err != nil {
			if err == io.EOF {
				err = nil
			} else {
				// Provide detailed context: file path, document number, and underlying error
				err = fmt.Errorf("failed to decode unit file %s (document %d): %w", filename, docNum, err)
			}
			return
		}

		if unit.Kind == "" {
			// Skip empty documents
			continue
		}

		units = append(units, unit)
	}
}

// LoadDir loads units from a directory
func LoadDir(dir string) (units []Unit, err error) {
	for _, ext := range []string{"*.yml", "*.yaml"} {
		var files []string
		if files, err = filepath.Glob(filepath.Join(dir, ext)); err != nil {
			err = fmt.Errorf("failed to glob directory %s with pattern %s: %w", dir, ext, err)
			return
		}
		sort.Strings(files)
		for _, file := range files {
			var _units []Unit
			if _units, err = LoadFile(file); err != nil {
				// Error already has context from LoadFile, just return it
				return
			}
			units = append(units, _units...)
		}
	}
	return
}

const (
	unitDirNone = "none"
)

func ParseUnitDirPattern(pattern string) (dirs []string) {
outerLoop:
	for _, dir := range strings.Split(pattern, ":") {
		dir = strings.TrimSpace(dir)

		if dir == "" {
			continue
		}

		if dir == unitDirNone {
			continue
		}

		if err := os.MkdirAll(dir, 0755); err != nil {
			continue
		}

		for _, existed := range dirs {
			if existed == dir {
				continue outerLoop
			}
		}

		dirs = append(dirs, dir)
	}

	return
}
