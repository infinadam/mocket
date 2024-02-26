package urlpath

import "regexp"

type Tree struct {
	Regexp *regexp.Regexp

	Children []*Tree
}

func makeTree(name string) *Tree {
	t := new(Tree)

	var err error
	if name != "" {
		if t.Regexp, err = regexp.Compile(name); err != nil {
			return nil
		}
	}
	t.Children = make([]*Tree, 0)

	return t
}

func (t *Tree) findChild(name string) *Tree {
	for _, child := range t.Children {
		if child.Regexp.String() == name || child.Regexp.MatchString(name) {
			return child
		}
	}

	return nil
}

func (t *Tree) addChild(child *Tree) {
	if child == nil {
		return
	}
	name := child.Regexp.String()
	if node := t.findChild(name); node == nil {
		t.Children = append(t.Children, child)
	}
}

func Add(tree *Tree, path []string) *Tree {
	if len(path) == 0 {
		return tree
	}

	if tree == nil {
		tree = makeTree("")
	}

	if tree.Regexp != nil && tree.Regexp.String() == path[0] {
		path = path[1:]
	}

	var node *Tree
	if node = tree.findChild(path[0]); node == nil {
		node = makeTree(path[0])
	}
	tree.addChild(node)
	Add(node, path[1:])

	return tree
}

func Find(tree *Tree, path []string) *Tree {
	if len(path) == 0 || tree == nil {
		return tree
	}

	if tree.Regexp != nil && tree.Regexp.String() == path[0] {
		path = path[1:]
	}

	return Find(tree.findChild(path[0]), path[1:])
}
