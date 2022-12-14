package eval

type variables map[string]float64

func (v variables) get(name string) (float64, error) {
	if val, ok := v[name]; ok {
		return val, nil
	}
	return 0, errVarNotFound
}

func (v variables) set(name string, value float64) {
	v[name] = value
}

func (v variables) delete(name string) {
	delete(v, name)
}
