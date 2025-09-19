package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/DesSolo/rtc/internal/generator"
	"github.com/DesSolo/rtc/internal/generator/parser"
)

var (
	outputPath         = flag.String("output", "internal/config/config.go", "path to output file")
	templatePath       = flag.String("template", "", "path to template file")
	yamlPath           = flag.String("yaml_path", ".", "path to config block in yaml file")
	yamlDescriptionKey = flag.String("yaml_description_key", "usage", "key for description")
)

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		fatal("file path is required")
	}

	filePath := args[0]

	generatorOptions := []generator.OptionFunc{
		generator.WithOutputPath(*outputPath),
	}

	if *templatePath != "" {
		fileContent, err := os.ReadFile(*templatePath)
		if err != nil {
			fatal("failed to read template file: %s", err.Error())
		}

		tpl, err := generator.ParseWithFunctions(string(fileContent))
		if err != nil {
			fatal("failed to parse template file: %s", err.Error())
		}

		generatorOptions = append(generatorOptions, generator.WithTemplate(tpl))
	}

	constGenerator := generator.New(
		parser.NewYaml(*yamlPath, *yamlDescriptionKey),
		generatorOptions...,
	)

	if err := constGenerator.Generate(filePath); err != nil {
		fatal("failed to generate code: %s", err.Error())
	}
}

func fatal(message string, args ...any) {
	fmt.Fprintf(os.Stderr, message+"\n", args...)
	os.Exit(1)
}
