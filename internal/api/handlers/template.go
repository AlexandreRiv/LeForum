package handlers

import (
	"html/template"
	"net/http"
)

type TemplateService struct {
	templates map[string]*template.Template
}

func NewTemplateService() *TemplateService {
	service := &TemplateService{
		templates: make(map[string]*template.Template),
	}

	// Précharger les templates communs
	templateFiles := []string{
		"authentification.html",
		"categories.html",
		"home_page.html",
		"post_page.html",
		"user.html",
	}

	for _, file := range templateFiles {
		tmpl := template.New(file)
		tmpl, err := tmpl.ParseFiles("web/templates/" + file)
		if err != nil {
			continue
		}

		// Ajouter les composants
		tmpl, err = tmpl.ParseGlob("web/templates/components/*.html")
		if err != nil {
			continue
		}

		service.templates[file] = tmpl
	}

	return service
}

func (s *TemplateService) RenderTemplate(w http.ResponseWriter, name string, data interface{}) error {
	tmpl, exists := s.templates[name]
	if !exists {
		// Si le template n'est pas préchargé, le charger à la demande
		tmpl = template.New(name)
		var err error

		tmpl, err = tmpl.ParseFiles("web/templates/" + name)
		if err != nil {
			return err
		}

		tmpl, err = tmpl.ParseGlob("web/templates/components/*.html")
		if err != nil {
			return err
		}

		s.templates[name] = tmpl
	}

	return tmpl.Execute(w, data)
}
