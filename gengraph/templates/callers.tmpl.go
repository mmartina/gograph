package templates

const CallersTemplate = `// Generated automatically by gengraph: do not edit manually

{{if .Tags -}}
// +build {{.Tags}}

{{end -}}
{{$keyName := .KeyName}}{{$nodeName := .NodeName -}}{{$graphName := .GraphName -}}
package {{.PackageName}}

{{range .Callers -}}
func (n *{{$nodeName}}) HasAny{{.SField}}() bool {
	return len(n.{{.Field}}) > 0
}

{{if .EdgeName -}}
func (n *{{$nodeName}}) Get{{.Field}}() {{.KeyName}}s {
	callers := make({{.KeyName}}s, 0)
	for caller := range n.{{.Field}} {
		callers = append(callers, caller)
	}
	sort.Sort({{.KeyName}}s(callers))
	return callers
}

func (n *{{$nodeName}}) addTo{{.Field}}(key {{.KeyName}}, edge {{.EdgeName}}) {
	n.{{.Field}}[key] = edge
}

func (n *{{$nodeName}}) removeFrom{{.Field}}(key {{.KeyName}}) {{.EdgeName}}{
	edge := n.{{.Field}}[key]
	delete(n.{{.Field}}, key)
	return edge
}

func (n *{{$nodeName}}) Connect{{.SField}}(caller *{{.NodeName}}, edge {{.EdgeName}}) {
	n.addTo{{.Field}}(caller.{{.KeyField}}, edge)
	caller.addTo{{.ReverseField}}(n.{{.KeyField}}, edge)
}

func (n *{{$nodeName}}) Disconnect{{.SField}}(caller *{{.NodeName}}) {{.EdgeName}} {
	edge := n.removeFrom{{.Field}}(caller.{{.KeyField}})
	caller.removeFrom{{.ReverseField}}(n.{{.KeyField}})
	return edge
}

{{else -}}
func (n *{{$nodeName}}) Get{{.Field}}() {{.KeyName}}s {
	return {{.KeyName}}s(n.{{.Field}})
}

func (n *{{$nodeName}}) addTo{{.Field}}(key {{.KeyName}}) {
	if {{.KeyName}}s(n.{{.Field}}).Contains(key) {
		return
	}
	callers := append(n.{{.Field}}, key)
	sort.Sort({{.KeyName}}s(callers))
	n.{{.Field}} = callers
}

func (n *{{$nodeName}}) removeFrom{{.Field}}(key {{.KeyName}}) {
	n.{{.Field}} = {{.KeyName}}s(n.{{.Field}}).Without(key)
}

func (n *{{$nodeName}}) Connect{{.SField}}(caller *{{.NodeName}}) {
	n.addTo{{.Field}}(caller.{{.KeyField}})
	caller.addTo{{.ReverseField}}(n.{{.KeyField}})
}

func (n *{{$nodeName}}) Disconnect{{.SField}}(caller *{{.NodeName}}) {
	n.removeFrom{{.Field}}(caller.{{.KeyField}})
	caller.removeFrom{{.ReverseField}}(n.{{.KeyField}})
}

{{end -}}

func (g *{{$graphName}}) Migrate{{.Field}}Of{{$nodeName}}(old *{{$nodeName}}, new *{{$nodeName}}, callers {{.KeyName}}s) {
	for _, caller{{.KeyField}} := range callers {
		caller := g.Get{{.NodeName}}(caller{{.KeyField}})
		{{if .EdgeName -}}
		edge := old.Disconnect{{.SField}}(caller)
		new.Connect{{.SField}}(caller, edge)
		{{else -}}
		old.Disconnect{{.SField}}(caller)
		new.Connect{{.SField}}(caller)
		{{end -}}
	}
}

func (g *{{$graphName}}) Common{{.Field}}Of{{$nodeName}}s(nodes {{$nodeName}}s) {{.KeyName}}s {
	common{{.Field}} := {{.KeyName}}s{}
	for i, node := range nodes {
		if i == 0 {
			common{{.Field}} = node.Get{{.Field}}()
		} else {
			common{{.Field}} = node.Get{{.Field}}().Intersect(common{{.Field}})
		}
	}
	return common{{.Field}}
}

type CliqueOf{{$nodeName}}sBy{{.Field}} struct {
	{{.KeyField}}s {{$keyName}}s
	{{.Field}}  {{.KeyName}}s
}

func (c CliqueOf{{$nodeName}}sBy{{.Field}}) Score() int {
	return (len(c.{{.KeyField}}s) - 1) * len(c.{{.Field}})
}

type CliquesOf{{$nodeName}}sBy{{.Field}} []CliqueOf{{$nodeName}}sBy{{.Field}}
func (c CliquesOf{{$nodeName}}sBy{{.Field}}) Len() int           { return len(c) }
func (c CliquesOf{{$nodeName}}sBy{{.Field}}) Less(i, j int) bool {
	return c[i].Score() > c[j].Score() || c[i].Score() == c[j].Score() && len(c[i].{{.Field}}) > len(c[j].{{.Field}})
}
func (c CliquesOf{{$nodeName}}sBy{{.Field}}) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }

func (g *XGraph) MaxCliquesOf{{$nodeName}}sWithCommon{{.Field}}(nodes {{$nodeName}}s, min{{.Field}} int) CliquesOf{{$nodeName}}sBy{{.Field}} {
	candidates := nodes.FilteredBy(func(node {{$nodeName}}) bool { return len(node.{{.Field}}) >= min{{.Field}} })
	maxCliques := g.maxCliquesOf{{$nodeName}}sBy{{.Field}}(nil, nil, candidates, nil, min{{.Field}})
	sort.Sort(maxCliques)
	return maxCliques
}

func (g *XGraph) maxCliquesOf{{$nodeName}}sBy{{.Field}}(
	incl {{$keyName}}s, excl {{$keyName}}s, candidates {{$nodeName}}s, common{{.Field}} {{.KeyName}}s, min{{.Field}} int) CliquesOf{{$nodeName}}sBy{{.Field}} {
	if len(candidates) == 0 && len(excl) == 0 {
		return CliquesOf{{$nodeName}}sBy{{.Field}}{ { {{.KeyField}}s: incl, {{.Field}}: common{{.Field}} } }
	}
	maxCliques := make(CliquesOf{{$nodeName}}sBy{{.Field}}, 0)
	for i, candidate := range candidates {
		remaining{{.Field}} := candidate.Get{{.Field}}()
		if len(common{{.Field}}) > 0 {
			remaining{{.Field}} = remaining{{.Field}}.Intersect(common{{.Field}})
		}
		if len(remaining{{.Field}}) >= min{{.Field}} {
			var remainingCandidates {{$nodeName}}s
			var remainingExcl {{$keyName}}s
			for j, c := range candidates {
				if len(c.Get{{.Field}}().Intersect(candidate.Get{{.Field}}())) >= min{{.Field}} {
					if i < j {
						remainingCandidates = append(remainingCandidates, c)
					} else if j < i {
						remainingExcl = append(remainingExcl, c.Key)
					}
				}
			}
			maxCliques = append(maxCliques, g.maxCliquesOf{{$nodeName}}sBy{{.Field}}(
				incl.With(candidate.Key), remainingExcl, remainingCandidates, remaining{{.Field}}, min{{.Field}})...)
		}
	}
	return maxCliques
}

{{end -}}
`
