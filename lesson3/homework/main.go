package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

type Options struct {
	From   string
	To     string
	Offset int64
	Limit  int64
}

func ParseFlags() (*Options, error) {
	var opts Options

	flag.StringVar(&opts.From, "from", "", "file to read. by default - stdin")
	flag.StringVar(&opts.To, "to", "", "file to write. by default - stdout")
	flag.Int64Var(&opts.Offset, "offset", 0, "how many bytes to skip. by default - 0")
	flag.Int64Var(&opts.Limit, "limit", 0, "how many bytes to read. by default - all until EOF")

	flag.Parse()

	return &opts, nil
}

func main() {
	opts, err := ParseFlags()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "can not parse flags:", err)
		os.Exit(1)
	}

	var src io.Reader = os.Stdin
	var dst io.Writer = os.Stdout

	if opts.From != "" {
		srcFile, err := os.Open(opts.From)
		if err != nil {
			log.Fatal(err)
		}
		defer srcFile.Close()
		src = srcFile
	}

	if opts.To != "" {
		dstFile, err := os.OpenFile(opts.To, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer dstFile.Close()
		dst = dstFile
	}
	if opts.Limit != 0 {
		src = io.LimitReader(src, opts.Offset+opts.Limit)
	}

	io.CopyN(io.Discard, src, opts.Offset)
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}
