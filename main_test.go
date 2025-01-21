package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestName(t *testing.T) {

	tests := []struct {
		name     string
		input    string
		expected [][]FileInfo
		wantErr  bool
	}{
		{
			name: "Group Files by Depth",
			input: `
file0.txt
folder1/file1.txt
folder1/subfolder1/file2.txt
folder1/subfolder1/subfolder2/file3.txt
folder1/subfolder1/subfolder2/subfolder3/file4.txt
`,
			expected: [][]FileInfo{
				{
					{Path: "file0.txt", IsDir: false, Depth: 0},
				},
				{
					{Path: "folder1/file1.txt", IsDir: false, Depth: 1},
				},
				{
					{Path: "folder1/subfolder1/file2.txt", IsDir: false, Depth: 2},
					{Path: "folder1/subfolder1/subfolder2/file3.txt", IsDir: false, Depth: 3},
					{Path: "folder1/subfolder1/subfolder2/subfolder3/file4.txt", IsDir: false, Depth: 4},
				},
			},

			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			files := parseInputString(tt.input)
			result := groupFilesByDepth(files, 3)
			assert.Equal(t, tt.expected, result)
		},
		)
	}
}
