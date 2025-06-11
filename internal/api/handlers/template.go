package handlers

import (
	"html/template"
	"net/http"
	"strconv"
	"time"
)

type TemplateService struct {
	templates map[string]*template.Template
}

func NewTemplateService() *TemplateService {
	service := &TemplateService{
		templates: make(map[string]*template.Template),
	}

	_ = template.FuncMap{
		"formatDate": formatRelativeTime,
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

		_ = template.FuncMap{
			"formatDate": formatRelativeTime,
		}

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

func formatRelativeTime(dateStr string) string {
	// Analyser la date (format MySQL: YYYY-MM-DD HH:MM:SS)
	t, err := time.Parse("2006-01-02 15:04:05", dateStr)
	if err != nil {
		return dateStr // En cas d'erreur, retourner la chaîne originale
	}

	// Calculer la différence entre maintenant et la date donnée
	diff := time.Since(t)

	// Calculer en heures
	hours := int(diff.Hours())

	if hours < 24 {
		// Moins de 24 heures: afficher en heures
		if hours < 1 {
			return "Il y a moins d'une heure"
		} else if hours == 1 {
			return "Il y a 1 heure"
		}
		return "Il y a " + strconv.Itoa(hours) + " heures"
	} else {
		// Plus de 24 heures: afficher en jours
		days := hours / 24
		if days == 1 {
			return "Il y a 1 jour"
		}
		return "Il y a " + strconv.Itoa(days) + " jours"
	}
}
