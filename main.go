package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	repl()
}

func repl() {
	var reader = bufio.NewReader(os.Stdin)
	env, _ := NewEnv(nil, nil, nil)
	registPrimitives()

	for {
		fmt.Print("user> ")

		ast, err := parse(reader, 0)
		if err == io.EOF {
			return
		}

		if err != nil {
			fmt.Println(err)
			continue
		}

		if ast == nil {
			fmt.Println("parseの戻り値がnil")
			continue
		}

		v, err := eval(ast, env)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if v == nil {
			fmt.Println("evalの戻り値がnil")
			continue
		}

		fmt.Println("=>", v)
	}
}
