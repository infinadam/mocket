package router

import "testing"

// should create a tree node
func TestMakeNode(t *testing.T) {
	node := makePath("test")
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
	root := makePath("")
	if root == nil {
		t.Fatal("failed to create new node")
	}
	if root.Regexp != nil {
		t.Errorf("unexpected regexp")
	}
}

// should find a child in a tree's immediate children
func TestFindChild(t *testing.T) {
	tree := makePath("")
	child := makePath("test")
	tree.Children = append(tree.Children, child)

	found := tree.findChild("test")
	if found == nil {
		t.Error("did not find child")
	}
}

// should not find a child that is missing from the immediate children
func TestFindMissingChild(t *testing.T) {
	tree := makePath("")
	child := makePath("test")
	tree.Children = append(tree.Children, child)

	found := tree.findChild("missing")
	if found != nil {
		t.Error("found a missing child")
	}
}

// should find a child using regular expression
func TestFindRegexChild(t *testing.T) {
	tree := makePath("")
	child := makePath(`\d`)
	tree.Children = append(tree.Children, child)

	found := tree.findChild("123")
	if found == nil {
		t.Error("did not find regular expression node")
	}
}

// should add a child to a tree's children list
func TestAddChild(t *testing.T) {
	tree := makePath("")
	child := makePath("child")

	tree.addChild(child)
	if len(tree.Children) != 1 {
		t.Error("failed to add child")
	}
}

// should not add a nil to a tree's children list
func TestAddNilChild(t *testing.T) {
	tree := makePath("")

	tree.addChild(nil)
	if len(tree.Children) != 0 {
		t.Error("added nil child")
	}
}

// should not add a child with the same name
func TestAddSameChild(t *testing.T) {
	tree := makePath("")
	child := makePath("child")
	tree.Children = append(tree.Children, child)

	tree.addChild(makePath("child"))
	if len(tree.Children) > 1 {
		t.Error("added child of the same name")
	}
}

// should have well-defined children
func TestMakeChildren(t *testing.T) {
	node := makePath("")
	if node == nil {
		t.Fatal("failed to create new node")
	}
	if node.Children == nil {
		t.Error("node children is nil")
	}
}

// should do nothing if no path is given
func TestNoPath(t *testing.T) {
	tree := makePath("")

	child := tree.add([]string{})

	if tree != child {
		t.Error("expected tree to be returned")
	}

	if len(tree.Children) > 0 {
		t.Error("children should be empty")
	}
}

// should make a tree with a path
func TestCreatePath(t *testing.T) {
	tree := makePath("")
	tree.add([]string{"test", "child"})

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
func TestAddToExisting(t *testing.T) {
	tree := makePath("")
	child := makePath("child1")
	tree.addChild(child)

	tree.add([]string{"child1", "child2"})

	if child.findChild("child2") == nil {
		t.Error("could not find \"child2\".")
	}
}

// should not add anything if the path exists
func TestIdentity(t *testing.T) {
	tree := makePath("")
	child1 := makePath("child1")
	child2 := makePath("child2")
	tree.addChild(child1)
	child1.addChild(child2)

	tree.add([]string{"child1", "child2"})

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
func TestAddAdjacent(t *testing.T) {
	tree := makePath("")
	child1 := makePath("child1")
	child2 := makePath("child2")
	tree.addChild(child1)
	child1.addChild(child2)

	tree.add([]string{"child1", "child3"})

	if len(tree.Children) != 1 {
		t.Errorf("len(tree children) != 1 (%d)", len(tree.Children))
	}

	if len(child1.Children) != 2 {
		t.Error("len(child1 children) != 2")
	}
}

// should find a simple path
func TestFindSimple(t *testing.T) {
	tree := makePath("")
	child := makePath("child")
	tree.addChild(child)

	found := tree.find([]string{"child"})

	if found == nil || found.Regexp.String() != "child" {
		t.Error("did not find \"child\"")
	}
}

// should find a deep path
func TestFindDeep(t *testing.T) {
	tree := makePath("")
	child1 := makePath("child1")
	child2 := makePath("child2")
	tree.addChild(child1)
	child1.addChild(child2)

	found := tree.find([]string{"child1", "child2"})

	if found == nil || found.Regexp.String() != "child2" {
		t.Error("did not find \"child2\"")
	}
}

// should find nothing in an empty tree
func TestFindEmpty(t *testing.T) {
	tree := makePath("")
	if found := tree.find([]string{"test"}); found != nil {
		t.Error("expected nil when searching empty tree")
	}
}

// should not find a missing path
func TestFindMissingPath(t *testing.T) {
	tree := makePath("")
	child := makePath("child")
	tree.addChild(child)

	if found := tree.find([]string{"missing"}); found != nil {
		t.Error("expected nil when searching missing path")
	}
}

// should find a node through Regexp application
func TestFindRegexp(t *testing.T) {
	tree := makePath("")
	child1 := makePath(`\d+`)
	child2 := makePath("test")
	tree.addChild(child1)
	child1.addChild(child2)

	found := tree.find([]string{"123", "test"})
	if found != child2 {
		t.Error("did not find \"test\"")
	}
}
