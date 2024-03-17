package utils

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"
)

const defaultProjectRoot = "unset"

func ProjectRoot() string {
	// TODO(you) This should be the location that you are running this code from - the root of the git repo on your local machine.
	// return "/path/to/where/you/put/this/repository"
	return defaultProjectRoot
}

func MemoizeOperation(fileName string, fn func() ([]byte, error)) ([]byte, error) {
	if ProjectRoot() == defaultProjectRoot {
		return nil, fmt.Errorf("project root is not set, please set it to the root of the git repo on your local machine - see `func ProjectRoot()` in `utils/utils.go`")
	}
	filePath := filepath.Join(ProjectRoot(), "data", fileName)
	result, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("memoization for %s not found, starting to do computation\n", fileName)
			fnResult, err := fn()
			if err != nil {
				return nil, fmt.Errorf("failed to compute fn: %w", err)
			}
			if err := os.WriteFile(filePath, fnResult, 0777); err != nil {
				return nil, fmt.Errorf("failed to write compute fn memoization: %w", err)
			}
			fmt.Printf("memoization for %s computed and saved\n", fileName)
			return fnResult, nil
		}
		return nil, fmt.Errorf("failed to read %s: %w", fileName, err)
	}
	return result, nil
}

func AlphanumericOnly(s string) string {
	result := strings.Builder{}
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			result.WriteRune(r)
		}
	}
	return result.String()
}

func WriteValueToTempJSONFile(value any, prefix string) (filename string, err error) {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshaling value to json: %w", err)
	}
	return WriteBytesToTempFile(data, fmt.Sprintf("%s-*.json", prefix))
}

func WriteErrorsToTempFile(value []error, prefix string) (filename string, err error) {
	strs := make([]string, len(value))
	for i := range value {
		strs[i] = value[i].Error()
	}
	sort.Strings(strs)
	return WriteBytesToTempFile([]byte(strings.Join(strs, "\n")), fmt.Sprintf("%s-*.text", prefix))
}

func WriteValueToTempXMLFile(value any, prefix string) (filename string, err error) {
	data, err := xml.MarshalIndent(value, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshaling value to xml: %w", err)
	}
	return WriteBytesToTempFile(data, fmt.Sprintf("%s-*.xml", prefix))
}

func WriteBytesToTempFile(data []byte, filePathFmt string) (filename string, err error) {
	tmpFile, err := os.CreateTemp("", filePathFmt)
	if err != nil {
		return "", fmt.Errorf("creating temp file: %w", err)
	}
	defer tmpFile.Close()
	_, err = tmpFile.Write(data)
	if err != nil {
		return "", fmt.Errorf("writing to temp file: %w", err)
	}
	path, err := filepath.Abs(tmpFile.Name())
	if err != nil {
		return "", fmt.Errorf("getting absolute path: %w", err)
	}
	return path, nil
}

func CloneJSON[T any](t1 T) (T, error) {
	var t2 T
	data, err := json.Marshal(t1)
	if err != nil {
		return t2, fmt.Errorf("marshaling %T: %w", t1, err)
	}
	if err := json.Unmarshal(data, &t2); err != nil {
		return t2, fmt.Errorf("unmarshaling %T: %w", t2, err)
	}
	return t2, nil
}

func CloneXML[T any](t1 T) (T, error) {
	var t2 T
	data, err := xml.Marshal(t1)
	if err != nil {
		return t2, fmt.Errorf("marshaling %T: %w", t1, err)
	}
	if err := xml.Unmarshal(data, &t2); err != nil {
		return t2, fmt.Errorf("unmarshaling %T: %w", t2, err)
	}
	return t2, nil
}

func Shuffle[T any](ts []T) {
	rand.Shuffle(len(ts), func(i, j int) {
		ts[i], ts[j] = ts[j], ts[i]
	})
}

func SplitIntoBatches[T any](ts []T, n int) [][]T {
	result := [][]T{}
	for i := 0; i < len(ts); i += n {
		start := i
		end := i + n
		if end > len(ts) {
			end = len(ts)
		}
		result = append(result, ts[start:end])
	}
	return result
}
