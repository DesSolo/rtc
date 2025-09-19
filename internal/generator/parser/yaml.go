package parser

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/DesSolo/rtc/internal/generator"
)

type node map[string]any

// Yaml ...
type Yaml struct {
	path           string
	descriptionKey string
}

// NewYaml ...
func NewYaml(path string, descriptionKey string) *Yaml {
	return &Yaml{path: path, descriptionKey: descriptionKey}
}

// Parse ...
func (p *Yaml) Parse(src []byte) ([]*generator.Config, error) {
	tree := make(node)
	if err := yaml.Unmarshal(src, &tree); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tree YAML: %w", err)
	}

	if p.path == "." {
		// there is flat struct
		// key is key, value is struct

		return p.configsFromNode(tree)
	}

	targetNode := walk(tree, strings.Split(p.path, "."))
	if len(targetNode) == 0 {
		return nil, fmt.Errorf("no target node found in %s", p.path)
	}

	return p.configsFromNode(targetNode)
}

func (p *Yaml) configsFromNode(root node) ([]*generator.Config, error) {
	configs := make([]*generator.Config, 0, len(root))

	for key, value := range root {
		val, ok := value.(node)
		if !ok {
			return nil, fmt.Errorf("value for key %s is not a map", key)
		}

		description, ok := val[p.descriptionKey].(string)
		if !ok {
			return nil, fmt.Errorf("value for key %s has no descriptionKey: %s", key, p.descriptionKey)
		}

		configs = append(configs, &generator.Config{
			Key:         key,
			Description: description,
		})
	}

	return configs, nil
}

func walk(root node, path []string) node {
	for _, partPath := range path {
		newRoot := findNode(root, partPath)
		if newRoot == nil {
			return nil
		}

		root = newRoot
	}

	return root
}

func findNode(root node, path string) node {
	for name, value := range root {
		if name == path {
			subNode, ok := value.(node)
			if !ok {
				return nil
			}

			return subNode
		}
	}

	return nil
}
