package graphvizutil

import . "github.com/awalterschulze/gographviz"

type Attributes map[string]string

func AddEdge(g *Graph, src, dst string, directed bool, attrs map[string]string, reversed bool) error {
	if reversed {
		dir := attrs[string(Dir)]
		lHead := attrs[string(LHead)]
		lTail := attrs[string(LTail)]
		arrowHead := attrs[string(ArrowHead)]
		arrowTail := attrs[string(ArrowTail)]

		if dir == "" || dir == "forward" {
			SetOrDelete(attrs, string(Dir), "back")
		}
		if lHead != "" || lTail != "" {
			SetOrDelete(attrs, string(LTail), lHead)
			SetOrDelete(attrs, string(LHead), lTail)
		}
		if arrowHead != "" || arrowTail != "" {
			SetOrDelete(attrs, string(ArrowTail), arrowHead)
			SetOrDelete(attrs, string(ArrowHead), arrowTail)
		}

		return g.AddEdge(dst, src, directed, attrs)
	} else {
		return g.AddEdge(src, dst, directed, attrs)
	}
}

func SetOrDelete(attrs map[string]string, key string, value string) {
	if value != "" {
		attrs[key] = value
	} else {
		delete(attrs, key)
	}
}
