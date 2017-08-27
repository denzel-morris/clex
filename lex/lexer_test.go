package lex

import (
	"strings"
	"testing"

	"github.com/denzel-morris/clex/lex/lexemes"
)

type fullMatchTestCase struct {
	input    string
	expected lexemes.Type
}

type partialMatchTestCase struct {
	input    string
	expected Lexeme
}

func TestLexerNext(t *testing.T) {
	for _, c := range fullMatchTestCases {
		lexer := makeLexer(c.input, &EmptyErrorPolicy{})
		lexeme, err := lexer.Next()
		if err != nil {
			t.Error("On case:", c, "got error", err)
		}

		if lexeme.Type != c.expected || lexeme.Value != c.input {
			t.Error("Expected", c.expected, "got", lexeme, "for", c.input)
		}
	}

	for _, c := range partialMatchTestCases {
		lexer := makeLexer(c.input, &EmptyErrorPolicy{})
		lexeme, err := lexer.Next()
		if err != nil {
			t.Error("On case:", c, "got error", err)
		}

		if lexeme.Type != c.expected.Type || lexeme.Value != c.expected.Value {
			t.Error("Expected", c.expected, "got", lexeme, "for", c.input)
		}
	}
}

func TestLexerLexicalErrors(t *testing.T) {
	policy := &CountingErrorPolicy{}
	for _, input := range errorTestCases {
		lexer := makeLexer(input, policy)
		lexeme, _ := lexer.Next()

		if lexeme.Type != lexemes.Invalid {
			t.Error("Expected Invalid lexeme, got", lexeme, "on", input)
		}

		if policy.count == 0 {
			t.Error("No lexical errors reported on", input)
		}
	}
}

func makeLexer(input string, policy ErrorPolicy) Lexer {
	return NewLexer(
		NewLineReader(NewReader(strings.NewReader(input))),
		policy,
	)
}

type EmptyErrorPolicy struct{}

func (ep EmptyErrorPolicy) ReportError(message string, line string, position Position) {}

type LogErrorPolicy struct {
	t *testing.T
}

func (ep LogErrorPolicy) ReportError(message string, line string, position Position) {
	ep.t.Log(message, line, position)
}

type CountingErrorPolicy struct {
	count int
}

func (ep *CountingErrorPolicy) ReportError(message string, line string, position Position) {
	ep.count++
}

