package templates

const HelperTemplate = `// Generated automatically by gengraph: do not edit manually

{{if .Tags -}}
// +build {{.Tags}}

{{end -}}
{{$keyName := .KeyName -}}{{$nodeName := .NodeName -}}
package {{.PackageName}}

type Migrate{{$nodeName}}s func(similarNodes []*{{$nodeName}}, similarKeys {{$keyName}}s, clusterNode *{{$nodeName}})

func (g *{{.GraphName}}) ClusterSimilar{{$nodeName}}s(minimumGroupSize int,
	keys {{$keyName}}s,
	areSimilarNodes {{$nodeName}}Matcher,
	createClusterNode func(similarNode {{$nodeName}}, similarKeys {{$keyName}}s) *{{$nodeName}},
	migrateNodes Migrate{{$nodeName}}s) {{$keyName}}s {

	migratedNodes := make({{$keyName}}s, 0)
	for idx1, key1 := range keys {
		if !migratedNodes.Contains(key1) {
			node1 := g.Get{{$nodeName}}(key1)
			similarNodes := []*{{$nodeName}}{node1}
			similarKeys := {{$keyName}}s{key1}
			for idx2, key2 := range keys {
				if idx1 < idx2 && !migratedNodes.Contains(key2) {
					node2 := g.Get{{$nodeName}}(key2)
					if AreEqual{{$keyName}}s(node1.{{.Parent.Field}}, node2.{{.Parent.Field}}) && areSimilarNodes(*node1, *node2) {
						similarNodes = append(similarNodes, node2)
						similarKeys = append(similarKeys, key2)
					}
				}
			}
			if len(similarNodes) >= minimumGroupSize {
				clusterNode := createClusterNode(*node1, similarKeys)
				migrateNodes(similarNodes, similarKeys, clusterNode)
				migratedNodes = append(migratedNodes, similarKeys...)
			}
		}
	}
	return migratedNodes
}

func (g *{{.GraphName}}) DefaultMigrate{{$nodeName}}s(connect bool)  Migrate{{$nodeName}}s {
	return func(sourceNodes []*{{$nodeName}}, sourceKeys {{$keyName}}s, targetNode *{{$nodeName}}) {
		if common{{.Parent.Field}} := g.Common{{.Parent.Field}}sOf{{$nodeName}}s(As{{$nodeName}}s(sourceNodes)); common{{.Parent.Field}} != nil {
			for _, sourceNode := range sourceNodes {
				g.Migrate{{.Parent.Field}}Of{{$nodeName}}(sourceNode, targetNode, common{{.Parent.Field}})
			}
		}
		g.MigrateEdgesOf{{$nodeName}}s(sourceNodes, targetNode)
		if connect {
			for _, sourceNode := range sourceNodes {
				sourceNode.Connect{{.Parent.Field}}(targetNode)
			}
		}
	}
}

func (g *{{.GraphName}}) DeepMigrate{{.Children.Field}}OfAll{{$nodeName}}s() {
	for _, node := range g.{{.GraphField}} {
		if !node.Has{{.Parent.Field}}() {
			g.DeepMigrate{{.Children.Field}}Of{{$nodeName}}(node)
		}
	}
}

func (g *{{.GraphName}}) DeepMigrate{{.Children.Field}}Of{{$nodeName}}(parentNode *{{$nodeName}}) {
	childNodes := g.Get{{$nodeName}}s(parentNode.{{.Children.Field}})
	for _, childNode := range childNodes {
		g.DeepMigrate{{.Children.Field}}Of{{$nodeName}}(childNode)
	}
	g.MigrateEdgesOf{{$nodeName}}s(childNodes, parentNode)
}

func (g *{{.GraphName}}) MigrateEdgesOf{{$nodeName}}s(sourceNodes []*{{$nodeName}}, targetNode *{{$nodeName}}) {
	{{range .Callers -}}
	{
		commonCallers := g.Common{{.Field}}Of{{$nodeName}}s(As{{$nodeName}}s(sourceNodes))
		for _, sourceNode := range sourceNodes {
			g.Migrate{{.Field}}Of{{$nodeName}}(sourceNode, targetNode, commonCallers)
		}
	}
	{{end -}}
	{{range .Callees -}}
	{
		commonCallees := g.Common{{.Field}}Of{{$nodeName}}s(As{{$nodeName}}s(sourceNodes))
		for _, sourceNode := range sourceNodes {
			g.Migrate{{.Field}}Of{{$nodeName}}(sourceNode, targetNode, commonCallees)
		}
	}
	{{end -}}
}

func (g *{{.GraphName}}) Merge{{.Children.Field}}Into{{$nodeName}}(parentNode *{{$nodeName}}, childFilter {{$keyName}}Filter) {
	for _, key := range parentNode.{{.Children.Field}} {
		if childFilter(key) {
			child := g.Get{{$nodeName}}(key)
			g.MigrateFromTo{{$nodeName}}(child, parentNode)
			g.DisconnectAndRemove{{$nodeName}}(child)
		}
	}
}

func (g *{{.GraphName}}) MigrateFromTo{{$nodeName}}(sourceNode *{{$nodeName}}, targetNode *{{$nodeName}}) {
	{{range .Callers -}}
	g.Migrate{{.Field}}Of{{$nodeName}}(sourceNode, targetNode, sourceNode.Get{{.Field}}())
	{{end -}}
	{{range .Callees -}}
	g.Migrate{{.Field}}Of{{$nodeName}}(sourceNode, targetNode, sourceNode.Get{{.Field}}())
	{{end -}}
}

func (g *{{.GraphName}}) Discover{{$nodeName}}Clusters(keys {{$keyName}}s, discover func(g *{{.GraphName}}, remaining {{$keyName}}s) {{$keyName}}s) {
	remaining := keys
	for {
		if clustered := discover(g, remaining); len(clustered) > 0 {
			remaining = remaining.WithoutAnyOf(clustered)
		} else {
			break
		}
	}
}
`
