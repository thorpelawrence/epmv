package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"text/template"

	anyascii "github.com/anyascii/go"
	"github.com/gen2brain/go-fitz"
	flag "github.com/spf13/pflag"
)

var (
	list   = flag.BoolP("list", "l", false, "list file metadata, without renaming")
	ascii  = flag.BoolP("ascii", "a", false, "transliterate file names to ASCII")
	format = flag.StringP("format", "f", "{{.title}} - {{.author}}", "format string for output file name, .epub will be ignored")
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
	doc, err := fitz.New(file)
	if err != nil {
		return fmt.Errorf("loading file: %w", err)
	}

	md := doc.Metadata()

	if *list {
		name, err := filepath.Rel(".", file)
		if err != nil {
			// non-critical, just print the basename
			name = filepath.Base(file)
		}
		fmt.Fprintf(
			os.Stdout,
			"----------------------\n%s\n----------------------\n",
			name,
		)

		keys := make([]string, 0, len(md))
		for k := range md {
			keys = append(keys, k)
		}
		slices.Sort(keys)

		for _, k := range keys {
			fmt.Fprintf(os.Stdout, "%-20s | %s\n", k, trimNul(md[k]))
		}

		return nil
	}

	for k, v := range md {
		md[k] = trimNul(v)
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

	if err := os.Rename(file, filepath.Join(dir, base)); err != nil {
		return fmt.Errorf("renaming file: %w", file, err)
	}

	return nil
}
