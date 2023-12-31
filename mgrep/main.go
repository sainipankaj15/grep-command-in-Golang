package main

import (
	"fmt"
	"grep-command-clone/worker"
	"grep-command-clone/worklist"
	"os"
	"path/filepath"
	"sync"

	"github.com/alexflint/go-arg"
)

func discoverDir(wl *worklist.Worklist, path string) {

	entries, err := os.ReadDir(path)

	if err != nil {
		fmt.Println("Reading directors error ", err)
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			nextPath := filepath.Join(path, entry.Name())
			discoverDir(wl, nextPath)
		} else {
			wl.Add(worklist.NewJob(filepath.Join(path, entry.Name())))
		}
	}
}

var args struct {
	SearchTerm string `arg:"positional, required"`
	SearchDir  string `arg:"positional"`
}

func main() {

	arg.MustParse(&args)

	var workersWg sync.WaitGroup

	wl := worklist.New(100)

	results := make(chan worker.Result, 100)

	numWorkers := 10

	workersWg.Add(1)

	go func() {
		defer workersWg.Done()
		discoverDir(&wl, args.SearchDir)
		wl.Finalize(numWorkers)
	}()

	for i := 0; i < numWorkers; i++ {

		workersWg.Add(1)

		go func() {
			defer workersWg.Done()
			for {
				workEntry := wl.Next()

				if workEntry.Path != "" {
					workerResult := worker.FindInFile(workEntry.Path, args.SearchTerm)

					for _, r := range workerResult.Inner {
						results <- r
					}
				} else {
					// When the path is empty, this indicates that there are no more jobs available,
					// so we quit and will return from there
					return
				}
			}
		}()
	}

	blockWorkersWg := make(chan struct{})
	go func() {
		workersWg.Wait()
		close(blockWorkersWg)
	}()

	var displayWg sync.WaitGroup

	displayWg.Add(1)
	go func() {
		for {
			select {
			case r := <-results:
				fmt.Printf("%v[%v]: %v\n ", r.Path, r.LineNum, r.Line)
			case <- blockWorkersWg :
				if len(results) == 0 {
					displayWg.Done()
					return 
				}
			}
		}
	}()

	displayWg.Wait()
}
