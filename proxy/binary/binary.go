package binary

import (
	"fmt"
	"math/rand"
	"time"
)

type Node struct {
	Key    int
	Height int
	Left   *Node
	Right  *Node
}

type AVLTree struct {
	Root *Node
}

func NewNode(key int) *Node {
	return &Node{Key: key, Height: 1}
}

func (t *AVLTree) Insert(key int) {
	t.Root = insert(t.Root, key)
}

func (t *AVLTree) ToMermaid() string {
	return fmt.Sprintf("{{< /columns >}}\n\n{{< mermaid >}} graph TD\n %s\n%s", nodeToMermaid(t.Root), "{{< /mermaid >}}")
}

func nodeToMermaid(node *Node) string {
	if node == nil {
		return ""
	}
	var result string
	if node.Left != nil {
		result += fmt.Sprintf("%d --> %d\n", node.Key, node.Left.Key)
		result += nodeToMermaid(node.Left)
	}
	if node.Right != nil {
		result += fmt.Sprintf("%d --> %d\n", node.Key, node.Right.Key)
		result += nodeToMermaid(node.Right)
	}
	return result
}

func height(node *Node) int {
	if node == nil {
		return 0
	}
	return node.Height
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func updateHeight(node *Node) {
	node.Height = max(height(node.Left), height(node.Right)) + 1
}

func getBalance(node *Node) int {
	if node == nil {
		return 0
	}
	return height(node.Left) - height(node.Right)
}

func leftRotate(x *Node) *Node {
	y := x.Right
	t2 := y.Left

	y.Left = x
	x.Right = t2

	updateHeight(x)
	updateHeight(y)

	return y
}

func rightRotate(y *Node) *Node {
	x := y.Left
	t2 := x.Right

	x.Right = y
	y.Left = t2

	updateHeight(y)
	updateHeight(x)

	return x
}

func insert(node *Node, key int) *Node {
	if node == nil {
		return NewNode(key)
	}

	if key < node.Key {
		node.Left = insert(node.Left, key)
	} else if key > node.Key {
		node.Right = insert(node.Right, key)
	} else {
		return node
	}

	updateHeight(node)

	balance := getBalance(node)

	// Left Left Case
	if balance > 1 && key < node.Left.Key {
		return rightRotate(node)
	}
	// Right Right Case
	if balance < -1 && key > node.Right.Key {
		return leftRotate(node)
	}
	// Left Right Case
	if balance > 1 && key > node.Left.Key {
		node.Left = leftRotate(node.Left)
		return rightRotate(node)
	}
	// Right Left Case
	if balance < -1 && key < node.Right.Key {
		node.Right = rightRotate(node.Right)
		return leftRotate(node)
	}

	return node
}

func GenerateTree(count int) *AVLTree {
	rand.Seed(time.Now().UnixNano())
	tree := AVLTree{}
	for i := 0; i < count; i++ {
		key := rand.Intn(100)
		tree.Insert(key)
	}
	return &tree
}
