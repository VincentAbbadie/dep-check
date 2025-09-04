package graph

import (
	"fmt"
)

// NODE definition
type NodeType string

const (
	External NodeType = "external"
	Common   NodeType = "common"
	Service  NodeType = "service"
	Utility  NodeType = "utility"
)

type Node struct {
	Name, Path string
	NodeType   NodeType
}

func (n *Node) String() string {
	return fmt.Sprintf("[%s : %s <%s>]", n.Path, n.Name, n.NodeType)
}

// EDGE definition
type Edge struct {
	From, To string
}

func (e *Edge) String() string {
	return fmt.Sprintf("%s --> %s", e.From, e.To)
}

func (e *Edge) IsValid() (bool, string) {
	retrunMessageTemplate := "%s can not import %s"

	switch Nodes[e.From].NodeType {
	case Common: // Common can not import Utility
		if Nodes[e.To].NodeType == Utility {
			return false, fmt.Sprintf(retrunMessageTemplate, Nodes[e.From], Nodes[e.To])
		}
	case Service: // Service cannot import Utility or Service
		if Nodes[e.To].NodeType == Utility || Nodes[e.To].NodeType == Service {
			return false, fmt.Sprintf(retrunMessageTemplate, Nodes[e.From], Nodes[e.To])
		}
	}
	return true, ""
}

// GRAPH definition
var (
	Nodes = make(map[string]*Node)
	Edges = make(map[string]*Edge)
)

type GraphBuilder interface {
	Init() error
	BuildGraph()
}

// func FindNodeFromFilePath(filePath string) (*Node, error) {
// 	for _, node := range Nodes {
// 		if strings.Contains(filePath, node.Path) {
// 			return node, nil
// 		}
// 	}
// 	return nil, fmt.Errorf("no node found for file %s", filePath)
// }

// func (n *Node) GetParents() (p []*Node) {
// 	queue := []*Node{n}
// 	visited := make(map[string]bool)

// 	for len(queue) > 0 {
// 		cur := queue[0]
// 		queue = queue[1:]
// 		if _, ok := visited[cur.Name]; !ok {
// 			visited[cur.Name] = true
// 			p = append(p, cur)
// 			for _, edge := range Edges {
// 				if edge.To == cur.Name {
// 					queue = append(queue, Nodes[edge.From])
// 				}
// 			}
// 		}
// 	}

// 	return
// }
