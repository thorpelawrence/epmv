package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	anyascii "github.com/anyascii/go"
	flag "github.com/spf13/pflag"
	"github.com/taylorskalyo/goreader/epub"
)

var (
	list     = flag.BoolP("list", "l", false, "list file metadata, without renaming")
	listJSON = flag.BoolP("json", "j", false, "list metadata in JSONL format, requires --list")
	ascii    = flag.BoolP("ascii", "a", false, "transliterate file names to ASCII")
	format   = flag.StringP("format", "f", "{{.Title}} - {{.Creator}}", "format string for output file name, .epub will be ignored")
	dry      = flag.BoolP("dry", "n", false, "dry run, only print renames")
)

func main() {
	flag.Parse()

	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	files := flag.Args()

	if *listJSON && !*list {
		return errors.New("cannot use --json without --list")
	}

	t, err := template.New("name").Parse(trimSuffixEPUB(*format))
	if err != nil {
		return fmt.Errorf("parsing format string: %w", err)
	}

	for _, file := range files {
		if err := processFile(t, file); err != nil {
			return fmt.Errorf("processing file '%s': %w", file, err)
		}
	}

	return nil
}

func processFile(t *template.Template, file string) error {
	rc, err := epub.OpenReader(file)
	if err != nil {
		return fmt.Errorf("loading file: %w", err)
	}
	defer rc.Close()

	if len(rc.Rootfiles) == 0 || rc.Rootfiles[0] == nil {
		return errors.New("reading metadata: no root files")
	}
	md := rc.Rootfiles[0].Package.Metadata

	if *list {
		path, err := filepath.Rel(".", file)
		if err != nil {
			// non-critical, just print the basename
			path = filepath.Base(file)
		}

		if *listJSON {
			if err := json.NewEncoder(os.Stdout).Encode(map[string]any{
				"path":     path,
				"metadata": md,
			}); err != nil {
				return fmt.Errorf("encoding json: %w", err)
			}
			return nil
		}

		fmt.Fprintf(os.Stdout, "%-10s| %s\n", "path", path)
		j, err := json.MarshalIndent(md, "          | ", "  ")
		if err != nil {
			return fmt.Errorf("marshaling json: %w", err)
		}
		fmt.Fprintf(os.Stdout, "%-10s| %s\n----------+\n", "metadata", string(j))

		return nil
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, md); err != nil {
		return fmt.Errorf("executing template: %w", err)
	}
	base := buf.String() + ".epub"

	if *ascii && !isASCII(base) {
		fmt.Fprintln(os.Stderr, "info: transliterating non-ASCII:", base)
		base = anyascii.Transliterate(base)
	}

	dir := filepath.Dir(file)
	renamed := filepath.Join(dir, base)

	if *dry {
		fmt.Fprintf(os.Stderr, "would rename '%s' → '%s'\n", file, renamed)
		return nil
	}

	if file == renamed {
		return nil
	}

	if fileExists(renamed) {
		return fmt.Errorf("file already exists: %s", renamed)
	}

	if err := os.Rename(file, renamed); err != nil {
		return fmt.Errorf("renaming file: %w", err)
	}

	return nil
}
