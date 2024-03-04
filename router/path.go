package router

import "regexp"

type Path struct {
	Regexp *regexp.Regexp

	Children []*Path
}

func makePath(name string) *Path {
	p := new(Path)

	var err error
	if name != "" {
		if p.Regexp, err = regexp.Compile(name); err != nil {
			return nil
		}
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

func (p *Path) add(path []string) *Path {
	if len(path) == 0 {
		return p
	}

	if p.Regexp != nil && p.Regexp.String() == path[0] {
		path = path[1:]
	}

	var node *Path
	if node = p.findChild(path[0]); node == nil {
		node = makePath(path[0])
	}
	p.addChild(node)

	return node.add(path[1:])
}

func (p *Path) find(path []string) *Path {
	if len(path) == 0 {
		return p
	}

	if p.Regexp != nil && p.Regexp.String() == path[0] {
		path = path[1:]
	}

	return p.findChild(path[0]).find(path[1:])
}
