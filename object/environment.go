package object

var ENV *Environment

type Environment struct {
	store map[string]Object
}

func newEnvironment() *Environment {
	return &Environment{store: map[string]Object{}}
}

func InitEnvironment() {
	ENV = newEnvironment()
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

func init() {
	InitEnvironment()
}
