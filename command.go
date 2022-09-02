package main

type HandlerFunc func([]string) error

type Command struct {
	Name    string
	Usage   string
	Handler HandlerFunc
}

func NewCommand(name, usage string, handler HandlerFunc) *Command {
	return &Command{Name: name, Usage: usage, Handler: handler}
}
