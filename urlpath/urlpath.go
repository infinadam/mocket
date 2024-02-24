package urlpath

import "regexp"

type Tree struct {
	Regexp *regexp.Regexp

	Children map[string]*Tree
}

func makeTree(name string) *Tree {
	t := new(Tree)

	var err error
	if name != "" {
		if t.Regexp, err = regexp.Compile(name); err != nil {
			return nil
		}
	}
	t.Children = make(map[string]*Tree)

	return t
}

func Add(path []string, tree *Tree) *Tree {
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
	if node = tree.Children[path[0]]; node == nil {
		node = makeTree(path[0])
	}
	tree.Children[path[0]] = node
	Add(path[1:], node)

	return tree
}

func Find(path []string, tree *Tree) *Tree {
	if len(path) == 0 || tree == nil {
		return tree
	}

	if tree.Regexp != nil && tree.Regexp.String() == path[0] {
		path = path[1:]
	}

	return Find(path[1:], tree.Children[path[0]])
}
