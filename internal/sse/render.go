package sse

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
)

//go:embed *.html
var tmpls embed.FS

// RenderToString renders a template document to a string
func RenderToString(name string, data interface{}) (string, error) {
	var buf []byte
	w := bytes.NewBuffer(buf)

	err := template.Must(template.New(name).ParseFS(tmpls, name)).Execute(w, data)
	if err != nil {
		fmt.Println("Error rendering template:", err)
		return "", err
	}

	return w.String(), nil
}
