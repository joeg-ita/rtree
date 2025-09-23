package src

import (
	"fmt"
	"strings"
)

const ROOT = "ROOT"

type Node struct {
	Key        string
	Value      string
	Children   map[string]*Node
	IsEnd      bool
	parentNode *Node
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
		Children:   map[string]*Node{},
		IsEnd:      true,
		parentNode: nil,
	}
}

func NewRTree() *RTree {
	return &RTree{
		Root: &Node{
			Key:      ROOT,
			Children: map[string]*Node{},
			IsEnd:    false,
		},
	}
}

func (r *RTree) AddNodesToChildren(parentNode *Node, nodes ...*Node) *Node {
	parentNode.Children = r.appendNodesToMap(parentNode.Children, nodes...)
	return parentNode
}

func (r *RTree) AddChildrenToNodeChildren(parentNode *Node, nodesToAdd ...map[string]*Node) *Node {
	for _, v := range nodesToAdd {
		parentNode.Children = r.appendChildrenMapToMap(parentNode.Children, v)
	}
	return parentNode
}

func (r *RTree) DeleteNodeFromChildren(parentNode *Node, key string) *Node {
	delete(parentNode.Children, key)
	return parentNode
}

func (tree *RTree) Add(key string, value string) bool {
	return tree.addHandler(key, value, tree.Root)
}

