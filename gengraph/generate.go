package main

import (
	"fmt"
	"strings"

	"github.com/MarcGrol/golangAnnotations/generator/generationUtil"
	"github.com/MarcGrol/golangAnnotations/model"
	"github.com/jinzhu/inflection"

	"github.com/mmartina/gograph/gengraph/graphAnnotation"
	"github.com/mmartina/gograph/gengraph/templates"
)

type FieldName string

type KeyName string
type NodeName string
type EdgeName string
type GraphName string

type keyField struct {
	Field   FieldName
	SField  string
	KeyName KeyName
}

type keyTypeField struct {
	Field    FieldName
	TypeName string
}

type nodeByKeyField struct {
	Field    FieldName
	KeyName  KeyName
	NodeName NodeName
}

type key struct {
	keyName KeyName
	keyType keyTypeField
}

type relatedNode struct {
	Field        FieldName
	SField       string
	KeyName      KeyName
	KeyField     FieldName // resolved
	NodeName     NodeName  // resolved
	EdgeName     EdgeName
	ReverseField FieldName // resolved
}

type node struct {
	Tags         string
	PackageName  string
	KeyField     FieldName
	KeyName      KeyName
	KeyTypeField FieldName // resolved
	KeyTypeName  string    // resolved
	NodeName     NodeName
	Parent       *keyField
	Children     *keyField
	Callers      []relatedNode
	Callees      []relatedNode
	GraphName    GraphName // resolved
	GraphField   FieldName // resolved
}

type graph struct {
	Tags        string
	PackageName string
	GraphName   GraphName
	nodeByKey   []nodeByKeyField
	Nodes       []*node // resolved
}

type Generator struct {
	keyByName         map[KeyName]key
	nodeByName        map[NodeName]*node
	nodeNameByKeyName map[KeyName]NodeName
	graphs            []graph
}

func NewGenerator() *Generator {
	return &Generator{
		keyByName:         make(map[KeyName]key),
		nodeByName:        make(map[NodeName]*node),
		nodeNameByKeyName: make(map[KeyName]NodeName),
		graphs:            make([]graph, 0),
	}
}

func (g *Generator) createKeys(parsedSources model.ParsedSources) error {
	for _, s := range parsedSources.Structs {
		if graphAnnotation.IsKeyStruct(s) {
			key, err := createKey(s)
			if err != nil {
				return fmt.Errorf("Error parsing key %s: %s", s.Name, err)
			}
			g.keyByName[key.keyName] = *key
		}
	}
	return nil
}

func createKey(s model.Struct) (*key, error) {
	key := &key{
		keyName: KeyName(s.Name),
	}
	for _, f := range s.Fields {
		if f.GetTagMap()[graphAnnotation.FieldTagGraph] == "type" {
			if f.IsSlice() || f.IsPointer() {
				return nil, fmt.Errorf("field %s invalid type", f.Name)
			}
			key.keyType = keyTypeField{Field: FieldName(f.Name), TypeName: f.TypeName}
			break
		}
	}
	return key, nil
}

func (g *Generator) createNodes(parsedSources model.ParsedSources) error {
	for _, s := range parsedSources.Structs {
		if graphAnnotation.IsNodeStruct(s) {
			node, err := createNode(s)
			if err != nil {
				return fmt.Errorf("Error parsing node %s: %s", s.Name, err)
			}
			if _, ok := g.keyByName[node.KeyName]; !ok {
				return fmt.Errorf("Invalid node %s: wrong key", s.Name)
			}
			g.nodeByName[node.NodeName] = node
			g.nodeNameByKeyName[node.KeyName] = node.NodeName
		}
	}
	return nil
}

