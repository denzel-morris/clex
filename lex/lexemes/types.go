package lexemes

type Type int

const (
	Invalid Type = iota
	EOF
	Identifier
	Keyword
	IntegerConstant
	FloatingConstant
	CharLiteral
	StringLiteral
	Comment
	Whitespace
	LeftBracket
	RightBracket
	LeftParenthesis
	RightParenthesis
	LeftCurlyBrace
	RightCurlyBrace
	Period
	Arrow
	Increment
	Decrement
	Ampersand
	Pipe
	Caret
	Tilde
	Plus
	Minus
	Star
	ForwardSlash
	Exclamation
	Percent
	LessThan
	GreaterThan
	DoubleLessThan
	DoubleGreaterThan
	LessThanOrEqual
	GreaterThanOrEqual
	DoubleEqual
	ExclamationEqual
	DoubleAmpersand
	DoublePipe
	QuestionMark
	Colon
	SemiColon
	Ellipsis
	Equal
	PlusEqual
	MinusEqual
	StarEqual
	ForwardSlashEqual
	PercentEqual
	DoubleLessThanEqual
	DoubleGreaterThanEqual
	AmpersandEqual
	PipeEqual
	Hash
	DoubleHash
	Comma
)

var typeToName = map[Type]string{
	Invalid:                "Invalid",
	EOF:                    "EOF",
	Identifier:             "Identifier",
	Keyword:                "Keyword",
	IntegerConstant:        "IntegerConstant",
	FloatingConstant:       "FloatingConstant",
	CharLiteral:            "CharLiteral",
	StringLiteral:          "StringLiteral",
	Comment:                "Comment",
	Whitespace:             "Whitespace",
	LeftBracket:            "LeftBracket",
	RightBracket:           "RightBracket",
	LeftParenthesis:        "LeftParenthesis",
	RightParenthesis:       "RightParenthesis",
	LeftCurlyBrace:         "LeftCurlyBrace",
	RightCurlyBrace:        "RightCurlyBrace",
	Period:                 "Period",
	Arrow:                  "Arrow",
	Increment:              "Increment",
	Decrement:              "Decrement",
	Ampersand:              "Ampersand",
	Pipe:                   "Pipe",
	Caret:                  "Caret",
	Tilde:                  "Tilde",
	Plus:                   "Plus",
	Minus:                  "Minus",
	Star:                   "Star",
	ForwardSlash:           "ForwardSlash",
	Exclamation:            "Exclamation",
	Percent:                "Percent",
	LessThan:               "LessThan",
	GreaterThan:            "GreaterThan",
	DoubleLessThan:         "DoubleLessThan",
	DoubleGreaterThan:      "DoubleGreaterThan",
	LessThanOrEqual:        "LessThanOrEqual",
	GreaterThanOrEqual:     "GreaterThanOrEqual",
	DoubleEqual:            "DoubleEqual",
	ExclamationEqual:       "ExclamationEqual",
	DoubleAmpersand:        "DoubleAmpersand",
	DoublePipe:             "DoublePipe",
	QuestionMark:           "QuestionMark",
	Colon:                  "Colon",
	SemiColon:              "SemiColon",
	Ellipsis:               "Ellipsis",
	Equal:                  "Equal",
	PlusEqual:              "PlusEqual",
	MinusEqual:             "MinusEqual",
	StarEqual:              "StarEqual",
	ForwardSlashEqual:      "ForwardSlashEqual",
	PercentEqual:           "PercentEqual",
	DoubleLessThanEqual:    "DoubleLessThanEqual",
	DoubleGreaterThanEqual: "DoubleGreaterThanEqual",
	AmpersandEqual:         "AmpersandEqual",
	PipeEqual:              "PipeEqual",
	Hash:                   "Hash",
	DoubleHash:             "DoubleHash",
	Comma:                  "Comma",
}

func (t Type) String() string {
	return typeToName[t]
}
