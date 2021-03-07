package main

import (
	"bufio"
	"errors"
	"io"
	"strconv"
)

const MAX_PARSE_DEPTH uint = 100

func nextToken(reader *bufio.Reader) (string, error) {
	err := skipSpaces(reader)
	if err != nil {
		return "", err
	}

	r, _, err := reader.ReadRune()
	if err != nil {
		return "", err
	}

	if r == '(' || r == ')' || r == '\'' {
		return string(r), nil
	}

	list := []rune{r}

	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			if len(list) != 0 && err == io.EOF {
				return string(list), nil
			}

			return string(list), err
		}

		if r == '(' || r == ')' || r == ' ' || r == '\t' || r == '\n' {
			reader.UnreadRune()
			return string(list), nil
		}

		list = append(list, r)
	}
}

func skipSpaces(reader *bufio.Reader) error {
	for {
		r, _, err := reader.ReadRune()

		if err != nil {
			return err
		}

		switch r {
		case ' ', '\t', '\n':
			// nop
		default:
			reader.UnreadRune()
			return nil
		}
	}
}

func parse(reader *bufio.Reader, depth uint) (*LType, error) {
	if depth >= MAX_PARSE_DEPTH {
		return nil, errors.New("exceeds max parse depth")
	}

	token, err := nextToken(reader)
	if err != nil {
		return nil, err
	}

	if token == "(" {
		var list []*LType

		for {
			err = skipSpaces(reader)
			if err != nil {
				return nil, err
			}

			r, _, err := reader.ReadRune()
			if err != nil {
				return nil, err
			}

			if r == ')' {
				return NewList(list), nil
			}

			reader.UnreadRune()

			v, err := parse(reader, depth+1)
			if err != nil {
				return nil, err
			}

			list = append(list, v)
		}
	} else if token == ")" {
		return nil, errors.New("Syntax Error: Unexped ')'")
	} else if token == "'" {
		form, err := parse(reader, depth+1)
		if err != nil {
			return nil, err
		}

		list := []*LType{NewSymbol("quote"), form}
		return NewList(list), nil
	} else {
		return atom(token)
	}
}

func atom(token string) (*LType, error) {
	val, err := strconv.Atoi(token)
	if err == nil {
		return NewInt(val), nil
	}

	return NewSymbol(token), nil
}
