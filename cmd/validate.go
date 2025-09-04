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

			// // Build Graph
			// for _, pkg := range goPackages {
			// 	pkgNode := createOrGetNodeFromGo(pkg.PkgPath)
			// 	if pkgNode != nil {
			// 		for _, imp := range pkg.Imports {
			// 			impNode := createOrGetNodeFromGo(imp.PkgPath)
			// 			if impNode != nil && pkgNode.Name != impNode.Name {
			// 				e := &graph.Edge{From: pkgNode.Name, To: impNode.Name}
			// 				if _, ok := graph.Edges[e.String()]; !ok {
			// 					graph.Edges[e.String()] = e
			// 				}
			// 			}
			// 		}
			// 	}
			// }

			// slog.Debug(fmt.Sprintf("Nodes num : %d", len(graph.Nodes)))
			// if DebugMode {
			// 	for _, n := range graph.Nodes {
			// 		slog.Debug(fmt.Sprintf(" %s : %s", n.Name, n.Path))
			// 	}
			// }
			// slog.Debug(fmt.Sprintf("Edges num : %d", len(graph.Edges)))
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			slog.Error(fmt.Sprintf("[Valide CND]Selected language %s", config.SelectedLanguage))
			// changedPackages := make(map[string]string)

			// scanner := bufio.NewScanner(os.Stdin)
			// for scanner.Scan() {
			// 	changedFile := scanner.Text()
			// 	slog.Info(fmt.Sprintf("Input file : %s", changedFile))

			// 	if filepath.Dir(changedFile) == "." {
			// 		continue
			// 	}

			// 	filePackage := createOrGetNodeFromGo(path.Join(moduleName, changedFile))
			// 	if filePackage == nil {
			// 		slog.Warn("No node found")
			// 		continue
			// 	}
			// 	slog.Debug("Corresponding node :")
			// 	slog.Debug(fmt.Sprintf("%s : %s", filePackage.Name, filePackage.Path))

			// 	if _, ok := changedPackages[filePackage.Name]; !ok {
			// 		for _, parent := range filePackage.GetParents() {
			// 			if _, ok := changedPackages[parent.Name]; !ok {
			// 				changedPackages[parent.Name] = strings.Join([]string{parent.Path, "/..."}, "")
			// 				fmt.Println(changedPackages[parent.Name])
			// 			}
			// 		}
			// 	}
			// }
		},
	}
)

func init() {
	rootCmd.AddCommand(validateCmd)
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
