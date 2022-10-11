package eval

import (
	"errors"
	"io/ioutil"
	"testing"
)

func Test_Print(t *testing.T) {
	e := New(ioutil.Discard)

	if err := e.parseInput(testInput); err != nil {
		t.Fatal(err)
	}

	if err := e.print(map[string]interface{}{"x": 1}); err != nil && !errors.Is(err, errVarNotFound) {
		t.Fatalf("expected errVarNotFound, got %v", err)
	}

	if err := e.print(map[string]interface{}{"value": 1}); err != nil && !errors.Is(err, errVarBadName) {
		t.Fatalf("expected errVarBadName, got %v", err)
	}

	if err := e.print(map[string]interface{}{"value": "x"}); err != nil && !errors.Is(err, errVarBadRef) {
		t.Fatalf("expected errVarBadRef, got %v", err)
	}
}

func Test_GetTargetVariable(t *testing.T) {
	e := New(ioutil.Discard)
	if err := e.parseInput(testInput); err != nil {
		t.Fatal(err)
	}

	if _, err := e.getTargetVariable(map[string]interface{}{"x": 1}, nil); err != nil && !errors.Is(err, errVarNotFound) {
		t.Fatalf("expected errVarNotFound, got %v", err)
	}

	if _, err := e.getTargetVariable(map[string]interface{}{"id": 1}, nil); err != nil && !errors.Is(err, errVarBadName) {
		t.Fatalf("expected errVarBadName, got %v", err)
	}
}

func Test_Update(t *testing.T) {
	e := New(ioutil.Discard)
	if err := e.parseInput(testInput); err != nil {
		t.Fatal(err)
	}

	if err := e.update(map[string]interface{}{"x": 1}); err == nil {
		t.Fatal("expected error")
	}

	if err := e.update(map[string]interface{}{"id": "var1", "x": 3.5}); err == nil {
		t.Fatal("expected error")
	}
}

func Test_Sub(t *testing.T) {
	e := New(ioutil.Discard)
	if err := e.parseInput(testInput); err != nil {
		t.Fatal(err)
	}

	if err := e.sub(map[string]interface{}{"id": "$id", "operand1": "$value1", "operand2": "$value2"},
		map[string]interface{}{"id": "var1", "value1": "#var1", "value2": "#var2"}); err != nil {
		t.Fatal(err)
	}

	if e.vars["var1"] != -1 {
		t.Fatal("wrong sub result")
	}
}

func Test_Multiply(t *testing.T) {
	e := New(ioutil.Discard)
	if err := e.parseInput(testInput); err != nil {
		t.Fatal(err)
	}

	if err := e.multiply(map[string]interface{}{"id": "$id", "operand1": "$value1", "operand2": "$value2"},
		map[string]interface{}{"id": "var1", "value1": "#var1", "value2": "#var2"}); err != nil {
		t.Fatal(err)
	}

	if e.vars["var1"] != 2 {
		t.Fatal("wrong multiply result")
	}
}

func Test_Div(t *testing.T) {
	e := New(ioutil.Discard)
	if err := e.parseInput(testInput); err != nil {
		t.Fatal(err)
	}

	if err := e.div(map[string]interface{}{"id": "$id", "operand1": "$value1", "operand2": "$value2"},
		map[string]interface{}{"id": "var1", "value1": "#var1", "value2": "#var2"}); err != nil {
		t.Fatal(err)
	}

	if e.vars["var1"] != .5 {
		t.Fatal("wrong div result")
	}
}
