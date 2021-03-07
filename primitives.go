package main

import (
	"errors"
)

type LFunc func([]*LType, *Env) (*LType, error)

var PrimMacros map[string]LFunc

var PrimFuncs map[string]LFunc

func registPrimitives() {
	PrimMacros = map[string]LFunc{
		"quote":  quote,
		"if":     primIf,
		"define": define,
		"lambda": lambda,
	}

	PrimFuncs = map[string]LFunc{
		"car":  car,
		"cdr":  cdr,
		"cons": cons,
		"=":    equal,
		">":    greaterThan,
		"<":    lessThan,
		">=":   greaterThanOrEqual,
		"<=":   lessThanOrEqual,
		"+":    add,
		"-":    sub,
		"*":    mul,
	}
}

func quote(args []*LType, env *Env) (*LType, error) {
	if len(args) != 1 {
		return nil, errors.New("quote: 引数の数が違う。quoteは1つの引数を取る")
	}

	return args[0], nil
}

func primIf(args []*LType, env *Env) (*LType, error) {
	if len(args) != 3 {
		return nil, errors.New("if: 引数の数が違う。ifは3つの引数を取る")
	}

	pred, err := eval(args[0], env)
	if err != nil {
		return nil, err
	}

	if pred == False {
		return eval(args[2], env)
	}

	return eval(args[1], env)
}

func define(args []*LType, env *Env) (*LType, error) {
	if len(args) != 2 {
		return nil, errors.New("define: 引数の数が違う。2つの引数を取る")
	}

	v := args[0]
	if v.Type == TypeSymbol {
		d, err := eval(args[1], env)
		if err != nil {
			return nil, err
		}

		env.Set(getSymbolName(v), d)

		return v, nil
	} else if v.Type == TypeList {
		list := getList(v)

		if len(list) == 0 {
			return nil, errors.New("define: 1つめの引数に空リストはだめ")
		}

		if list[0].Type != TypeSymbol {
			return nil, errors.New("define: 1つめの引数のリストの１番目の値はシンボルじゃないとだめ")
		}

		varName := getSymbolName(list[0])
		vars := make([]string, 0)

		for _, d := range list[1:] {
			if d.Type != TypeSymbol {
				return nil, errors.New("define: 引数リストがシンボルじゃない")
			}

			vars = append(vars, getSymbolName(d))
		}

		closure := NewClosure(vars, args[1])

		env.Set(varName, closure)

		return closure, nil
	} else {
		return nil, errors.New("define: 1つめの引数はシンボルかリストじゃないとだめ")
	}
}

func lambda(args []*LType, env *Env) (*LType, error) {
	if len(args) != 2 {
		return nil, errors.New("lambda: 引数の数が違う。2つの引数を取る")
	}

	list := args[0]
	if list.Type != TypeList {
		return nil, errors.New("lambda: 1つめの引数はリストじゃないとだめ")
	}

	vars := make([]string, 0)
	for _, v := range getList(list) {
		if v.Type != TypeSymbol {
			return nil, errors.New("lambda: 引数リストの中はシンボルじゃないとだめ")
		}

		vars = append(vars, getSymbolName(v))
	}

	return NewClosure(vars, args[1]), nil
}

func car(args []*LType, env *Env) (*LType, error) {
	if len(args) != 1 {
		return nil, errors.New("car: 引数の数が違う。1つのリストを引数に取る")
	}

	v := args[0]
	if v.Type != TypeList {
		return nil, errors.New("car: 引数がリストじゃない")
	}

	list := getList(v)
	if len(list) == 0 {
		return nil, errors.New("car: 引数のリストがnil")
	}

	return list[0], nil
}

func cdr(args []*LType, env *Env) (*LType, error) {
	if len(args) != 1 {
		return nil, errors.New("cdr: 引数の数が違う。１つのリストを引数に取る")
	}

	v := args[0]
	if v.Type != TypeList {
		return nil, errors.New("cdr: 引数がリストじゃない")
	}

	list := getList(v)
	if len(list) == 0 {
		return nil, errors.New("cdr: 引数のリストがnil")
	}

	if len(list) == 1 {
		return Nil, nil
	}

	return NewList(list[1:]), nil
}

