package eval

import (
	"bytes"
	"errors"
	"io/ioutil"
	"testing"
)

func Test_RunOK(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	e := New(buf)
	if err := e.Run(testInput); err != nil {
		t.Fatal(err)
	}

	if buf.String() != testOutput {
		t.Errorf("got %q, want %q", buf.String(), testOutput)
	}
}

func Test_RunErrors(t *testing.T) {
	e := New(ioutil.Discard)
	if err := e.parseInput([]byte(`{`)); err == nil {
		t.Fatal("expected unmarshal error")
	}

	if err := e.parseInput([]byte(`{"init": {"cmd" : "#setup" }}`)); err != nil && !errors.Is(err, errParseBadType) {
		t.Fatalf("expected errParseBadType, got %v", err)
	}

	if err := e.parseInput([]byte(`{"init": [1,2]}`)); err != nil && !errors.Is(err, errParseBadInput) {
		t.Fatalf("expected errParseBadInput, got %v", err)
	}
}

var testOutput = "3.5\n5.5\nvariable not found\n2\n5\n"
var testInput = []byte(`{
  "var1":1,
  "var2":2,
  
  "init": [
    {"cmd" : "#setup" }
  ],
  
  "setup": [
    {"cmd":"update", "id": "var1", "value":3.5},
    {"cmd":"print", "value": "#var1"},
    {"cmd":"#sum", "id": "var1", "value1":"#var1", "value2":"#var2"},
    {"cmd":"print", "value": "#var1"},
    {"cmd":"create", "id": "var3", "value":5},
    {"cmd":"delete", "id": "var1"},
    {"cmd":"#printAll"}
  ],
  
  "sum": [
      {"cmd":"add", "id": "$id", "operand1":"$value1", "operand2":"$value2"}
  ],

  "printAll":
  [
    {"cmd":"print", "value": "#var1"},
    {"cmd":"print", "value": "#var2"},
    {"cmd":"print", "value": "#var3"}
  ]
}`)
