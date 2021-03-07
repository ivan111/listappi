package main

import (
	"errors"
)

type Env struct {
	data  map[string]*LType
	outer *Env
	depth uint
}

const MAX_ENV_DEPTH = 1000

// key が現れる一番内側のEnvを探す
func (env Env) Find(key string) *Env {
	if _, ok := env.data[key]; ok {
		return &env
	}

	if env.outer != nil {
		return env.outer.Find(key)
	}

	return nil
}

func (env Env) Set(key string, value *LType) *LType {
	env.data[key] = value
	return value
}

func (e Env) Get(key string) (*LType, error) {
	env := e.Find(key)
	if env == nil {
		return nil, errors.New("'" + key + "' not found")
	}

	return env.data[key], nil
}

func NewEnv(vars []string, vals []*LType, outer *Env) (*Env, error) {
	data := make(map[string]*LType)

	if len(vars) != len(vals) {
		return nil, errors.New("引数の数が合わない")
	}

	for i, v := range vars {
		data[v] = vals[i]
	}

	var depth uint
	if outer != nil {
		depth = outer.depth + 1
	}

	if depth >= MAX_ENV_DEPTH {
		return nil, errors.New("スタック制限")
	}

	return &Env{data, outer, depth}, nil
}