func createNode(s model.Struct) (*node, error) {
	node := &node{
		PackageName: s.PackageName,
		NodeName:    NodeName(s.Name),
		Callers:     []relatedNode{},
		Callees:     []relatedNode{},
	}
	for _, f := range s.Fields {
		fieldName := FieldName(f.Name)
		tagMap := f.GetTagMap()
		switch tagMap[graphAnnotation.FieldTagGraph] {
		case "key":
			if f.IsSlice() || f.IsPointer() {
				return nil, fmt.Errorf("field %s invalid type", fieldName)
			}
			node.KeyField = fieldName
			node.KeyName = KeyName(f.TypeName)
		case "parent":
			if !f.IsPointer() {
				return nil, fmt.Errorf("field %s invalid type", fieldName)
			}
			node.Parent = &keyField{Field: fieldName, KeyName: KeyName(f.DereferencedTypeName())}
		case "children":
			if !f.IsSlice() {
				return nil, fmt.Errorf("field %s invalid type", fieldName)
			}
			node.Children = &keyField{Field: fieldName, KeyName: KeyName(f.SliceElementTypeName())}
		case "callers":
			if f.IsSlice() {
				node.Callers = append(node.Callers, relatedNode{Field: fieldName, KeyName: KeyName(f.SliceElementTypeName())})
			} else if f.IsMap() {
				keyName, edgeName := f.SplitMapTypeNames()
				if isPointer(keyName) || isSlice(keyName) {
					return nil, fmt.Errorf("field %s invalid key", fieldName)
				}
				if isPointer(edgeName) || isSlice(edgeName) {
					return nil, fmt.Errorf("field %s invalid value", fieldName)
				}
				node.Callers = append(node.Callers, relatedNode{Field: fieldName, KeyName: KeyName(keyName), EdgeName: EdgeName(edgeName)})
			} else {
				return nil, fmt.Errorf("field %s invalid type", fieldName)
			}
		case "callees":
			if f.IsSlice() {
				node.Callees = append(node.Callees, relatedNode{Field: fieldName, KeyName: KeyName(f.SliceElementTypeName())})
			} else if f.IsMap() {
				keyName, edgeName := f.SplitMapTypeNames()
				if isPointer(keyName) || isSlice(keyName) {
					return nil, fmt.Errorf("field %s invalid key", fieldName)
				}
				if isPointer(edgeName) || isSlice(edgeName) {
					return nil, fmt.Errorf("field %s invalid value", fieldName)
				}
				node.Callees = append(node.Callees, relatedNode{Field: fieldName, KeyName: KeyName(keyName), EdgeName: EdgeName(edgeName)})
			} else {
				return nil, fmt.Errorf("field %s invalid type", fieldName)
			}
		}
	}
	return node, nil
}

func isPointer(name string) bool {
	return strings.HasPrefix(name, "*")
}

func isSlice(name string) bool {
	return strings.HasPrefix(name, "[]")
}

func (g *Generator) validateNodes() error {
	for _, n := range g.nodeByName {
		if n.Parent != nil {
			if n.Parent.KeyName != n.KeyName {
				return fmt.Errorf("Invalid node %s: wrong parent", n.NodeName)
			}
			if n.Children == nil {
				return fmt.Errorf("Invalid node %s: parent without children", n.NodeName)
			}
		}
		if n.Children != nil {
			if n.Children.KeyName != n.KeyName {
				return fmt.Errorf("Invalid node %s: wrong child", n.NodeName)
			}
			if n.Parent == nil {
				return fmt.Errorf("Invalid node %s: children without parent", n.NodeName)
			}
		}
	}
	return nil
}

