package main

import (
	"errors"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"

	. "github.com/shellyln/takenoco/base"
	"github.com/shellyln/takenoco/extra"
	objparser "github.com/shellyln/takenoco/object"
	. "github.com/shellyln/takenoco/string"
)

var rootParser ParserFn

func init() {
	rootParser = program()
}

// Remove the resulting AST.
func erase(fn ParserFn) ParserFn {
	return Trans(fn, Erase)
}

// Whitespaces
func sp0() ParserFn {
	return erase(ZeroOrMoreTimes(Whitespace()))
}

func number() ParserFn {
	return Trans(
		FlatGroup(
			First(
				SeqI("pi"),
				extra.FloatNumberStr(),
				extra.IntegerNumberStr(),
			),
			WordBoundary(),
			erase(sp0()),
		),
		func(ctx ParserContext, asts AstSlice) (AstSlice, error) {
			var v float64
			if strings.ToLower(asts[0].Value.(string)) == "pi" {
				v = math.Pi
				asts = AstSlice{{
					Type:      AstType_Float,
					ClassName: "Number",
					Value:     v,
				}}
			} else {
				v, err := strconv.ParseFloat(asts[0].Value.(string), 64)
				if err != nil {
					return nil, err
				}
				asts = AstSlice{{
					Type:      AstType_Float,
					ClassName: "Number",
					Value:     v,
				}}
			}
			return asts, nil
		},
	)
}

func symbol() ParserFn {
	return Trans(
		FlatGroup(
			First(
				SeqI("sqrt"),
			),
			erase(CharClass("(")),
			erase(sp0()),
			number(),
			erase(sp0()),
			erase(CharClass(")")),
		),
		func(ctx ParserContext, asts AstSlice) (AstSlice, error) {
			v := asts[1].Value.(float64)
			var result float64
			switch asts[0].Value {
			case "sqrt":
				result = math.Sqrt(v)
			}

			asts = AstSlice{{
				Type:      AstType_Float,
				ClassName: "Number",
				Value:     result,
			}}
			return asts, nil
		},
	)
}

// Unary operators
func unaryOperator() ParserFn {
	return Trans(
		FlatGroup(
			CharClass("-"),
			erase(sp0()),
		),
		ChangeClassName("UnaryOperator"),
	)
}

// Binary operators
func binaryOperator() ParserFn {
	return Trans(
		FlatGroup(
			CharClass("+", "-", "*", "/"),
			erase(sp0()),
		),
		ChangeClassName("BinaryOperator"),
	)
}

// Expression without parentheses
func simpleExpression() ParserFn {
	return FlatGroup(
		number(),
		ZeroOrMoreTimes(
			binaryOperator(),
			number(),
		),
	)
}

// Expression without parentheses
func functionExpression() ParserFn {
	return FlatGroup(
		symbol(),
	)
}

// Expression enclosed in parentheses
func groupedExpression() ParserFn {
	return FlatGroup(
		erase(CharClass("(")),
		First(
			FlatGroup(
				erase(sp0()),
				expression(),
				erase(CharClass(")")),
				erase(sp0()),
			),
			Error("Error in grouped expression"),
		),
	)
}

// Expression before applying production rules
func expressionInner() ParserFn {
	return FlatGroup(
		ZeroOrMoreTimes(unaryOperator()),
		First(
			functionExpression(),
			simpleExpression(),
			Indirect(groupedExpression),
			Error("Value required"),
		),
		ZeroOrMoreTimes(
			binaryOperator(),
			First(
				FlatGroup(
					ZeroOrMoreTimes(unaryOperator()),
					First(
						simpleExpression(),
						Indirect(groupedExpression),
					),
				),
				Error("Error in the expression after the binary operator"),
			),
		),
	)
}

// Single expression
func expression() ParserFn {
	return Trans(
		expressionInner(),
		formulaProductionRules(),
	)
}

// Entire program
func program() ParserFn {
	return FlatGroup(
		Start(),
		erase(sp0()),
		expression(),
		End(),
	)
}

