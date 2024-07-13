package storage

import (
	"context"
	"sync"
)

// Result represents the Size function result
type Result struct {
	// Total Size of File objects
	Size int64
	// Count is a count of File objects processed
	Count int64
}

type DirSizer interface {
	// Size calculate a size of given Dir, receive a ctx and the root Dir instance
	// will return Result or error if happened
	Size(ctx context.Context, d Dir) (Result, error)
}

// sizer implement the DirSizer interface
type sizer struct {
	// maxWorkersCount number of workers for asynchronous run
	maxWorkersCount int

	// TODO: add other fields as you wish
}

// NewSizer returns new DirSizer instance
func NewSizer() DirSizer {
	return &sizer{}
}

func (a *sizer) Size(ctx context.Context, d Dir) (r Result, e error) {
	r = Result{
		Size:  0,
		Count: 0,
	}

	dirs, files, err := d.Ls(ctx)
	if err != nil {
		e = err
		return
	}

	var wg sync.WaitGroup
	wg.Add(len(dirs))
	for i := 0; i < len(dirs); i++ {
		dir := dirs[i]
		go func() {
			defer wg.Done()

			subResult, err := a.Size(ctx, dir)
			if err != nil {
				e = err
				return
			}
			r.Size += subResult.Size
			r.Count += subResult.Count
		}()
	}

	for i := 0; i < len(files); i++ {
		file := files[i]
		size, err := file.Stat(ctx)
		if err != nil {
			e = err
			return
		}
		r.Size += size
		r.Count += 1
	}
	wg.Wait()

	return
}
