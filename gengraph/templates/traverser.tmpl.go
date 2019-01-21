package templates

const TraverserTemplate = `// Generated automatically by gengraph: do not edit manually

{{if .Tags -}}
// +build {{.Tags}}

{{end -}}
package {{.PackageName}}

type {{.GraphName}}Visitor interface {
	{{range .Nodes -}}
		{{.NodeName}}Visitor
	{{end -}}
}

type {{.GraphName}}Traverser struct {
	Graph {{.GraphName}}
}

func (t {{.GraphName}}Traverser) TraverseAll(v {{.GraphName}}Visitor, parentID string) {
	{{range .Nodes -}}
		t.Traverse{{.NodeName}}s(v, parentID)
	{{end -}}
}
`
