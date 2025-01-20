package ptlbuilder

type Field struct {
	Name string
	Type string
}

type Command struct {
	Name   string
	Fields []Field
}

type ProtocolSpec struct {
	Package  string
	Commands []Command
}