package main

import (
	"fmt"
	"strconv"
	"strings"
)

type ValueType uint8

const (
	TypeNil ValueType = iota + 1
	TypeInt
	TypeBool
	TypeList
	TypeSymbol
	TypeClosure
)

type LType struct {
	Type ValueType
	Val  LValue
}

type LValue interface {
}

type LClosure struct {
	args []string
	body *LType
}

var Nil = &LType{TypeNil, nil}
var True = &LType{TypeBool, true}
var False = &LType{TypeBool, false}

func (this *LType) String() string {
	switch this.Type {
	case TypeNil:
		return "nil"
	case TypeBool:
		val, ok := this.Val.(bool)
		if ok == false {
			return ""
		}

		if val {
			return "#t"
		} else {
			return "#f"
		}
	case TypeInt:
		return strconv.Itoa(getIntValue(this))
	case TypeSymbol:
		return getSymbolName(this)
	case TypeList:
		list, ok := this.Val.([]*LType)
		if ok == false {
			return ""
		}

		var strList []string

		for _, v := range list {
			strList = append(strList, v.String())
		}

		return fmt.Sprintf("(%s)", strings.Join(strList, " "))
	case TypeClosure:
		return "closure"
	}

	return ""
}

func NewInt(val int) *LType {
	return &LType{TypeInt, val}
}

func NewSymbol(name string) *LType {
	return &LType{TypeSymbol, name}
}

func NewList(list []*LType) *LType {
	return &LType{TypeList, list}
}

func NewClosure(args []string, body *LType) *LType {
	return &LType{TypeClosure, &LClosure{args, body}}
}

func getIntValue(obj *LType) int {
	if obj.Type != TypeInt {
		fmt.Print("getIntValue: 来ないはず1")
		return 0
	}

	val, ok := obj.Val.(int)
	if ok == false {
		fmt.Print("getIntValue: 来ないはず2")
		return 0
	}

	return val
}

func getSymbolName(obj *LType) string {
	if obj.Type != TypeSymbol {
		fmt.Print("getSymbolName: 来ないはず1")
		return ""
	}

	val, ok := obj.Val.(string)
	if ok == false {
		fmt.Print("getSymbolName: 来ないはず2")
		return ""
	}

	return val
}

func getList(obj *LType) []*LType {
	if obj.Type != TypeList {
		fmt.Print("getList: 来ないはず1")
		return nil
	}

	val, ok := obj.Val.([]*LType)
	if ok == false {
		fmt.Print("getList: 来ないはず2")
		return nil
	}

	return val
}

func getClosureStruct(obj *LType) *LClosure {
	if obj.Type != TypeClosure {
		fmt.Print("getClosureStruct: 来ないはず1")
		return nil
	}

	val, ok := obj.Val.(*LClosure)
	if ok == false {
		fmt.Print("getClosureStruct: 来ないはず2")
		return nil
	}

	return val
}
