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

func (p *Path) findChild(name string) *Path {
	for _, child := range p.Children {
		re := child.Regexp
		if re.String() == name || re.MatchString(name) {
			return child
		}
	}

	return nil
}

func (p *Path) addChild(child *Path) {
	if child == nil {
		return
	}
	name := child.Regexp.String()
	if node := p.findChild(name); node == nil {
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

	node := p.findChild(path[0].String())
	if node == nil {
		node = new(Path)
		node.Regexp = path[0]
		node.Children = make([]*Path, 0)
	}
	p.addChild(node)

	return node.Add(path[1:])
}

func (p *Path) Find(path []string) *Path {
	if len(path) == 0 {
		return p
	}

	if child := p.findChild(path[0]); child != nil {
		return child.Find(path[1:])
	}

	return nil
}
