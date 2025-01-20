package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type DirStats struct {
	Path         string
	FileCount    int
	SubdirCount  int
	TotalEntries int
}

func main() {
	// homeDir, err := os.UserHomeDir()
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Could not get home directory: %v\n", err)
	// 	os.Exit(1)
	// }

	dirStats := make(map[string]*DirStats)
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		path := scanner.Text()
		dir := filepath.Dir(path)

		// Corrected initialization check
		if _, exists := dirStats[dir]; !exists {
			dirStats[dir] = &DirStats{Path: dir}
		}

		// Update stats
		if filepath.Dir(path) != path {
			if filepath.IsAbs(path) || strings.Contains(path, string(os.PathSeparator)) {
				if filepath.Ext(path) == "" {
					dirStats[dir].SubdirCount++
				} else {
					dirStats[dir].FileCount++
				}
				dirStats[dir].TotalEntries++
			}
		}
	}

	// Convert map to slice for sorting
	statsList := make([]*DirStats, 0, len(dirStats))
	for _, stat := range dirStats {
		statsList = append(statsList, stat)
	}

	// Sort by total entries in descending order
	sort.Slice(statsList, func(i, j int) bool {
		return statsList[i].TotalEntries > statsList[j].TotalEntries
	})

	// Print top 20 directories with most entries
	fmt.Println("Top Directories by Number of Entries:")
	fmt.Println("Path\t\tTotal Entries\tFiles\tSubdirs")
	for i, stat := range statsList {
		if i >= 40 || stat.TotalEntries < 10 {
			break
		}
		fmt.Printf("%d\t\t%d\t\t%d\t%s\n",
			stat.TotalEntries,
			stat.FileCount,
			stat.SubdirCount,
			stat.Path)
	}
}
