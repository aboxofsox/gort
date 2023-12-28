package gort

import (
	"os"
	"path/filepath"
)

// readDir recursively reads the directory specified by 'dir' and calls the 'handler' function for each file or directory found.
// The 'handler' function is called with the full path of the file or directory and its corresponding os.DirEntry.
// If an error occurs while reading the directory, the function returns without processing any files.
func readDir(dir string, handler func(path string, entry os.DirEntry)) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	for _, file := range files {
		fp := filepath.Join(dir, file.Name())
		if file.IsDir() {
			readDir(fp, handler)
		} else {
			handler(fp, file)
		}
	}
}
