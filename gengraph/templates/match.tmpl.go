package templates

const MatchTemplate = `// Generated automatically by gengraph: do not edit manually

{{if .Tags -}}
// +build {{.Tags}}

{{end -}}
{{$keyName := .KeyName -}}{{$nodeName := .NodeName -}}
package {{.PackageName}}

import "reflect"

type {{$keyName}}Matcher = func(key1 {{$keyName}}, key2 {{$keyName}}) bool 

type {{$nodeName}}Matcher = func(node1 {{$nodeName}}, node2 {{$nodeName}}) bool 

var (
	Any{{$keyName}}s = func(_ {{$keyName}}, _ {{$keyName}}) bool { return true }
	Any{{$nodeName}}s = func(_ {{$nodeName}}, _ {{$nodeName}}) bool { return true }
)

func AreEqual{{$keyName}}s(key1 *{{$keyName}}, key2 *{{$keyName}}) bool {
	return key1 == nil && key2 == nil || key1 != nil && key2 != nil && *key1 == *key2
}

func (g *{{.GraphName}}) DeepEqual{{$nodeName}}s(keyMatcher {{$keyName}}Matcher) {{$nodeName}}Matcher {
	return func (node1 {{$nodeName}}, node2 {{$nodeName}}) bool {
		if AreEqual{{$keyName}}s(node1.{{.Parent.Field}}, node2.{{.Parent.Field}}) {
			return g.areDeepEqual{{$nodeName}}s(keyMatcher, node1, node2)
		}
		return false
	}
}

func (g *{{.GraphName}}) areDeepEqual{{$nodeName}}s(keyMatcher {{$keyName}}Matcher, node1 {{$nodeName}}, node2 {{$nodeName}}) bool {
	if !keyMatcher(node1.{{.KeyField}}, node2.{{.KeyField}}) ||	len(node1.{{.Children.Field}}) != len(node2.{{.Children.Field}}) {
		return false
	}
	{{range .Callers -}}
	if !reflect.DeepEqual(node1.{{.Field}}, node2.{{.Field}}) {
		return false
	}
	{{end -}}
	{{range .Callees -}}
	if !reflect.DeepEqual(node1.{{.Field}}, node2.{{.Field}}) {
		return false
	}
	{{end -}}
	for i := 0; i < len(node1.{{.Children.Field}}); i++ {
		if !g.areDeepEqual{{$nodeName}}s(keyMatcher, *g.Get{{$nodeName}}(node1.{{.Children.Field}}[i]), *g.Get{{$nodeName}}(node2.{{.Children.Field}}[i])) {
			return false
		}
	}
	return true
}
`
