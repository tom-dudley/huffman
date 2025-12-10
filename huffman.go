package huffman

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Node struct {
	char byte
	name string
	freq int
	l    *Node
	r    *Node
}

func buildBinaryTree(frequencies []*Node) *Node {
	for {
		if len(frequencies) <= 1 {
			break
		}

		l := frequencies[0]
		r := frequencies[1]

		p := &Node{
			l:    l,
			r:    r,
			freq: l.freq + r.freq,
			name: l.name + r.name,
		}

		frequencies = frequencies[2:]

		if len(frequencies) == 0 {
			// frequencies = append(frequencies, *p)
			return p
		}

		if p.freq > frequencies[len(frequencies)-1].freq {
			frequencies = append(frequencies, p)
		} else if p.freq < frequencies[0].freq {
			frequencies = slices.Insert(frequencies, 0, p)
		} else {
			maxIndex := len(frequencies) - 1
			for i := 0; i < maxIndex; i++ {
				if frequencies[i].freq <= p.freq && p.freq <= frequencies[i+1].freq {
					frequencies = slices.Insert(frequencies, i+1, p)
					break
				}
			}
		}
	}

	return nil
}

func printTree(node *Node) {
	if node == nil {
		return
	}
	fmt.Printf("Node: %s : %d\n", string(node.name), node.freq)
	printTree(node.l)
	printTree(node.r)
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

type codebookEntry struct {
	b    byte
	code string
}

func sortCodes(codes map[byte]string) []codebookEntry {
	codebook := []codebookEntry{}
	for k, v := range codes {
		codebook = append(codebook, codebookEntry{
			b:    k,
			code: v,
		})
	}
	for {
		resort := false
		for i := 0; i < len(codebook)-1; i++ {
			this := codebook[i]
			next := codebook[i+1]
			if len(this.code) > len(next.code) ||
				(len(this.code) == len(next.code) && this.b > next.b) {
				tmp := codebook[i]
				codebook[i] = codebook[i+1]
				codebook[i+1] = tmp

				resort = true
				break
			}
		}

		if !resort {
			break
		}
	}

	return codebook
}

func constructCanoncial(codes map[byte]string) map[byte]string {
	codebook := sortCodes(codes)
	for i := range codebook {
		if i == 0 {
			codebook[0].code = strings.Repeat("0", len(codebook[0].code))
			continue
		}

		// Increment the bitstring by 1
		prev := codebook[i-1].code
		asInt, err := strconv.ParseInt(prev, 2, 8)
		if err != nil {
			fmt.Println(err)
		}

		asInt++

		// If the length of this bitstring is longer than the previous, then append a 0 to the RHS

		fmtStr := "%0" + strconv.Itoa(len(prev)) + "b"
		s := fmt.Sprintf(fmtStr, asInt)

		if len(codebook[i].code) > len(prev) {
			s += "0"
		}

		codebook[i].code = s
	}

	newCodes := map[byte]string{}
	for _, code := range codebook {
		newCodes[code.b] = code.code
	}

	return newCodes
}

func Encode(input []byte) []byte {
	input = append(input, 0x00)

	frequenciesMap := map[byte]int{}
	for _, b := range input {
		f, ok := frequenciesMap[b]
		if ok {
			frequenciesMap[b] = f + 1
		} else {
			frequenciesMap[b] = 1
		}
	}

	frequencies := []*Node{}

	for k, v := range frequenciesMap {
		frequencies = append(frequencies, &Node{char: k, freq: v, name: string(k)})
	}

	// for _, charAndFreq := range frequencies {
	// 	fmt.Println(charAndFreq)
	// }

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
	// printTree(rootNode)
	codes := buildCodes(rootNode, "", map[byte]string{})
	codes = constructCanoncial(codes)
	// fmt.Println("Generated canonical huffman codes:")
	// for k, code := range codes {
	// 	fmt.Printf("    %s : %s\n", string(k), code)
	// }

	// TODO: Switch from treating codes as strings to bits

	encodedInputBytes := encode(input, codes)
	encodedHuffmanBytes := encodeHuffman(codes)

	encoded := append(encodedHuffmanBytes, encodedInputBytes...)
	return encoded
}

func Decode(encoded []byte) []byte {
	numberOfSymbols := encoded[0]
	var huffmanBytesLength int
	if numberOfSymbols%2 == 0 {
		huffmanBytesLength = int(numberOfSymbols) + int(numberOfSymbols)/2
	} else {
		huffmanBytesLength = int(numberOfSymbols) + int(numberOfSymbols)/2 + 1
	}

	codes := decodeHuffman(encoded[:huffmanBytesLength+1])
	// for symbol, code := range codes {
	// 	fmt.Printf("%s : %s\n", string(symbol), code)
	// }

	rootNode := buildHuffmanTree(codes)
	return decodeWithTree(encoded[huffmanBytesLength+1:], rootNode)
}

func encode(input []byte, codes map[byte]string) []byte {
	// TODO: String builder would be better here
	// TODO: Even better bit shifting/construction
	var encoded string
	for _, b := range input {
		encoded += fmt.Sprintf(codes[b])
	}

	// TODO: This might want amending. given that 00 is a valid code.
	padding := 8 - (len(encoded) % 8)
	for i := 0; i < padding; i++ {
		encoded += "0"
	}
	// fmt.Printf("Encoded as %d bits\n", len(encoded))
	// fmt.Printf("Encoded as %d bytes\n", len(encoded)/8)

	encodedInputBytes := []byte{}

	for i := 0; i < len(encoded); i += 8 {
		// fmt.Printf("Got bitstring: %s\n", encoded[i:i+8])
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
		// fmt.Printf("Constructed byte: %08b\n", x)
		encodedInputBytes = append(encodedInputBytes, x)
	}

	return encodedInputBytes
}

// TODO: Why aren't we just passing the ordered slice rather than a map?
func encodeHuffman(codes map[byte]string) []byte {
	encoded := []byte{
		byte(len(codes)), // First byte is the number of symbols
	}

	codebook := sortCodes(codes)
	for _, code := range codebook {
		encoded = append(encoded, code.b)
	}

	// Build packed nibbles of code lengths
	for i := 0; i < len(codebook); i += 2 {
		var packed byte
		if i == len(codebook)-1 {
			// fmt.Printf("Packing: %s\n", codebook[i].code)
			packed = byte(len(codebook[i].code)) << 4
		} else {
			// fmt.Printf("Packing: %s\n", codebook[i].code)
			// fmt.Printf("Packing: %s\n", codebook[i+1].code)
			packed = byte(len(codebook[i].code))<<4 | byte(len(codebook[i+1].code))
		}
		encoded = append(encoded, packed)
	}

	os.WriteFile("test", encoded, 0o600)

	return encoded
}

func decodeHuffman(encoded []byte) map[byte]string {
	numberOfSymbols := encoded[0]
	// fmt.Printf("Decoding %d symbols\n", numberOfSymbols)
	symbols := encoded[1 : numberOfSymbols+1]

	lengths := []int{}

	packedLengths := encoded[numberOfSymbols+1:]
	for _, b := range packedLengths {
		upperNibble := b >> 4
		lengths = append(lengths, int(upperNibble))

		lowerNibble := b & 0b00001111
		lengths = append(lengths, int(lowerNibble))
	}

	codes := map[byte]string{}
	for i := 0; i < len(symbols); i++ {
		// fmt.Printf("%s : %d\n", string(symbols[i]), lengths[i])
		codes[symbols[i]] = strings.Repeat("0", lengths[i])
	}

	return constructCanoncial(codes)
}

func buildHuffmanTree(codes map[byte]string) *Node {
	codebook := sortCodes(codes)

	rootNode := &Node{}

	var currentNode *Node
	for _, code := range codebook {
		currentNode = rootNode
		lastCodeIndex := len(code.code) - 1
		for i, digit := range code.code {
			if digit == '0' {
				if currentNode.l == nil {
					currentNode.l = &Node{}
				}
				currentNode = currentNode.l
			} else {
				if currentNode.r == nil {
					currentNode.r = &Node{}
				}
				currentNode = currentNode.r
			}

			if i == lastCodeIndex {
				currentNode.char = code.b
			}
		}
	}

	return rootNode
}

func decodeWithTree(encoded []byte, rootNode *Node) []byte {
	digitsString := ""
	for _, b := range encoded {
		digitsString += fmt.Sprintf("%08b", b)
	}

	currentNode := rootNode
	decoded := []byte{}
	for _, digit := range digitsString {
		if digit == '0' {
			currentNode = currentNode.l
		} else {
			currentNode = currentNode.r
		}

		if currentNode.l == nil && currentNode.r == nil {
			decoded = append(decoded, currentNode.char)
			// If we hit NUL that's the end of the data
			if currentNode.char == 0x0 {
				return decoded
			}
			currentNode = rootNode
		}
	}

	return decoded
}