var fullMatchTestCases = []fullMatchTestCase{
	{"a", lexemes.Identifier},
	{"B", lexemes.Identifier},
	{"Cd", lexemes.Identifier},
	{"E1f", lexemes.Identifier},
	{`\u12Ab`, lexemes.Identifier},
	{`\U5678CdEf`, lexemes.Identifier},
	{`\u90AB4gh5`, lexemes.Identifier},
	{`\U1234567890ij6`, lexemes.Identifier},
	{"u8", lexemes.Identifier},
	{"u8ab", lexemes.Identifier},
	{"auto", lexemes.Keyword},
	{"break", lexemes.Keyword},
	{"case", lexemes.Keyword},
	{"char", lexemes.Keyword},
	{"const", lexemes.Keyword},
	{"continue", lexemes.Keyword},
	{"default", lexemes.Keyword},
	{"do", lexemes.Keyword},
	{"double", lexemes.Keyword},
	{"else", lexemes.Keyword},
	{"enum", lexemes.Keyword},
	{"extern", lexemes.Keyword},
	{"float", lexemes.Keyword},
	{"for", lexemes.Keyword},
	{"goto", lexemes.Keyword},
	{"if", lexemes.Keyword},
	{"inline", lexemes.Keyword},
	{"int", lexemes.Keyword},
	{"long", lexemes.Keyword},
	{"register", lexemes.Keyword},
	{"restrict", lexemes.Keyword},
	{"return", lexemes.Keyword},
	{"short", lexemes.Keyword},
	{"signed", lexemes.Keyword},
	{"sizeof", lexemes.Keyword},
	{"static", lexemes.Keyword},
	{"struct", lexemes.Keyword},
	{"switch", lexemes.Keyword},
	{"typedef", lexemes.Keyword},
	{"union", lexemes.Keyword},
	{"unsigned", lexemes.Keyword},
	{"void", lexemes.Keyword},
	{"volatile", lexemes.Keyword},
	{"while", lexemes.Keyword},
	{"_Alignas", lexemes.Keyword},
	{"_Alignof", lexemes.Keyword},
	{"_Atomic", lexemes.Keyword},
	{"_Bool", lexemes.Keyword},
	{"_Complex", lexemes.Keyword},
	{"_Generic", lexemes.Keyword},
	{"_Imaginary", lexemes.Keyword},
	{"_Noreturn", lexemes.Keyword},
	{"_Static_assert", lexemes.Keyword},
	{"_Thread_local", lexemes.Keyword},
	{"1", lexemes.IntegerConstant},
	{"20123456789", lexemes.IntegerConstant},
	{"0", lexemes.IntegerConstant},
	{"01234567", lexemes.IntegerConstant},
	{"0x0123456789ABCDEF", lexemes.IntegerConstant},
	{"0X0123456789ABCDEF", lexemes.IntegerConstant},
	{"0u", lexemes.IntegerConstant},
	{"1U", lexemes.IntegerConstant},
	{"2ul", lexemes.IntegerConstant},
	{"3uL", lexemes.IntegerConstant},
	{"4Ul", lexemes.IntegerConstant},
	{"5UL", lexemes.IntegerConstant},
	{"6ull", lexemes.IntegerConstant},
	{"7uLL", lexemes.IntegerConstant},
	{"8Ull", lexemes.IntegerConstant},
	{"9ULL", lexemes.IntegerConstant},
	{"10l", lexemes.IntegerConstant},
	{"11L", lexemes.IntegerConstant},
	{"12lu", lexemes.IntegerConstant},
	{"13lU", lexemes.IntegerConstant},
	{"14Lu", lexemes.IntegerConstant},
	{"15LU", lexemes.IntegerConstant},
	{"16llu", lexemes.IntegerConstant},
	{"17llU", lexemes.IntegerConstant},
	{"18LLu", lexemes.IntegerConstant},
	{"19LLU", lexemes.IntegerConstant},
	{".0", lexemes.FloatingConstant},
	{"1.", lexemes.FloatingConstant},
	{"123.0", lexemes.FloatingConstant},
	{".0e1", lexemes.FloatingConstant},
	{"1.E0", lexemes.FloatingConstant},
	{".2e+1", lexemes.FloatingConstant},
	{"3.E-20", lexemes.FloatingConstant},
	{".4e5f", lexemes.FloatingConstant},
	{"6.E7F", lexemes.FloatingConstant},
	{".8e+9l", lexemes.FloatingConstant},
	{"0.E-10L", lexemes.FloatingConstant},
	{"1e+0", lexemes.FloatingConstant},
	{"2E-1f", lexemes.FloatingConstant},
	{"0x0123ABCDEF.p01", lexemes.FloatingConstant},
	{"0xa.p1F", lexemes.FloatingConstant},
	{"0xAbCdp-1L", lexemes.FloatingConstant},
	{"'a'", lexemes.CharLiteral},
	{"'bc'", lexemes.CharLiteral},
	{`'\''`, lexemes.CharLiteral},
	{`'\"'`, lexemes.CharLiteral},
	{`'\\'`, lexemes.CharLiteral},
	{`'\a'`, lexemes.CharLiteral},
	{`'\b'`, lexemes.CharLiteral},
	{`'\f'`, lexemes.CharLiteral},
	{`'\n'`, lexemes.CharLiteral},
	{`'\r'`, lexemes.CharLiteral},
	{`'\t'`, lexemes.CharLiteral},
	{`'\v'`, lexemes.CharLiteral},
	{`'\0'`, lexemes.CharLiteral},
	{`'\12'`, lexemes.CharLiteral},
	{`'\345'`, lexemes.CharLiteral},
	{`'\x0'`, lexemes.CharLiteral},
	{`'\x0123456789ABCDEFabcdef'`, lexemes.CharLiteral},
	{`'\u12Ab'`, lexemes.CharLiteral},
	{`'\U5678CdEf'`, lexemes.CharLiteral},
	{`'\u90AB4gh5'`, lexemes.CharLiteral},
	{`'\U1234567890ij6'`, lexemes.CharLiteral},
	{"L'a'", lexemes.CharLiteral},
	{"L'bc'", lexemes.CharLiteral},
	{`L'\\'`, lexemes.CharLiteral},
	{`L'\a'`, lexemes.CharLiteral},
	{`L'\0'`, lexemes.CharLiteral},
	{`L'\12'`, lexemes.CharLiteral},
	{`L'\345'`, lexemes.CharLiteral},
	{`L'\x0'`, lexemes.CharLiteral},
	{`L'\x0123456789ABCDEFabcdef'`, lexemes.CharLiteral},
	{`L'\u12Ab'`, lexemes.CharLiteral},
	{`L'\U5678CdEf'`, lexemes.CharLiteral},
	{`L'\u90AB4gh5'`, lexemes.CharLiteral},
	{`L'\U1234567890ij6'`, lexemes.CharLiteral},
	{`u'\U5678CdEf'`, lexemes.CharLiteral},
	{`U'\U5678CdEf'`, lexemes.CharLiteral},
	{`"a"`, lexemes.StringLiteral},
	{`"bc"`, lexemes.StringLiteral},
	{`"\""`, lexemes.StringLiteral},
	{`"\""`, lexemes.StringLiteral},
	{`"\\"`, lexemes.StringLiteral},
	{`"\a"`, lexemes.StringLiteral},
	{`"\b"`, lexemes.StringLiteral},
	{`"\f"`, lexemes.StringLiteral},
	{`"\n"`, lexemes.StringLiteral},
	{`"\r"`, lexemes.StringLiteral},
	{`"\t"`, lexemes.StringLiteral},
	{`"\v"`, lexemes.StringLiteral},
	{`"\0"`, lexemes.StringLiteral},
	{`"\12"`, lexemes.StringLiteral},
	{`"\345"`, lexemes.StringLiteral},
	{`"\x0"`, lexemes.StringLiteral},
	{`"\x0123456789ABCDEFabcdef"`, lexemes.StringLiteral},
	{`"\u12Ab"`, lexemes.StringLiteral},
	{`"\U5678CdEf"`, lexemes.StringLiteral},
	{`"\u90AB4gh5"`, lexemes.StringLiteral},
	{`"\U1234567890ij6"`, lexemes.StringLiteral},
	{`L"a"`, lexemes.StringLiteral},
	{`L"bc"`, lexemes.StringLiteral},
	{`L"\\"`, lexemes.StringLiteral},
	{`L"\a"`, lexemes.StringLiteral},
	{`L"\0"`, lexemes.StringLiteral},
	{`L"\12"`, lexemes.StringLiteral},
	{`L"\345"`, lexemes.StringLiteral},
	{`L"\x0"`, lexemes.StringLiteral},
	{`L"\x0123456789ABCDEFabcdef"`, lexemes.StringLiteral},
	{`L"\u12Ab"`, lexemes.StringLiteral},
	{`L"\U5678CdEf"`, lexemes.StringLiteral},
	{`L"\u90AB4gh5"`, lexemes.StringLiteral},
	{`L"\U1234567890ij6"`, lexemes.StringLiteral},
	{`"hello world"`, lexemes.StringLiteral},
	{`L"hello world"`, lexemes.StringLiteral},
	{`u"hello world"`, lexemes.StringLiteral},
	{`U"hello world"`, lexemes.StringLiteral},
	{`u8"hello world"`, lexemes.StringLiteral},
	{"// introducing a comment", lexemes.Comment},
	{"/* multi-\nline\ncomment */", lexemes.Comment},
	{"[", lexemes.LeftBracket},
	{"]", lexemes.RightBracket},
	{"(", lexemes.LeftParenthesis},
	{")", lexemes.RightParenthesis},
	{"{", lexemes.LeftCurlyBrace},
	{"}", lexemes.RightCurlyBrace},
	{".", lexemes.Period},
	{"->", lexemes.Arrow},
	{"++", lexemes.Increment},
	{"--", lexemes.Decrement},
	{"&", lexemes.Ampersand},
	{"|", lexemes.Pipe},
	{"^", lexemes.Caret},
	{"~", lexemes.Tilde},
	{"+", lexemes.Plus},
	{"-", lexemes.Minus},
	{"*", lexemes.Star},
	{"/", lexemes.ForwardSlash},
	{"!", lexemes.Exclamation},
	{"%", lexemes.Percent},
	{"<<", lexemes.DoubleLessThan},
	{">>", lexemes.DoubleGreaterThan},
	{"<", lexemes.LessThan},
	{">", lexemes.GreaterThan},
	{"<=", lexemes.LessThanOrEqual},
	{">=", lexemes.GreaterThanOrEqual},
	{"==", lexemes.DoubleEqual},
	{"!=", lexemes.ExclamationEqual},
	{"&&", lexemes.DoubleAmpersand},
	{"||", lexemes.DoublePipe},
	{"?", lexemes.QuestionMark},
	{":", lexemes.Colon},
	{";", lexemes.SemiColon},
	//{"...", lexemes.Ellipsis},
	{"=", lexemes.Equal},
	{"+=", lexemes.PlusEqual},
	{"-=", lexemes.MinusEqual},
	{"*=", lexemes.StarEqual},
	{"/=", lexemes.ForwardSlashEqual},
	{"%=", lexemes.PercentEqual},
	{"<<=", lexemes.DoubleLessThanEqual},
	{">>=", lexemes.DoubleGreaterThanEqual},
	{"&=", lexemes.AmpersandEqual},
	{"|=", lexemes.PipeEqual},
	{"#", lexemes.Hash},
	{"##", lexemes.DoubleHash},
	{"<:", lexemes.LeftBracket},
	{":>", lexemes.RightBracket},
	{"<%", lexemes.LeftCurlyBrace},
	{"%>", lexemes.RightCurlyBrace},
	{"%:", lexemes.Hash},
	//{"%:%:", lexemes.DoubleHash},
	{",", lexemes.Comma},
	{"    \f\n\r\t\v\f\n\r\t\v ", lexemes.Whitespace},
	{"1.0E", lexemes.Invalid},
}

var partialMatchTestCases = []partialMatchTestCase{
	{"1.0+", Lexeme{lexemes.FloatingConstant, "1.0"}},
	{"..", Lexeme{lexemes.Period, "."}},
	{"...", Lexeme{lexemes.Period, "."}},
	{"%:%:", Lexeme{lexemes.Hash, "%:"}},
}

var errorTestCases = []string{
	`\u000`,
	`\U0000000`,
	`\b`,
	`"\u000"`,
	`"\U0000"`,
	`"\z"`,
	`"\x"`,
	"$",
	`'\z'`,
	"\"hello\n",
	"'h\n",
	"2E",
	"0x",
}
