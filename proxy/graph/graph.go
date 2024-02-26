package graph

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type Node struct {
	ID      int
	Name    string
	Form    string // "circle", "rect", "square", "ellipse", "round-rect", "rhombus"
	Links   []*Node
	Visited bool
}

func (n *Node) AddLink(node *Node) {
	n.Links = append(n.Links, node)
}

func NewNode(id int, name, form string) *Node {
	return &Node{ID: id, Name: name, Form: form}
}
func getRandomNodes(nodes []*Node, count int) []*Node {
	if count > len(nodes) {
		count = len(nodes)
	}
	rand.Shuffle(len(nodes), func(i, j int) { nodes[i], nodes[j] = nodes[j], nodes[i] })
	return nodes[:count]
}

func ToMermaid(nodes []*Node) string {
	var sb strings.Builder
	for _, node := range nodes {
		for _, linkedNode := range node.Links {
			if !node.Visited {
				sb.WriteString(fmt.Sprintf("%s%s --> %s%s\n", node.Name, node.Form, linkedNode.Name, linkedNode.Form))
				node.Visited = true
			}
			if !linkedNode.Visited {
				sb.WriteString(fmt.Sprintf("%s --> %s%s\n", node.Name, linkedNode.Name, linkedNode.Form))
				linkedNode.Visited = true
			}
		}
	}
	return sb.String()
}
func GenerateMermaid() string {
	rand.Seed(time.Now().UnixNano())
	forms := []string{"[Square Rect]", "((Circle))", "(Round Rect)", "{Rhombus}", "[Ellipse]"}
	nodes := make([]*Node, 0)
	nodeCount := rand.Intn(26) + 5
	for i := 0; i < nodeCount; i++ {
		id := i + 1
		name := "A" + strconv.Itoa(id)
		form := forms[rand.Intn(len(forms)-1)]
		nodes = append(nodes, NewNode(id, name, form))
	}

	for i := 0; i < len(nodes); i++ {
		linkCount := rand.Intn(len(nodes)-1) + 1
		for _, linkedNode := range getRandomNodes(nodes, linkCount) {
			nodes[i].AddLink(linkedNode)
		}
	}

	return fmt.Sprintf("{{< mermaid >}} graph LR \n%s\n%s",
		ToMermaid(nodes), "{{< /mermaid >}}\n\n{{< /columns >}}")

}
