package store

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/piotr-m-jurek/roadmap-personal-blog/internal/models"
)

const contentDir = "content"

type Store struct {
	contentDir string
}

func New() *Store {
	return &Store{
		contentDir: contentDir,
	}
}

func (s *Store) GetEntries() (map[string]models.Entry, error) {
	entries := map[string]models.Entry{}
	globPath := filepath.Join(s.contentDir, "*.md")
	files, err := filepath.Glob(globPath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			log.Printf("error reading file %s: %v", file, err)
			continue
		}

		parts := strings.SplitN(string(content), "\n\n", 2)
		if len(parts) != 2 {
			log.Printf("malformed content in %s", file)
			continue
		}
		ID := strings.TrimSuffix(filepath.Base(file), ".md")

		entries[ID] = models.Entry{ID: ID, Title: parts[0], Content: parts[1]}
	}
	return entries, nil
}

func (s *Store) SaveEntry(entry models.Entry) error {
	fileName := filepath.Join(s.contentDir, entry.ID+".md")
	data := []byte(entry.Title + "\n\n" + entry.Content)
	return os.WriteFile(fileName, data, 0644)
}

func (s *Store) DeleteEntry(id string) error {
	fileName := filepath.Join(s.contentDir, id+".md")
	return os.Remove(fileName)
}
