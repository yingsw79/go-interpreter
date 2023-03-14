package evaluator

import (
	"go-interpreter/lexer"
	"go-interpreter/object"
	"go-interpreter/parser"
	"testing"
)

type test struct {
	input    string
	expected any
}

func TestEvalIntegerExpression(t *testing.T) {
	tests := []test{
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

	autoTest(t, tests)
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []test{
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

	autoTest(t, tests)
}

func TestStringLiteral(t *testing.T) {
	tests := []test{
		{`"Hello World!"`, "Hello World!"},
		{`"Hello" + " " + "World!"`, "Hello World!"},
		{`"12" * 0`, ""},
		{`"12" * 1`, "12"},
		{`"12" * 10`, "12121212121212121212"},
	}

	autoTest(t, tests)
}

func TestPrefixOperator(t *testing.T) {
	tests := []test{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
		{"!!5", true},
		{"~0", -1},
		{"~-1", 0},
		{"~65535", -65536},
	}

	autoTest(t, tests)
}

func TestInfixOperator(t *testing.T) {
	tests := []test{
		{"2 >= 2", true},
		{"1 >= 2", false},
		{"1 <= 1", true},
		{`"ab" < "abc"`, true},
		{"123124 % 13", 1},
		{"234 >> 5 & 1", 1},
		{"1 << 10", 1024},
		{"123 | 321", 379},
		{"0 && 2", 0},
		{"0 || 2", 2},
		{"0 && 1 || 0 && 2", 0},
		{"let a = b = [1, 1, 1]; a[0] = 2; b", []any{2, 1, 1}},
		{"let a = [[1, 2], [3, 4]]; a[0][0] = 5; a[0]", []any{5, 2}},
	}

	autoTest(t, tests)
}

func TestIfElseExpressions(t *testing.T) {
	tests := []test{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	autoTest(t, tests)
}

func TestReturnStatements(t *testing.T) {
	tests := []test{
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

	autoTest(t, tests)
}

func TestLetStatements(t *testing.T) {
	tests := []test{
		{"let a; a", nil},
		{"let a = 5;", nil},
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
		{"let a = 1; let b, a = 2, 2;", "identifier 'a' has already been declared"},
		{"let a, b, c, d = 1, 2, 3, 4; [a, b, c, d]", []any{1, 2, 3, 4}},
	}

	autoTest(t, tests)
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
	tests := []test{
		{"let f = fn(x) { x; }; f(5);", 5},
		{"let f = fn(x) { return x; }; f(5);", 5},
		{"let f = fn(x) { x * 2; }; f(5);", 10},
		{"let f = fn(x, y) { x + y; }; f(5, 5);", 10},
		{"let f = fn(x, y) { x + y; }; f(5 + 5, f(5, 5));", 20},
		{"let f = fn(n) { if (n < 2) { return n } return f(n - 1) + f(n - 2)}; f(10)", 55},
		{"fn(x) { x; }(5)", 5},
	}

	autoTest(t, tests)
}

func TestEnvironment(t *testing.T) {
	tests := []test{
		{"let a = 1; a = 2; a", 2},
		{"let a = 1; fn() { fn() { a = 2 }() }(); a", 2},
		{`
		let a = 10
		let b = 10
		let c = 10
		let f = fn(a) {
			let b = 20
		
			return a + b + c
		}
		f(20) + a + b;`, 70},
	}

	autoTest(t, tests)
}

func TestClosures(t *testing.T) {
	tests := []test{
		{`let newAdder = fn(x) { return fn(y) { x + y }; }; let addTwo = newAdder(2); addTwo(2);`, 4},
	}

	autoTest(t, tests)
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []test{
		{`let a = []; append(a, 1, 2, 3); a`, []any{1, 2, 3}},
		{`let a = []; let b = append(a, 1, 2, 3); b`, []any{1, 2, 3}},
		{`let a = [1, 2, 3, 4]; pop(a); a`, []any{1, 2, 3}},
		{`let a = [1, 2, 3, 4]; reverse(a); a`, []any{4, 3, 2, 1}},
		{`let a = [1, 3, 2, 4, 6, 5, 7]; sort(a); a`, []any{1, 2, 3, 4, 5, 6, 7}},
		{`len("")`, 0},
		{`len([])`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len([1, 2, 3, 4, 5, 6, 7])`, 7},
		{`sum([])`, 0},
		{`sum([], 2)`, 2},
		{`sum([1, 2, 3, 4, 5, 6, 7])`, 28},
		{`sum([1, 2, 3, 4, 5, 6, 7], 8)`, 36},
		{`sum(["1", "23", "456"])`, "123456"},
	}

	autoTest(t, tests)
}

func TestArrayLiterals(t *testing.T) {
	tests := []test{
		{"[1, 2 * 2, 3 + 3]", []any{1, 4, 6}},
		{`[1, 2, false, "hello", [3, 4, true]]`, []any{1, 2, false, "hello", []any{3, 4, true}}},
		{`[1, 2, [3, 4, [5, 6, 7]]]`, []any{1, 2, []any{3, 4, []any{5, 6, 7}}}},
		{"[1, 2, 3] + [4, 5, 6]", []any{1, 2, 3, 4, 5, 6}},
		{"[1, 2, 3] + 4", []any{1, 2, 3, 4}},
		{`[1, 2, 3] + "hello"`, []any{1, 2, 3, "hello"}},
		{`4 + [1, 2, 3]`, []any{4, 1, 2, 3}},
		{`[1, 2] * 0`, []any{}},
		{`[1, 2] * 4`, []any{1, 2, 1, 2, 1, 2, 1, 2}},
	}

	autoTest(t, tests)
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []test{
		{"[1, 2, 3][0]", 1},
		{"[1, 2, 3][1]", 2},
		{"[1, 2, 3][2]", 3},
		{"let i = 0; [1][i];", 1},
		{"[1, 2, 3][1 + 1];", 3},
		{"let a = [1, 2, 3]; a[2];", 3},
		{"let a = [1, 2, 3]; a[0] + a[1] + a[2];", 6},
		{"let a = [1, 2, 3]; let i = a[0]; a[i]", 2},
	}

	autoTest(t, tests)
}

func autoTest(t *testing.T, tests []test) {
	for _, tt := range tests {
		res, err := testEval(tt.input)
		if err != nil {
			testObject(t, object.NewString(err.Error()), tt.expected)
		} else if res != nil {
			testObject(t, res, tt.expected)
		} else if tt.expected != nil {
			t.Error("missing return value")
		}
	}
}

func testEval(input string) (object.Object, error) {
	p := parser.NewParser(lexer.NewLexer(input))
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}

func testObject(t *testing.T, obj object.Object, expected any) {
	switch e := expected.(type) {
	case int:
		testIntegerObject(t, obj, int64(e))
	case int64:
		testIntegerObject(t, obj, e)
	case string:
		testStringObject(t, obj, e)
	case bool:
		testBooleanObject(t, obj, e)
	case []any:
		testArrayObject(t, obj, e)
	case nil:
		testNullObject(t, obj)
	default:
		t.Errorf("missing type")
	}
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

func testArrayObject(t *testing.T, obj object.Object, expected []any) {
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
		testObject(t, o, expected[i])
	}
}

func testNullObject(t *testing.T, obj object.Object) {
	if obj != object.NULL {
		t.Errorf("object is not NULL. got=%T (%v)", obj, obj)
	}
}
