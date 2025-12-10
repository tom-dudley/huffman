package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Println("Arg should be encode or decode")
		os.Exit(1)
	}
	switch args[0] {
	case "encode":
		buf, err := os.ReadFile("input")
		if err != nil {
			panic("Error reading file")
		}

		encodeInput(buf)
	case "decode":
		encoded, err := os.ReadFile("encoded")
		if err != nil {
			panic("Error reading file")
		}
		output := decode(encoded)
		fmt.Println(string(output))
	default:
		fmt.Println("first arg should be encode or decode")
		os.Exit(1)
	}
}
