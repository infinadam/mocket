package router

import (
	"regexp"
)

type Path struct {
	Regexp *regexp.Regexp
	// this should be a generic in the future.
	Action *HTTPAction

	Children []*Path
}

func (p *Path) findChild(name string) (*Path, map[string]string) {
	for _, child := range p.Children {
		re := child.Regexp
		if re.String() == name {
			return child, nil
		}
		if matched, groups := match(re, name); matched {
			return child, groups
		}
	}

	return nil, nil
}

func (p *Path) addChild(child *Path) {
	if child == nil {
		return
	}
	name := child.Regexp.String()
	if node, _ := p.findChild(name); node == nil {
		p.Children = append(p.Children, child)
	}
}

func (p *Path) Add(path []*regexp.Regexp) *Path {
	if len(path) == 0 {
		return p
	}

	if p.Regexp != nil && p.Regexp.String() == path[0].String() {
		path = path[1:]
	}

	node, _ := p.findChild(path[0].String())
	if node == nil {
		node = new(Path)
		node.Regexp = path[0]
		node.Children = make([]*Path, 0)
	}
	p.addChild(node)

	return node.Add(path[1:])
}

func (p *Path) Find(path []string, groups map[string]string) (*Path, map[string]string) {
	if len(path) == 0 {
		return p, groups
	}

	if child, g := p.findChild(path[0]); child != nil {
		if groups == nil {
			groups = make(map[string]string)
		}
		for k, v := range g {
			groups[k] = v
		}
		return child.Find(path[1:], groups)
	}

	return nil, groups
}
