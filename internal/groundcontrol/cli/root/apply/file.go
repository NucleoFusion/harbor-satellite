package apply

import (
	"fmt"
	"os"
	"path/filepath"
)

func getAllFiles(argFiles []string, recursive bool) ([]string, error) {
	files := argFiles

	for _, v := range files {
		isDir, err := isDirectory(v)
		if err != nil {
			return nil, err
		}

		if isDir {
			dirFiles, err := WalkDirectory(v, recursive)
			if err != nil {
				return nil, err
			}

			files = append(files, dirFiles...)
		}
	}

	return files, nil
}

func isDirectory(name string) (bool, error) {
	info, err := os.Stat(name)
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}

func WalkDirectory(dir string, recursive bool) ([]string, error) {
	isDir, err := isDirectory(dir)
	if err != nil {
		return nil, err
	} else if isDir == false {
		return nil, fmt.Errorf("not a directory")
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	files := []string{}
	for _, v := range entries {
		path := filepath.Join(dir, v.Name())
		if !v.IsDir() {
			files = append(files, path)
			continue
		}
		if !recursive {
			continue
		}

		// We can ignore the error since we _know_ its a directory
		nestedFiles, _ := WalkDirectory(path, recursive)
		files = append(files, nestedFiles...)
	}

	return files, nil
}

func readFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading file %s: %v", path, err)
	}

	return data, nil
}
