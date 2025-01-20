package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

type FileInfo struct {
	Path  string
	IsDir bool
	Depth int
}

func main() {
	files := readInput()
	sortFiles(files)

	// supposed to be the fastest. TODO: deeper analysis
	//
	// buf := make([]byte, 65536)
	// for _, file := range files {
	// 	n := copy(buf, file.Path)
	// 	buf[n] = '\n'
	// 	syscall.Write(1, buf[:n+1])
	// }

	// This feels faster and is slightly faster in the terminal output
	writer := bufio.NewWriterSize(os.Stdout, 65536)
	defer writer.Flush()

	for _, file := range files {
		writer.WriteString(file.Path)
		writer.WriteByte('\n')
	}
}

func readInput() []FileInfo {
	var files []FileInfo
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		path := strings.TrimSpace(scanner.Text())
		if path == "" {
			continue // Skip empty lines
		}

		depth := strings.Count(path, string(os.PathSeparator))
		isDir := strings.HasSuffix(path, string(os.PathSeparator))

		// Remove trailing slash for directories
		if isDir {
			path = strings.TrimSuffix(path, string(os.PathSeparator))
		}

		files = append(files, FileInfo{
			Path:  path,
			IsDir: isDir,
			Depth: depth,
		})
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}

	return files
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
