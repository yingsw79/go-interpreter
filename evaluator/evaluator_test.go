package evaluator

import (
	"go-interpreter/lexer"
	"go-interpreter/object"
	"go-interpreter/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated, err := testEval(tt.input)
		if err != nil {
			t.Fatal(err)
		}

		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"true < 2", true},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for _, tt := range tests {
		evaluated, err := testEval(tt.input)
		if err != nil {
			t.Fatal(err)
		}

		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestStringLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"Hello World!"`, "Hello World!"},
		{`"Hello" + " " + "World!"`, "Hello World!"},
		{`"12" * 0`, ""},
		{`"12" * 1`, "12"},
		{`"12" * 10`, "12121212121212121212"},
		{`sum(["1", "23", "456"])`, "123456"},
	}

	for _, tt := range tests {
		evaluated, err := testEval(tt.input)
		if err != nil {
			t.Fatal(err)
		}

		testStringObject(t, evaluated, tt.expected)
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated, err := testEval(tt.input)
		if err != nil {
			t.Fatal(err)
		}

		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		evaluated, err := testEval(tt.input)
		if err != nil {
			t.Fatal(err)
		}

		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{"if (10 > 1) { return 10; }", 10},
		{
			`
if (10 > 1) {
  if (10 > 1) {
    return 10;
  }

  return 1;
}
`, 10,
		},
		{
			`
let f = fn(x) {
  return x;
  x + 10;
};
f(10);`, 10,
		},
		{
			`
let f = fn(x) {
   let res = x + 10;
   return res;
   return 10;
};
f(10);`, 20,
		},
	}

	for _, tt := range tests {
		evaluated, err := testEval(tt.input)
		if err != nil {
			t.Fatal(err)
		}

		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		evaluated, err := testEval(tt.input)
		if err != nil {
			t.Fatal(err)
		}

		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2; };"

	evaluated, err := testEval(input)
	if err != nil {
		t.Fatal(err)
	}

	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%v", fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
	}

	expectedBody := "(x + 2)"

	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let double = fn(x) { x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"let fib = fn(n) { if (n < 2) { return n } return fib(n - 1) + fib(n - 2)}; fib(10)", 55},
		{"fn(x) { x; }(5)", 5},
		{`len("")`, 0},
		{`len([])`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len([1, 2, 3, 4, 5, 6, 7])`, 7},
		{`sum([])`, 0},
		{`sum([], 2)`, 2},
		{`sum([1, 2, 3, 4, 5, 6, 7])`, 28},
		{`sum([1, 2, 3, 4, 5, 6, 7], 8)`, 36},
	}

	for _, tt := range tests {
		evaluated, err := testEval(tt.input)
		if err != nil {
			t.Fatal(err)
		}

		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestClosures(t *testing.T) {
	input := `
let newAdder = fn(x) {
  return fn(y) { x + y };
};

let addTwo = newAdder(2);
addTwo(2);`

	evaluated, err := testEval(input)
	if err != nil {
		t.Fatal(err)
	}

	testIntegerObject(t, evaluated, 4)
}

func TestArrayLiterals(t *testing.T) {
	tests := []struct {
		input    string
		expected []interface{}
	}{
		{"[1, 2 * 2, 3 + 3]", []interface{}{1, 4, 6}},
		{`[1, 2, false, "hello", [3, 4, true]]`, []interface{}{1, 2, false, "hello", []interface{}{3, 4, true}}},
		{`[1, 2, [3, 4, [5, 6, 7]]]`, []interface{}{1, 2, []interface{}{3, 4, []interface{}{5, 6, 7}}}},
		{"[1, 2, 3] + [4, 5, 6]", []interface{}{1, 2, 3, 4, 5, 6}},
		{"[1, 2, 3] + 4", []interface{}{1, 2, 3, 4}},
		{`[1, 2, 3] + "hello"`, []interface{}{1, 2, 3, "hello"}},
		{`4 + [1, 2, 3]`, []interface{}{4, 1, 2, 3}},
		{`[1, 2] * 0`, []interface{}{}},
		{`[1, 2] * 4`, []interface{}{1, 2, 1, 2, 1, 2, 1, 2}},
		{`let a = []; push(a, 1, 2, 3); a`, []interface{}{1, 2, 3}},
		{`let a = []; let b = push(a, 1, 2, 3); b`, []interface{}{1, 2, 3}},
		{`let a = [1, 2, 3, 4]; pop(a); a`, []interface{}{1, 2, 3}},
		{`let a = [1, 2, 3, 4]; reverse(a); a`, []interface{}{4, 3, 2, 1}},
		{`let a = [1, 3, 2, 4, 6, 5, 7]; sort(a); a`, []interface{}{1, 2, 3, 4, 5, 6, 7}},
	}

	for _, tt := range tests {
		evaluated, err := testEval(tt.input)
		if err != nil {
			t.Fatal(err)
		}

		testArrayObject(t, evaluated, tt.expected)
	}
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"[1, 2, 3][0]", 1},
		{"[1, 2, 3][1]", 2},
		{"[1, 2, 3][2]", 3},
		{"let i = 0; [1][i];", 1},
		{"[1, 2, 3][1 + 1];", 3},
		{"let myArray = [1, 2, 3]; myArray[2];", 3},
		{"let myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];", 6},
		{"let myArray = [1, 2, 3]; let i = myArray[0]; myArray[i]", 2},
	}

	for _, tt := range tests {
		evaluated, err := testEval(tt.input)
		if err != nil {
			t.Fatal(err)
		}

		testIntegerObject(t, evaluated, tt.expected)
	}
}

func testEval(input string) (object.Object, error) {
	p := parser.NewParser(lexer.NewLexer(input))
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) {
	res, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%v)", obj, obj)
		return
	}

	if res.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", res.Value, expected)
	}
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) {
	res, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%v)", obj, obj)
		return
	}

	if res.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t", res.Value, expected)
	}
}

func testStringObject(t *testing.T, obj object.Object, expected string) {
	res, ok := obj.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%v)", obj, obj)
		return
	}

	if res.Value != expected {
		t.Errorf("object has wrong value. got='%s', want='%s'", res.Value, expected)
	}
}

func testArrayObject(t *testing.T, obj object.Object, expected []interface{}) {
	res, ok := obj.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%v)", obj, obj)
		return
	}

	if len(res.Elements) != len(expected) {
		t.Errorf("wrong num of elements. want=%d, got=%d", len(expected), len(res.Elements))
		return
	}

	for i, o := range res.Elements {
		switch e := expected[i].(type) {
		case int:
			testIntegerObject(t, o, int64(e))
		case int64:
			testIntegerObject(t, o, e)
		case string:
			testStringObject(t, o, e)
		case bool:
			testBooleanObject(t, o, e)
		case []interface{}:
			testArrayObject(t, o, e)
		case nil:
			testNullObject(t, o)
		default:
			t.Errorf("missing type")
		}
	}
}

func testNullObject(t *testing.T, obj object.Object) {
	if obj != object.NULL {
		t.Errorf("object is not NULL. got=%T (%v)", obj, obj)
	}
}
