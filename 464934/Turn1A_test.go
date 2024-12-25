package main

import "testing"

func doubleSlice(input []int) []int {
	var result []int
	for _, value := range input {
		result = append(result, value*2)
	}
	return result
}

func TestDoubleSlice(t *testing.T) {
	// Test case 1: Basic functionality with positive integers
	input1 := []int{1, 2, 3, 4, 5}
	expected1 := []int{2, 4, 6, 8, 10}
	output1 := doubleSlice(input1)
	if len(output1) != len(expected1) {
		t.Errorf("TestCase 1: Length mismatch. Output length = %d, Expected length = %d", len(output1), len(expected1))
	}
	for i, v := range expected1 {
		if output1[i] != v {
			t.Errorf("TestCase 1: Expected %d at index %d, but got %d", v, i, output1[i])
		}
	}

	// Test case 2: Zero input slice
	input2 := []int{}
	expected2 := []int{}
	output2 := doubleSlice(input2)
	if len(output2) != len(expected2) {
		t.Errorf("TestCase 2: Length mismatch. Output length = %d, Expected length = %d", len(output2), len(expected2))
	}

	// Test case 3: Slice with negative integers
	input3 := []int{-1, -2, -3, -4, -5}
	expected3 := []int{-2, -4, -6, -8, -10}
	output3 := doubleSlice(input3)
	if len(output3) != len(expected3) {
		t.Errorf("TestCase 3: Length mismatch. Output length = %d, Expected length = %d", len(output3), len(expected3))
	}
	for i, v := range expected3 {
		if output3[i] != v {
			t.Errorf("TestCase 3: Expected %d at index %d, but got %d", v, i, output3[i])
		}
	}
}
