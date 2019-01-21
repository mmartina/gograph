package templates

const TraverseTemplate = `// Generated automatically by gengraph: do not edit manually

{{if .Tags -}}
// +build {{.Tags}}

{{end -}}
{{$keyName := .KeyName -}}{{$nodeName := .NodeName -}}{{$graphName := .GraphName -}}{{$children := .Children -}}
package {{.PackageName}}

type {{$keyName}}Pair struct {
	Sealed{{.KeyField}} {{$keyName}}
	Logical{{.KeyField}} *{{$keyName}}
}

func (kp {{$keyName}}Pair) {{.KeyField}}() {{$keyName}} {
	if kp.Logical{{.KeyField}} != nil {
		return *kp.Logical{{.KeyField}}
	}
	return kp.Sealed{{.KeyField}}
}

type {{$nodeName}}Visitor interface {
	IsSealed{{$keyName}}(key {{$keyName}}) bool
	VisitSealed{{$nodeName}}(parentID string, node {{$nodeName}}, depth int)
	VisitUnsealed{{$nodeName}}(parentID string, node {{$nodeName}}) string
	{{range .Callees -}}
	{{if .EdgeName -}}
	Visit{{$nodeName}}To{{.NodeName}}Edge(caller {{$keyName}}Pair, callee {{.KeyName}}Pair, edge {{.EdgeName}})
	{{else -}}
	Visit{{$nodeName}}To{{.NodeName}}Edge(caller {{$keyName}}Pair, callee {{.KeyName}}Pair)
	{{end -}}
	{{end -}}
}

func (g {{$graphName}}) InheritedDegreeOf{{$nodeName}}(v {{.GraphName}}Visitor, node {{$nodeName}}) int {
	degree := node.Degree()
	for node.Has{{.Parent.Field}}() {
		node = *g.Get{{$nodeName}}(*node.{{.Parent.Field}})
		degree += node.Degree() / len(g.CollectDeep{{.Children.Field}}Of{{$nodeName}}(node.{{.KeyField}}, v.IsSealed{{$keyName}}))
	}
	return degree
}

func (t {{.GraphName}}Traverser) Traverse{{$nodeName}}s(v {{.GraphName}}Visitor, parentID string) {
	for _, node := range t.Graph.Sorted{{$nodeName}}s() {
		if !node.Has{{.Parent.Field}}() {
			t.traverse{{$nodeName}}(v, parentID, node, 0)
			{{range .Callees -}}
			t.deepTraverseAll{{$nodeName}}To{{.NodeName}}Edges(v, node)
			{{end -}}
		}
	}
}

func (t {{.GraphName}}Traverser) traverse{{$nodeName}}(v {{$nodeName}}Visitor, parentID string, node {{$nodeName}}, depth int) {
	if v.IsSealed{{$keyName}}(node.{{.KeyField}}) {
		v.VisitSealed{{$nodeName}}(parentID, node, depth)
	} else {
		id := v.VisitUnsealed{{$nodeName}}(parentID, node)
		for _, child{{.KeyField}} := range node.{{.Children.Field}} {
			t.traverse{{$nodeName}}(v, id, *t.Graph.Get{{$nodeName}}(child{{.KeyField}}), depth + 1)
		}
	}
}

func (t {{.GraphName}}Traverser) get{{$keyName}}Pair(v {{$nodeName}}Visitor, key {{$keyName}}) {{$keyName}}Pair {
	if v.IsSealed{{$keyName}}(key) {
		return {{$keyName}}Pair{Sealed{{.KeyField}}: key}
	} else {
		return {{$keyName}}Pair{
			Sealed{{.KeyField}}: *t.Graph.AnyDeep{{.Children.SField}}Of{{$nodeName}}(key, v.IsSealed{{$keyName}}),
			Logical{{.KeyField}}: &key,
		}
	}
}

{{range .Callees -}}
func (t {{$graphName}}Traverser) deepTraverseAll{{$nodeName}}To{{.NodeName}}Edges(v {{$graphName}}Visitor, node {{$nodeName}}) {
	{{if .EdgeName -}}
	for _, callee{{.KeyField}} := range node.Get{{.Field}}() {
		v.Visit{{$nodeName}}To{{.NodeName}}Edge(t.get{{$keyName}}Pair(v, node.{{.KeyField}}), t.get{{.KeyName}}Pair(v, callee{{.KeyField}}), node.{{.Field}}[callee{{.KeyField}}])
	}
	{{else -}}
	for _, callee{{.KeyField}} := range node.{{.Field}} {
		v.Visit{{$nodeName}}To{{.NodeName}}Edge(t.get{{$keyName}}Pair(v, node.{{.KeyField}}), t.get{{.KeyName}}Pair(v, callee{{.KeyField}}))
	}
	{{end -}}
	for _, child{{.KeyField}} := range node.{{$children.Field}} {
		t.deepTraverseAll{{$nodeName}}To{{.NodeName}}Edges(v, *t.Graph.Get{{$nodeName}}(child{{.KeyField}}))
	}
}

{{end -}}
`
