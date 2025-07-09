package parser

import (
	"bytes"
	"io"
	"log"

	"github.com/dslipak/pdf"
)

type File struct {
	Name      string
	Content   []byte
	URL       string
	CreatedAt string
}

type ResumeEntry struct {
	Content string
	URL     string
}

// Map the name of the resume to its content and url
type ParsedResumeResults map[string]ResumeEntry

func UpdateResumeContent(files []*File) {
	// Iterate over each file
	for _, file := range files {
		// Extract the content as bytes into a reader
		reader := bytes.NewReader(file.Content)

		// Parse the PDF content
		parsedContent, err := ReadPDF(reader, int64(len(file.Content)))
		if err != nil {
			log.Printf("error reading PDF %s: %v", file.Name, err)
			continue
		}

		file.Content = []byte(parsedContent) // Update the file content with parsed text
	}
}

func ReadPDF(contentReader io.ReaderAt, size int64) (string, error) {
	r, err := pdf.NewReader(contentReader, size)
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
