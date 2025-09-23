package test

import (
	"fmt"
	r "rtree/src"
	"testing"
)

func TestAddNodesToChildren(t *testing.T) {

	rtree := r.NewRTree()

	node := r.NewNode("test", "test value")

	rtree.AddNodesToChildren(rtree.Root, node)

	if len(rtree.Root.Children) != 1 {
		t.Errorf(`rtree.Root.Children want 1 match  %d, nil`, len(rtree.Root.Children))
	}
}

func TestDeleteNodeFromChildren(t *testing.T) {

	rtree := r.NewRTree()

	node1 := r.NewNode("test1", "test value 1")
	node2 := r.NewNode("test2", "test value 2")
	node3 := r.NewNode("test3", "test value 3")

	nodes := []*r.Node{node1, node2, node3}
	rtree.AddNodesToChildren(rtree.Root, nodes...)

	rtree.DeleteNodeFromChildren(rtree.Root, node1.Key)

	_, exists1 := rtree.Root.Children[node1.Key]
	value3, _ := rtree.Root.Children[node3.Key]

	if len(rtree.Root.Children) != 2 ||
		exists1 ||
		value3.Key != "test3" {
		t.Errorf(`rtree.Root.Children error len=%d`, len(rtree.Root.Children))
	}
}

func TestAdd(t *testing.T) {

	rtree := r.NewRTree()

	keys := []string{"ciao", "ciaone", "ciauz", "help", "helper", "cia", "test"}

	for _, k := range keys {
		rtree.Add(k, fmt.Sprintf("val of %s", k))
	}
	r.PrintNode(rtree.Root, true)

	if len(rtree.Root.Children) != 3 {
		t.Errorf(`rtree.Root.Children error len=%d`, len(rtree.Root.Children))
	}
	value, _ := rtree.Root.Children[keys[5]]
	if len(value.Children) != 2 {
		t.Errorf(`rtree.Root.Children[0].Children error len=%d`, len(value.Children))
	}
}

func TestAdd2(t *testing.T) {

	rtree := r.NewRTree()

	// keys := []string{"test", "testing", "tea"}
	keys := []string{"test", "te", "testing", "tea"}

	for _, k := range keys {
		rtree.Add(k, fmt.Sprintf("val of %s", k))
	}
	r.PrintNode(rtree.Root, true)

	if len(rtree.Root.Children) != 1 {
		t.Errorf(`rtree.Root.Children error len=%d`, len(rtree.Root.Children))
	}
}

func TestAdd3(t *testing.T) {

	rtree := r.NewRTree()

	keys := []string{
		"bef9e715-9e22-441a-9029-612bfa335e86",
		"e7712def-ae1f-4f72-9c01-4d0bc64905d4",
		"41d192c5-fd52-4209-bbda-8e1318b5c935",
		"b325a99d-dc17-4358-887f-8b686cf6eeb8",
		"e7d833c4-3e1a-49e9-8338-f7a2f9eeddf1"}

	for _, k := range keys {
		rtree.Add(k, fmt.Sprintf("val of %s", k))
	}
	r.PrintNode(rtree.Root, true)

	if len(rtree.Root.Children) != 4 {
		t.Errorf(`rtree.Root.Children error len=%d`, len(rtree.Root.Children))
	}
	value, _ := rtree.Root.Children[keys[0]]
	if len(value.Children) != 2 {
		t.Errorf(`rtree.Root.Children[0].Children error len=%d`, len(value.Children))
	}
}

func TestSearchKeyIsPresent(t *testing.T) {

	rtree := r.NewRTree()

	keys := []string{"ciao", "ciaone", "ciauz", "help", "helper", "cia", "test"}

	for _, k := range keys {
		rtree.Add(k, fmt.Sprintf("val of %s", k))
	}

	key := "ciauz"
	node := rtree.Search(key)
	if node != nil {
		r.PrintNode(node, false)
	} else {
		t.Errorf(`Not Found expected key %s`, key)
	}
}

func TestSearchKeyIsNotPresent(t *testing.T) {

	rtree := r.NewRTree()

	keys := []string{"ciao", "ciaone", "ciauz", "help", "helper", "cia", "test"}

	for _, k := range keys {
		rtree.Add(k, fmt.Sprintf("val of %s", k))
	}

	key := "hello"
	nodeHello := rtree.Search(key)
	if nodeHello != nil {
		t.Errorf(`Found unexpected key %s`, key)
	} else {
		fmt.Printf("Key [%s] not found\n", key)
	}

}

func TestSearchInternmediateKeyIsPresent(t *testing.T) {

	rtree := r.NewRTree()

	keys := []string{"ciao", "ciaone", "ciauz", "help", "helper", "cia", "test"}

	for _, k := range keys {
		rtree.Add(k, fmt.Sprintf("val of %s", k))
	}

	r.PrintNode(rtree.Root, true)

	key := "cia"
	node := rtree.Search(key)
	if node != nil {
		r.PrintNode(node, false)
	} else {
		t.Errorf(`Not Found expected key %s`, key)
	}

}

func TestSearchChunkOfPresentKey(t *testing.T) {

	rtree := r.NewRTree()

	keys := []string{"ciao", "ciaone", "ciauz", "help", "helper", "cia", "test"}

	for _, k := range keys {
		rtree.Add(k, fmt.Sprintf("val of %s", k))
	}

	key := "uz"
	node := rtree.Search(key)
	if node != nil {
		t.Errorf(`Found unexpected key %s`, key)
	} else {
		fmt.Printf("Key [%s] not found\n", key)
	}

}

func TestCompact(t *testing.T) {
	rtree := r.NewRTree()

	keys := []string{"ciao", "ciaone", "ciauz", "help", "helper", "cia", "test"}
	for _, k := range keys {
		rtree.Add(k, fmt.Sprintf("val of %s", k))
	}
	r.PrintNode(rtree.Root, true)
	key := "ciauz"
	result := rtree.Delete(key)
	r.PrintNode(rtree.Root, true)
	if !result {
		t.Errorf(`Fail to delete key %s`, key)
	}
}

func TestDeleteNodeWithoutChildren(t *testing.T) {

	rtree := r.NewRTree()

	keys := []string{"ciao", "ciaone", "ciauz", "help", "helper", "cia", "test"}
	for _, k := range keys {
		rtree.Add(k, fmt.Sprintf("val of %s", k))
	}

	keyToDelete := "test"
	deleteResult := rtree.Delete(keyToDelete)

	if !deleteResult {
		r.PrintNode(rtree.Root, true)
		t.Errorf(`Error deleting key %s`, keyToDelete)
	}
}

func TestDeleteNodeWithChildrenAndCompact(t *testing.T) {

	rtree := r.NewRTree()

	keys := []string{"ciao", "ciaone", "ciauz", "help", "helper", "cia", "test"}
	for _, k := range keys {
		rtree.Add(k, fmt.Sprintf("val of %s", k))
	}

	keyToDelete := "ciao"
	deleteResult := rtree.Delete(keyToDelete)

	if !deleteResult {
		r.PrintNode(rtree.Root, true)
		t.Errorf(`Error deleting key %s`, keyToDelete)
	}
}
