package template

import (
	"embed"
	"fmt"
	_ "io/fs"
	"log"
	"path/filepath"
	"strings"
	"text/template"
)

type TemplateManager struct {
	templates map[string]map[string]map[string]*template.Template // [language][channel][eventType]*template.Template
}

func NewTemplateManager() *TemplateManager {
	return &TemplateManager{
		templates: make(map[string]map[string]map[string]*template.Template),
	}
}

func (m *TemplateManager) LoadAllTemplates(fs embed.FS, basePath string) error {
	// Read the base directory
	languages, err := fs.ReadDir(basePath)
	if err != nil {
		return err
	}

	for _, langDir := range languages {
		if !langDir.IsDir() {
			continue
		}
		language := langDir.Name()

		// Read the language subdirectory
		channels, err := fs.ReadDir(filepath.Join(basePath, language))
		if err != nil {
			return err
		}

		for _, channelDir := range channels {
			if !channelDir.IsDir() {
				continue
			}
			channel := channelDir.Name()

			// Read the channel subdirectory
			templates, err := fs.ReadDir(filepath.Join(basePath, language, channel))
			if err != nil {
				return err
			}

			for _, tmplFile := range templates {
				if tmplFile.IsDir() {
					continue
				}

				fileName := tmplFile.Name()
				eventType := strings.TrimSuffix(fileName, filepath.Ext(fileName))

				// Read the template file
				content, err := fs.ReadFile(filepath.Join(basePath, language, channel, fileName))
				if err != nil {
					return err
				}

				// Parse the template
				tmpl, err := template.New(fileName).Parse(string(content))
				if err != nil {
					return err
				}

				// Initialize the maps if necessary
				if m.templates[language] == nil {
					m.templates[language] = make(map[string]map[string]*template.Template)
				}
				if m.templates[language][channel] == nil {
					m.templates[language][channel] = make(map[string]*template.Template)
				}

				// Store the template
				m.templates[language][channel][eventType] = tmpl
			}
		}
	}

	return nil
}

func (m *TemplateManager) GetTemplate(language, channel, eventType string) (*template.Template, error) {
	if tmpl, ok := m.templates[language][channel][eventType]; ok {
		return tmpl, nil
	}
	if tmpl, ok := m.templates["en"][channel][eventType]; ok {
		log.Printf("Template for language '%s', channel '%s', event '%s' not found. Falling back to English.", language, channel, eventType)
		return tmpl, nil
	}
	return nil, fmt.Errorf("template not found for language '%s' (or default), channel '%s', and event '%s'", language, channel, eventType)
}
