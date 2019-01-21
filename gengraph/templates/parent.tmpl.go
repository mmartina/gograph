package templates

const ParentTemplate = `// Generated automatically by gengraph: do not edit manually

{{if .Tags -}}
// +build {{.Tags}}

{{end -}}
{{$keyName := .KeyName -}}{{$nodeName := .NodeName -}}
package {{.PackageName}}

func (n *{{$nodeName}}) Has{{.Parent.Field}}() bool {
	return n.{{.Parent.Field}} != nil
}

func (n *{{$nodeName}}) set{{.Parent.Field}}(key {{$keyName}}) {
	n.{{.Parent.Field}} = &key
}

func (n *{{$nodeName}}) remove{{.Parent.Field}}() {
	n.{{.Parent.Field}} = nil
}

func (n *{{$nodeName}}) Connect{{.Parent.Field}}(parent *{{$nodeName}}) {
	if n.{{.Parent.Field}} == nil && parent.{{.KeyField}} != n.{{.KeyField}}{
		n.set{{.Parent.Field}}(parent.{{.KeyField}})
		parent.add{{.Children.SField}}(n.{{.KeyField}})
	}
}

func (n *{{$nodeName}}) Disconnect{{.Parent.Field}}(parent *{{$nodeName}}) {
	if n.{{.Parent.Field}} != nil && *n.{{.Parent.Field}} == parent.{{.KeyField}} {
		n.remove{{.Parent.Field}}()
		parent.remove{{.Children.SField}}(n.{{.KeyField}})
	}
}

func (n *{{$nodeName}}) Switch{{.Parent.Field}}(old *{{$nodeName}}, new *{{$nodeName}}) {
	if n.{{.Parent.Field}} != nil && *n.{{.Parent.Field}} == old.{{.KeyField}}  && new.{{.KeyField}} != n.{{.KeyField}} {
		n.set{{.Parent.Field}}(new.{{.KeyField}})
		old.remove{{.Children.SField}}(n.{{.KeyField}})
		new.add{{.Children.SField}}(n.{{.KeyField}})
	}
}

func (g *{{.GraphName}}) Migrate{{.Parent.Field}}Of{{$nodeName}}(old *{{$nodeName}}, new *{{$nodeName}}, parent *{{$keyName}}) {
	if parent != nil {
		parentNode := g.Get{{$nodeName}}(*parent)
		new.Connect{{.Parent.Field}}(parentNode)
		old.Disconnect{{.Parent.Field}}(parentNode)
	}
}

func (g *{{.GraphName}}) Common{{.Parent.Field}}sOf{{$nodeName}}s(nodes {{$nodeName}}s) *{{$keyName}} {
	var common{{.Parent.Field}} *{{$keyName}}
	for i, node := range nodes {
		if node.{{.Parent.Field}} == nil {
			return nil
		}
		if i == 0 {
			common{{.Parent.Field}} = node.{{.Parent.Field}}
		} else if *node.{{.Parent.Field}} != *common{{.Parent.Field}} {
			return nil
		}
	}
	return common{{.Parent.Field}}
}

func (g *{{.GraphName}}) RootOf{{$keyName}}(key {{$keyName}}) {{$keyName}} {
	return g.RootOf{{$nodeName}}(g.Get{{$nodeName}}(key)).{{.KeyField}}
}

func (g *{{.GraphName}}) RootOf{{$nodeName}}(node *{{$nodeName}}) *{{$nodeName}} {
	if node.{{.Parent.Field}} != nil {
		return g.RootOf{{$nodeName}}(g.Get{{$nodeName}}(*node.{{.Parent.Field}}))
	}
	return node
}
`
