package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/output"
	"github.com/yourorg/logslice/internal/parser"
	"github.com/yourorg/logslice/internal/reader"
)

func main() {
	var (
		start      = flag.String("start", "", "start timestamp (inclusive)")
		end        = flag.String("end", "", "end timestamp (inclusive)")
		fields     = flag.String("fields", "", "field filters as key=value pairs, comma-separated")
		format     = flag.String("format", "raw", "output format: raw or json")
		countOnly  = flag.Bool("count", false, "print only the count of matched lines")
		bufferSize = flag.Int("buffer", 0, "read buffer size in bytes (0 = default)")
	)
	flag.Parse()

	args := flag.Args()

	var r *reader.Reader
	var err error
	if len(args) == 0 {
		r, err = reader.NewStdinReader()
	} else {
		opts := []reader.Option{}
		if *bufferSize > 0 {
			opts = append(opts, reader.WithBufferSize(*bufferSize))
		}
		r, err = reader.NewFileReader(args[0], opts...)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening input: %v\n", err)
		os.Exit(1)
	}
	defer r.Close()

	tr, err := filter.NewTimeRange(*start, *end)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid time range: %v\n", err)
		os.Exit(1)
	}

	var fieldFilters []parser.FieldFilter
	if *fields != "" {
		for _, pair := range strings.Split(*fields, ",") {
			ff, err := parser.ParseFieldFilter(strings.TrimSpace(pair))
			if err != nil {
				fmt.Fprintf(os.Stderr, "invalid field filter %q: %v\n", pair, err)
				os.Exit(1)
			}
			fieldFilters = append(fieldFilters, ff)
		}
	}

	lf := filter.NewLineFilter(tr, fieldFilters)
	lp := parser.NewLineParser()
	w, err := output.NewStdoutWriter(*format)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid output format: %v\n", err)
		os.Exit(1)
	}

	for r.Scan() {
		line := r.Text()
		parsed := lp.Parse(line)
		if lf.Accepts(parsed) {
			if err := w.WriteLine(parsed); err != nil {
				fmt.Fprintf(os.Stderr, "write error: %v\n", err)
				os.Exit(1)
			}
		}
	}
	if err := r.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "read error: %v\n", err)
		os.Exit(1)
	}

	if *countOnly {
		fmt.Println(w.Count())
	}
}