func (g *Generator) enrichNodes() error {

	// phase 1: resolve node.keyType, node.CallerByField and node.CalleeByField
	for _, nodeName := range g.nodeNameByKeyName {
		n := g.nodeByName[nodeName]
		if n.Children != nil {
			n.Children.SField = inflection.Singular(string(n.Children.Field))
		}
		keyType := g.keyByName[n.KeyName].keyType
		n.KeyTypeField = keyType.Field
		n.KeyTypeName = keyType.TypeName
		for idx, callerField := range n.Callers {
			callerNodeName, ok := g.nodeNameByKeyName[callerField.KeyName]
			if !ok {
				return fmt.Errorf("Invalid node %s: wrong caller", n.NodeName)
			}
			caller := g.nodeByName[callerNodeName]
			n.Callers[idx].SField = inflection.Singular(string(callerField.Field))
			n.Callers[idx].KeyField = caller.KeyField
			n.Callers[idx].NodeName = caller.NodeName
		}
		for idx, calleeField := range n.Callees {
			calleeNodeName, ok := g.nodeNameByKeyName[calleeField.KeyName]
			if !ok {
				return fmt.Errorf("Invalid node %s: wrong callee", n.NodeName)
			}
			callee := g.nodeByName[calleeNodeName]
			n.Callees[idx].SField = inflection.Singular(string(calleeField.Field))
			n.Callees[idx].KeyField = callee.KeyField
			n.Callees[idx].NodeName = callee.NodeName
		}
	}

	// phase 2: resolve relatedNode.ReverseField
	for _, n := range g.nodeByName {
		for idx, calleeField := range n.Callees {
			callee := g.nodeByName[calleeField.NodeName]
			var reverseField FieldName
			for _, callerField := range callee.Callers {
				if callerField.NodeName == n.NodeName {
					if reverseField != "" {
						return fmt.Errorf("Invalid callee %s -> %s: ambiguous reverse caller", n.NodeName, calleeField.NodeName)
					}
					if callerField.EdgeName != calleeField.EdgeName {
						return fmt.Errorf("Invalid callee %s -> %s: wrong reverse edgeName", n.NodeName, calleeField.NodeName)
					}
					reverseField = callerField.Field
				}
			}
			if reverseField == "" {
				return fmt.Errorf("Invalid callee %s -> %s: missing reverse caller", n.NodeName, calleeField.NodeName)
			}
			n.Callees[idx].ReverseField = reverseField
		}
		for idx, callerField := range n.Callers {
			caller := g.nodeByName[callerField.NodeName]
			var reverseField FieldName
			for _, calleeField := range caller.Callees {
				if calleeField.NodeName == n.NodeName {
					if reverseField != "" {
						return fmt.Errorf("Invalid caller %s -> %s: ambiguous reverse callee", n.NodeName, callerField.NodeName)
					}
					if calleeField.EdgeName != callerField.EdgeName {
						return fmt.Errorf("Invalid callee %s -> %s: wrong reverse edgeName", n.NodeName, calleeField.NodeName)
					}
					reverseField = calleeField.Field
				}
			}
			if reverseField == "" {
				return fmt.Errorf("Invalid caller %s -> %s: missing reverse callee", n.NodeName, callerField.NodeName)
			}
			n.Callers[idx].ReverseField = reverseField
		}
	}
	return nil
}

func (g *Generator) createGraphs(parsedSources model.ParsedSources) error {
	for _, s := range parsedSources.Structs {
		if graphAnnotation.IsGraphStruct(s) {
			graph, err := createGraph(s)
			if err != nil {
				return fmt.Errorf("Error parsing graph %s: %s", s.Name, err)
			}
			for _, keyNode := range graph.nodeByKey {
				graph.Nodes = append(graph.Nodes, g.nodeByName[keyNode.NodeName])
			}
			g.graphs = append(g.graphs, *graph)

			// set node.GraphName and node.GraphField
			for _, ntField := range graph.nodeByKey {
				nt, ok := g.nodeByName[ntField.NodeName]
				if !ok {
					return fmt.Errorf("Invalid graph %s: wrong node map", graph.GraphName, ntField.Field)
				}
				nt.GraphName = GraphName(graph.GraphName)
				nt.GraphField = ntField.Field
			}
		}
	}
	return nil
}

func createGraph(s model.Struct) (*graph, error) {
	graph := graph{
		PackageName: s.PackageName,
		GraphName:   GraphName(s.Name),
		nodeByKey:   make([]nodeByKeyField, 0),
		Nodes:       make([]*node, 0),
	}
	for _, f := range s.Fields {
		fieldName := FieldName(f.Name)
		tagMap := f.GetTagMap()
		switch tagMap[graphAnnotation.FieldTagGraph] {
		case "nodes":
			if !f.IsMap() {
				return nil, fmt.Errorf("field %s invalid type", fieldName)
			}
			keyName, rawNodeName := f.SplitMapTypeNames()
			if isPointer(keyName) || isSlice(keyName) {
				return nil, fmt.Errorf("field %s invalid key", fieldName)
			}
			if !isPointer(rawNodeName) {
				return nil, fmt.Errorf("field %s invalid value", fieldName)
			}
			nodeName := NodeName(strings.TrimPrefix(rawNodeName, "*"))
			graph.nodeByKey = append(graph.nodeByKey, nodeByKeyField{Field: fieldName, KeyName: KeyName(keyName), NodeName: nodeName})
		}
	}
	return &graph, nil
}

