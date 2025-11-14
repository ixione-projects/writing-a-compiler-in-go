package object

type Environment struct {
	Values    map[string]Object
	Enclosing *Environment
}

func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		Values:    make(map[string]Object),
		Enclosing: enclosing,
	}
}

func (env *Environment) Get(ident string) (Object, bool) {
	value, found := env.Values[ident]
	if !found && env.Enclosing != nil {
		value, found = env.Enclosing.Get(ident)
	}
	return value, found
}

func (env *Environment) Set(ident string, value Object) {
	env.Values[ident] = value
}

func (env *Environment) Length() int {
	return len(env.Values)
}
