package sse

import (
	"bytes"
	"html/template"
)

// RenderToString renders a template document to a string
func RenderToString(name string, data interface{}) (string, error) {
	var buf []byte
	w := bytes.NewBuffer(buf)

	err := template.Must(template.New(name).ParseFiles("internal/sse/"+name)).Execute(w, data)
	if err != nil {
		return "", err
	}

	return w.String(), nil
}
