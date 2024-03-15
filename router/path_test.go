package router

import (
	"regexp"
	"testing"
)

// should find a child in a tree's immediate children
func TestPathFindChild(t *testing.T) {
	var tree, child Path

	child.Regexp, _ = regexp.Compile("test")
	tree.Children = append(tree.Children, &child)

	found := tree.findChild("test")

	if found == nil {
		t.Error("did not find child")
	}
}

// should not find a child that is missing from the immediate children
func TestPathFindMissingChild(t *testing.T) {
	var tree, child Path

	child.Regexp, _ = regexp.Compile("test")
	tree.Children = append(tree.Children, &child)

	found := tree.findChild("missing")

	if found != nil {
		t.Error("found a missing child")
	}
}

// should find a child using regular expression
func TestPathFindRegexChild(t *testing.T) {
	var tree, child Path

	child.Regexp, _ = regexp.Compile(`\d+`)
	tree.Children = append(tree.Children, &child)

	found := tree.findChild("123")

	if found == nil {
		t.Error("did not find regular expression node")
	}
}

// should add a child to a tree's children list
func TestPathAddChild(t *testing.T) {
	var tree, child Path

	child.Regexp, _ = regexp.Compile("child")

	tree.addChild(&child)

	if len(tree.Children) != 1 {
		t.Error("failed to add child")
	}
}

// should not add a nil to a tree's children list
func TestPathAddNilChild(t *testing.T) {
	var tree Path

	tree.addChild(nil)

	if len(tree.Children) != 0 {
		t.Error("added nil child")
	}
}

// should not add a child with the same name
func TestPathAddSameChild(t *testing.T) {
	var tree, child Path

	child.Regexp, _ = regexp.Compile("child")
	tree.Children = append(tree.Children, &child)

	re, _ := regexp.Compile("child")
	tree.addChild(&Path{re, nil, nil})

	if len(tree.Children) > 1 {
		t.Error("added child of the same name")
	}
}

// should do nothing if no path is given
func TestPathNoPath(t *testing.T) {
	var tree Path

	child := tree.Add([]*regexp.Regexp{})

	if &tree != child {
		t.Error("expected tree to be returned")
	}

	if len(tree.Children) > 0 {
		t.Error("children should be empty")
	}
}

// should make a tree with a path
func TestPathCreatePath(t *testing.T) {
	var tree Path

	var res []*regexp.Regexp
	re, _ := regexp.Compile("test")
	res = append(res, re)
	re, _ = regexp.Compile("child")
	res = append(res, re)

	tree.Add(res)

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
	var tree, child Path

	child.Regexp, _ = regexp.Compile("child1")
	tree.addChild(&child)

	var res []*regexp.Regexp
	re, _ := regexp.Compile("child1")
	res = append(res, re)
	re, _ = regexp.Compile("child2")
	res = append(res, re)

	tree.Add(res)

	if child.findChild("child2") == nil {
		t.Error("could not find \"child2\".")
	}
}

// should not add anything if the path exists
func TestPathIdentity(t *testing.T) {
	var tree, child1, child2 Path
	child1.Regexp, _ = regexp.Compile("child1")
	child2.Regexp, _ = regexp.Compile("child2")

	tree.addChild(&child1)
	child1.addChild(&child2)

	var res []*regexp.Regexp
	re, _ := regexp.Compile("child1")
	res = append(res, re)
	re, _ = regexp.Compile("child2")
	res = append(res, re)

	tree.Add(res)

	if tree.findChild("child1") != &child1 {
		t.Error("unexpected child1")
	}

	if child1.findChild("child2") != &child2 {
		t.Error("unexpected child2")
	}

	if len(child2.Children) != 0 {
		t.Error("len(child2 children) != 0")
	}
}

// should add an adjacent node
func TestPathAddAdjacent(t *testing.T) {
	var tree, child1, child2 Path
	child1.Regexp, _ = regexp.Compile("child1")
	child2.Regexp, _ = regexp.Compile("child2")

	tree.addChild(&child1)
	child1.addChild(&child2)

	var res []*regexp.Regexp
	re, _ := regexp.Compile("child1")
	res = append(res, re)
	re, _ = regexp.Compile("child3")
	res = append(res, re)

	tree.Add(res)

	if len(tree.Children) != 1 {
		t.Errorf("len(tree children) != 1 (%d)", len(tree.Children))
	}

	if len(child1.Children) != 2 {
		t.Error("len(child1 children) != 2")
	}
}

// should find a simple path
func TestPathFindSimple(t *testing.T) {
	var tree, child Path
	child.Regexp, _ = regexp.Compile("child")

	tree.addChild(&child)

	found := tree.Find([]string{"child"})

	if found == nil || found.Regexp.String() != "child" {
		t.Error("did not find \"child\"")
	}
}

// should find a deep path
func TestPathFindDeep(t *testing.T) {
	var tree, child1, child2 Path
	child1.Regexp, _ = regexp.Compile("child1")
	child2.Regexp, _ = regexp.Compile("child2")

	tree.addChild(&child1)
	child1.addChild(&child2)

	found := tree.Find([]string{"child1", "child2"})

	if found == nil || found.Regexp.String() != "child2" {
		t.Error("did not find \"child2\"")
	}
}

// should find nothing in an empty tree
func TestPathFindEmpty(t *testing.T) {
	var tree Path
	if found := tree.Find([]string{"test"}); found != nil {
		t.Error("expected nil when searching empty tree")
	}
}

// should not find a missing path
func TestPathFindMissingPath(t *testing.T) {
	var tree, child Path
	child.Regexp, _ = regexp.Compile("child1")

	tree.addChild(&child)

	if found := tree.Find([]string{"missing"}); found != nil {
		t.Error("expected nil when searching missing path")
	}
}

// should find a non-leaf node
func TestFindShortPath(t *testing.T) {
	var tree, child1, child2 Path
	child1.Regexp, _ = regexp.Compile("child1")
	child2.Regexp, _ = regexp.Compile("child2")

	tree.addChild(&child1)
	child1.addChild(&child2)

	if found := tree.Find([]string{"child1"}); found != &child1 {
		t.Error("expected to find child1")
	}
}

// should find a node through Regexp application
func TestPathFindRegexp(t *testing.T) {
	var tree, child1, child2 Path
	child1.Regexp, _ = regexp.Compile(`\d+`)
	child2.Regexp, _ = regexp.Compile("test")

	tree.addChild(&child1)
	child1.addChild(&child2)

	found := tree.Find([]string{"123", "test"})
	if found != &child2 {
		t.Error("did not find \"test\"")
	}
}
