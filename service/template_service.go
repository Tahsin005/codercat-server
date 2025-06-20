package service

import (
	"bytes"
	"html/template"
	"path/filepath"
)

type TemplateService interface {
	RenderEmailTemplate(templateName string, data interface{}) (string, error)
}

type templateService struct {
	templatesDir string
}

func NewTemplateService(templatesDir string) TemplateService {
	return &templateService{
		templatesDir: templatesDir,
	}
}

func (s *templateService) RenderEmailTemplate(templateName string, data interface{}) (string, error) {
	templatePath := filepath.Join(s.templatesDir, "email", templateName+".html")

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
