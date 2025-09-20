package src

import (
	"fmt"
	"strings"
)

const ROOT = "ROOT"

type Node struct {
	Key        string
	Value      string
	Children   []*Node
	IsEnd      bool
	parentNode *Node
	index      int
}

type RTree struct {
	Root *Node
}

func PrintNode(node *Node, printChildren bool) {
	fmt.Println("node", &node)
	fmt.Println("key", node.Key)
	fmt.Println("value", node.Value)
	fmt.Println("isEnd", node.IsEnd)
	fmt.Println("parentNode", node.parentNode)
	fmt.Println("children len", len(node.Children))
	if printChildren {
		for _, n := range node.Children {
			PrintNode(n, printChildren)
		}
	}
}

func NewNode(key string, value string) *Node {
	return &Node{
		Key:        key,
		Value:      value,
		Children:   []*Node{},
		IsEnd:      true,
		parentNode: nil,
	}
}

func NewRTree() *RTree {
	return &RTree{
		Root: &Node{
			Key:      ROOT,
			Children: []*Node{},
			IsEnd:    false,
		},
	}
}

func AddNodesToChildren(parentNode *Node, nodesToAdd ...*Node) *Node {
	parentNode.Children = append(parentNode.Children, nodesToAdd...)
	return parentNode
}

func DeleteNodeFromChildren(parentNode *Node, index int) *Node {
	fmt.Printf("Deleting children node with key=[%s] from key=[%s]\n", parentNode.Children[index].Key, parentNode.Key)
	parentNode.Children = append(parentNode.Children[:index], parentNode.Children[index+1:]...)
	return parentNode
}

func Add(key string, value string, tree *RTree) bool {
	fmt.Printf("Add key=[%s] value=[%s]\n", key, value)
	return addHandler(key, value, tree.Root)
}

func addHandler(key string, value string, node *Node) bool {
	result := false

	if key == node.Key {
		fmt.Printf("key %s already present n.Key %s\n", key, node.Key)
		return false
	}

	// Add when is empty
	if len(node.Children) == 0 {
		newNode := NewNode(key, value)
		AddNodesToChildren(node, newNode)
		return true
	}

	// Is not empty check the key
	tmpKey := ""
	tmpKeyAlreadyPresent := false
	tmpKeyOffset := ""
	tmpKeyOrphan := ""
	childrenIndex := -1
	var selectedNode *Node
	for i, n := range node.Children {
		selectedNode = n
		childrenIndex = i
		tmpKey = ""
		tmpKeyOffset = ""
		tmpKeyOrphan = ""
		if n.Key[0] != key[0] {
			selectedNode = nil
			childrenIndex = -1
			continue
		}
		for j := 0; j < len(n.Key) && j < len(key); j++ {
			if key[j] == n.Key[j] {
				tmpKey = key[:j+1]
			} else {
				break
			}
		}
		if tmpKey == n.Key {
			tmpKeyAlreadyPresent = true
		}
		if len(key) > len(n.Key) {
			tmpKeyOffset = key[len(tmpKey):]
		}
		if len(tmpKey) < len(n.Key) {
			tmpKeyOrphan = n.Key[len(tmpKey):]
		}
		if tmpKeyOrphan != "" {
			fmt.Println("tmpKeyOrphan", tmpKeyOrphan)
		}
		if childrenIndex != -1 {
			fmt.Println("childrenIndex", childrenIndex)
		}
		if tmpKey != "" {
			fmt.Println("tmpKey", tmpKey)
		}
		if tmpKeyOffset != "" {
			fmt.Println("tmpKeyOffset", tmpKeyOffset)
		}
		if tmpKeyAlreadyPresent {
			fmt.Println("tmpKeyAlreadyPresent", tmpKeyAlreadyPresent)
		}
		fmt.Println("--------------------------")
		break
	}
	var newNode *Node
	if selectedNode != nil && len(selectedNode.Children) > 0 {
		fmt.Println("len(currentNode.Children)", len(selectedNode.Children))
	}
	if tmpKeyAlreadyPresent && tmpKeyOffset == "" {
		currentNode := node.Children[childrenIndex]
		currentNode.IsEnd = true
		currentNode.Value = value
		return true
	}
	if selectedNode != nil && tmpKeyAlreadyPresent && tmpKeyOffset != "" {
		return addHandler(tmpKeyOffset, value, selectedNode)
	}
	if tmpKeyOrphan == "" && tmpKey == "" && tmpKeyOffset == "" {
		newNode = NewNode(key, value)
		AddNodesToChildren(node, newNode)
		return true
	} else if tmpKeyOrphan != "" && tmpKey != "" && tmpKeyOffset != "" {
		currentNode := node.Children[childrenIndex]

		originalIsEnd := currentNode.IsEnd
		originalValue := currentNode.Value

		currentNode.Key = tmpKey
		currentNode.IsEnd = false
		currentNode.Value = ""

		orphanNode := NewNode(tmpKeyOrphan, "")
		if originalIsEnd {
			orphanNode.IsEnd = true
			orphanNode.Value = originalValue
		}

		AddNodesToChildren(orphanNode, currentNode.Children...)
		currentNode.Children = []*Node{}
		AddNodesToChildren(currentNode, orphanNode)

		newNode := NewNode(tmpKeyOffset, value)
		AddNodesToChildren(currentNode, newNode)

		return true
	} else if tmpKeyOrphan != "" && tmpKey != "" && tmpKeyOffset == "" {
		currentNode := node.Children[childrenIndex]
		currentNode.Key = tmpKey

		orphanNode := NewNode(tmpKeyOrphan, currentNode.Value)
		if currentNode.IsEnd {
			orphanNode.IsEnd = true
		}

		currentNode.IsEnd = true
		currentNode.Value = value

		AddNodesToChildren(orphanNode, currentNode.Children...)
		currentNode.Children = []*Node{}
		AddNodesToChildren(currentNode, orphanNode)

		return true
	} else if tmpKeyOffset != "" && tmpKeyOrphan == "" {
		newNode = NewNode(tmpKeyOffset, value)
		AddNodesToChildren(node.Children[childrenIndex], newNode)
		return true
	}

	return result
}

