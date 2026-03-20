package cmd

import (
	"fmt"
	"os"

	"barong-cli/cmd/user"

	"github.com/spf13/cobra"
)

var baseURL string

var rootCmd = &cobra.Command{
	Use:   "barong-cli",
	Short: "Command line tool to interact with Barong APIs",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&baseURL, "url", "", "Barong server base URL (overrides BARONG_URL)")
	rootCmd.AddCommand(user.NewUserCmd(getBaseURL))
}

func getBaseURL() string {
	if baseURL != "" {
		return baseURL
	}
	if u := os.Getenv("BARONG_URL"); u != "" {
		return u
	}
	return "http://localhost:9090"
}
