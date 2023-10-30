package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
)

// Compress takes a string and compresses it using huffman coding
// The result is a byte array that can be stored in a file
// The first 8 bytes stores the lenght of the serialized tree.
// The next 4 bytes stores the added padding.
// The rest of the file is the compressed content
// The header is used to decompress the file
func Compress(content string) ([]byte, error) {

	inputContent := strings.Split(content, "")

	// split files
	freqTable := getFreq(inputContent)

	// create root node
	root := createRootNode(freqTable)

	// generate precfix codes
	prefixCodes := prefixCodes(root)

	bitString, paddingAdded := stringToBitString(inputContent, prefixCodes)

	// convert bitString to bytes
	compressed, err := bitStringToByteArray(bitString)
	if err != nil {
		return []byte{}, fmt.Errorf("failed to convert bitstring to bytes: %w", err)
	}

	// create the metadata for the compresstion
	metadata := encodeHeader(*root, paddingAdded)

	// join metadata and compressed file content
	res := append(metadata, compressed...)

	return res, nil
}

// Decompress takes a byte array and decompresses it using huffman coding
func Decompress(content []byte) (string, error) {
	pad, decodedTree, byteString, err := decodeContent(content)
	if err != nil {
		return "", fmt.Errorf("failed to decode content: %w", err)
	}

	decodedUncompressed := byteArrayToBitString(byteString)
	decoded := decode(decodedUncompressed, decodedTree, int(pad))

	return decoded, nil
}

// prefixCodes generates a map with prefix codes for every character
// i.e. a = 01, h = 0101
func prefixCodes(item *Item) map[string]string {
	prefixCodes := map[string]string{}
	var prefix string
	getPrefixMap(item, &prefix, prefixCodes)
	return prefixCodes
}

func createRootNode(f FrequencyTable) *Item {
	pq := NewPriorityQueue(f)
	root := pq.generateNodeTree()
	return &root
}

// encodeHeader constructs a header section that should be
// stored first thing in the compressed file.
// The first 8 bytes stores the lenght of the serialized tree.
// The next 4 bytes stores the added padding.
func encodeHeader(root Item, padding int) []byte {
	serilizedTree := Serialize(&root)

	var treeLength uint64
	treeLength = uint64(len(serilizedTree))

	header := []byte{}
	header = binary.LittleEndian.AppendUint64(header, treeLength)

	padHeader := binary.LittleEndian.AppendUint32([]byte{}, uint32(padding))
	header = append(header, padHeader...)

	header = append(header, serilizedTree...)

	return header
}

func decodeContent(content []byte) (uint32, *Item, []byte, error) {
	var treeLength uint64
	var paddingLength uint32
	err := binary.Read(bytes.NewReader(content[:8]), binary.LittleEndian, &treeLength)
	if err != nil {
		return 0, nil, []byte{}, fmt.Errorf("failed to read binary: %w", err)
	}

	err = binary.Read(bytes.NewReader(content[8:12]), binary.LittleEndian, &paddingLength)
	if err != nil {
		return 0, nil, []byte{}, fmt.Errorf("failed to read binary: %w", err)
	}

	ftStartIdx := 12
	ftEndIdx := treeLength

	treeText := content[ftStartIdx : ftEndIdx+uint64(ftStartIdx)]

	tree := Deserialize(string(treeText))

	return paddingLength, tree, content[ftEndIdx+uint64(ftStartIdx):], nil

}

// decode walks down the tree for each character in the "binary string" 0010001110010010
// we know each character is node when we reach a leaf node.
// then currentNode is reseted to root node
// Walk 0 - 0 = A
// Walk 0 - 1 = B
//
//	     	  0
//				/   \
//		   	  0       1
//			/           \
//		   a              b
//
// -----------------------------
func decode(s string, node *Item, padding int) string {
	s = s[:len(s)-padding]
	var result string

	currentNode := node
	if currentNode.LeftNode == nil || currentNode.RightNode == nil {
		return ""
	}
	for _, v := range s {
		if v == '0' {
			currentNode = currentNode.LeftNode
		} else if v == '1' {
			currentNode = currentNode.RightNode
		}

		if currentNode == nil {
			continue
		}
		if currentNode.LeftNode == nil && currentNode.RightNode == nil {
			result += currentNode.Value
			currentNode = node
		}
	}
	return result
}

// stringToBitString iterates over each character in the string, fetching the prefix code from
// frequency table where the prefix is store for a key
//
// map["a"] = 002
func stringToBitString(in []string, prefixCodes map[string]string) (string, int) {
	bitString := ""

	for _, c := range in {
		bitString += prefixCodes[string(c)]
	}

	// caluculate if padding is needed
	// The bitString must be dividle by 8 (as a byte is 8 bits)
	// Therefor we need to pad the string with zeros in the end
	// The padding must later be removed when decoding
	paddingAdded := 0
	for len(bitString)%8 != 0 {
		bitString += "0"
		paddingAdded++
	}

	return bitString, paddingAdded
}

// bitStringToByteArray takes a string of bits (0001010010) and transforms it to a byte array
func bitStringToByteArray(bitString string) ([]byte, error) {
	lenBits := len(bitString) / 8
	if len(bitString)%8 != 0 {
		return []byte{}, fmt.Errorf("provided string: %s is not valid 8 length", bitString)
	}

	// creates a new byte slice to store results
	out := make([]byte, lenBits)

	for i := 0; i < lenBits; i++ {
		start := i * 8
		end := start + 8

		byteValue := bitString[start:end]

		int, err := strconv.ParseUint(byteValue, 2, 8)
		if err != nil {
			return []byte{}, fmt.Errorf("failed to parse to uint: %w", err)
		}

		out[i] = byte(int)
	}
	return out, nil
}

// byteArrayToBitString takes a byte array and transforms it to a string of bits (0001010010)
// This is used when decoding the file
func byteArrayToBitString(b []byte) string {
	var stringByte string
	for _, char := range b {
		stringByte += fmt.Sprintf("%08b", char)
	}
	return stringByte
}

// getPrefixMap is a recursive function that walks down the tree and stores the prefix code for each character
func getPrefixMap(i *Item, prefix *string, prexfixCodes map[string]string) {
	if i == nil {
		return
	} else {
		zero := fmt.Sprintf("%s0", *prefix)
		one := fmt.Sprintf("%s1", *prefix)
		getPrefixMap(i.LeftNode, &zero, prexfixCodes)
		getPrefixMap(i.RightNode, &one, prexfixCodes)

		if len(i.Value) > 0 {
			prexfixCodes[i.Value] = *prefix
		}
	}
}

// getFreq takes a string and returns a map with the frequency of each character
func getFreq(s []string) FrequencyTable {
	m := make(map[string]int, len(s))
	// Iterate over all chars and add frequency for character to map
	for _, c := range s {
		m[c]++
	}
	return m
}
