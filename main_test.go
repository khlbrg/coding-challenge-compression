package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestBitStringToByte(t *testing.T) {

	input := "000000110000000100000000"
	result, _ := bitStringToByteArray(input)

	var stringByte string
	for _, char := range result {
		stringByte += fmt.Sprintf("%08b", char)
	}

	if input != stringByte {
		t.Errorf("failed: expected: %s, got: %s", input, stringByte)
	}
}

func TestBitStringToByteError(t *testing.T) {
	input := "0000001100000001000000" // incorret numver of bits should return error
	_, err := bitStringToByteArray(input)
	if err == nil {
		t.Errorf("failed: expected: %s, got: %s", input, err)
	}
}

func TestGetFrequenceTable(t *testing.T) {
	got := getFreq([]string{"a", "b", "c", "a", "b", "a", " "})
	want := FrequencyTable{
		"a": 3,
		"b": 2,
		"c": 1,
		" ": 1,
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v, got: %v", want, got)
	}
}

func TestPriorityQueue(t *testing.T) {
	freqTable := getFreq([]string{"a", "b", "c", "a", "b", "a", " "})

	pq := NewPriorityQueue(freqTable)
	root := pq.generateNodeTree()

	if root.LeftNode.Value != "a" ||
		root.RightNode.LeftNode.Value != "b" ||
		root.RightNode.RightNode.RightNode.Value != "c" ||
		root.RightNode.RightNode.LeftNode.Value != " " {
		t.Fatalf("failed")
	}
}

func TestTreeSerialize(t *testing.T) {
	testItem := Item{
		LeftNode: &Item{
			LeftNode: &Item{
				Value: "B",
			},
			RightNode: &Item{
				Value: "C",
			},
		},
		RightNode: &Item{
			Value: "A",
		},
	}
	serializedRoot := Serialize(&testItem)

	deseralizedRoot := Deserialize(serializedRoot)

	if testItem.LeftNode.RightNode.Value != deseralizedRoot.LeftNode.RightNode.Value ||
		testItem.LeftNode.LeftNode.Value != deseralizedRoot.LeftNode.LeftNode.Value ||
		testItem.RightNode.Value != deseralizedRoot.RightNode.Value {
		t.Errorf("failed")
	}
}

func TestCompressionAndDecompression(t *testing.T) {
	input := "Hello World. This is a #test, string. I hope it works. 😱"

	compressed, err := Compress(input)
	if err != nil {
		t.Errorf("failed to compress: %v", err)
	}

	decompressed, err := Decompress(compressed)
	if err != nil {
		t.Errorf("failed to decompress: %v", err)
	}

	if decompressed != string(input) {
		t.Errorf("failed: expected: %s, got: %s", input, decompressed)
	}
	t.Logf("Input: %s \n", input)
	t.Logf("Output: %s \n", decompressed)
}

func TestPrefixCodes(t *testing.T) {
	testItem := Item{
		LeftNode: &Item{
			LeftNode: &Item{
				Value: "B",
			},
			RightNode: &Item{
				Value: "C",
			},
		},
		RightNode: &Item{
			Value: "A",
		},
	}

	got := prefixCodes(&testItem)
	want := map[string]string{
		"A": "1",
		"B": "00",
		"C": "01",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v, got: %v", want, got)
	}
}
