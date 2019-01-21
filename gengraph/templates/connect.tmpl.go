package templates

const ConnectTemplate = `// Generated automatically by gengraph: do not edit manually

{{if .Tags -}}
// +build {{.Tags}}

{{end -}}
{{$keyName := .KeyName -}}{{$nodeName := .NodeName -}}{{$graphName := .GraphName -}}
package {{.PackageName}}

func (g *{{.GraphName}}) migrateSubgraphFor{{$nodeName}}(sg *{{.GraphName}}, key {{.KeyName}}) {
	if sg.Contains{{$nodeName}}(key) {
		return
	}
	node := g.Get{{$nodeName}}(key)
	sg.add{{$nodeName}}(node)
	if node.Has{{.Parent.Field}}() {
		g.migrateSubgraphFor{{$nodeName}}(sg, *node.{{.Parent.Field}})
	}
	for _, child := range node.{{.Children.Field}} {
		g.migrateSubgraphFor{{$nodeName}}(sg, child)
	}
	{{range .Callers -}}{{if .EdgeName -}}
	for callerKey, _ := range node.{{.Field}} {
		g.migrateSubgraphFor{{.NodeName}}(sg, callerKey)
	}
	{{else -}}
	for _, callerKey := range node.{{.Field}} {
		g.migrateSubgraphFor{{.NodeName}}(sg, callerKey)
	}
	{{end -}}{{end -}}	
	{{range .Callees -}}{{if .EdgeName -}}
	for calleeKey, _ := range node.{{.Field}} {
		g.migrateSubgraphFor{{.NodeName}}(sg, calleeKey)
	}
	{{else -}}
	for _, calleeKey := range node.{{.Field}} {
		g.migrateSubgraphFor{{.NodeName}}(sg, calleeKey)
	}
	{{end -}}{{end -}}
	g.remove{{$nodeName}}(key)
}
`
