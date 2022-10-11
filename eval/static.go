package eval

import (
	"fmt"
)

func (e *Eval) add(code, called map[string]interface{}) error {
	return e.callback(code, called, func(a, b float64) float64 {
		return a + b
	})
}

func (e *Eval) sub(code, called map[string]interface{}) error {
	return e.callback(code, called, func(a, b float64) float64 {
		return a - b
	})
}

func (e *Eval) multiply(code, called map[string]interface{}) error {
	return e.callback(code, called, func(a, b float64) float64 {
		return a * b
	})
}

func (e *Eval) div(code, called map[string]interface{}) error {
	return e.callback(code, called, func(a, b float64) float64 {
		return a / b
	})
}

func (e *Eval) update(code map[string]interface{}) error {
	name, err := e.getTargetVariable(code, nil)
	if err != nil {
		return err
	}

	value, err := e.getVarValue(valueKey, code, nil)
	if err != nil {
		return err
	}

	e.vars.set(name, value)

	return nil
}

func (e *Eval) create(code map[string]interface{}) error {
	name, err := e.getTargetVariable(code, nil)
	if err != nil {
		return err
	}

	if _, err := e.vars.get(name); err == nil {
		return errVarExists
	}

	value, err := e.getVarValue(valueKey, code, nil)
	if err != nil {
		return err
	}

	e.vars.set(name, value)
	return nil
}

func (e *Eval) delete(code map[string]interface{}) error {
	name, err := e.getTargetVariable(code, nil)
	if err != nil {
		return err
	}

	if _, err := e.vars.get(name); err != nil {
		return errVarNotFound
	}

	e.vars.delete(name)
	return nil
}

func (e *Eval) callback(code, called map[string]interface{}, f func(float64, float64) float64) error {
	name, err := e.getTargetVariable(code, called)
	if err != nil {
		return err
	}
	a, err := e.getVarValue(operand1Key, code, called)
	if err != nil {
		return err
	}
	b, err := e.getVarValue(operand2Key, code, called)
	if err != nil {
		return err
	}
	e.vars.set(name, f(a, b))
	return nil
}

func (e *Eval) print(code map[string]interface{}) error {
	if _, ok := code[valueKey]; !ok {
		return errVarNotFound
	}

	name, ok := code[valueKey].(string)
	if !ok {
		return errVarBadName
	}

	if name[0] != '#' {
		return errVarBadRef
	}

	value, err := e.vars.get(name[1:])
	if err != nil {
		fmt.Fprintln(e.out, err)
		return nil
	}

	fmt.Fprintln(e.out, value)
	return nil
}

func (e *Eval) getTargetVariable(code, called map[string]interface{}) (string, error) {
	if _, ok := code[varKey]; !ok {
		return "", errVarNotFound
	}

	name, ok := code[varKey].(string)
	if !ok {
		return "", errVarBadName
	}

	if name[0] == '$' {
		return e.getTargetVariable(called, nil)
	}

	return name, nil
}

func (e *Eval) getVarValue(name string, code, called map[string]interface{}) (float64, error) {
	v, ok := code[name]
	if !ok {
		return 0, errVarNotFound
	}

	switch val := v.(type) {
	case float64:
		return val, nil
	case string:
		name = val
	default:
		return 0, errVarBadType
	}

	if len(name) == 0 {
		return 0, errVarBadName
	}

	switch name[0] {
	case '#':
		name = name[1:]
	case '$':
		return e.getVarValue(name[1:], called, nil)
	default:
		return 0, errVarBadName
	}
	return e.vars.get(name)
}
