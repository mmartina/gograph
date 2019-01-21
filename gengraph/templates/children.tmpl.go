package templates

const ChildrenTemplate = `// Generated automatically by gengraph: do not edit manually

{{if .Tags -}}
// +build {{.Tags}}

{{end -}}
{{$keyName := .KeyName -}}{{$nodeName := .NodeName -}}
package {{.PackageName}}

func (n *{{$nodeName}}) HasAny{{.Children.SField}}() bool {
	return len(n.{{.Children.Field}}) > 0
}

func (n *{{$nodeName}}) Get{{.Children.Field}}() {{.Children.KeyName}}s {
	return {{.Children.KeyName}}s(n.{{.Children.Field}})
}

func (n *{{$nodeName}}) add{{.Children.SField}}(key {{$keyName}}) {
	if {{$keyName}}s(n.{{.Children.Field}}).Contains(key) {
		return
	}
	children := append(n.{{.Children.Field}}, key)
	sort.Sort({{.Children.KeyName}}s(children))
	n.{{.Children.Field}} = children
}

func (n *{{$nodeName}}) remove{{.Children.SField}}(key {{$keyName}}) {
	n.{{.Children.Field}} = {{.Children.KeyName}}s(n.{{.Children.Field}}).Without(key)
}

func (n *{{$nodeName}}) Connect{{.Children.SField}}(child *{{$nodeName}}) {
	if child.{{.Parent.Field}} == nil && child.{{.KeyField}} != n.{{.KeyField}} {
		n.add{{.Children.SField}}(child.{{.KeyField}})
		child.set{{.Parent.Field}}(n.{{.KeyField}})
	}
}

func (n *{{$nodeName}}) Disconnect{{.Children.SField}}(child *{{$nodeName}}) {
	if child.{{.Parent.Field}} != nil && *child.{{.Parent.Field}} == n.{{.KeyField}} {
		n.remove{{.Children.SField}}(child.{{.KeyField}})
		child.remove{{.Parent.Field}}()
	}
}

func (g *{{.GraphName}}) DeepDisconnectAndRemove{{$nodeName}}(node *{{$nodeName}}) {
	for _, child{{.KeyField}} := range node.{{.Children.Field}} {
		g.DeepDisconnectAndRemove{{$nodeName}}(g.Get{{$nodeName}}(child{{.KeyField}}))
	}
	g.DisconnectAndRemove{{$nodeName}}(node)
}

func (g *{{.GraphName}}) AnyDeep{{.Children.SField}}Of{{$nodeName}}(key {{$keyName}}, filter {{$keyName}}Filter) *{{$keyName}} {
	deepChildren := g.CollectDeep{{.Children.Field}}Of{{$nodeName}}(key, filter)
	if len(deepChildren) == 0 {
		return nil
	}
	childKey := deepChildren[rand.Intn(len(deepChildren))]
	return &childKey
}

func (g XGraph) CollectDeep{{.Children.Field}}Of{{$nodeName}}(key {{$keyName}}, filter {{$keyName}}Filter) {{$keyName}}s {
	if filter(key) {
		return {{$keyName}}s{key}
	}
	node := g.Get{{$nodeName}}(key)
	deepChildren := make({{$keyName}}s,0)
	for _, child{{.KeyField}} := range node.{{.Children.Field}} {
		deepChildren = append(deepChildren, g.CollectDeep{{.Children.Field}}Of{{$nodeName}}(child{{.KeyField}}, filter)...)
	}
	return deepChildren
}
`
