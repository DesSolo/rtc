package ctl

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func newConfigsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "configs",
		Short: "manage configs",
	}

	cmd.AddCommand(newUpsertConfigCommand())

	return cmd
}

type upsertConfigYaml map[string]struct {
	Value    string `yaml:"value"`
	Type     string `yaml:"type"`
	Usage    string `yaml:"usage"`
	Group    string `yaml:"group"`
	Writable bool   `yaml:"writable"`
}

func newUpsertConfigCommand() *cobra.Command {
	var (
		projectName    string
		envName        string
		releaseName    string
		valuesFilePath string
	)

	cmd := &cobra.Command{
		Use:   "upsert",
		Short: "upsert configs",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()

			slog.DebugContext(ctx, "reading configs", "valuesFilePath", valuesFilePath)
			data, err := os.ReadFile(valuesFilePath) // nolint:gosec
			if err != nil {
				return fmt.Errorf("os.ReadFile: %w", err)
			}

			slog.DebugContext(ctx, "parsing configs")
			var configs upsertConfigYaml
			if err := yaml.Unmarshal(data, &configs); err != nil {
				return fmt.Errorf("yaml.Unmarshal: %w", err)
			}

			slog.DebugContext(ctx, "upserting configs", "projectName", projectName, "envName", envName, "releaseName", releaseName, "configs", len(configs))
			if err := clientFromContext(ctx).UpsertConfigs(ctx, projectName, envName, releaseName, convertConfigsToReq(configs)); err != nil {
				return fmt.Errorf("client.UpsertConfigs: %w", err)
			}

			slog.InfoContext(ctx, "upserted configs", "projectName", projectName, "envName", envName, "releaseName", releaseName, "configs", len(configs))

			return nil
		},
	}

	cmd.Flags().StringVarP(&projectName, "project", "p", "", "Project name")
	cmd.Flags().StringVarP(&envName, "env", "e", "", "Environment name")
	cmd.Flags().StringVarP(&releaseName, "release", "r", "", "Release")
	cmd.Flags().StringVarP(&valuesFilePath, "values", "v", "", "Values file path")

	return cmd
}
