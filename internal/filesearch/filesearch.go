package filesearch

import (
	"bytes"
	"net/http"
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

	if !isTextFile(path) {
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

// isTextFile reads the first 512 bytes of a file and uses http.DetectContentType
// to determine if the file is a text file.
func isTextFile(filePath string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer file.Close()

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return false
	}

	contentType := http.DetectContentType(buffer)

	// MIME types starting with "text/" are text files.
	return strings.HasPrefix(contentType, "text/")
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

	excludes = append(excludes, ".gitignore", ".git", "LICENSE")

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
		buffer.WriteString("\n\n")
	}
	return buffer.String(), nil
}

func SearchAndConcatenateFiles(extensions, excludes []string) ([]string, string, error) {
	files, err := findFiles(extensions, excludes)
	if err != nil {
		return []string{}, "", err
	}

	concatenatedContent, err := concatenateFiles(files)
	if err != nil {
		return []string{}, "", err 
	}

	return files, concatenatedContent, nil
}
