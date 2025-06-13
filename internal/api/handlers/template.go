package handlers

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
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

	// Associer les chemins de fichiers aux noms de définition
	templateDefinitions := map[string]string{
		"admin/dashboard.html":      "dashboard",
		"admin/users.html":          "users",
		"admin/categories.html":     "categories",
		"admin/reports.html":        "reports",
		"moderation/dashboard.html": "moderation-dashboard",
	}

	templateFiles := []string{
		"authentification.html",
		"categories.html",
		"home_page.html",
		"post_page.html",
		"user.html",
		"admin/dashboard.html",
		"moderation/dashboard.html",
	}

	for _, file := range templateFiles {
		// Créer un nouveau template avec le bon nom de définition s'il existe
		var tmplName string
		if defName, exists := templateDefinitions[file]; exists {
			tmplName = defName
		} else {
			tmplName = file
		}

		tmpl := template.New(tmplName).Funcs(funcMap)
		var err error

		// Utiliser le chemin complet pour le parsing
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
	// Déterminer le nom de la définition du template
	templateDefinitions := map[string]string{
		"admin/dashboard.html":      "dashboard",
		"admin/users.html":          "users",
		"admin/categories.html":     "categories",
		"admin/reports.html":        "reports",
		"moderation/dashboard.html": "moderation-dashboard",
	}

	// Récupérer le nom de la définition
	definitionName, exists := templateDefinitions[name]
	if !exists {
		definitionName = filepath.Base(name)
	}

	tmpl, exists := s.templates[name]
	if !exists {
		// Charger le template s'il n'existe pas déjà
		funcMap := template.FuncMap{
			"formatDate": formatRelativeTime,
		}

		// IMPORTANT: Créer le template avec le nom de la définition, pas le nom du fichier
		tmpl = template.New(definitionName).Funcs(funcMap)
		var err error

		// Parser le fichier principal
		tmpl, err = tmpl.ParseFiles("web/templates/" + name)
		if err != nil {
			return err
		}

		// Parser les composants
		tmpl, err = tmpl.ParseGlob("web/templates/components/*.html")
		if err != nil {
			return err
		}

		s.templates[name] = tmpl
	}

	log.Printf("Exécution du template %s avec définition %s", name, definitionName)

	// Exécuter le template avec le nom de la définition
	return tmpl.ExecuteTemplate(w, definitionName, data)
}

func formatRelativeTime(dateStr string) string {
	// Analyser la date (format MySQL: YYYY-MM-DD HH:MM:SS)
	t, err := time.Parse("2006-01-02 15:04:05", dateStr)
	if err != nil {
		return dateStr // En cas d'erreur, retourner la chaîne originale
	}

	// Calculer la différence entre maintenant et la date donnée
	diff := time.Since(t)

	// Ajouter deux heures à la différence
	diff = diff + 2*time.Hour

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
