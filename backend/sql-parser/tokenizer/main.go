package tokenizer

import (
	"fmt"
	"strings"
)

type Token struct {
	kind  int
	value string
}

func (t Token) String() string {
	return fmt.Sprintf("Token{kind: %d, value: %q}", t.kind, t.value)
}

func TokenizeSQL(sqlFile string) []Token {
	splitSQL := strings.Split(sqlFile, " ")
	tokens := []Token{}
	for _, token := range splitSQL {
		tokens = append(tokens, Token{0, token})
	}

	return tokens
}
