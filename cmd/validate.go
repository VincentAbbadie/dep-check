package cmd

import (
	"fmt"
	"log/slog"

	"github.com/moveaxlab/dep-check/config"
	"github.com/moveaxlab/dep-check/graph"
	"github.com/moveaxlab/dep-check/languages"
	"github.com/spf13/cobra"
)

var (
	validateCmd = &cobra.Command{
		Use:   "validate",
		Short: "validate the dependencies graph based on the seleced language",
		Long:  fmt.Sprintf("validate the dependencies graph based on the seleced language and the configuration provide in the %s", config.DepCheckFileName),
		PreRunE: func(cmd *cobra.Command, args []string) error {

			var graphBuilder graph.GraphBuilder
			switch config.SelectedLanguage {
			case languages.Go:
				graphBuilder = &languages.GoGraphBuilder{}
			}

			if err := graphBuilder.Init(); err != nil {
				slog.Error("failure when initialize the language graph builder")
				return err
			}

			slog.Info("creating dependencies graph")
			graphBuilder.BuildGraph()
			slog.Info(fmt.Sprintf("graph contains %d nodes", len(graph.Nodes)))
			slog.Info(fmt.Sprintf("graph contains %d edges", len(graph.Edges)))

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			validationErrors := make([]string, 0)
			for _, e := range graph.Edges {
				if isValid, reason := e.IsValid(); !isValid {
					validationErrors = append(validationErrors, reason)
				}
			}
			if len(validationErrors) > 0 {
				slog.Error("Graph is not Valide :")
				for _, reason := range validationErrors {
					slog.Error(reason)
				}
				cmd.SilenceUsage = true
				return fmt.Errorf("invalid imports detected, check error log for more details")
			}
			return nil
		},
	}
)

func init() {
	RootCmd.AddCommand(validateCmd)
}

// func createOrGetNodeFromGo(pkgPath string) (n *graph.Node) {
// 	for _, artif := range config.DepCheckConfig.Artifacts {
// 		prefix := path.Join(moduleName, artif)
// 		hasWildcard := strings.HasSuffix(prefix, "*")
// 		if hasWildcard {
// 			prefix = strings.TrimSuffix(prefix, "*")
// 		}
// 		if strings.HasPrefix(pkgPath, prefix) {
// 			name := artif
// 			if hasWildcard {
// 				name = strings.Replace(artif, "*", strings.Split(strings.Replace(pkgPath, prefix, "", 1), "/")[0], 1)
// 				prefix = path.Join(moduleName, name)
// 			}
// 			if _, ok := graph.Nodes[name]; !ok {
// 				n = &graph.Node{Name: name, Path: prefix}
// 				graph.Nodes[name] = n
// 			} else {
// 				n = graph.Nodes[name]
// 			}
// 			return n
// 		}
// 	}
// 	return nil
// }
