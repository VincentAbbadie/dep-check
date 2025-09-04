package detect

import (
	"fmt"
	"strings"

	"github.com/moveaxlab/dep-check/cmd"
	"github.com/moveaxlab/dep-check/config"
	"github.com/spf13/cobra"
)

const (
	outputFormatFlag             = "format"
	outputFormatFlagDefaultValue = "multiline"
)

var (
	outputFormatFlagList = []string{"oneline", outputFormatFlagDefaultValue}
	detectCmd            = &cobra.Command{
		Use: "detect",
	}
)

func init() {
	detectCmd.PersistentFlags().StringVar(&config.Format, outputFormatFlag, outputFormatFlagDefaultValue, fmt.Sprintf("Language to parse. Available languages are : %s", strings.Join(languageFlagList, ", ")))

	cmd.RootCmd.AddCommand(detectCmd)
}
