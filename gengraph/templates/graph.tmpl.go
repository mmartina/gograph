package templates

const GraphTemplate = `// Generated automatically by gengraph: do not edit manually

{{if .Tags -}}
// +build {{.Tags}}

{{end -}}
{{$keyName := .KeyName -}}{{$nodeName := .NodeName -}}
package {{.PackageName}}

func (g *{{.GraphName}}) add{{$nodeName}}(node *{{$nodeName}}) {
	g.{{.GraphField}}[node.{{.KeyField}}] = node
}

func (g *{{.GraphName}}) remove{{$nodeName}}(key {{$keyName}}) {
	delete(g.{{.GraphField}}, key)
}

func (g *{{.GraphName}}) Contains{{$nodeName}}(key {{$keyName}}) bool {
	_, ok := g.{{.GraphField}}[key]
	return ok
}

func (g *{{.GraphName}}) Get{{$nodeName}}(key {{$keyName}}) *{{$nodeName}} {
	if node, ok := g.{{.GraphField}}[key]; ok {
		return node
	}
	return nil
}

func (g *{{.GraphName}}) Get{{$nodeName}}s(keys {{$keyName}}s) []*{{$nodeName}} {
	nodes := make([]*{{$nodeName}}, 0)
	for _, key := range keys {
		nodes = append(nodes, g.Get{{$nodeName}}(key))
	}
	return nodes
}

func (g *{{.GraphName}}) Filtered{{$nodeName}}s(keys {{$keyName}}s) {{$nodeName}}s {
	return As{{$nodeName}}s(g.Get{{$nodeName}}s(keys))
}

func (g {{.GraphName}}) Sorted{{$keyName}}s() {{$keyName}}s {
	keys := make({{$keyName}}s, 0)
	for key := range g.{{.GraphField}} {
		keys = append(keys, key)
	}
	sort.Sort(keys)
	return keys
}

func (g {{.GraphName}}) Sorted{{$nodeName}}s() {{$nodeName}}s {
	nodes := make({{$nodeName}}s, 0)
	for _, node := range g.{{.GraphField}} {
		nodes = append(nodes, *node)
	}
	sort.Sort(nodes)
	return nodes
}
`
