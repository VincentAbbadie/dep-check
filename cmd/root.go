package cmd

import (
	"github.com/moveaxlab/dep-check/config"
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/spf13/cobra"
)

const (
	debugFlag    = "debug"
	languageFlag = "language"
)

var (
	// DebugMode        = false
	languageFlagList = []string{"go", "js", "java"}
	// SelectedLanguage string
	rootCmd = &cobra.Command{
		Use:   "dep-check",
		Short: "Compute the dependency graph of a monorepo",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if config.DebugMode {
				slog.SetLogLoggerLevel(slog.LevelDebug)
				slog.Debug("Debug mode activated")
			}

			isValid := false
			if slices.Contains(languageFlagList, config.SelectedLanguage) {
				isValid = true
				slog.Info(fmt.Sprintf("Language selected : %s", config.SelectedLanguage))
			}

			if !isValid {
				return fmt.Errorf("langugage %s is not supported", config.SelectedLanguage)
			}

			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if config.DepCheckConfig.IsEmpty() {
				return fmt.Errorf("%s is empty", config.DepCheckFileName)
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Usage()
		},
	}
)

func init() {
	rootCmd.PersistentFlags().BoolVar(&config.DebugMode, debugFlag, false, "Add aditional information for debuging purpose")
	rootCmd.PersistentFlags().StringVarP(&config.SelectedLanguage, languageFlag, "l", "", fmt.Sprintf("Language to parse. Available languages are : %s", strings.Join(languageFlagList, ", ")))
	rootCmd.MarkPersistentFlagRequired(languageFlag)
}

func Execute() {
	rootCmd.Execute()
}
