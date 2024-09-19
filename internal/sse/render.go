package sse

import (
	"bytes"
	"html/template"
	"io"
)

var path = "internal/sse"

var templates = template.Must(template.ParseGlob(path + "/*.html"))

// Render renders a template document
func render(w io.Writer, name string, data interface{}) error {
	return templates.ExecuteTemplate(w, name, data)
}

// RenderToString renders a template document to a string
func RenderToString(name string, data interface{}) (string, error) {
	var buf []byte
	w := bytes.NewBuffer(buf)

	err := render(w, name, data)
	if err != nil {
		return "", err
	}

	return w.String(), nil
}
