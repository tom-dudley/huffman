package main

import (
	"fmt"
	"os"

	"github.com/tom-dudley/huffman"
)

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Println("Arg should be encode or decode")
		os.Exit(1)
	}
	switch args[0] {
	case "encode":
		input, err := os.ReadFile("input")
		if err != nil {
			panic("Error reading file")
		}

		encoded := huffman.Encode(input)
		err = os.WriteFile("encoded", encoded, 0o600)
		if err != nil {
			fmt.Printf("Error writing to file: %s\n", err.Error())
		}

		fmt.Printf("Input: %d bytes\n", len(input))
		fmt.Printf("Encoded: %d bytes\n", len(encoded))
		fmt.Printf("%d bytes saved\n", len(input)-len(encoded))
		fmt.Println("Encoded file saved to: encoded")
	case "decode":
		encoded, err := os.ReadFile("encoded")
		if err != nil {
			panic("Error reading file")
		}
		output := huffman.Decode(encoded)
		fmt.Println(string(output))
	default:
		fmt.Println("first arg should be encode or decode")
		os.Exit(1)
	}
}
