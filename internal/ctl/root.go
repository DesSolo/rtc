package ctl

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/DesSolo/rtc/internal/ctl/client"
)

// Execute ...
func Execute() {
	if err := newRootCommand().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rtcctl",
		Short: "rtcctl is a tool to manage RTC",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			url := cmd.Flags().Lookup("url").Value.String()

			token := cmd.Flags().Lookup("token").Value.String()
			if token == "" {
				token = os.Getenv("RTCCTL_TOKEN")
			}

			logLevel, err := parseLogLevel(cmd.Flags().Lookup("log-level").Value.String())
			if err != nil {
				return fmt.Errorf("invalid log level: %w", err)
			}

			slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				Level: logLevel,
			})))

			rtcClient := client.NewClient(url, token)

			cmd.SetContext(clientToContext(cmd.Context(), rtcClient))
			return nil
		},
	}

	cmd.PersistentFlags().StringP("url", "u", "http://localhost:8080/api/v1", "RTC server url")
	cmd.PersistentFlags().StringP("token", "t", "", "RTC server token")
	cmd.PersistentFlags().StringP("log-level", "l", "0", "log level info=0 debug=-4")

	cmd.AddCommand(newConfigsCommand())

	return cmd
}

func parseLogLevel(logLevel string) (slog.Level, error) {
	val, err := strconv.Atoi(logLevel)
	if err != nil {
		return 0, fmt.Errorf("strconv.Atoi: %w", err)
	}

	return slog.Level(val), nil
}
