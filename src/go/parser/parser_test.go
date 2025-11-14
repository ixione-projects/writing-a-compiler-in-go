package parser

import (
	"fmt"
	"testing"

	"github.com/ixione-projects/writing-a-compiler-in-go/src/go/ast"
)

type ParserTest struct {
	input   string
	errors  []string
	debug   bool
	program ProgramTest
}

type NodeTest interface {
	node()
}

func (p ProgramTest) node()                {}
func (ls LetDeclarationTest) node()        {}
func (rs ReturnStatementTest) node()       {}
func (es ExpressionStatementTest) node()   {}
func (bs BlockStatementTest) node()        {}
func (ue UnaryExpressionTest) node()       {}
func (be BinaryExpressionTest) node()      {}
func (ce ConditionalExpressionTest) node() {}
func (fl FunctionLiteralTest) node()       {}
func (ce CallExpressionTest) node()        {}
func (ae AssignmentExpressionTest) node()  {}
func (se SubscriptExpressionTest) node()   {}
func (i IdentifierTest) node()             {}
func (nl NumberLiteralTest) node()         {}
func (bl BooleanLiteralTest) node()        {}
func (sl StringLiteralTest) node()         {}
func (al ArrayLiteralTest) node()          {}
func (hl HashLiteralTest) node()           {}

type ProgramTest struct {
	Statements []StatementTest
}

type StatementTest interface {
	NodeTest
	statementNode()
}

func (ls LetDeclarationTest) statementNode()      {}
func (rs ReturnStatementTest) statementNode()     {}
func (es ExpressionStatementTest) statementNode() {}
func (bs BlockStatementTest) statementNode()      {}

type LetDeclarationTest struct {
	Name   IdentifierTest
	Value  ExpressionTest
	String string
}

type ReturnStatementTest struct {
	ReturnValue ExpressionTest
	String      string
}

type ExpressionStatementTest struct {
	Expression ExpressionTest
	String     string
}

type BlockStatementTest struct {
	Statements []StatementTest
}

type ExpressionTest interface {
	NodeTest
	expressionNode()
}

func (ue UnaryExpressionTest) expressionNode()       {}
func (be BinaryExpressionTest) expressionNode()      {}
func (ce ConditionalExpressionTest) expressionNode() {}
func (fl FunctionLiteralTest) expressionNode()       {}
func (ce CallExpressionTest) expressionNode()        {}
func (ae AssignmentExpressionTest) expressionNode()  {}
func (se SubscriptExpressionTest) expressionNode()   {}
func (i IdentifierTest) expressionNode()             {}
func (nl NumberLiteralTest) expressionNode()         {}
func (nl BooleanLiteralTest) expressionNode()        {}
func (sl StringLiteralTest) expressionNode()         {}
func (al ArrayLiteralTest) expressionNode()          {}
func (hl HashLiteralTest) expressionNode()           {}

type UnaryExpressionTest struct {
	Operator string
	Right    ExpressionTest
}

type BinaryExpressionTest struct {
	Left     ExpressionTest
	Operator string
	Right    ExpressionTest
}

type ConditionalExpressionTest struct {
	Condition   ExpressionTest
	Consequence BlockStatementTest
	Alternative BlockStatementTest
}

type FunctionLiteralTest struct {
	Parameters []IdentifierTest
	Body       BlockStatementTest
}

type CallExpressionTest struct {
	Callee    ExpressionTest
	Arguments []ExpressionTest
}

type AssignmentExpressionTest struct {
	LValue ExpressionTest
	RValue ExpressionTest
}

type SubscriptExpressionTest struct {
	Base      ExpressionTest
	Subscript ExpressionTest
}

type (
	IdentifierTest     string
	NumberLiteralTest  float64
	BooleanLiteralTest bool
	StringLiteralTest  string
)

type ArrayLiteralTest struct {
	Elements []ExpressionTest
}

type HashLiteralTest struct {
	Keys  []ExpressionTest
	Pairs map[ExpressionTest]ExpressionTest
}

func TestLetStatement(t *testing.T) {
	tests := []ParserTest{
		{
			input: `
			let x = 5;
			let y = 10;
			let foobar = 838383;`,
			program: ProgramTest{
				[]StatementTest{
					LetDeclarationTest{
						IdentifierTest("x"), NumberLiteralTest(5),
						"let x=5;",
					},
					LetDeclarationTest{
						IdentifierTest("y"), NumberLiteralTest(10),
						"let y=10;",
					},
					LetDeclarationTest{
						IdentifierTest("foobar"), NumberLiteralTest(838383),
						"let foobar=838383;",
					},
				},
			},
		},
		{
			input: `
			let x 5;
			let y 10;
			let 838383;`,
			errors: []string{
				"expected next token to be <ASSIGN> but was <NUMBER>",
				"expected next token to be <ASSIGN> but was <NUMBER>",
				"expected next token to be <IDENT> but was <NUMBER>",
			},
		},
	}

	for i, test := range tests {
		testParser(t, i, test)
	}
}

func TestReturnStatement(t *testing.T) {
	tests := []ParserTest{
		{
			input: `
			return 5;
			return 10;
			return 993322;`,
			program: ProgramTest{
				[]StatementTest{
					ReturnStatementTest{
						NumberLiteralTest(5),
						"return 5;",
					},
					ReturnStatementTest{
						NumberLiteralTest(10),
						"return 10;",
					},
					ReturnStatementTest{
						NumberLiteralTest(993322),
						"return 993322;",
					},
				},
			},
		},
	}

	for i, test := range tests {
		testParser(t, i, test)
	}
}

