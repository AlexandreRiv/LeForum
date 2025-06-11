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

	funcMap := template.FuncMap{
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
		tmpl := template.New(file).Funcs(funcMap)
		var err error
		tmpl, err = tmpl.ParseFiles("web/templates/" + file)
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
		funcMap := template.FuncMap{
			"formatDate": formatRelativeTime,
		}

		tmpl = template.New(name).Funcs(funcMap)
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

	// Convertir en différentes unités de temps
	seconds := int(diff.Seconds())
	minutes := int(diff.Minutes())
	hours := int(diff.Hours())
	days := hours / 24

	switch {
	case seconds < 60:
		// Moins d'une minute
		if seconds < 5 {
			return "À l'instant"
		} else if seconds == 1 {
			return "Il y a 1 seconde"
		}
		return "Il y a " + strconv.Itoa(seconds) + " secondes"

	case minutes < 60:
		// Moins d'une heure
		if minutes == 1 {
			return "Il y a 1 minute"
		}
		return "Il y a " + strconv.Itoa(minutes) + " minutes"

	case hours < 24:
		// Moins d'un jour
		if hours == 1 {
			return "Il y a 1 heure"
		}
		return "Il y a " + strconv.Itoa(hours) + " heures"

	case days < 30:
		// Moins d'un mois
		if days == 1 {
			return "Il y a 1 jour"
		}
		return "Il y a " + strconv.Itoa(days) + " jours"

	default:
		// Plus d'un mois, afficher la date formatée
		return t.Format("02/01/2006")
	}
}
