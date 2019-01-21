package templates

const KeyTemplate = `// Generated automatically by gengraph: do not edit manually

{{if .Tags -}}
// +build {{.Tags}}

{{end -}}
package {{.PackageName}}

func (keys {{.KeyName}}s) Contains(other {{.KeyName}}) bool {
	for _, k := range keys {
		if k == other {
			return true
		}
	}
	return false
}

func (keys {{.KeyName}}s) ContainsAll(other {{.KeyName}}s) bool {
	for _, k := range other {
		if !keys.Contains(k) {
			return false
		}
	}
	return true
}

func (keys {{.KeyName}}s) Intersect(other {{.KeyName}}s) {{.KeyName}}s {
	remaining := make({{.KeyName}}s, 0)
	for _, k := range keys {
		if other.Contains(k) {
			remaining = append(remaining, k)
		}
	}
	return remaining
}

func (keys {{.KeyName}}s) Union(other {{.KeyName}}s) {{.KeyName}}s {
	for _, k := range other {
		if !keys.Contains(k) {
			keys = append(keys, k)
		}
	}
	return keys
}

func (keys {{.KeyName}}s) With(other {{.KeyName}}) {{.KeyName}}s {
	if keys.Contains(other) {
		return keys
	}
	return append(keys, other)
}

func (keys {{.KeyName}}s) Without(other {{.KeyName}}) {{.KeyName}}s {
	remaining := make({{.KeyName}}s, 0)
	for _, k := range keys {
		if k != other {
			remaining = append(remaining, k)
		}
	}
	return remaining
}

func (keys {{.KeyName}}s) WithoutAnyOf(other {{.KeyName}}s) {{.KeyName}}s {
	remaining := make({{.KeyName}}s, 0)
	for _, k := range keys {
		if !other.Contains(k) {
			remaining = append(remaining, k)
		}
	}
	return remaining
}

func (keys {{.KeyName}}s) LessThan(other {{.KeyName}}s) bool {
	for i := 0; i < len(keys) && i < len(other); i++ {
		if keys[i] != other[i] {
			return keys[i].LessThan(other[i])
		}
	}
	return len(keys) < len(other)
}

func (keys {{.KeyName}}s) Any() *{{.KeyName}} {
	if len(keys) > 0 {
		return &keys[rand.Intn(len(keys))]
	}
	return nil
}

type {{.KeyName}}s []{{.KeyName}}
func (k {{.KeyName}}s) Len() int           { return len(k) }
func (k {{.KeyName}}s) Less(i, j int) bool { return k[i].LessThan(k[j]) }
func (k {{.KeyName}}s) Swap(i, j int)      { k[i], k[j] = k[j], k[i] }
`