func TestExpressionStatement(t *testing.T) {
	suites := []struct {
		name  string
		tests []ParserTest
	}{
		{
			name: "TestIdentifier",
			tests: []ParserTest{
				{
					input: `foobar;`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								IdentifierTest("foobar"),
								"foobar;",
							},
						},
					},
				},
			},
		},
		{
			name: "TestNumberLiteral",
			tests: []ParserTest{
				{
					input: `5;`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								NumberLiteralTest(5),
								"5;",
							},
						},
					},
				},
			},
		},
		{
			name: "TestPrefixExpression",
			tests: []ParserTest{
				{
					input: `!5;`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								UnaryExpressionTest{"!", NumberLiteralTest(5)},
								"(!5);",
							},
						},
					},
				},
				{
					input: `-15;`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								UnaryExpressionTest{"-", NumberLiteralTest(15)},
								"(-15);",
							},
						},
					},
				},
				{
					input: `!foobar;`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								UnaryExpressionTest{"!", IdentifierTest("foobar")},
								"(!foobar);",
							},
						},
					},
				},
				{
					input: `-foobar;`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								UnaryExpressionTest{"-", IdentifierTest("foobar")},
								"(-foobar);",
							},
						},
					},
				},
				{
					input: `!true;`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								UnaryExpressionTest{"!", BooleanLiteralTest(true)},
								"(!true);",
							},
						},
					},
				},
				{
					input: `!false;`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								UnaryExpressionTest{"!", BooleanLiteralTest(false)},
								"(!false);",
							},
						},
					},
				},
				{
					input: `!-a;`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								UnaryExpressionTest{
									"!",
									UnaryExpressionTest{"-", IdentifierTest("a")},
								},
								"(!(-a));",
							},
						},
					},
				},
				{
					input: `-(5 + 5)`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								UnaryExpressionTest{
									"-",
									BinaryExpressionTest{
										NumberLiteralTest(5),
										"+",
										NumberLiteralTest(5),
									},
								},
								"(-(5+5));",
							},
						},
					},
				},
				{
					input: `!(true == true)`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								UnaryExpressionTest{
									"!",
									BinaryExpressionTest{
										BooleanLiteralTest(true),
										"==",
										BooleanLiteralTest(true),
									},
								},
								"(!(true==true));",
							},
						},
					},
				},
			},
		},
		{
			name: "TestInfixExpression",
			tests: []ParserTest{
				{
					input: `5 + 5;`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								BinaryExpressionTest{NumberLiteralTest(5), "+", NumberLiteralTest(5)},
								"(5+5);",
							},
						},
					},
				},
				{
					input: `5 - 5;`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								BinaryExpressionTest{NumberLiteralTest(5), "-", NumberLiteralTest(5)},
								"(5-5);",
							},
						},
					},
				},
				{
					input: `5 * 5;`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								BinaryExpressionTest{NumberLiteralTest(5), "*", NumberLiteralTest(5)},
								"(5*5);",
							},
						},
					},
				},
				{
					input: `5 / 5;`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								BinaryExpressionTest{NumberLiteralTest(5), "/", NumberLiteralTest(5)},
								"(5/5);",
							},
						},
					},
				},
				{
					input: `5 > 5;`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								BinaryExpressionTest{NumberLiteralTest(5), ">", NumberLiteralTest(5)},
								"(5>5);",
							},
						},
					},
				},
				{
					input: `5 < 5;`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								BinaryExpressionTest{NumberLiteralTest(5), "<", NumberLiteralTest(5)},
								"(5<5);",
							},
						},
					},
				},
				{
					input: `5 == 5;`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								BinaryExpressionTest{NumberLiteralTest(5), "==", NumberLiteralTest(5)},
								"(5==5);",
							},
						},
					},
				},
				{
					input: `5 != 5;`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								BinaryExpressionTest{NumberLiteralTest(5), "!=", NumberLiteralTest(5)},
								"(5!=5);",
							},
						},
					},
				},
				{
					input: `foobar + barfoo;`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								BinaryExpressionTest{IdentifierTest("foobar"), "+", IdentifierTest("barfoo")},
								"(foobar+barfoo);",
							},
						},
					},
				},
				{
					input: `foobar - barfoo;`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								BinaryExpressionTest{IdentifierTest("foobar"), "-", IdentifierTest("barfoo")},
								"(foobar-barfoo);",
							},
						},
					},
				},
				{
					input: `foobar * barfoo;`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								BinaryExpressionTest{IdentifierTest("foobar"), "*", IdentifierTest("barfoo")},
								"(foobar*barfoo);",
							},
						},
					},
				},
				{
					input: `foobar / barfoo;`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								BinaryExpressionTest{IdentifierTest("foobar"), "/", IdentifierTest("barfoo")},
								"(foobar/barfoo);",
							},
						},
					},
				},
				{
					input: `foobar > barfoo;`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								BinaryExpressionTest{IdentifierTest("foobar"), ">", IdentifierTest("barfoo")},
								"(foobar>barfoo);",
							},
						},
					},
				},
				{
					input: `foobar < barfoo;`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								BinaryExpressionTest{IdentifierTest("foobar"), "<", IdentifierTest("barfoo")},
								"(foobar<barfoo);",
							},
						},
					},
				},
				{
					input: `foobar == barfoo;`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								BinaryExpressionTest{IdentifierTest("foobar"), "==", IdentifierTest("barfoo")},
								"(foobar==barfoo);",
							},
						},
					},
				},
				{
					input: `foobar != barfoo;`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								BinaryExpressionTest{IdentifierTest("foobar"), "!=", IdentifierTest("barfoo")},
								"(foobar!=barfoo);",
							},
						},
					},
				},
				{
					input: `true == true`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								BinaryExpressionTest{BooleanLiteralTest(true), "==", BooleanLiteralTest(true)},
								"(true==true);",
							},
						},
					},
				},
				{
					input: `true != false`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								BinaryExpressionTest{BooleanLiteralTest(true), "!=", BooleanLiteralTest(false)},
								"(true!=false);",
							},
						},
					},
				},
				{
					input: `false == false`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								BinaryExpressionTest{BooleanLiteralTest(false), "==", BooleanLiteralTest(false)},
								"(false==false);",
							},
						},
					},
				},
				{
					input: `-a * b;`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								BinaryExpressionTest{
									UnaryExpressionTest{"-", IdentifierTest("a")},
									"*",
									IdentifierTest("b"),
								},
								"((-a)*b);",
							},
						},
					},
				},
				{
					input: `a + b + c;`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								BinaryExpressionTest{
									BinaryExpressionTest{IdentifierTest("a"), "+", IdentifierTest("b")},
									"+",
									IdentifierTest("c"),
								},
								"((a+b)+c);",
							},
						},
					},
				},
				{
					input: `a + b - c;`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								BinaryExpressionTest{
									BinaryExpressionTest{IdentifierTest("a"), "+", IdentifierTest("b")},
									"-",
									IdentifierTest("c"),
								},
								"((a+b)-c);",
							},
						},
					},
				},
				{
					input: `a * b * c;`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								BinaryExpressionTest{
									BinaryExpressionTest{IdentifierTest("a"), "*", IdentifierTest("b")},
									"*",
									IdentifierTest("c"),
								},
								"((a*b)*c);",
							},
						},
					},
				},
				{
					input: `a * b / c;`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								BinaryExpressionTest{
									BinaryExpressionTest{IdentifierTest("a"), "*", IdentifierTest("b")},
									"/",
									IdentifierTest("c"),
								},
								"((a*b)/c);",
							},
						},
					},
				},
				{
					input: `a + b / c;`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								BinaryExpressionTest{
									IdentifierTest("a"),
									"+",
									BinaryExpressionTest{IdentifierTest("b"), "/", IdentifierTest("c")},
								},
								"(a+(b/c));",
							},
						},
					},
				},
				{
					input: `a + b * c + d / e - f;`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								BinaryExpressionTest{
									BinaryExpressionTest{
										BinaryExpressionTest{
											IdentifierTest("a"),
											"+",
											BinaryExpressionTest{IdentifierTest("b"), "*", IdentifierTest("c")},
										},
										"+",
										BinaryExpressionTest{IdentifierTest("d"), "/", IdentifierTest("e")},
									},
									"-",
									IdentifierTest("f"),
								},
								"(((a+(b*c))+(d/e))-f);",
							},
						},
					},
				},
				{
					input: `3 + 4; -5 * 5`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								BinaryExpressionTest{
									NumberLiteralTest(3),
									"+",
									NumberLiteralTest(4),
								},
								"(3+4);",
							},
							ExpressionStatementTest{
								BinaryExpressionTest{
									UnaryExpressionTest{"-", NumberLiteralTest(5)},
									"*",
									NumberLiteralTest(5),
								},
								"((-5)*5);",
							},
						},
					},
				},
				{
					input: `5 > 4 == 3 < 4`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								BinaryExpressionTest{
									BinaryExpressionTest{NumberLiteralTest(5), ">", NumberLiteralTest(4)},
									"==",
									BinaryExpressionTest{NumberLiteralTest(3), "<", NumberLiteralTest(4)},
								},
								"((5>4)==(3<4));",
							},
						},
					},
				},
				{
					input: `5 < 4 != 3 > 4`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								BinaryExpressionTest{
									BinaryExpressionTest{NumberLiteralTest(5), "<", NumberLiteralTest(4)},
									"!=",
									BinaryExpressionTest{NumberLiteralTest(3), ">", NumberLiteralTest(4)},
								},
								"((5<4)!=(3>4));",
							},
						},
					},
				},
				{
					input: `3 + 4 * 5 == 3 * 1 + 4 * 5`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								BinaryExpressionTest{
									BinaryExpressionTest{
										NumberLiteralTest(3),
										"+",
										BinaryExpressionTest{NumberLiteralTest(4), "*", NumberLiteralTest(5)},
									},
									"==",
									BinaryExpressionTest{
										BinaryExpressionTest{NumberLiteralTest(3), "*", NumberLiteralTest(1)},
										"+",
										BinaryExpressionTest{NumberLiteralTest(4), "*", NumberLiteralTest(5)},
									},
								},
								"((3+(4*5))==((3*1)+(4*5)));",
							},
						},
					},
				},
				{
					input: `3 > 5 == false`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								BinaryExpressionTest{
									BinaryExpressionTest{NumberLiteralTest(3), ">", NumberLiteralTest(5)},
									"==",
									BooleanLiteralTest(false),
								},
								"((3>5)==false);",
							},
						},
					},
				},
				{
					input: `3 < 5 == true`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								BinaryExpressionTest{
									BinaryExpressionTest{NumberLiteralTest(3), "<", NumberLiteralTest(5)},
									"==",
									BooleanLiteralTest(true),
								},
								"((3<5)==true);",
							},
						},
					},
				},
				{
					input: `1 + (2 + 3) + 4`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								BinaryExpressionTest{
									BinaryExpressionTest{
										NumberLiteralTest(1),
										"+",
										BinaryExpressionTest{NumberLiteralTest(2), "+", NumberLiteralTest(3)},
									},
									"+",
									NumberLiteralTest(4),
								},
								"((1+(2+3))+4);",
							},
						},
					},
				},
				{
					input: `(5 + 5) * 2`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								BinaryExpressionTest{
									BinaryExpressionTest{
										NumberLiteralTest(5),
										"+",
										NumberLiteralTest(5),
									},
									"*",
									NumberLiteralTest(2),
								},
								"((5+5)*2);",
							},
						},
					},
				},
				{
					input: `2 / (5 + 5)`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								BinaryExpressionTest{
									NumberLiteralTest(2),
									"/",
									BinaryExpressionTest{
										NumberLiteralTest(5),
										"+",
										NumberLiteralTest(5),
									},
								},
								"(2/(5+5));",
							},
						},
					},
				},
				{
					input: `a + add(b * c) + d`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								BinaryExpressionTest{
									BinaryExpressionTest{
										IdentifierTest("a"),
										"+",
										CallExpressionTest{
											IdentifierTest("add"),
											[]ExpressionTest{
												BinaryExpressionTest{IdentifierTest("b"), "*", IdentifierTest("c")},
											},
										},
									},
									"+",
									IdentifierTest("d"),
								},
								"((a+add((b*c)))+d);",
							},
						},
					},
				},
				{
					input: `a * [1, 2, 3, 4][b * c] * d`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								BinaryExpressionTest{
									BinaryExpressionTest{
										IdentifierTest("a"),
										"*",
										SubscriptExpressionTest{
											ArrayLiteralTest{
												[]ExpressionTest{
													NumberLiteralTest(1),
													NumberLiteralTest(2),
													NumberLiteralTest(3),
													NumberLiteralTest(4),
												},
											},
											BinaryExpressionTest{
												IdentifierTest("b"),
												"*",
												IdentifierTest("c"),
											},
										},
									},
									"*",
									IdentifierTest("d"),
								},
								"((a*([1,2,3,4][(b*c)]))*d);",
							},
						},
					},
				},
			},
		},
		{
			name: "TestBooleanLiteral",
			tests: []ParserTest{
				{
					input: `true;`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								BooleanLiteralTest(true),
								"true;",
							},
						},
					},
				},
				{
					input: `false;`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								BooleanLiteralTest(false),
								"false;",
							},
						},
					},
				},
			},
		},
		{
			name: "TestConditionalExpression",
			tests: []ParserTest{
				{
					input: `if (x < y) { x }`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								ConditionalExpressionTest{
									BinaryExpressionTest{IdentifierTest("x"), "<", IdentifierTest("y")},
									BlockStatementTest{
										[]StatementTest{
											ExpressionStatementTest{
												IdentifierTest("x"),
												"x;",
											},
										},
									},
									BlockStatementTest{},
								},
								"if (x<y) {x;};",
							},
						},
					},
				},
				{
					input: `if (x < y) { x } else { y }`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								ConditionalExpressionTest{
									BinaryExpressionTest{IdentifierTest("x"), "<", IdentifierTest("y")},
									BlockStatementTest{
										[]StatementTest{
											ExpressionStatementTest{
												IdentifierTest("x"),
												"x;",
											},
										},
									},
									BlockStatementTest{
										[]StatementTest{
											ExpressionStatementTest{
												IdentifierTest("y"),
												"y;",
											},
										},
									},
								},
								"if (x<y) {x;} else {y;};",
							},
						},
					},
				},
			},
		},
		{
			name: "TestFunctionLiteral",
			tests: []ParserTest{
				{
					input: `fn(x, y) { x + y; }`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								FunctionLiteralTest{
									[]IdentifierTest{
										IdentifierTest("x"),
										IdentifierTest("y"),
									},
									BlockStatementTest{
										[]StatementTest{
											ExpressionStatementTest{
												BinaryExpressionTest{IdentifierTest("x"), "+", IdentifierTest("y")},
												"(x+y);",
											},
										},
									},
								},
								"fn(x,y){(x+y);};",
							},
						},
					},
				},
				{
					input: `fn() {};`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								FunctionLiteralTest{
									[]IdentifierTest{},
									BlockStatementTest{
										[]StatementTest{},
									},
								},
								"fn(){};",
							},
						},
					},
				},
				{
					input: `fn(x) {};`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								FunctionLiteralTest{
									[]IdentifierTest{
										IdentifierTest("x"),
									},
									BlockStatementTest{
										[]StatementTest{},
									},
								},
								"fn(x){};",
							},
						},
					},
				},
				{
					input: `fn(x, y, z) {};`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								FunctionLiteralTest{
									[]IdentifierTest{
										IdentifierTest("x"),
										IdentifierTest("y"),
										IdentifierTest("z"),
									},
									BlockStatementTest{
										[]StatementTest{},
									},
								},
								"fn(x,y,z){};",
							},
						},
					},
				},
			},
		},
		{
			name: "TestCallExpression",
			tests: []ParserTest{
				{
					input: `add();`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								CallExpressionTest{
									IdentifierTest("add"),
									[]ExpressionTest{},
								},
								"add();",
							},
						},
					},
				},
				{
					input: `add(1);`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								CallExpressionTest{
									IdentifierTest("add"),
									[]ExpressionTest{
										NumberLiteralTest(1),
									},
								},
								"add(1);",
							},
						},
					},
				},
				{
					input: `add(1, 2 * 3, 4 + 5);`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								CallExpressionTest{
									IdentifierTest("add"),
									[]ExpressionTest{
										NumberLiteralTest(1),
										BinaryExpressionTest{NumberLiteralTest(2), "*", NumberLiteralTest(3)},
										BinaryExpressionTest{NumberLiteralTest(4), "+", NumberLiteralTest(5)},
									},
								},
								"add(1,(2*3),(4+5));",
							},
						},
					},
				},
				{
					input: `add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								CallExpressionTest{
									IdentifierTest("add"),
									[]ExpressionTest{
										IdentifierTest("a"),
										IdentifierTest("b"),
										NumberLiteralTest(1),
										BinaryExpressionTest{NumberLiteralTest(2), "*", NumberLiteralTest(3)},
										BinaryExpressionTest{NumberLiteralTest(4), "+", NumberLiteralTest(5)},
										CallExpressionTest{
											IdentifierTest("add"),
											[]ExpressionTest{
												NumberLiteralTest(6),
												BinaryExpressionTest{NumberLiteralTest(7), "*", NumberLiteralTest(8)},
											},
										},
									},
								},
								"add(a,b,1,(2*3),(4+5),add(6,(7*8)));",
							},
						},
					},
				},
				{
					input: `add(a + b + c * d / f + g);`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								CallExpressionTest{
									IdentifierTest("add"),
									[]ExpressionTest{
										BinaryExpressionTest{
											BinaryExpressionTest{
												BinaryExpressionTest{IdentifierTest("a"), "+", IdentifierTest("b")},
												"+",
												BinaryExpressionTest{
													BinaryExpressionTest{IdentifierTest("c"), "*", IdentifierTest("d")},
													"/",
													IdentifierTest("f"),
												},
											},
											"+",
											IdentifierTest("g"),
										},
									},
								},
								"add((((a+b)+((c*d)/f))+g));",
							},
						},
					},
				},
				{
					input: `add(a * b[2], b[1], 2 * [1, 2][1])`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								CallExpressionTest{
									IdentifierTest("add"),
									[]ExpressionTest{
										BinaryExpressionTest{
											IdentifierTest("a"),
											"*",
											SubscriptExpressionTest{
												IdentifierTest("b"),
												NumberLiteralTest(2),
											},
										},
										SubscriptExpressionTest{
											IdentifierTest("b"),
											NumberLiteralTest(1),
										},
										BinaryExpressionTest{
											NumberLiteralTest(2),
											"*",
											SubscriptExpressionTest{
												ArrayLiteralTest{
													[]ExpressionTest{
														NumberLiteralTest(1),
														NumberLiteralTest(2),
													},
												},
												NumberLiteralTest(1),
											},
										},
									},
								},
								"add((a*(b[2])),(b[1]),(2*([1,2][1])));",
							},
						},
					},
				},
			},
		},
		{
			name: "TestStringLiteral",
			tests: []ParserTest{
				{
					input: `"hello world";`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								StringLiteralTest("hello world"),
								"\"hello world\";",
							},
						},
					},
				},
			},
		},
		{
			name: "TestArrayLiteral",
			tests: []ParserTest{
				{
					input: `[1, 2 * 2, 3 + 3]`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								ArrayLiteralTest{
									[]ExpressionTest{
										NumberLiteralTest(1),
										BinaryExpressionTest{
											NumberLiteralTest(2),
											"*",
											NumberLiteralTest(2),
										},
										BinaryExpressionTest{
											NumberLiteralTest(3),
											"+",
											NumberLiteralTest(3),
										},
									},
								},
								"[1,(2*2),(3+3)];",
							},
						},
					},
				},
				{
					input: `[]`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								ArrayLiteralTest{
									[]ExpressionTest{},
								},
								"[];",
							},
						},
					},
				},
				{
					input: `array[1 + 1]`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								SubscriptExpressionTest{
									IdentifierTest("array"),
									BinaryExpressionTest{
										NumberLiteralTest(1),
										"+",
										NumberLiteralTest(1),
									},
								},
								"(array[(1+1)]);",
							},
						},
					},
				},
			},
		},
		{
			name: "TestHashLiteral",
			tests: []ParserTest{
				{
					input: `{"one": 1, "two": 2, "three": 3}`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								HashLiteralTest{
									[]ExpressionTest{
										StringLiteralTest("one"),
										StringLiteralTest("two"),
										StringLiteralTest("three"),
									},
									map[ExpressionTest]ExpressionTest{
										StringLiteralTest("one"):   NumberLiteralTest(1),
										StringLiteralTest("two"):   NumberLiteralTest(2),
										StringLiteralTest("three"): NumberLiteralTest(3),
									},
								},
								"{\"one\":1,\"two\":2,\"three\":3};",
							},
						},
					},
				},
				{
					input: `{};`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								HashLiteralTest{},
								"{};",
							},
						},
					},
				},
				{
					input: `{true: 1, false: 2}`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								HashLiteralTest{
									[]ExpressionTest{
										BooleanLiteralTest(true),
										BooleanLiteralTest(false),
									},
									map[ExpressionTest]ExpressionTest{
										BooleanLiteralTest(true):  NumberLiteralTest(1),
										BooleanLiteralTest(false): NumberLiteralTest(2),
									},
								},
								"{true:1,false:2};",
							},
						},
					},
				},
				{
					input: `{1: 1, 2: 2, 3: 3}`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								HashLiteralTest{
									[]ExpressionTest{
										NumberLiteralTest(1),
										NumberLiteralTest(2),
										NumberLiteralTest(3),
									},
									map[ExpressionTest]ExpressionTest{
										NumberLiteralTest(1): NumberLiteralTest(1),
										NumberLiteralTest(2): NumberLiteralTest(2),
										NumberLiteralTest(3): NumberLiteralTest(3),
									},
								},
								"{1:1,2:2,3:3};",
							},
						},
					},
				},
				{
					input: `{"one": 0 + 1, "two": 10 - 8, "three": 15 / 5}`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								HashLiteralTest{
									[]ExpressionTest{
										StringLiteralTest("one"),
										StringLiteralTest("two"),
										StringLiteralTest("three"),
									},
									map[ExpressionTest]ExpressionTest{
										StringLiteralTest("one"): BinaryExpressionTest{
											NumberLiteralTest(0),
											"+",
											NumberLiteralTest(1),
										},
										StringLiteralTest("two"): BinaryExpressionTest{
											NumberLiteralTest(10),
											"-",
											NumberLiteralTest(8),
										},
										StringLiteralTest("three"): BinaryExpressionTest{
											NumberLiteralTest(15),
											"/",
											NumberLiteralTest(5),
										},
									},
								},
								"{\"one\":(0+1),\"two\":(10-8),\"three\":(15/5)};",
							},
						},
					},
				},
			},
		},
	}

	for _, suite := range suites {
		t.Run(suite.name, func(t *testing.T) {
			for i, test := range suite.tests {
				testParser(t, i, test)
			}
		})
	}
}

func testParser(t *testing.T, i int, test ParserTest) {
	p := NewParser(test.input, test.debug)
	program := p.ParseProgram()

	if len(test.errors) != len(p.Errors()) {
		t.Errorf("test[%d] - len(p.Errors()) ==> expected: <%d> but was: <%d>", i, len(test.errors), len(p.Errors()))
		for j, msg := range p.Errors() {
			t.Errorf("--------- p.Errors()[%d]: %s", j, msg)
		}
		t.Fatalf("test[%d] - %s", i, test.input)
	}

	for j, msg := range test.errors {
		fail := false
		if msg != p.Errors()[j] {
			t.Errorf("test[%d] - p.Errors()[%d] ==> expected: <%s> but was: <%s>", i, j, msg, p.Errors()[j])
			fail = true
		}

		if fail {
			t.FailNow()
		}
	}

	if len(test.program.Statements) > 0 {
		testProgram(t, i, test.program, program)
	}
}

func testProgram(t *testing.T, i int, expected ProgramTest, actual *ast.Program) {
	if actual == nil {
		t.Fatalf("test[%d] - ParseProgram() ==> expected: not <%#v>", i, actual)
	}

	if len(expected.Statements) != len(actual.Statements) {
		t.Fatalf("test[%d] - len(program.Statements) ==> expected: <%d> but was: <%d>", i, len(expected.Statements), len(actual.Statements))
	}

	for j, stmt := range expected.Statements {
		if !testStatement(t, i, j, stmt, actual.Statements[j]) {
			t.Fatalf("test[%d][%d] - %s", i, j, actual.Statements[j].String())
		}
	}
}

func testStatement(t *testing.T, i, j int, expected StatementTest, actual ast.Statement) bool {
	switch expected := expected.(type) {
	case LetDeclarationTest:
		if !testLetDeclaration(t, i, j, expected, actual) {
			return false
		}
	case ReturnStatementTest:
		if !testReturnStatement(t, i, j, expected, actual) {
			return false
		}
	case ExpressionStatementTest:
		if !testExpressionStatement(t, i, j, expected, actual) {
			return false
		}
	case BlockStatementTest:
		if !testBlockStatement(t, i, j, expected, actual) {
			return false
		}
	default:
		t.Fatalf("test[%d][%d] - unexpected type <%T>", i, j, expected)
	}
	return true
}

func testLetDeclaration(t *testing.T, i, j int, expected LetDeclarationTest, actual ast.Statement) bool {
	if "let" != actual.TokenLiteral() {
		t.Errorf("test[%d][%d] - *ast.LetStatement.TokenLiteral() ==> expected: <%s> but was: <%s>", i, j, "let", actual.TokenLiteral())
		return false
	}

	stmt, ok := actual.(*ast.LetDeclaration)
	if !ok {
		t.Errorf("test[%d][%d] - actual.(*ast.LetStatement) ==> unexpected type, expected: <%T> but was: <%T>", i, j, &ast.LetDeclaration{}, actual)
		return false
	}

	if !testIdentifier(t, i, j, expected.Name, stmt.Name) {
		return false
	}

	if !testExpression(t, i, j, expected.Value, stmt.Value) {
		return false
	}

	if expected.String != actual.String() {
		t.Errorf("test[%d][%d] - *ast.LetStatement.String() ==> expected: <%s> but was: <%s>", i, j, expected.String, actual.String())
		return false
	}

	return true
}

func testReturnStatement(t *testing.T, i, j int, expected ReturnStatementTest, actual ast.Statement) bool {
	if "return" != actual.TokenLiteral() {
		t.Errorf("test[%d][%d] - *ast.ReturnStatement.TokenLiteral() ==> expected: <%s> but was: <%s>", i, j, "return", actual.TokenLiteral())
		return false
	}

	stmt, ok := actual.(*ast.ReturnStatement)
	if !ok {
		t.Errorf("test[%d][%d] - actual.(*ast.ReturnStatement) ==> unexpected type, expected: <%T> but was: <%T>", i, j, &ast.ReturnStatement{}, actual)
		return false
	}

	if !testExpression(t, i, j, expected.ReturnValue, stmt.ReturnValue) {
		return false
	}

	if expected.String != actual.String() {
		t.Errorf("test[%d][%d] - *ast.ReturnStatement.String() ==> expected: <%s> but was: <%s>", i, j, expected.String, actual.String())
		return false
	}

	return true
}

func testExpressionStatement(t *testing.T, i, j int, expected ExpressionStatementTest, actual ast.Statement) bool {
	stmt, ok := actual.(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("test[%d][%d] - actual.(*ast.ExpressionStatement) ==> unexpected type, expected: <%T> but was: <%T>", i, j, &ast.ExpressionStatement{}, actual)
		return false
	}

	if !testExpression(t, i, j, expected.Expression, stmt.Expression) {
		return false
	}

	if expected.String != actual.String() {
		t.Errorf("test[%d][%d] - *ast.ExpressionStatement.String() ==> expected: <%s> but was: <%s>", i, j, expected.String, actual.String())
		return false
	}

	return true
}

func testBlockStatement(t *testing.T, i, j int, expected BlockStatementTest, actual ast.Statement) bool {
	if "{" != actual.TokenLiteral() {
		t.Errorf("test[%d][%d] - *ast.BlockStatement.TokenLiteral() ==> expected: <%s> but was: <%s>", i, j, "{", actual.TokenLiteral())
		return false
	}

	block, ok := actual.(*ast.BlockStatement)
	if !ok {
		t.Errorf("test[%d][%d] - actual.(*ast.BlockStatement) ==> unexpected type, expected: <%T> but was: <%T>", i, j, &ast.BlockStatement{}, actual)
		return false
	}

	if len(expected.Statements) != len(block.Statements) {
		t.Errorf("test[%d] - len(block.Statements) ==> expected: <%d> but was: <%d>", i, len(expected.Statements), len(block.Statements))
		return false
	}

	for k, stmt := range expected.Statements {
		if !testStatement(t, i, j, stmt, block.Statements[k]) {
			return false
		}
	}

	return true
}

func testExpression(t *testing.T, i, j int, expected ExpressionTest, actual ast.Expression) bool {
	switch expected := expected.(type) {
	case UnaryExpressionTest:
		if !testUnaryExpression(t, i, j, expected, actual) {
			return false
		}
	case BinaryExpressionTest:
		if !testBinaryExpression(t, i, j, expected, actual) {
			return false
		}
	case ConditionalExpressionTest:
		if !testConditionalExpression(t, i, j, expected, actual) {
			return false
		}
	case FunctionLiteralTest:
		if !testFunctionLiteral(t, i, j, expected, actual) {
			return false
		}
	case CallExpressionTest:
		if !testCallExpression(t, i, j, expected, actual) {
			return false
		}
	case AssignmentExpressionTest:
		if !testAssignmentExpression(t, i, j, expected, actual) {
			return false
		}
	case SubscriptExpressionTest:
		if !testSubscriptExpression(t, i, j, expected, actual) {
			return false
		}
	case IdentifierTest:
		if !testIdentifier(t, i, j, expected, actual) {
			return false
		}
	case NumberLiteralTest:
		if !testNumberLiteral(t, i, j, expected, actual) {
			return false
		}
	case BooleanLiteralTest:
		if !testBooleanLiteral(t, i, j, expected, actual) {
			return false
		}
	case StringLiteralTest:
		if !testStringLiteral(t, i, j, expected, actual) {
			return false
		}
	case ArrayLiteralTest:
		if !testArrayLiteral(t, i, j, expected, actual) {
			return false
		}
	case HashLiteralTest:
		if !testHashLiteral(t, i, j, expected, actual) {
			return false
		}
	default:
		t.Fatalf("test[%d][%d] - unexpected type <%T>", i, j, expected)
	}
	return true
}

func testUnaryExpression(t *testing.T, i, j int, expected UnaryExpressionTest, actual ast.Expression) bool {
	if expected.Operator != actual.TokenLiteral() {
		t.Errorf("test[%d][%d] - *ast.PrefixExpression.TokenLiteral() ==> expected: <%s> but was: <%s>", i, j, expected.Operator, actual.TokenLiteral())
		return false
	}

	expr, ok := actual.(*ast.UnaryExpression)
	if !ok {
		t.Errorf("test[%d][%d] - actual.(*ast.PrefixExpression) ==> unexpected type, expected: <%T> but was: <%T>", i, j, &ast.UnaryExpression{}, actual)
		return false
	}

	if !testExpression(t, i, j, expected.Right, expr.Right) {
		return false
	}

	return true
}

func testBinaryExpression(t *testing.T, i, j int, expected BinaryExpressionTest, actual ast.Expression) bool {
	if expected.Operator != actual.TokenLiteral() {
		t.Errorf("test[%d][%d] - *ast.InfixExpression.TokenLiteral() ==> expected: <%s> but was: <%s>", i, j, expected.Operator, actual.TokenLiteral())
		return false
	}

	expr, ok := actual.(*ast.BinaryExpression)
	if !ok {
		t.Errorf("test[%d][%d] - actual.(*ast.InfixExpression) ==> unexpected type, expected: <%T> but was: <%T>", i, j, &ast.BinaryExpression{}, actual)
		return false
	}

	if !testExpression(t, i, j, expected.Left, expr.Left) {
		return false
	}

	if !testExpression(t, i, j, expected.Right, expr.Right) {
		return false
	}

	return true
}

func testConditionalExpression(t *testing.T, i, j int, expected ConditionalExpressionTest, actual ast.Expression) bool {
	if "if" != actual.TokenLiteral() {
		t.Errorf("test[%d][%d] - *ast.IfExpression.TokenLiteral() ==> expected: <%s> but was: <%s>", i, j, "if", actual.TokenLiteral())
		return false
	}

	expr, ok := actual.(*ast.ConditionalExpression)
	if !ok {
		t.Errorf("test[%d][%d] - actual.(*ast.IfExpression) ==> unexpected type, expected: <%T> but was: <%T>", i, j, &ast.ConditionalExpression{}, actual)
		return false
	}

	if !testExpression(t, i, j, expected.Condition, expr.Condition) {
		return false
	}

	if !testBlockStatement(t, i, j, expected.Consequence, expr.Consequence) {
		return false
	}

	if len(expected.Alternative.Statements) > 0 && !testBlockStatement(t, i, j, expected.Alternative, expr.Alternative) {
		return false
	}

	return true
}

func testFunctionLiteral(t *testing.T, i, j int, expected FunctionLiteralTest, actual ast.Expression) bool {
	if "fn" != actual.TokenLiteral() {
		t.Errorf("test[%d][%d] - *ast.FunctionLiteral.TokenLiteral() ==> expected: <%s> but was: <%s>", i, j, "fn", actual.TokenLiteral())
		return false
	}

	fn, ok := actual.(*ast.FunctionLiteral)
	if !ok {
		t.Errorf("test[%d][%d] - actual.(*ast.FunctionLiteral) ==> unexpected type, expected: <%T> but was: <%T>", i, j, &ast.FunctionLiteral{}, actual)
		return false
	}

	if len(expected.Parameters) != len(fn.Parameters) {
		t.Errorf("test[%d] - len(fn.Parameters) ==> expected: <%d> but was: <%d>", i, len(expected.Parameters), len(fn.Parameters))
		return false
	}

	for k, parameter := range expected.Parameters {
		if !testIdentifier(t, i, j, parameter, fn.Parameters[k]) {
			return false
		}
	}

	if !testBlockStatement(t, i, j, expected.Body, fn.Body) {
		return false
	}

	return true
}

func testCallExpression(t *testing.T, i, j int, expected CallExpressionTest, actual ast.Expression) bool {
	if "(" != actual.TokenLiteral() {
		t.Errorf("test[%d][%d] - *ast.CallExpression.TokenLiteral() ==> expected: <%s> but was: <%s>", i, j, "(", actual.TokenLiteral())
		return false
	}

	call, ok := actual.(*ast.CallExpression)
	if !ok {
		t.Errorf("test[%d][%d] - actual.(*ast.CallExpression) ==> unexpected type, expected: <%T> but was: <%T>", i, j, &ast.CallExpression{}, actual)
		return false
	}

	if !testExpression(t, i, j, expected.Callee, call.Callee) {
		return false
	}

	if len(expected.Arguments) != len(call.Arguments) {
		t.Errorf("test[%d] - len(call.Arguments) ==> expected: <%d> but was: <%d>", i, len(expected.Arguments), len(call.Arguments))
		return false
	}

	for k, argument := range expected.Arguments {
		if !testExpression(t, i, j, argument, call.Arguments[k]) {
			return false
		}
	}

	return true
}

func testAssignmentExpression(t *testing.T, i, j int, expected AssignmentExpressionTest, actual ast.Expression) bool {
	if "=" != actual.TokenLiteral() {
		t.Errorf("test[%d][%d] - *ast.AssignmentExpression.TokenLiteral() ==> expected: <%s> but was: <%s>", i, j, "=", actual.TokenLiteral())
		return false
	}

	expr, ok := actual.(*ast.AssignmentExpression)
	if !ok {
		t.Errorf("test[%d][%d] - actual.(*ast.AssignmentExpression) ==> unexpected type, expected: <%T> but was: <%T>", i, j, &ast.AssignmentExpression{}, actual)
		return false
	}

	if !testExpression(t, i, j, expected.LValue, expr.LValue) {
		return false
	}

	if !testExpression(t, i, j, expected.RValue, expr.RValue) {
		return false
	}

	return true
}

func testSubscriptExpression(t *testing.T, i, j int, expected SubscriptExpressionTest, actual ast.Expression) bool {
	if "[" != actual.TokenLiteral() {
		t.Errorf("test[%d][%d] - *ast.SubscriptExpression.TokenLiteral() ==> expected: <%s> but was: <%s>", i, j, "[", actual.TokenLiteral())
		return false
	}

	expr, ok := actual.(*ast.SubscriptExpression)
	if !ok {
		t.Errorf("test[%d][%d] - actual.(*ast.SubscriptExpression) ==> unexpected type, expected: <%T> but was: <%T>", i, j, &ast.SubscriptExpression{}, actual)
		return false
	}

	if !testExpression(t, i, j, expected.Base, expr.Base) {
		return false
	}

	if !testExpression(t, i, j, expected.Subscript, expr.Subscript) {
		return false
	}

	return true
}

func testIdentifier(t *testing.T, i, j int, expected IdentifierTest, actual ast.Expression) bool {
	if string(expected) != actual.TokenLiteral() {
		t.Errorf("test[%d][%d] - *ast.Identifier.TokenLiteral() ==> expected: <%s> but was: <%s>", i, j, string(expected), actual.TokenLiteral())
		return false
	}

	expr, ok := actual.(*ast.Identifier)
	if !ok {
		t.Errorf("test[%d][%d] - actual.(*ast.Identifier) ==> unexpected type, expected: <%T> but was: <%T>", i, j, &ast.Identifier{}, actual)
		return false
	}

	if string(expected) != expr.Value {
		t.Errorf("test[%d][%d] - *ast.Identifier.Value ==> expected: <%s> but was: <%s>", i, j, string(expected), expr.Value)
		return false
	}

	return true
}

func testNumberLiteral(t *testing.T, i, j int, expected NumberLiteralTest, actual ast.Expression) bool {
	if fmt.Sprintf("%.0f", float64(expected)) != actual.TokenLiteral() {
		t.Errorf("test[%d][%d] - *ast.NumberLiteral.TokenLiteral() ==> expected: <%f> but was: <%s>", i, j, float64(expected), actual.TokenLiteral())
		return false
	}

	expr, ok := actual.(*ast.NumberLiteral)
	if !ok {
		t.Errorf("test[%d][%d] - actual.(*ast.NumberLiteral) ==> unexpected type, expected: <%T> but was: <%T>", i, j, &ast.NumberLiteral{}, actual)
		return false
	}

	if float64(expected) != expr.Value {
		t.Errorf("test[%d][%d] - *ast.NumberLiteral.Value ==> expected: <%f> but was: <%f>", i, j, float64(expected), expr.Value)
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, i, j int, expected BooleanLiteralTest, actual ast.Expression) bool {
	if fmt.Sprintf("%t", bool(expected)) != actual.TokenLiteral() {
		t.Errorf("test[%d][%d] - *ast.BooleanLiteral.TokenLiteral() ==> expected: <%t> but was: <%s>", i, j, bool(expected), actual.TokenLiteral())
		return false
	}

	expr, ok := actual.(*ast.BooleanLiteral)
	if !ok {
		t.Errorf("test[%d][%d] - actual.(*ast.BooleanLiteral) ==> unexpected type, expected: <%T> but was: <%T>", i, j, &ast.NumberLiteral{}, actual)
		return false
	}

	if bool(expected) != expr.Value {
		t.Errorf("test[%d][%d] - *ast.BooleanLiteral.Value ==> expected: <%t> but was: <%t>", i, j, bool(expected), expr.Value)
		return false
	}

	return true
}

func testStringLiteral(t *testing.T, i, j int, expected StringLiteralTest, actual ast.Expression) bool {
	if "\""+string(expected)+"\"" != actual.TokenLiteral() {
		t.Errorf("test[%d][%d] - *ast.StringLiteral.TokenLiteral() ==> expected: <%s> but was: <%s>", i, j, "\""+string(expected)+"\"", actual.TokenLiteral())
		return false
	}

	expr, ok := actual.(*ast.StringLiteral)
	if !ok {
		t.Errorf("test[%d][%d] - actual.(*ast.StringLiteral) ==> unexpected type, expected: <%T> but was: <%T>", i, j, &ast.StringLiteral{}, actual)
		return false
	}

	if string(expected) != expr.Value {
		t.Errorf("test[%d][%d] - *ast.StringLiteral.Value ==> expected: <%s> but was: <%s>", i, j, string(expected), expr.Value)
		return false
	}

	return true
}

func testArrayLiteral(t *testing.T, i, j int, expected ArrayLiteralTest, actual ast.Expression) bool {
	if "[" != actual.TokenLiteral() {
		t.Errorf("test[%d][%d] - *ast.ArrayLiteral.TokenLiteral() ==> expected: <%s> but was: <%s>", i, j, "[", actual.TokenLiteral())
		return false
	}

	array, ok := actual.(*ast.ArrayLiteral)
	if !ok {
		t.Errorf("test[%d][%d] - actual.(*ast.ArrayLiteral) ==> unexpected type, expected: <%T> but was: <%T>", i, j, &ast.ArrayLiteral{}, actual)
		return false
	}

	if len(expected.Elements) != len(array.Elements) {
		t.Errorf("test[%d][%d] - len(array.Elements) ==> expected: <%d> but was: <%d>", i, j, len(expected.Elements), len(array.Elements))
		return false
	}

	for k, element := range expected.Elements {
		if !testExpression(t, i, j, element, array.Elements[k]) {
			return false
		}
	}

	return true
}

func testHashLiteral(t *testing.T, i, j int, expected HashLiteralTest, actual ast.Expression) bool {
	if "{" != actual.TokenLiteral() {
		t.Errorf("test[%d][%d] - *ast.HashLiteral.TokenLiteral() ==> expected: <%s> but was: <%s>", i, j, "{", actual.TokenLiteral())
		return false
	}

	hash, ok := actual.(*ast.HashLiteral)
	if !ok {
		t.Errorf("test[%d][%d] - actual.(*ast.HashLiteral) ==> unexpected type, expected: <%T> but was: <%T>", i, j, &ast.ArrayLiteral{}, actual)
		return false
	}

	if len(expected.Pairs) != len(hash.Pairs) {
		t.Errorf("test[%d][%d] - len(hash.Pairs) ==> expected: <%d> but was: <%d>", i, j, len(expected.Pairs), len(hash.Pairs))
		return false
	}

	for k, key := range expected.Keys {
		if !testExpression(t, i, j, key, hash.Keys[k]) {
			return false
		}

		if !testExpression(t, i, j, expected.Pairs[key], hash.Pairs[hash.Keys[k]]) {
			return false
		}
	}

	return true
}