func (g *Generator) build(parsedSources model.ParsedSources) error {

	if err := g.createKeys(parsedSources); err != nil {
		return err
	}

	if err := g.createNodes(parsedSources); err != nil {
		return err
	}

	if err := g.validateNodes(); err != nil {
		return err
	}

	if err := g.enrichNodes(); err != nil {
		return err
	}

	if err := g.createGraphs(parsedSources); err != nil {
		return err
	}

	return nil
}

func (g *Generator) Generate(inputDir string, tags string, parsedSources model.ParsedSources) error {

	packageName, err := generationUtil.GetPackageNameForStructs(parsedSources.Structs)
	if err != nil {
		return err
	}

	targetDir, err := generationUtil.DetermineTargetPath(inputDir, packageName)
	if err != nil {
		return err
	}

	if err := g.build(parsedSources); err != nil {
		return err
	}

	for _, node := range g.nodeByName {
		node.Tags = tags
		if err = g.doGenerateForNode(packageName, targetDir, *node); err != nil {
			return err
		}
	}

	for _, graph := range g.graphs {
		graph.Tags = tags
		if err = g.doGenerateForGraph(packageName, targetDir, graph); err != nil {
			return err
		}
	}

	return nil
}

func (g *Generator) doGenerateForNode(packageName string, targetDir string, node node) error {
	if err := doGenerateFor(templates.KeyTemplate, "key", packageName, targetDir, string(node.NodeName), node); err != nil {
		return err
	}
	if err := doGenerateFor(templates.NodeTemplate, "node", packageName, targetDir, string(node.NodeName), node); err != nil {
		return err
	}
	if err := doGenerateFor(templates.ParentTemplate, "parent", packageName, targetDir, string(node.NodeName), node); err != nil {
		return err
	}
	if err := doGenerateFor(templates.ChildrenTemplate, "children", packageName, targetDir, string(node.NodeName), node); err != nil {
		return err
	}
	if len(node.Callers) > 0 {
		if err := doGenerateFor(templates.CallersTemplate, "callers", packageName, targetDir, string(node.NodeName), node); err != nil {
			return err
		}
	}
	if len(node.Callees) > 0 {
		if err := doGenerateFor(templates.CalleesTemplate, "callees", packageName, targetDir, string(node.NodeName), node); err != nil {
			return err
		}
	}
	if err := doGenerateFor(templates.GraphTemplate, "graph", packageName, targetDir, string(node.NodeName), node); err != nil {
		return err
	}
	if err := doGenerateFor(templates.FilterTemplate, "filter", packageName, targetDir, string(node.NodeName), node); err != nil {
		return err
	}
	if err := doGenerateFor(templates.MatchTemplate, "match", packageName, targetDir, string(node.NodeName), node); err != nil {
		return err
	}
	if err := doGenerateFor(templates.HelperTemplate, "helper", packageName, targetDir, string(node.NodeName), node); err != nil {
		return err
	}
	if err := doGenerateFor(templates.TraverseTemplate, "traverse", packageName, targetDir, string(node.NodeName), node); err != nil {
		return err
	}
	if err := doGenerateFor(templates.ConnectTemplate, "connect", packageName, targetDir, string(node.NodeName), node); err != nil {
		return err
	}
	return nil
}

func (g *Generator) doGenerateForGraph(packageName string, targetDir string, graph graph) error {
	if err := doGenerateFor(templates.TraverserTemplate, "traverser", packageName, targetDir, string(graph.GraphName), graph); err != nil {
		return err
	}
	if err := doGenerateFor(templates.ConnectedTemplate, "connected", packageName, targetDir, string(graph.GraphName), graph); err != nil {
		return err
	}
	return nil
}

func doGenerateFor(template string, templateName string, packageName string, targetDir string, name string, data interface{}) error {

	filename := fmt.Sprintf("%s_%s.go", name, templateName)
	target := generationUtil.Prefixed(fmt.Sprintf("%s/%s", targetDir, filename))

	err := generationUtil.Generate(generationUtil.Info{
		Src:            packageName,
		TargetFilename: target,
		TemplateName:   templateName,
		TemplateString: template,
		Data:           data,
	})
	if err != nil {
		return fmt.Errorf("Error generating %s: %s", filename, err)
	}

	return nil
}
