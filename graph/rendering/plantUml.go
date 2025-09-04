package rendering

import (
	"fmt"

	"github.com/moveaxlab/dep-check/graph"
)

const (
	plantUmlFileName    = "plantUml_graph.txt"
	plantUmlGraphHeader = "@startuml"
	plantUmlGraphFooter = "@enduml"
)

func PlantUmlNode(n *graph.Node) string {
	return fmt.Sprintf("object %s", n.Name)
}

func PlantUmlEdge(e *graph.Edge) string {
	return e.String()
}
