package robo

import (
	"fmt"
	"strconv"
	"strings"
)

/*
Instruction Syntax
R=Right, L=Left, U=Up, D=Down
<dir> = Direction
<dis> = Displacement
<inst> = Instruction

Notation:
<dir> = "U" | "R" | "D" | "L" | "u" | "r" | "d" | "l"
<dis> = "0" | "1" | "2" | "3" | ... | "9"
<inst> =  "" | <dir> | <dir><dis> | <inst>
*/

type TokenType int

const (
	TokenUnknown TokenType = -1
	TokenDir     TokenType = iota + 1
	TokenDis
)

type Token struct {
	Type  TokenType
	Match string
}

type Program string

type MoveCommand struct {
	Direction    FaceDirection
	Displacement int
}

func (p Program) Compile() ([]MoveCommand, error) {
	// lexical analysis (lexer)
	var tokenizer Tokenizer = Tokenizer(p)
	tokens, err := tokenizer.Lexify()
	if err != nil {
		return nil, err
	}

	// parsing logic (parser)
	var parser Parser = tokens
	return parser.Parse()
}

// Tokenizer performs lexical analysis of instruction.
type Tokenizer string

func (t Tokenizer) Lexify() ([]Token, error) {
	var tokens []Token
	for i, chr := range t {
		ttype := t.getType(chr)
		if ttype == TokenUnknown {
			return nil, fmt.Errorf("syntax error: invalid instruction at index %d", i)
		}

		if i == 0 && ttype != TokenDir {
			return nil, fmt.Errorf("syntax error: program must begin with a directional instruction")
		}

		tokens = append(tokens, Token{ttype, strings.ToLower(string(chr))})
	}

	return tokens, nil
}
func (t Tokenizer) getType(ch rune) TokenType {
	switch ch {
	case 'U', 'R', 'D', 'L', 'u', 'r', 'd', 'l':
		return TokenDir
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return TokenDis
	default:
		return TokenUnknown
	}
}

type Parser []Token

func (p Parser) Parse() ([]MoveCommand, error) {
	var dirs = map[string]FaceDirection{
		"u": FaceNorth,
		"r": FaceEast,
		"d": FaceSouth,
		"l": FaceWest,
	}
	var cmds []MoveCommand
	i, n := 0, len(p)-1
	for i <= n {
		token := p[i]
		cmd := MoveCommand{dirs[token.Match], 1}
		dis := ""
		for i+1 <= n && p[i+1].Type == TokenDis {
			i++
			dis += p[i].Match
		}

		if dis != "" {
			cmd.Displacement, _ = strconv.Atoi(dis)
		}
		cmds = append(cmds, cmd)
		i++
	}

	return cmds, nil
}
