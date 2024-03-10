package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

var largeOffsetError = errors.New("offset is larger than input")

type Options struct {
	From   string
	To     string
	Offset uint
	Limit  uint
	Conv   string
}

// OffsetReader returns a Reader that reads from r
// but skips the first n bytes.
// The underlying implementation is an *OffsettedReader.
func OffsetReader(r io.Reader, n int64) io.Reader { return &OffsettedReader{r, n} }

// An OffsettedReader reads from R but discards
// the first N bytes. Each call to Read
// updates N to reflect the new amount to discard remaining.
// Read returns EOF when N <= 0 or when the underlying R returns EOF.
type OffsettedReader struct {
	R io.Reader // underlying reader
	N int64     // bytes to discard
}

func (o *OffsettedReader) Read(p []byte) (n int, err error) {
	if o.N > 0 {
		discarded, err := io.CopyN(io.Discard, o.R, o.N)
		if err != nil {
			if err == io.EOF && discarded != 0 {
				return 0, largeOffsetError
			}
			return 0, err
		}
	}
	n, err = o.R.Read(p)
	return
}

func NewTrimSpaceWriter(w io.Writer) io.Writer { return &TrimSpaceWriter{w} }

type TrimSpaceWriter struct {
	W io.Writer // underlying writer
}

func (t *TrimSpaceWriter) Write(p []byte) (n int, err error) {
	p = bytes.TrimSpace(p)
	n, err = t.W.Write(p)
	return
}

func ParseFlags() (*Options, error) {
	var opts Options

	flag.StringVar(&opts.From, "from", "", "file to read. by default - stdin")
	flag.StringVar(&opts.To, "to", "", "file to write. by default - stdout")
	flag.UintVar(&opts.Offset, "offset", 0, "how many bytes to skip. by default - 0")
	flag.UintVar(&opts.Limit, "limit", 0, "how many bytes to read. by default - all until EOF")
	flag.StringVar(&opts.Conv, "conv", "", "how many bytes to read. by default - all until EOF")

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
		dstFile, err := os.OpenFile(opts.To, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer dstFile.Close()
		dst = dstFile
	}
	if opts.Limit != 0 {
		src = io.LimitReader(src, int64(opts.Offset+opts.Limit))
	}
	if opts.Offset != 0 {
		src = OffsetReader(src, int64(opts.Offset))
	}
	if opts.Conv != "" {
		dst = NewTrimSpaceWriter(dst)
	}

	tee := io.TeeReader(src, dst)
	if _, err := io.ReadAll(tee); err != nil {
		log.Fatal(err)
	}
}
