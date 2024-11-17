package utils

import "testing"

func TestIsDigit(t *testing.T) {

	for i := 0; i < 10; i++ {
		if !IsDigit(byte(i) + '0') {
			t.Errorf("Expected %d to be a digit", i)
		}
	}
	for i := 0; i < 26; i++ {
		if IsDigit(byte(i) + 'a') {
			t.Errorf("Expected %c to not be a digit", byte(i) + 'a')
		}
	}
}

func TestIsAlpha(t *testing.T) {
	
	for i := 0; i < 26; i++ {
		if !IsAlpha(byte(i) + 'a') {
			t.Errorf("Expected %c to be an alpha", byte(i) + 'a')
		}
	}
	for i := 0; i < 26; i++ {
		if !IsAlpha(byte(i) + 'A') {
			t.Errorf("Expected %c to be an alpha", byte(i) + 'A')
		}
	}
	for i := 0; i < 10; i++ {
		if IsAlpha(byte(i) + '0') {
			t.Errorf("Expected %d to not be an alpha", i)
		}
	}
}