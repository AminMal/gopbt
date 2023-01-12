package gopbt

type typeGenMapping struct {
	// todo, add named generators in addition to type generators. name priority should be higher than type name
	generatorMapping map[string]anyGen
}

func (mapping *typeGenMapping) setGenerator(typeName string, g anyGen) {
	mapping.generatorMapping[typeName] = g
}
