package main

import (
	"fmt"
	"os"
	"slices"
)

type Node struct {
	char byte
	name string
	freq int
	l    *Node
	r    *Node
	p    *Node
}

func buildBinaryTree(frequencies []Node) *Node {
	// Construct binary tree
	for {
		if len(frequencies) <= 1 {
			break
		}

		l := frequencies[0]
		r := frequencies[1]

		p := &Node{
			l:    &l,
			r:    &r,
			freq: l.freq + r.freq,
			name: l.name + r.name,
		}

		// TODO: Do we need to know the parent?
		l.p = p
		r.p = p

		frequencies = frequencies[2:]

		if len(frequencies) == 0 {
			// frequencies = append(frequencies, *p)
			return p
		}

		if p.freq > frequencies[len(frequencies)-1].freq {
			frequencies = append(frequencies, *p)
		} else {
			maxIndex := len(frequencies) - 1
			for i := 0; i < maxIndex; i++ {
				if frequencies[i].freq <= p.freq && p.freq <= frequencies[i+1].freq {
					frequencies = slices.Insert(frequencies, i+1, *p)
					break
				}
			}
		}
	}

	return nil
}

func traverseBinaryTree(node *Node) {
	if node == nil {
		return
	}
	fmt.Printf("Node: %s : %d\n", string(node.name), node.freq)
	traverseBinaryTree(node.l)
	traverseBinaryTree(node.r)
}

func buildCodes(node *Node, code string, codes map[byte]string) map[byte]string {
	if node.l == nil && node.r == nil {
		codes[node.char] = code
		return codes
	}
	codes = buildCodes(node.l, code+"0", codes)
	codes = buildCodes(node.r, code+"1", codes)
	return codes
}

func main() {
	buf, err := os.ReadFile("input")
	if err != nil {
		panic("Error reading file")
	}

	fmt.Printf("Got %d bytes\n", len(buf))
	frequenciesMap := map[byte]int{}
	for _, b := range buf {
		f, ok := frequenciesMap[b]
		if ok {
			frequenciesMap[b] = f + 1
		} else {
			frequenciesMap[b] = 1
		}
	}

	frequencies := []Node{}

	for k, v := range frequenciesMap {
		frequencies = append(frequencies, Node{char: k, freq: v, name: string(k)})
	}

	for _, charAndFreq := range frequencies {
		fmt.Println(charAndFreq)
	}

	// Now sort the frequencies
	for {
		resort := false
		// Walk the slice, checking if we're ordered from lowest to highest freq
		for i := 0; i < len(frequencies)-1; i++ {
			if frequencies[i].freq > frequencies[i+1].freq {
				tmp := frequencies[i]
				frequencies[i] = frequencies[i+1]
				frequencies[i+1] = tmp

				resort = true
				break
			}
		}

		if !resort {
			break
		}
	}

	// Build the tree and generate the codes
	rootNode := buildBinaryTree(frequencies)
	traverseBinaryTree(rootNode)
	codes := buildCodes(rootNode, "", map[byte]string{})
	for k, code := range codes {
		fmt.Printf("%s : %s\n", string(k), code)
	}

	// TODO: Switch from treating codes as strings to bits

	// Encode
	// TODO: String builder would be better here
	// TODO: Even better bit shifting/construction
	var encoded string
	for _, b := range buf {
		encoded += fmt.Sprintf(codes[b])
	}

	padding := 8 - (len(encoded) % 8)
	for i := 0; i < padding; i++ {
		encoded += "0"
	}
	fmt.Printf("Encoded as %d bits\n", len(encoded))
	fmt.Printf("Encoded as %d bytes\n", len(encoded)/8)

	encodedBytes := []byte{}

	for i := 0; i < len(encoded); i += 8 {
		fmt.Printf("Got bitstring: %s\n", encoded[i:i+8])
		b := encoded[i : i+8]
		var x byte
		for i := 0; i < 8; i++ {
			if b[i] == '1' {
				switch i {
				case 0:
					x += 1 << 7
				case 1:
					x += 1 << 6
				case 2:
					x += 1 << 5
				case 3:
					x += 1 << 4
				case 4:
					x += 1 << 3
				case 5:
					x += 1 << 2
				case 6:
					x += 1 << 1
				case 7:
					x += 1 << 0
				}
			}
		}
		fmt.Printf("Constructed byte: %08b\n", x)
		encodedBytes = append(encodedBytes, x)
	}

	fmt.Println(encodedBytes)
}
