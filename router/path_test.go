package router

import (
	"regexp"
	"testing"
)

// should create a tree node
func TestPathMakeNode(t *testing.T) {
	re, _ := regexp.Compile("test")

	node := makePath(re)

	if node == nil {
		t.Fatal("failed to create new node")
	}
	if node.Regexp.String() != "test" {
		t.Errorf("unexpected name (%s)", node.Regexp.String())
	}
	if node.Children == nil {
		t.Error("node children is nil")
	}
}

// should create a root node
func TestPathMakeRoot(t *testing.T) {
	re, _ := regexp.Compile("")

	root := makePath(re)

	if root == nil {
		t.Fatal("failed to create new node")
	}
	if root.Regexp != nil {
		t.Errorf("unexpected regexp")
	}
}

// should find a child in a tree's immediate children
func TestPathFindChild(t *testing.T) {
	re, _ := regexp.Compile("")
	tree := makePath(re)
	re, _ = regexp.Compile("test")
	child := makePath(re)
	tree.Children = append(tree.Children, child)

	found := tree.findChild("test")

	if found == nil {
		t.Error("did not find child")
	}
}

// should not find a child that is missing from the immediate children
func TestPathFindMissingChild(t *testing.T) {
	re, _ := regexp.Compile("")
	tree := makePath(re)
	re, _ = regexp.Compile("test")
	child := makePath(re)
	tree.Children = append(tree.Children, child)

	found := tree.findChild("missing")

	if found != nil {
		t.Error("found a missing child")
	}
}

// should find a child using regular expression
func TestPathFindRegexChild(t *testing.T) {
	re, _ := regexp.Compile("")
	tree := makePath(re)
	re, _ = regexp.Compile(`\d`)
	child := makePath(re)
	tree.Children = append(tree.Children, child)

	found := tree.findChild("123")

	if found == nil {
		t.Error("did not find regular expression node")
	}
}

// should add a child to a tree's children list
func TestPathAddChild(t *testing.T) {
	re, _ := regexp.Compile("")
	tree := makePath(re)
	re, _ = regexp.Compile("child")
	child := makePath(re)

	tree.addChild(child)

	if len(tree.Children) != 1 {
		t.Error("failed to add child")
	}
}

// should not add a nil to a tree's children list
func TestPathAddNilChild(t *testing.T) {
	re, _ := regexp.Compile("")
	tree := makePath(re)

	tree.addChild(nil)

	if len(tree.Children) != 0 {
		t.Error("added nil child")
	}
}

// should not add a child with the same name
func TestPathAddSameChild(t *testing.T) {
	re, _ := regexp.Compile("")
	tree := makePath(re)
	re, _ = regexp.Compile("child")
	child := makePath(re)
	tree.Children = append(tree.Children, child)

	re, _ = regexp.Compile("child")
	tree.addChild(makePath(re))

	if len(tree.Children) > 1 {
		t.Error("added child of the same name")
	}
}

// should have well-defined children
func TestPathMakeChildren(t *testing.T) {
	re, _ := regexp.Compile("")
	node := makePath(re)
	if node == nil {
		t.Fatal("failed to create new node")
	}
	if node.Children == nil {
		t.Error("node children is nil")
	}
}

// should do nothing if no path is given
func TestPathNoPath(t *testing.T) {
	re, _ := regexp.Compile("")
	tree := makePath(re)

	child := tree.add([]*regexp.Regexp{})

	if tree != child {
		t.Error("expected tree to be returned")
	}

	if len(tree.Children) > 0 {
		t.Error("children should be empty")
	}
}

// should make a tree with a path
func TestPathCreatePath(t *testing.T) {
	re, _ := regexp.Compile("")
	tree := makePath(re)

	var res []*regexp.Regexp
	re, _ = regexp.Compile("test")
	res = append(res, re)
	re, _ = regexp.Compile("child")
	res = append(res, re)

	tree.add(res)

	if tree.Children == nil || len(tree.Children) != 1 {
		t.Error("invalid tree children")
	}

	test := tree.findChild("test")
	if test == nil {
		t.Fatal("could not find \"test\".")
	}
	if test.findChild("child") == nil {
		t.Error("could not find \"child\".")
	}
}

