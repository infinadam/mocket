package urlpath

import "testing"

// should create a tree node
func TestMakeNode(t *testing.T) {
	node := makeTree("test")
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
func TestMakeRoot(t *testing.T) {
	root := makeTree("")
	if root == nil {
		t.Fatal("failed to create new node")
	}
	if root.Regexp != nil {
		t.Errorf("unexpected regexp")
	}
}

// should have well-defined children
func TestMakeChildren(t *testing.T) {
	node := makeTree("test")
	if node == nil {
		t.Fatal("failed to create new node")
	}
	if node.Children == nil {
		t.Error("node children is nil")
	}
}

// should do nothing if no path is given
func TestNoPath(t *testing.T) {
	tree := makeTree("test")

	tree = Add([]string{}, tree)

	if tree == nil {
		t.Fatal("tree is nil")
	}

	if name := tree.Regexp.String(); tree.Regexp.String() != "test" {
		t.Errorf("expected tree name \"test\", got %q", name)
	}

	if len(tree.Children) > 0 {
		t.Error("children should be empty")
	}
}

// should make a rooted tree from a path
func TestCreateTree(t *testing.T) {
	tree := Add([]string{"test"}, nil)

	if tree == nil {
		t.Fatal("tree is nil")
	}

	if tree.Children == nil || len(tree.Children) != 1 {
		t.Error("invalid tree children")
	}
}

// should make a tree with a path
func TestCreatePath(t *testing.T) {
	tree := Add([]string{"test", "child"}, nil)

	if tree == nil {
		t.Fatal("tree is nil")
	}

	if tree.Children == nil || len(tree.Children) != 1 {
		t.Error("invalid tree children")
	}

	test := tree.Children["test"]
	if test.Children == nil || test.Children["child"] == nil {
		t.Error("could not find \"child\".")
	}
}

// should add a child to an existing tree
func TestAddToExisting(t *testing.T) {
	tree := makeTree("")
	child := makeTree("child1")
	tree.Children["child1"] = child

	Add([]string{"child1", "child2"}, tree)

	if child.Children["child2"] == nil {
		t.Error("could not find \"child2\".")
	}
}

// should not add anything if the path exists
func TestIdentity(t *testing.T) {
	tree := makeTree("")
	child1 := makeTree("child1")
	child2 := makeTree("child2")
	tree.Children["child1"] = child1
	child1.Children["child2"] = child2

	Add([]string{"child1", "child2"}, tree)

	if len(tree.Children) != 1 {
		t.Error("len(tree children) != 1")
	}

	if len(child1.Children) != 1 {
		t.Error("len(child1 children) != 1")
	}

	if len(child2.Children) != 0 {
		t.Error("len(child2 children) != 0")
	}
}

// should add an adjacent node
func TestAddAdjacent(t *testing.T) {
	tree := makeTree("")
	child1 := makeTree("child1")
	child2 := makeTree("child2")
	tree.Children["child1"] = child1
	child1.Children["child2"] = child2

	Add([]string{"child1", "child3"}, tree)

	if len(tree.Children) != 1 {
		t.Error("len(tree children) != 1")
	}

	if len(child1.Children) != 2 {
		t.Error("len(child1 children) != 2")
	}
}

// should find a simple path
func TestFindSimple(t *testing.T) {
	tree := makeTree("")
	child := makeTree("child")
	tree.Children["child"] = child

	found := Find([]string{"child"}, tree)

	if found == nil || found.Regexp.String() != "child" {
		t.Error("did not find \"child\"")
	}
}

// should find a deep path
func TestFindDeep(t *testing.T) {
	tree := makeTree("")
	child1 := makeTree("child1")
	child2 := makeTree("child2")
	tree.Children["child1"] = child1
	child1.Children["child2"] = child2

	found := Find([]string{"child1", "child2"}, tree)

	if found == nil || found.Regexp.String() != "child2" {
		t.Error("did not find \"child2\"")
	}
}

// should find nothing in a nil tree
func TestFindNil(t *testing.T) {
	if found := Find([]string{"test"}, nil); found != nil {
		t.Error("expected nil when searching nil")
	}
}

// should find nothing in an empty tree
func TestFindEmpty(t *testing.T) {
	if found := Find([]string{"test"}, makeTree("")); found != nil {
		t.Error("expected nil when searching empty tree")
	}
}

// should not find a missing path
func TestFindMissingPath(t *testing.T) {
	tree := makeTree("")
	child := makeTree("child")
	tree.Children["child"] = child

	if found := Find([]string{"missing"}, tree); found != nil {
		t.Error("expected nil when searching missing path")
	}
}