func cons(args []*LType, env *Env) (*LType, error) {
	if len(args) != 2 {
		return nil, errors.New("cons: 引数の数が違う。2つの引数に取る")
	}

	v := args[0]
	if v.Type != TypeList {
		return nil, errors.New("car: 引数がリストじゃない")
	}

	list := getList(v)
	if len(list) == 0 {
		return nil, errors.New("car: 引数のリストがnil")
	}

	return list[0], nil
}

func equal(args []*LType, env *Env) (*LType, error) {
	if len(args) < 2 {
		return nil, errors.New("=: 引数の数が違う。2つ以上の引数に取る")
	}

	prev := args[0]
	if prev.Type != TypeInt {
		return False, nil
	}
	prevVal := getIntValue(prev)

	for _, v := range args[1:] {
		if v.Type != TypeInt {
			return False, nil
		}

		nextVal := getIntValue(v)
		if prevVal != nextVal {
			return False, nil
		}

		prev = v
		prevVal = nextVal
	}

	return True, nil
}

func greaterThan(args []*LType, env *Env) (*LType, error) {
	if len(args) < 2 {
		return nil, errors.New(">: 引数の数が違う。2つ以上の引数に取る")
	}

	prev := args[0]
	if prev.Type != TypeInt {
		return False, nil
	}
	prevVal := getIntValue(prev)

	for _, v := range args[1:] {
		if v.Type != TypeInt {
			return False, nil
		}

		nextVal := getIntValue(v)
		if prevVal <= nextVal {
			return False, nil
		}

		prev = v
		prevVal = nextVal
	}

	return True, nil
}

func lessThan(args []*LType, env *Env) (*LType, error) {
	if len(args) < 2 {
		return nil, errors.New("<: 引数の数が違う。2つ以上の引数に取る")
	}

	prev := args[0]
	if prev.Type != TypeInt {
		return False, nil
	}
	prevVal := getIntValue(prev)

	for _, v := range args[1:] {
		if v.Type != TypeInt {
			return False, nil
		}

		nextVal := getIntValue(v)
		if prevVal >= nextVal {
			return False, nil
		}

		prev = v
		prevVal = nextVal
	}

	return True, nil
}

func greaterThanOrEqual(args []*LType, env *Env) (*LType, error) {
	if len(args) < 2 {
		return nil, errors.New(">=: 引数の数が違う。2つ以上の引数に取る")
	}

	prev := args[0]
	if prev.Type != TypeInt {
		return False, nil
	}
	prevVal := getIntValue(prev)

	for _, v := range args[1:] {
		if v.Type != TypeInt {
			return False, nil
		}

		nextVal := getIntValue(v)
		if prevVal < nextVal {
			return False, nil
		}

		prev = v
		prevVal = nextVal
	}

	return True, nil
}

func lessThanOrEqual(args []*LType, env *Env) (*LType, error) {
	if len(args) < 2 {
		return nil, errors.New("<=: 引数の数が違う。2つ以上の引数に取る")
	}

	prev := args[0]
	if prev.Type != TypeInt {
		return False, nil
	}
	prevVal := getIntValue(prev)

	for _, v := range args[1:] {
		if v.Type != TypeInt {
			return False, nil
		}

		nextVal := getIntValue(v)
		if prevVal > nextVal {
			return False, nil
		}

		prev = v
		prevVal = nextVal
	}

	return True, nil
}

func add(args []*LType, env *Env) (*LType, error) {
	var res int

	for _, v := range args {
		if v.Type != TypeInt {
			return nil, errors.New("+: 引数にint以外がある")
		}

		res += getIntValue(v)
	}

	return NewInt(res), nil
}

func sub(args []*LType, env *Env) (*LType, error) {
	var res int

	for i, v := range args {
		if v.Type != TypeInt {
			return nil, errors.New("-: 引数にint以外がある")
		}

		if i == 0 {
			res = getIntValue(v)
		} else {
			res -= getIntValue(v)
		}
	}

	return NewInt(res), nil
}

func mul(args []*LType, env *Env) (*LType, error) {
	res := 1

	for _, v := range args {
		if v.Type != TypeInt {
			return nil, errors.New("*: 引数にint以外がある")
		}

		res *= getIntValue(v)
	}

	return NewInt(res), nil
}
