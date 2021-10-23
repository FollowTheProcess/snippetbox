package main

import (
	"html/template"
	"path/filepath"

	"github.com/FollowTheProcess/snippetbox/pkg/models"
)

type templateData struct {
	CurrentYear int
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := make(map[string]*template.Template)

	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		// Parse the page template into a set
		ts, err := template.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Add any layout templates to the set
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		// Add any partial templates to the set
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		// Add the template set to the cache with the name of the page
		// as the key
		cache[name] = ts
	}

	return cache, nil
}