// should add a child to an existing tree
func TestPathAddToExisting(t *testing.T) {
	re, _ := regexp.Compile("")
	tree := makePath(re)
	re, _ = regexp.Compile("child1")
	child := makePath(re)

	tree.addChild(child)

	var res []*regexp.Regexp
	re, _ = regexp.Compile("child1")
	res = append(res, re)
	re, _ = regexp.Compile("child2")
	res = append(res, re)

	tree.add(res)

	if child.findChild("child2") == nil {
		t.Error("could not find \"child2\".")
	}
}

// should not add anything if the path exists
func TestPathIdentity(t *testing.T) {
	re, _ := regexp.Compile("")
	tree := makePath(re)
	re, _ = regexp.Compile("child1")
	child1 := makePath(re)
	re, _ = regexp.Compile("child2")
	child2 := makePath(re)

	tree.addChild(child1)
	child1.addChild(child2)

	var res []*regexp.Regexp
	re, _ = regexp.Compile("child1")
	res = append(res, re)
	re, _ = regexp.Compile("child2")
	res = append(res, re)

	tree.add(res)

	if tree.findChild("child1") != child1 {
		t.Error("unexpected child1")
	}

	if child1.findChild("child2") != child2 {
		t.Error("unexpected child2")
	}

	if len(child2.Children) != 0 {
		t.Error("len(child2 children) != 0")
	}
}

// should add an adjacent node
func TestPathAddAdjacent(t *testing.T) {
	re, _ := regexp.Compile("")
	tree := makePath(re)
	re, _ = regexp.Compile("child1")
	child1 := makePath(re)
	re, _ = regexp.Compile("child2")
	child2 := makePath(re)

	tree.addChild(child1)
	child1.addChild(child2)

	var res []*regexp.Regexp
	re, _ = regexp.Compile("child1")
	res = append(res, re)
	re, _ = regexp.Compile("child3")
	res = append(res, re)

	tree.add(res)

	if len(tree.Children) != 1 {
		t.Errorf("len(tree children) != 1 (%d)", len(tree.Children))
	}

	if len(child1.Children) != 2 {
		t.Error("len(child1 children) != 2")
	}
}

// should find a simple path
func TestPathFindSimple(t *testing.T) {
	re, _ := regexp.Compile("")
	tree := makePath(re)
	re, _ = regexp.Compile("child")
	child := makePath(re)

	tree.addChild(child)

	found := tree.find([]string{"child"})

	if found == nil || found.Regexp.String() != "child" {
		t.Error("did not find \"child\"")
	}
}

// should find a deep path
func TestPathFindDeep(t *testing.T) {
	re, _ := regexp.Compile("")
	tree := makePath(re)
	re, _ = regexp.Compile("child1")
	child1 := makePath(re)
	re, _ = regexp.Compile("child2")
	child2 := makePath(re)

	tree.addChild(child1)
	child1.addChild(child2)

	found := tree.find([]string{"child1", "child2"})

	if found == nil || found.Regexp.String() != "child2" {
		t.Error("did not find \"child2\"")
	}
}

// should find nothing in an empty tree
func TestPathFindEmpty(t *testing.T) {
	re, _ := regexp.Compile("")
	tree := makePath(re)
	if found := tree.find([]string{"test"}); found != nil {
		t.Error("expected nil when searching empty tree")
	}
}

// should not find a missing path
func TestPathFindMissingPath(t *testing.T) {
	re, _ := regexp.Compile("")
	tree := makePath(re)
	re, _ = regexp.Compile("child1")
	child := makePath(re)

	tree.addChild(child)

	if found := tree.find([]string{"missing"}); found != nil {
		t.Error("expected nil when searching missing path")
	}
}

// should find a node through Regexp application
func TestPathFindRegexp(t *testing.T) {
	re, _ := regexp.Compile("")
	tree := makePath(re)
	re, _ = regexp.Compile(`\d+`)
	child1 := makePath(re)
	re, _ = regexp.Compile("test")
	child2 := makePath(re)
	tree.addChild(child1)

	tree.addChild(child1)
	child1.addChild(child2)

	found := tree.find([]string{"123", "test"})
	if found != child2 {
		t.Error("did not find \"test\"")
	}
}
