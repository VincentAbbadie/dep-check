package languages

import (
	"fmt"
	"log/slog"
	"os"
	"path"
	"strings"

	"github.com/moveaxlab/dep-check/config"
	"github.com/moveaxlab/dep-check/graph"

	"golang.org/x/mod/modfile"
	"golang.org/x/tools/go/packages"
)

const (
	goModFileName = "go.mod"
	Go            = "go"
)

var (
	moduleName string
	goPackages []*packages.Package
)

type GoGraphBuilder struct{}

// var (
// 	goCmd      = &cobra.Command{
// 		Use:   "go",
// 		Short: "build dependencies graph based on the go language",
// 		PreRun: func(cmd *cobra.Command, args []string) {
//

// 			// Build Graph
// 			for _, pkg := range goPackages {
// 				pkgNode := createOrGetNodeFromGo(pkg.PkgPath)
// 				if pkgNode != nil {
// 					for _, imp := range pkg.Imports {
// 						impNode := createOrGetNodeFromGo(imp.PkgPath)
// 						if impNode != nil && pkgNode.Name != impNode.Name {
// 							e := &graph.Edge{From: pkgNode.Name, To: impNode.Name}
// 							if _, ok := graph.Edges[e.String()]; !ok {
// 								graph.Edges[e.String()] = e
// 							}
// 						}
// 					}
// 				}
// 			}

// 			slog.Debug(fmt.Sprintf("Nodes num : %d", len(graph.Nodes)))
// 			if DebugMode {
// 				for _, n := range graph.Nodes {
// 					slog.Debug(fmt.Sprintf(" %s : %s", n.Name, n.Path))
// 				}
// 			}
// 			slog.Debug(fmt.Sprintf("Edges num : %d", len(graph.Edges)))
// 		},
// 		Run: func(cmd *cobra.Command, args []string) {
// 			changedPackages := make(map[string]string)

// 			scanner := bufio.NewScanner(os.Stdin)
// 			for scanner.Scan() {
// 				changedFile := scanner.Text()
// 				slog.Info(fmt.Sprintf("Input file : %s", changedFile))

// 				if filepath.Dir(changedFile) == "." {
// 					continue
// 				}

// 				filePackage := createOrGetNodeFromGo(path.Join(moduleName, changedFile))
// 				if filePackage == nil {
// 					slog.Warn("No node found")
// 					continue
// 				}
// 				slog.Debug("Corresponding node :")
// 				slog.Debug(fmt.Sprintf("%s : %s", filePackage.Name, filePackage.Path))

// 				if _, ok := changedPackages[filePackage.Name]; !ok {
// 					for _, parent := range filePackage.GetParents() {
// 						if _, ok := changedPackages[parent.Name]; !ok {
// 							changedPackages[parent.Name] = strings.Join([]string{parent.Path, "/..."}, "")
// 							fmt.Println(changedPackages[parent.Name])
// 						}
// 					}
// 				}
// 			}
// 		},
// 	}
// )

func (g *GoGraphBuilder) Init() error {
	//Load Module Name
	data, err := os.ReadFile(goModFileName)
	if err != nil {
		slog.Error(fmt.Sprintf("No %s file", goModFileName))
		return err
	}
	modFile, err := modfile.Parse(goModFileName, data, nil)
	if err != nil {
		slog.Error(fmt.Sprintf("%s file unparsable", goModFileName))
		return err
	}
	moduleName = modFile.Module.Mod.Path
	slog.Info(fmt.Sprintf("ModuleName : %s", moduleName))

	// Load packages
	packagesConfig := &packages.Config{Mode: packages.NeedImports | packages.NeedName}

	goPackages, err = packages.Load(packagesConfig, "./...")
	if err != nil {
		slog.Error(fmt.Sprintf("%s file unparsable", goModFileName))
		return err
	}

	slog.Info(fmt.Sprintf("%d go packages found", len(goPackages)))
	if config.DebugMode {
		for _, n := range goPackages {
			slog.Debug(fmt.Sprintf(" %s : %s", n.Name, n.PkgPath))
		}
	}

	return nil
}

func (g *GoGraphBuilder) BuildGraph() {
	for _, pkg := range goPackages {
		pkgNode, newCreation := tryToCreateNode(pkg.PkgPath)
		if pkgNode != nil {
			if newCreation {
				slog.Debug(fmt.Sprintf("create node %v", pkgNode))
			}
			for _, imp := range pkg.Imports {
				impNode, newCreation := tryToCreateNode(imp.PkgPath)
				if impNode != nil {
					if newCreation {
						slog.Debug(fmt.Sprintf("create node %v", impNode))
					}

					if pkgNode.Name != impNode.Name {
						e := &graph.Edge{From: pkgNode.Name, To: impNode.Name}
						if _, ok := graph.Edges[e.String()]; !ok {
							graph.Edges[e.String()] = e
							slog.Debug(fmt.Sprintf("create edge %s", e))
						}
					}
				} else {
					slog.Debug(fmt.Sprintf("can not create node with path %s according to %s content", imp.PkgPath, config.DepCheckFileName))
				}
			}
		} else {
			slog.Debug(fmt.Sprintf("can not create node with path %s according to %s content", pkg.PkgPath, config.DepCheckFileName))
		}
	}

}

func tryToCreateNode(pkgPath string) (n *graph.Node, create bool) {
	// 1 external
	n, create = createOrGetNode(pkgPath, config.DepCheckConfig.External, graph.External)
	if n == nil {
		// 2 utility
		n, create = createOrGetNode(pkgPath, config.DepCheckConfig.Utility, graph.Utility)
		if n == nil {
			// 3 common
			n, create = createOrGetNode(pkgPath, config.DepCheckConfig.Common, graph.Common)
			if n == nil {
				//4 service
				n, create = createOrGetNode(pkgPath, config.DepCheckConfig.Service, graph.Service)
			}
		}
	}
	return
}

func createOrGetNode(pkgPath string, artifactOptions []string, nodeType graph.NodeType) (*graph.Node, bool) {
	slog.Debug(fmt.Sprintf("try to create node with path %s of type %s", pkgPath, nodeType))
	slog.Debug(fmt.Sprintf("considering artifactOptions %v", artifactOptions))

	for _, artif := range artifactOptions {
		prefix := path.Join(moduleName, artif)
		hasWildcard := strings.HasSuffix(prefix, "*")
		if hasWildcard {
			prefix = strings.TrimSuffix(prefix, "*")
		}

		if strings.HasPrefix(pkgPath, prefix) {
			name := artif
			if hasWildcard {
				name = strings.Replace(artif, "*", strings.Split(strings.Replace(pkgPath, prefix, "", 1), "/")[0], 1)
				prefix = path.Join(moduleName, name)
			}

			if _, ok := graph.Nodes[name]; !ok {
				n := &graph.Node{Name: name, Path: prefix, NodeType: nodeType}
				slog.Debug(fmt.Sprintf(" creating node %v", n))
				graph.Nodes[name] = n
				return n, true
			} else {
				n := graph.Nodes[name]
				slog.Debug(fmt.Sprintf(" node %v already exist", n))
				return n, false
			}
		}
	}
	slog.Debug(" no node created")
	return nil, false
}
