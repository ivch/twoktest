package eval

import (
	"encoding/json"
	"errors"
	"io"
)

const (
	funcInit     = "init"
	funcCreate   = "create"
	funcDelete   = "delete"
	funcUpdate   = "update"
	funcPrint    = "print"
	funcAdd      = "add"
	funcSub      = "sub"
	funcDiv      = "div"
	funcMultiply = "multiply"
	callKey      = "cmd"
	varKey       = "id"
	valueKey     = "value"
	operand1Key  = "operand1"
	operand2Key  = "operand2"
)

type Eval struct {
	out   io.Writer
	vars  variables
	funcs map[string][]cmd
}

type cmd map[string]interface{}

// New creates new instance of Eval
func New(out io.Writer) *Eval {
	return &Eval{
		out:   out,
		vars:  make(map[string]float64),
		funcs: make(map[string][]cmd),
	}
}

// Run runs evaluation of the given code
func (e *Eval) Run(in []byte) error {
	if err := e.parseInput(in); err != nil {
		return err
	}
	return e.run(funcInit, nil, nil)
}

func (e *Eval) run(name string, commands []cmd, called map[string]interface{}) error {
	_, ok := e.funcs[name]
	if !ok {
		return errors.New("no function: " + name)
	}

	for i := range commands {
		call, err := e.evalCommand(commands[i])
		if err != nil {
			return err
		}

		var callErr error
		switch call {
		case funcUpdate:
			callErr = e.update(commands[i])
		case funcPrint:
			callErr = e.print(commands[i])
		case funcAdd:
			callErr = e.add(commands[i], called)
		case funcSub:
			callErr = e.sub(commands[i], called)
		case funcDiv:
			callErr = e.div(commands[i], called)
		case funcMultiply:
			callErr = e.multiply(commands[i], called)
		case funcCreate:
			callErr = e.create(commands[i])
		case funcDelete:
			callErr = e.delete(commands[i])
		}

		if callErr != nil {
			return callErr
		}
	}
	return nil
}

func (e *Eval) evalCommand(in map[string]interface{}) (string, error) {
	c, ok := in[callKey]
	if !ok {
		return "", errors.New("wrong function command call")
	}

	command, ok := c.(string)
	if !ok {
		return "", errors.New("wrong function command")
	}

	if len(command) == 0 {
		return "", errors.New("wrong function command size")
	}

	switch command[0] {
	case '#':
		cmds, ok := e.funcs[command[1:]]
		if !ok {
			return "", errors.New("no function: " + command[1:])
		}
		return "", e.run(command[1:], cmds, in)
	case '$':
	default:
		return command, nil
	}
	return "", nil
}

func (e *Eval) parseInput(in []byte) error {
	code := make(map[string]interface{})
	if err := json.Unmarshal(in, &code); err != nil {
		return err
	}

	for k := range code {
		switch v := code[k].(type) {
		case float64:
			e.vars[k] = v
		case []interface{}:
			functions := make([]cmd, len(v))
			for name, f := range v {
				fn, ok := f.(map[string]interface{})
				if !ok {
					return errors.New("incorrect cmd input")
				}
				functions[name] = fn
			}
			e.funcs[k] = functions
		default:
			return errors.New("unexpected type")
		}
	}
	return nil
}
