package object

type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnvironment() *Environment {
	return &Environment{store: map[string]Object{}}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	return &Environment{store: map[string]Object{}, outer: outer}
}

func (e *Environment) Get(name string) (Object, *Environment) {
	if v, ok := e.store[name]; ok {
		return v, e
	} else if e.outer != nil {
		return e.outer.Get(name)
	}
	return nil, nil
}

func (e *Environment) IsExist(name string) bool {
	_, ok := e.store[name]
	return ok
}

func (e *Environment) Set(name string, val Object) {
	e.store[name] = val
}
