package generator

import "text/template"

// OptionFunc ...
type OptionFunc func(*Generator)

// WithOutputPath set a specific output path
func WithOutputPath(path string) OptionFunc {
	return func(g *Generator) {
		g.targetPath = path
	}
}

// WithTemplate returns a OptionFunc for set custom template
// note: use ParseWithFunctions for access some rich functions
func WithTemplate(template *template.Template) OptionFunc {
	return func(g *Generator) {
		g.template = template
	}
}
