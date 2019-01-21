package templates

const FilterTemplate = `// Generated automatically by gengraph: do not edit manually

{{if .Tags -}}
// +build {{.Tags}}

{{end -}}
{{$keyName := .KeyName -}}{{$nodeName := .NodeName -}}
package {{.PackageName}}

type {{$keyName}}Filter func({{$keyName}}) bool

type {{$nodeName}}Filter func({{$nodeName}}) bool

func FilterBy{{.KeyTypeName}}s(types ...{{.KeyTypeName}}) {{$keyName}}Filter {
	return func(key {{$keyName}}) bool {
		for _, t := range types {
			if key.{{.KeyTypeField}} == t {
				return true
			}
		}
		return false
	}
}

func (keys {{$keyName}}s) FilteredBy(filter {{$keyName}}Filter) {{$keyName}}s {
	remaining := make({{$keyName}}s, 0)
	for _, k := range keys {
		if filter(k) {
			remaining = append(remaining, k)
		}
	}
	return remaining
}

func (nodes {{$nodeName}}s) FilteredBy(filter {{$nodeName}}Filter) {{$nodeName}}s {
	remaining := make({{$nodeName}}s, 0)
	for _, node := range nodes {
		if filter(node) {
			remaining = append(remaining, node)
		}
	}
	return remaining
}

func (keys {{$keyName}}s) FilteredBy{{.KeyTypeName}}s(types ...{{.KeyTypeName}}) {{$keyName}}s {
	return keys.FilteredBy(FilterBy{{.KeyTypeName}}s(types...))
}

func (keys {{$keyName}}s) FilteredByNode(g *{{.GraphName}}, filter {{$nodeName}}Filter) {{$keyName}}s {
	remaining := make({{$keyName}}s, 0)
	for _, k := range keys {
		if node := g.Get{{$nodeName}}(k); node != nil && filter(*node) {
			remaining = append(remaining, k)
		}
	}
	return remaining
}
`
