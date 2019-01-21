package templates

const NodeTemplate = `// Generated automatically by gengraph: do not edit manually

{{if .Tags -}}
// +build {{.Tags}}

{{end -}}
{{$keyName := .KeyName -}}{{$nodeName := .NodeName -}}
package {{.PackageName}}

func (g *{{.GraphName}}) GetOrCreate{{$nodeName}}(key {{$keyName}}) *{{$nodeName}} {
	node := g.Get{{$nodeName}}(key)
	if node == nil {
		node = new{{$nodeName}}(key)
		g.add{{$nodeName}}(node)
	}
	return node
}

func (g *{{.GraphName}}) DisconnectAndRemove{{$nodeName}}(node *{{$nodeName}}) {
	if node.{{.Parent.Field}} != nil {
		parent := g.Get{{$nodeName}}(*node.{{.Parent.Field}})
		node.Disconnect{{.Parent.Field}}(parent)
	}
	for _, child{{.KeyField}} := range node.{{.Children.Field}} {
		child := g.Get{{$nodeName}}(child{{.KeyField}})
		node.Disconnect{{.Children.SField}}(child)
	}
	{{range .Callers -}}
	for _, caller{{.KeyField}} := range node.Get{{.Field}}() {
		caller := g.Get{{.NodeName}}(caller{{.KeyField}})
		node.Disconnect{{.SField}}(caller)
	}
	{{end -}}
	{{range .Callees -}}
	for _, callee{{.KeyField}} := range node.Get{{.Field}}() {
		callee := g.Get{{.NodeName}}(callee{{.KeyField}})
		node.Disconnect{{.SField}}(callee)
	}
	{{end -}}
	g.remove{{$nodeName}}(node.{{.KeyField}})
}

func (node {{$nodeName}}) Degree() int {
	degree := 0
	{{range .Callers -}}
	degree += len(node.{{.Field}})
	{{end -}}
	{{range .Callees -}}
	degree += len(node.{{.Field}})
	{{end -}}
	return degree
}

func As{{$nodeName}}s(nodes []*{{$nodeName}}) {{$nodeName}}s {
	values := {{$nodeName}}s{}
	for _, node := range nodes {
		values = append(values, *node)
	}
	return values
}

func (nodes {{$nodeName}}s) {{.KeyField}}s() {{$keyName}}s {
	keys := {{$keyName}}s{}
	for _, node := range nodes {
		keys = append(keys, node.{{.KeyField}})
	}
	return keys
}

func (nodes {{$nodeName}}s) Any() *{{$nodeName}} {
	if len(nodes) > 0 {
		return &nodes[rand.Intn(len(nodes))]
	}
	return nil
}

type {{$nodeName}}s []{{$nodeName}}
func (n {{$nodeName}}s) Len() int           { return len(n) }
func (n {{$nodeName}}s) Less(i, j int) bool { return n[i].{{.KeyField}}.LessThan(n[j].{{.KeyField}}) }
func (n {{$nodeName}}s) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
`
