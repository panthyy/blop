package cache

import (
	"fmt"
	"os"
	"path/filepath"
)

type Cache struct {
	Dir string
}

func New(dir string) *Cache {
	return &Cache{Dir: dir}
}

func (c *Cache) Get(key string) ([]byte, error) {
	path := filepath.Join(c.Dir, key)
	return os.ReadFile(path)
}

func (c *Cache) Set(key string, data []byte) error {
	if err := os.MkdirAll(c.Dir, 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	path := filepath.Join(c.Dir, key)
	return os.WriteFile(path, data, 0644)
}

func (c *Cache) Remove(key string) error {
	path := filepath.Join(c.Dir, key)
	return os.Remove(path)
}

func (c *Cache) List() ([]string, error) {
	files, err := os.ReadDir(c.Dir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to read cache directory: %w", err)
	}

	var templates []string
	for _, f := range files {
		if !f.IsDir() {
			templates = append(templates, f.Name())
		}
	}

	return templates, nil
}
