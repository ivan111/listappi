package main

import (
	"errors"
)

func eval(x *LType, env *Env) (*LType, error) {
	switch x.Type {
	case TypeSymbol:
		return env.Get(getSymbolName(x))
	case TypeNil, TypeInt, TypeClosure:
		return x, nil
	case TypeList:
		list := getList(x)
		car := list[0]
		cdr := list[1:]
		var prim LFunc
		var closure *LType

		if car.Type == TypeSymbol {
			symbol := getSymbolName(car)

			var ok bool
			prim, ok = PrimMacros[symbol]
			if ok {
				return prim(cdr, env)
			}

			prim, ok = PrimFuncs[symbol]
			if ok == false {
				var err error
				closure, err = eval(car, env)
				if err != nil {
					return nil, err
				}

				if closure.Type != TypeClosure {
					return nil, errors.New("リストの先頭シンボルの参照先がクロージャじゃない")
				}
			}
		} else if car.Type == TypeClosure {
			closure = car
		}

		if prim == nil && closure == nil {
			return nil, errors.New("リストの先頭要素が処理できない")
		}

		var exps []*LType

		for _, exp := range cdr {
			result, err := eval(exp, env)
			if err != nil {
				return nil, err
			}

			exps = append(exps, result)
		}

		if prim != nil {
			return prim(exps, env)
		} else {
			return callClosure(getClosureStruct(closure), exps, env)
		}
	}

	return Nil, nil
}

func callClosure(closure *LClosure, args []*LType, outerEnv *Env) (*LType, error) {
	env, err := NewEnv(closure.args, args, outerEnv)
	if err != nil {
		return Nil, err
	}

	return eval(closure.body, env)
}
