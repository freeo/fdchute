package main

import (
	"bufio"
	"flag"
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

type Config struct {
	debug bool
}

var CFG Config

// TODO /home/freeo/wb edge case
func main() {
	CFG = Config{
		debug: false,
	}
	flag.BoolVar(&CFG.debug, "debug", false, "Enable debug output")
	flag.Parse()

	if CFG.debug {
		fmt.Printf("DEBUG: %t \n", CFG.debug)
	}

	files := readInput()
	process(files)

	// supposed to be the fastest. TODO: deeper analysis
	//
	// buf := make([]byte, 65536)
	// for _, file := range files {
	// 	n := copy(buf, file.Path)
	// 	buf[n] = '\n'
	// 	syscall.Write(1, buf[:n+1])
	// }

	// writer := bufio.NewWriterSize(os.Stdout, 65536)
	// defer writer.Flush()

	// for i, file := range files {
	// 	writer.WriteString(file.Path)
	// 	writer.WriteByte('\n')

	// 	// Flush the buffer every 1000 lines or when the buffer is close to full
	// 	if (i+1)%1000 == 0 || writer.Available() < 1024 {
	// 		writer.WriteString("FLUSH: " + strconv.Itoa(i) + "\n")
	// 		writer.Flush()
	// 	}
	// }
}

func process(files []FileInfo) {
	groups := groupFilesByDepth(files, 3)
	for i, g := range groups {
		if CFG.debug {
			fmt.Printf("FD LEVEL %d: %d \n", i, len(g))
		}
		if i == 0 {
			sortFiles(g)
			writeStream(g)
		} else {
			writeStream(g)
		}
	}
}

func writeStream(files []FileInfo) {
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

		fi := FileInfo{
			Path:  path,
			IsDir: isDir,
			Depth: depth,
		}
		files = append(files, fi)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}

	return files
}

func filterFilesByDepth(files []FileInfo, depthThreshold int) ([]FileInfo, []FileInfo) {
	shallowFiles := make([]FileInfo, 0)
	deepFiles := make([]FileInfo, 0)

	for _, file := range files {
		if file.Depth <= depthThreshold {
			shallowFiles = append(shallowFiles, file)
		} else {
			deepFiles = append(deepFiles, file)
		}
	}

	return shallowFiles, deepFiles
}

func groupFilesByDepth(files []FileInfo, countGroups int) [][]FileInfo {
	// Initialize the result slice with countGroups number of inner slices
	result := make([][]FileInfo, countGroups)

	for _, file := range files {
		// Calculate the depth of the file
		depth := strings.Count(file.Path, string(os.PathSeparator))

		// Determine which group the file belongs to
		groupIndex := depth
		if groupIndex >= countGroups-1 {
			groupIndex = countGroups - 1 // Put in the last group if depth exceeds or equals countGroups-1
		}

		// Add the FileInfo to the appropriate group
		result[groupIndex] = append(result[groupIndex], file)
	}

	// Remove empty slices
	finalResult := [][]FileInfo{}
	for _, group := range result {
		if len(group) > 0 {
			finalResult = append(finalResult, group)
		}
	}

	return finalResult
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

func parseInputString(input string) []FileInfo {
	var files []FileInfo
	scanner := bufio.NewScanner(strings.NewReader(input))

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

		fi := FileInfo{
			Path:  path,
			IsDir: isDir,
			Depth: depth,
		}
		files = append(files, fi)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing input string: %v\n", err)
		return nil
	}

	return files
}
