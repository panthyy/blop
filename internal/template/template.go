package template

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/panthyy/blop/internal/manifest"
	"github.com/panthyy/blop/internal/utils"
)

func Render(m *manifest.Manifest, outputDir string, variables map[string]string) error {
	for _, file := range m.Files {
		var content string
		var err error

		if file.Content != "" {
			content = file.Content
		} else if file.Source != "" {
			fc, err := utils.DownloadFile(file.Source)
			if err != nil {
				return fmt.Errorf("failed to download file %s: %w", file.Source, err)
			}
			content = string(fc)
		} else {
			return fmt.Errorf("file %s has no content or source", file.Path)
		}

		pathTmpl, err := template.New("path").Parse(file.Path)
		if err != nil {
			return fmt.Errorf("failed to parse path template for %s: %w", file.Path, err)
		}

		var pathBuf bytes.Buffer
		if err := pathTmpl.Execute(&pathBuf, variables); err != nil {
			return fmt.Errorf("failed to execute path template for %s: %w", file.Path, err)
		}
		evaluatedPath := pathBuf.String()

		contentTmpl, err := template.New(evaluatedPath).Parse(content)
		if err != nil {
			return fmt.Errorf("failed to parse template for %s: %w", evaluatedPath, err)
		}

		outputPath := filepath.Join(outputDir, evaluatedPath)
		if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
			return fmt.Errorf("failed to create directory for %s: %w", outputPath, err)
		}

		f, err := os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("failed to create file %s: %w", outputPath, err)
		}
		defer f.Close()

		if err := contentTmpl.Execute(f, variables); err != nil {
			return fmt.Errorf("failed to execute template for %s: %w", evaluatedPath, err)
		}
	}

	return nil
}
