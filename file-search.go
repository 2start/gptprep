package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"

	gitignore "github.com/sabhiram/go-gitignore"
)

func loadGitIgnorePatterns() (*gitignore.GitIgnore, error) {
	gitIgnorePath := ".gitignore"
	if _, err := os.Stat(gitIgnorePath); os.IsNotExist(err) {
		return nil, nil // No .gitignore file, so no patterns to load.
	}

	patterns, err := os.ReadFile(gitIgnorePath)
	if err != nil {
		return nil, err
	}

	return gitignore.CompileIgnoreLines(strings.Split(string(patterns), "\n")...), nil
}

func isExcludedDir(path string, excludes []string, ignoreObject *gitignore.GitIgnore) bool {
    dirName := filepath.Base(path)
    for _, exclude := range excludes {
        if matched, err := filepath.Match(exclude, dirName); err == nil && matched {
            return true
        }
    }
    if ignoreObject != nil && ignoreObject.MatchesPath(path) {
        return true
    }
    return false
}

func isExcludedFile(path string, extensions, excludes []string, ignoreObject *gitignore.GitIgnore) bool {
    fileName := filepath.Base(path)
    if len(extensions) > 0 && !hasValidExtension(path, extensions) {
        return true
    }
    for _, exclude := range excludes {
        if matched, err := filepath.Match(exclude, fileName); err == nil && matched {
            return true
        }
    }
    if ignoreObject != nil && ignoreObject.MatchesPath(path) {
        return true
    }
    return false
}

func hasValidExtension(path string, extensions []string) bool {
	for _, ext := range extensions {
		if strings.HasSuffix(path, ext) {
			return true
		}
	}
	return false
}

func walkFileFunc(files *[]string, extensions, excludes []string, ignoreObject *gitignore.GitIgnore) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if isExcludedDir(path, excludes, ignoreObject) {
				return filepath.SkipDir
			}
		} else if !isExcludedFile(path, extensions, excludes, ignoreObject) {
			*files = append(*files, path)
		}
		return nil
	}
}

func findFiles(extensions, excludes []string) ([]string, error) {
	var files []string
	ignoreObject, err := loadGitIgnorePatterns()
	if err != nil {
		return nil, err
	}

	// Ensure '.gitignore' is always in the excludes list
	excludes = append(excludes, ".gitignore")

	err = filepath.Walk(".", walkFileFunc(&files, extensions, excludes, ignoreObject))
	if err != nil {
		return nil, err
	}

	return files, nil
}

func concatenateFiles(files []string) (string, error) {
	var buffer bytes.Buffer
	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return "", err
		}
		buffer.WriteString("### filename: " + file + "\n")
		buffer.Write(content)
		buffer.WriteString("\n\n") // Add some spacing between files
	}
	return buffer.String(), nil
}

func SearchAndConcatenateFiles(extensions, excludes []string) (string, error) {
	files, err := findFiles(extensions, excludes)
	if err != nil {
		return "", err
	}
	return concatenateFiles(files)
}
