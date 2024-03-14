package router

import "regexp"

type Path struct {
	Regexp *regexp.Regexp
	// this should be a generic in the future.
	Action *HTTPAction

	Children []*Path
}

func makePath(re *regexp.Regexp) *Path {
	p := new(Path)

	if re != nil && re.String() != "" {
		p.Regexp = re
	}
	p.Children = make([]*Path, 0)

	return p
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

	var node *Path
	if node = p.findChild(path[0].String()); node == nil {
		node = makePath(path[0])
	}
	p.addChild(node)

	return node.Add(path[1:])
}

func (p *Path) Find(path []string) *Path {
	if len(path) == 0 {
		return p
	}

	if p.Regexp != nil && p.Regexp.MatchString(path[0]) {
		path = path[1:]
	}

	return p.findChild(path[0]).Find(path[1:])
}
