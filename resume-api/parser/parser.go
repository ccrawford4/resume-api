package parser

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dslipak/pdf"
)

func ParseResume() map[string]string {
	resumesDir := "../resumes"

	files, err := os.ReadDir(resumesDir)
	if err != nil {
		fmt.Println("Error reading resumes directory:", err)
		return map[string]string{}
	}

	result := map[string]string{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		path := filepath.Join(resumesDir, file.Name())
		content, err := ReadPDF(path)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", file.Name(), err)
			continue
		}
		result[file.Name()] = string(content)
	}

	return result
}

func ReadPDF(path string) (string, error) {
	r, err := pdf.Open(path)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		return "", err
	}
	buf.ReadFrom(b)
	return buf.String(), nil
}
