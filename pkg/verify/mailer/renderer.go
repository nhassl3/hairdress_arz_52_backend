package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
)

const (
	basePathToTemplates = "pkg/verify/mailer/templates/"
	EmailConfirmation   = "email_confirmation.html"
)

// Render rendering new html message on email client
func Render(name string, data interface{}) (string, error) {
	if !strings.HasSuffix(name, ".html") {
		name = name + ".html"
	}
	t, err := template.ParseFiles(
		basePathToTemplates+"base.html",
		basePathToTemplates+name,
	)
	if err != nil {
		return "", fmt.Errorf("renderer.Render: failed to parse template: %w", err)
	}

	var body bytes.Buffer
	if err = t.Execute(&body, data); err != nil {
		return "", fmt.Errorf("renderer.Render: failed to execute template: %w", err)
	}

	return body.String(), nil
}
