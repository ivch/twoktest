# 2K TestTask

This is an implementation of the test task for the FSL language parser

The app receives input file(s) with the FSL code.
The FSL code can include variables and function calls. The FSL expects that the parser app implements basic operations:

- add: `{"cmd": "add", "id": "var1", "operand1": 1, "operand2": 1}`
- subtract: `{"cmd": "subtract", "id": "var1", "operand1": 1, "operand2": 1}`
- multiply: `{"cmd": "multiply", "id": "var1", "operand1": 1, "operand2": 2}`
- divide: `{"cmd": "divide", "id": "var1", "operand1": 1, "operand2": 2}`
- create: `{"cmd": "create", "id": "var", "value": 1}`
- delete: `{"cmd": "delete", "id": "var1"}`
- update: `{"cmd": "update", "id": "var1", "value": 1}`
- print: `{"cmd": "print", "id": "var1"}`

The application expects all variable to be numeric

The application supports several files to be passed as input.

### Test Run

`make run`