func (r *RTree) addHandler(key string, value string, node *Node) bool {
	result := false

	if key == node.Key {
		return false
	}

	// Add when is empty
	if len(node.Children) == 0 {
		newNode := NewNode(key, value)
		r.AddNodesToChildren(node, newNode)
		return true
	}

	// Is not empty check the key
	tmpKey := ""
	tmpKeyAlreadyPresent := false
	tmpKeyOffset := ""
	tmpKeyOrphan := ""
	childKey := ""
	var selectedNode *Node
	for k, n := range node.Children {
		selectedNode = n
		childKey = k
		tmpKey = ""
		tmpKeyOffset = ""
		tmpKeyOrphan = ""
		if len(n.Key) == 0 || len(key) == 0 || n.Key[0] != key[0] {
			selectedNode = nil
			childKey = ""
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
		if len(tmpKey) < len(key) {
			tmpKeyOffset = key[len(tmpKey):]
		}
		if len(tmpKey) < len(n.Key) {
			tmpKeyOrphan = n.Key[len(tmpKey):]
		}
		break
	}

	if tmpKeyAlreadyPresent && tmpKeyOffset == "" {
		currentNode := node.Children[childKey]
		currentNode.IsEnd = true
		currentNode.Value = value
		return true
	}
	if selectedNode != nil && tmpKeyAlreadyPresent && tmpKeyOffset != "" {
		return r.addHandler(tmpKeyOffset, value, selectedNode)
	}
	if tmpKeyOrphan == "" && tmpKey != "" && tmpKeyOffset != "" {
		return r.addHandler(tmpKeyOffset, value, selectedNode)
	} else if tmpKeyOrphan == "" && tmpKey == "" && tmpKeyOffset == "" {
		newNode := NewNode(key, value)
		r.AddNodesToChildren(node, newNode)
		return true
	} else if tmpKeyOrphan != "" && tmpKey != "" && tmpKeyOffset != "" {
		currentNode := node.Children[childKey]

		originalIsEnd := currentNode.IsEnd
		originalValue := currentNode.Value

		currentNode.Key = tmpKey
		currentNode.IsEnd = false
		currentNode.Value = ""

		delete(node.Children, childKey)
		node.Children[tmpKey] = currentNode

		orphanNode := NewNode(tmpKeyOrphan, "")
		if originalIsEnd {
			orphanNode.IsEnd = true
			orphanNode.Value = originalValue
		}

		r.AddChildrenToNodeChildren(orphanNode, currentNode.Children)
		currentNode.Children = map[string]*Node{}
		r.AddNodesToChildren(currentNode, orphanNode)
		newNode := NewNode(tmpKeyOffset, value)
		r.AddNodesToChildren(currentNode, newNode)

		return true
	} else if tmpKeyOrphan != "" && tmpKey != "" && tmpKeyOffset == "" {
		currentNode := node.Children[childKey]
		currentNode.Key = tmpKey

		orphanNode := NewNode(tmpKeyOrphan, currentNode.Value)
		if currentNode.IsEnd {
			orphanNode.IsEnd = true
		}

		currentNode.IsEnd = true
		currentNode.Value = value

		delete(node.Children, childKey)
		node.Children[tmpKey] = currentNode

		r.AddChildrenToNodeChildren(orphanNode, currentNode.Children)
		currentNode.Children = map[string]*Node{}
		r.AddNodesToChildren(currentNode, orphanNode)

		return true
	} else if tmpKeyOffset != "" && tmpKeyOrphan == "" {
		newNode := NewNode(tmpKeyOffset, value)
		r.AddNodesToChildren(node.Children[childKey], newNode)
		return true
	}

	return result
}

func (tree *RTree) Search(key string) *Node {
	return tree.searchHandler(key, "", tree.Root, nil, 0)
}

func (r *RTree) searchHandler(key string, foundedKeyPart string, node *Node, parentNode *Node, level int) *Node {

	search := true
	for search {
		search = false

		keyToCheck := fmt.Sprintf("%s%s", foundedKeyPart, key)
		nod, exists := node.Children[keyToCheck]
		if exists && nod.IsEnd {
			nod.parentNode = node
			return nod
		} else {

			for k, child := range node.Children {

				nodeKey := fmt.Sprintf("%s%s", foundedKeyPart, k)
				if keyToCheck == nodeKey && child.IsEnd {
					child.parentNode = node
					return child
				}

				tmpFoundedKeyParts := ""
				tmpKey := keyToCheck

				if strings.HasPrefix(tmpKey, nodeKey) {
					tmpFoundedKeyParts = nodeKey
					tmpKey = keyToCheck[len(nodeKey):]
				} else {
					continue
				}

				key = tmpKey
				foundedKeyPart = tmpFoundedKeyParts
				search = true
				parentNode = node
				node = child
				level = level + 1
				break
			}

		}

	}
	return nil
}

func (tree *RTree) Delete(key string) bool {
	node := tree.Search(key)
	if node != nil && node.parentNode != nil && node.IsEnd && len(node.Children) == 0 {
		tree.DeleteNodeFromChildren(node.parentNode, node.Key)
		return true
	} else if node != nil && node.IsEnd && len(node.Children) > 0 {
		node.IsEnd = false
		// compactHandler(node)
		return true
	}
	return false
}

func (tree *RTree) Compact() {
	root := tree.Root
	for _, n := range root.Children {
		tree.compactHandler(n)
	}

}

func (r *RTree) compactHandler(node *Node) {
	if len(node.Children) == 1 && !node.IsEnd {
		var child *Node
		for _, value := range node.Children {
			child = value
		}
		node.Key = fmt.Sprintf("%s%s", node.Key, child.Key)
		node.IsEnd = child.IsEnd
		node.Value = child.Value
		node.Children = child.Children
	}
}

func (r *RTree) appendToMap(m1 map[string]*Node, m2 map[string]*Node) map[string]*Node {
	for key, value := range m2 {
		m1[key] = value
	}
	return m1
}

func (r *RTree) appendChildrenMapToMap(m1 map[string]*Node, nodesToAdd ...map[string]*Node) map[string]*Node {
	for _, mapToAdd := range nodesToAdd {
		r.appendToMap(m1, mapToAdd)
	}
	return m1
}

func (r *RTree) appendNodesToMap(m1 map[string]*Node, nodes ...*Node) map[string]*Node {
	for _, node := range nodes {
		m1[node.Key] = node
	}
	return m1
}
