package assets

import (
	"embed"
	"io/fs"
)

//go:embed mates/*.chatmode.md
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

// GetEmbeddedMatesList returns a list of all embedded chatmate filenames
func GetEmbeddedMatesList() ([]string, error) {
	matesFS := GetEmbeddedMates()

	entries, err := fs.ReadDir(matesFS, ".")
	if err != nil {
		return nil, err
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() && entry.Name()[len(entry.Name())-12:] == ".chatmode.md" {
			files = append(files, entry.Name())
		}
	}

	return files, nil
}

// GetEmbeddedMateContent returns the content of a specific embedded chatmate file
func GetEmbeddedMateContent(filename string) ([]byte, error) {
	matesFS := GetEmbeddedMates()
	return fs.ReadFile(matesFS, filename)
}
