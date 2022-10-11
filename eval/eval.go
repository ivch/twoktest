package eval

import (
	"encoding/json"
	"errors"
	"fmt"
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

var (
	errVarExists    = errors.New("variable already exists")
	errVarNotFound  = errors.New("variable not found")
	errVarBadType   = errors.New("bad variable type")
	errVarBadName   = errors.New("bad variable name")
	errVarBadRef    = errors.New("bad variable reference")
	errFuncNotFound = errors.New("function not found")
	errFuncBadCall  = errors.New("function bad call")
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
		return fmt.Errorf("%f: %s", errFuncNotFound, name)
	}

	if commands == nil {
		commands = e.funcs[name]
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
		return "", fmt.Errorf("%w: no command call", errFuncBadCall)
	}

	command, ok := c.(string)
	if !ok {
		return "", fmt.Errorf("%w: wrong function command", errFuncBadCall)
	}

	if len(command) == 0 {
		return "", fmt.Errorf("%w: wrong function command size", errFuncBadCall)
	}

	if command[0] != '#' {
		return command, nil
	}

	cmds, ok := e.funcs[command[1:]]
	if !ok {
		return "", fmt.Errorf("%w: %s", errFuncNotFound, command[1:])
	}
	return "", e.run(command[1:], cmds, in)
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
