package graphAnnotation

import (
	"github.com/MarcGrol/golangAnnotations/generator/annotation"
	"github.com/MarcGrol/golangAnnotations/model"
)

const (
	TypeKey       = "Key"
	TypeNode      = "Node"
	TypeGraph     = "Graph"
	FieldTagGraph = "graph"
)

var register annotation.AnnotationRegister

func init() {
	register = annotation.NewRegistry([]annotation.AnnotationDescriptor{
		{
			Name:       TypeKey,
			ParamNames: []string{},
			Validator:  func(a annotation.Annotation) bool { return a.Name == TypeKey },
		},
		{
			Name:       TypeNode,
			ParamNames: []string{},
			Validator:  func(a annotation.Annotation) bool { return a.Name == TypeNode },
		},
		{
			Name:       TypeGraph,
			ParamNames: []string{},
			Validator:  func(a annotation.Annotation) bool { return a.Name == TypeGraph },
		},
	})
}

func IsKeyStruct(s model.Struct) bool {
	_, ok := register.ResolveAnnotationByName(s.DocLines, TypeKey)
	return ok
}

func IsNodeStruct(s model.Struct) bool {
	_, ok := register.ResolveAnnotationByName(s.DocLines, TypeNode)
	return ok
}

func IsGraphStruct(s model.Struct) bool {
	_, ok := register.ResolveAnnotationByName(s.DocLines, TypeGraph)
	return ok
}