func Parse(s string) (float64, error) {
	out, err := rootParser(*NewStringParserContext(s))
	if err != nil {
		pos := GetLineAndColPosition(s, out.SourcePosition, 4)
		return 0, errors.New(
			err.Error() +
				"\n --> Line " + strconv.Itoa(pos.Line) +
				", Col " + strconv.Itoa(pos.Col) + "\n" +
				pos.ErrSource)
	}

	if out.MatchStatus == MatchStatus_Matched {
		return out.AstStack[0].Value.(float64), nil
	} else {
		pos := GetLineAndColPosition(s, out.SourcePosition, 4)
		return 0, errors.New(
			"Parse failed" +
				"\n --> Line " + strconv.Itoa(pos.Line) +
				", Col " + strconv.Itoa(pos.Col) + "\n" +
				pos.ErrSource + "\n" + out.MatchStatus.String())
	}
}

// Production rule (Precedence = 3)
var expressionRuleUnaryOp = Precedence{
	Rules: []ParserFn{
		Trans(
			FlatGroup(
				isOperator("UnaryOperator", []string{"-"}),
				anyOperand(),
			),
			func(ctx ParserContext, asts AstSlice) (AstSlice, error) {
				opcode := asts[0].Value.(string)
				op1 := asts[1].Value.(float64)

				var v float64
				switch opcode {
				case "-":
					v = -op1
				}

				return AstSlice{{
					ClassName: "Number",
					Value:     v,
				}}, nil
			},
		),
	},
	Rtol: true,
}

// Production rule (Precedence = 2)
var expressionRuleMulDiv = Precedence{
	Rules: []ParserFn{
		Trans(
			FlatGroup(
				anyOperand(),
				isOperator("BinaryOperator", []string{"*", "/"}),
				anyOperand(),
			),
			func(ctx ParserContext, asts AstSlice) (AstSlice, error) {
				opcode := asts[1].Value.(string)
				op1 := asts[0].Value.(float64)
				op2 := asts[2].Value.(float64)

				var v float64
				switch opcode {
				case "*":
					v = op1 * op2
				case "/":
					v = op1 / op2
				}

				return AstSlice{{
					ClassName: "Number",
					Value:     v,
				}}, nil
			},
		),
	},
}

// Production rule (Precedence = 1)
var expressionRulePlusMinus = Precedence{
	Rules: []ParserFn{
		Trans(
			FlatGroup(
				anyOperand(),
				isOperator("BinaryOperator", []string{"+", "-"}),
				anyOperand(),
			),
			func(ctx ParserContext, asts AstSlice) (AstSlice, error) {
				opcode := asts[1].Value.(string)
				op1 := asts[0].Value.(float64)
				op2 := asts[2].Value.(float64)

				var v float64
				switch opcode {
				case "+":
					v = op1 + op2
				case "-":
					v = op1 - op2
				}

				return AstSlice{{
					ClassName: "Number",
					Value:     v,
				}}, nil
			},
		),
	},
}

// Production rules
var precedences = []Precedence{
	expressionRuleUnaryOp,
	expressionRuleMulDiv,
	expressionRulePlusMinus,
}

// Production rules
func formulaProductionRules() TransformerFn {
	return ProductionRule(
		precedences,
		FlatGroup(Start(), objparser.Any(), objparser.End()),
	)
}

func unwrapOperandItem(ctx ParserContext, asts AstSlice) (AstSlice, error) {
	return AstSlice{asts[0].Value.(Ast)}, nil
}

func makeOpMatcher(className string, ops []string) func(c interface{}) bool {
	return func(c interface{}) bool {
		ast, ok := c.(Ast)
		if !ok || ast.ClassName != className {
			return false
		}
		val := ast.Value.(string)
		for _, op := range ops {
			if op == val {
				return true
			}
		}
		return false
	}
}

// An assertion that matches all single tokens
func anyOperand() ParserFn {
	return Trans(objparser.Any(), unwrapOperandItem)
}

// An assertion matching a single token that matches the class name
func isOperator(className string, ops []string) ParserFn {
	return Trans(objparser.ObjClassFn(makeOpMatcher(className, ops)), unwrapOperandItem)
}

func main() {
	testCases := []string{
		//"pi",
		//"1",
		//"1 + 2 + 3",
		//"(123 + 456 ) + pi",
		//"10 + (100 + 1)",
		//"((1 + 2) + (3 + 4)) + 5 + 6",
		//"6.6",
		//"((1.1 + 2.2) + (3.3 + 4.4 )) + 5.5 + 6.6",
		"sqrt(100)",
	}
	for _, input := range testCases {
		data, err := Parse(input)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("input: %v, result: %v\n", input, data)
	}
}
