package project

import (
	"bytes"
	"text/template"
)

func renderTemplate(c string, args map[string]interface{}) (string, error) {
	tmpl, err := template.New("main").Parse(c)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, args); err != nil {
		return "", err
	}

	return buf.String(), nil
}
