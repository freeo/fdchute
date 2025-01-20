package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"sort"
	"strings"
	"sync"
)

type FileInfo struct {
	Path  string
	IsDir bool
	Depth int
}

func main() {
	cpuProfile := flag.String("cpuprofile", "", "write cpu profile to file")
	memProfile := flag.String("memprofile", "", "write memory profile to file")

	// Parse flags
	flag.Parse()

	// Check if there's at least one non-flag argument (the input path)
	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] <path>\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Get the input path (it will be the first non-flag argument)
	inputPath := flag.Arg(0)

	if *cpuProfile != "" {
		f, _ := os.Create(*cpuProfile)
		defer f.Close()
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if *memProfile != "" {
		f, _ := os.Create(*memProfile)
		defer f.Close()
		pprof.WriteHeapProfile(f)
	}

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <path>\n", os.Args[0])
		os.Exit(1)
	}

	if strings.HasPrefix(inputPath, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
			os.Exit(1)
		}
		inputPath = filepath.Join(home, inputPath[1:])
	}

	files := parallelWalk(inputPath)
	sortFiles(files)

	for _, file := range files {
		fmt.Println(file.Path)
	}
}

func parallelWalk(root string) []FileInfo {
	var wg sync.WaitGroup
	numWorkers := runtime.NumCPU()
	jobs := make(chan string, numWorkers)
	results := make(chan FileInfo, numWorkers)

	// Start worker goroutines
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(jobs, results, &wg)
	}

	// Start a goroutine to close results channel when all workers are done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Start a goroutine to feed jobs
	go func() {
		defer close(jobs)
		jobs <- root
	}()

	var files []FileInfo
	for result := range results {
		files = append(files, result)
	}

	return files
}

func worker(jobs <-chan string, results chan<- FileInfo, wg *sync.WaitGroup) {
	defer wg.Done()
	for path := range jobs {
		processPath(path, results)
	}
}

func processPath(path string, results chan<- FileInfo) {
	entries, err := os.ReadDir(path)
	if err != nil {
		if os.IsPermission(err) {
			return
		}
		// Silently ignore other errors
		return
	}

	for _, entry := range entries {
		fullPath := filepath.Join(path, entry.Name())
		relPath, _ := filepath.Rel(path, fullPath)
		depth := strings.Count(relPath, string(os.PathSeparator))

		results <- FileInfo{
			Path:  relPath,
			IsDir: entry.IsDir(),
			Depth: depth,
		}

		if entry.IsDir() {
			processPath(fullPath, results)
		}
	}
}

func sortFiles(files []FileInfo) {
	sort.Slice(files, func(i, j int) bool {
		if files[i].Depth != files[j].Depth {
			return files[i].Depth < files[j].Depth
		}
		if files[i].IsDir != files[j].IsDir {
			return !files[i].IsDir
		}
		return files[i].Path < files[j].Path
	})
}