func Search(key string, tree RTree) *Node {
	fmt.Printf("Search key=[%s]\n", key)
	return searchHandler(key, "", tree.Root, nil, 0, -1)
}

func searchHandler(key string, foundedKeyPart string, node *Node, parentNode *Node, level int, index int) *Node {
	// fmt.Printf("key=[%s] foundedKeyPart=[%s] node.Key=[%s]\n", key, foundedKeyPart, node.Key)
	keyToCheck := fmt.Sprintf("%s%s", foundedKeyPart, key)

	nodeKey := fmt.Sprintf("%s%s", foundedKeyPart, node.Key)
	// fmt.Printf("keyToCheck[%s]==[%s]nodeKey node.IsEnd=[%v]\n", keyToCheck, nodeKey, node.IsEnd)
	if keyToCheck == nodeKey && node.IsEnd {
		// fmt.Printf("Found at level %d\n", level)
		node.parentNode = parentNode
		node.index = index
		return node
	}

	tmpFoundedKeyParts := ""
	tmpKey := keyToCheck
	// fmt.Printf("strings.HasPrefix([%s], node.Key[%s])\n", keyToCheck, node.Key)
	if strings.HasPrefix(keyToCheck, node.Key) {
		tmpFoundedKeyParts = keyToCheck[:len(node.Key)]
		tmpKey = keyToCheck[len(node.Key):]
	}

	for i, n := range node.Children {
		node := searchHandler(tmpKey, tmpFoundedKeyParts, n, node, (level + 1), i)
		if node != nil {
			return node
		}
		if level == 1 && tmpFoundedKeyParts == "" {
			break
		}
	}
	return nil
}

func Delete(key string, tree *RTree) bool {
	fmt.Printf("Delete key=[%s]\n", key)
	node := Search(key, *tree)
	if node != nil && node.IsEnd && len(node.Children) == 0 {
		DeleteNodeFromChildren(node.parentNode, node.index)
		return true
	} else if node != nil && node.IsEnd && len(node.Children) > 0 {
		node.IsEnd = false
		compactHandler(node)
		return true
	}
	return false
}

func Compact(tree *RTree) {
	fmt.Printf("Compact tree\n")
	root := tree.Root
	for _, n := range root.Children {
		compactHandler(n)
	}

}

func compactHandler(node *Node) {
	if len(node.Children) == 1 && !node.IsEnd {
		child := node.Children[0]
		node.Key = fmt.Sprintf("%s%s", node.Key, child.Key)
		node.IsEnd = child.IsEnd
		node.Value = child.Value
		node.Children = child.Children
	}
}
