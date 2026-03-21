package assets

import (
	"embed"
	"errors"
	"io/fs"
)

//go:embed mates/*.chatmode.md mates/legacy/*.chatmode.md
var embeddedMates embed.FS

// GetEmbeddedMates returns the embedded mates filesystem
func GetEmbeddedMates() fs.FS {
	matesFS, err := fs.Sub(embeddedMates, "mates")
	if err != nil {
		// This should never happen with valid embed
		panic("failed to access embedded mates: " + err.Error())
	}
	return matesFS
}

// getEmbeddedList enumerates chatmate files within the given directory prefix.
func getEmbeddedList(dir string) ([]string, error) {
	matesFS := GetEmbeddedMates()

	entries, err := fs.ReadDir(matesFS, dir)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if len(name) >= len(".chatmode.md") && name[len(name)-12:] == ".chatmode.md" {
			files = append(files, name)
		}
	}

	return files, nil
}

// GetEmbeddedModernMatesList returns embedded chatmates that are not legacy.
func GetEmbeddedModernMatesList() ([]string, error) {
	return getEmbeddedList(".")
}

// GetEmbeddedLegacyMatesList returns embedded legacy chatmate filenames.
func GetEmbeddedLegacyMatesList() ([]string, error) {
	files, err := getEmbeddedList("legacy")
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}

	return files, nil
}

// GetEmbeddedMatesList returns all embedded chatmate filenames across modern and legacy sets.
func GetEmbeddedMatesList() ([]string, error) {
	modern, err := GetEmbeddedModernMatesList()
	if err != nil {
		return nil, err
	}

	legacy, err := GetEmbeddedLegacyMatesList()
	if err != nil {
		return nil, err
	}

	return append(modern, legacy...), nil
}

// GetEmbeddedMateContent returns the content of a specific embedded chatmate file
// regardless of whether it lives in the main or legacy directory.
func GetEmbeddedMateContent(filename string) ([]byte, error) {
	matesFS := GetEmbeddedMates()

	if data, err := fs.ReadFile(matesFS, filename); err == nil {
		return data, nil
	}

	return fs.ReadFile(matesFS, "legacy/"+filename)
}
