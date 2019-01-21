package templates

const ConnectedTemplate = `// Generated automatically by gengraph: do not edit manually

{{if .Tags -}}
// +build {{.Tags}}

{{end -}}
{{$graphName := .GraphName -}}
package {{.PackageName}}

func (g *{{$graphName}}) PartitionIntoSubgraphs() []{{$graphName}} {
	subgraphs := make([]{{$graphName}},0)
	{{range .Nodes -}}
	for len(g.{{.GraphField}}) > 0 {
		for key := range g.{{.GraphField}} {
			sg := New{{$graphName}}()
			g.migrateSubgraphFor{{.NodeName}}(sg, key)
			subgraphs = append(subgraphs, *sg)
			break
		}
	}
	{{end -}}
	return subgraphs
}
`
