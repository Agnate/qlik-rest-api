package util

import (
	"testing"
)

func TestIsPalindromeTrue(t *testing.T) {
	test := "radar"
	want := true
	result := IsPalindrome(test)
	if result != want {
		t.Fatalf(`IsPalindrome("%s") = %t, should be %t, nil`, test, result, want)
	}
}

func TestIsPalindromeFalse(t *testing.T) {
	test := "sword"
	want := false
	result := IsPalindrome(test)
	if result != want {
		t.Fatalf(`IsPalindrome("%s") = %t, should be %t, nil`, test, result, want)
	}
}

func TestIsPalindromeNil(t *testing.T) {
	var test string
	want := false
	result := IsPalindrome(test)
	if result != want {
		t.Fatalf(`IsPalindrome("%s") = %t, should be %t, nil`, test, result, want)
	}
}

func TestIsPalindromeEmpty(t *testing.T) {
	test := ""
	want := false
	result := IsPalindrome(test)
	if result != want {
		t.Fatalf(`IsPalindrome("%s") = %t, should be %t, nil`, test, result, want)
	}
}

func TestIsPalindromeSingle(t *testing.T) {
	test := "a"
	want := true
	result := IsPalindrome(test)
	if result != want {
		t.Fatalf(`IsPalindrome("%s") = %t, should be %t, nil`, test, result, want)
	}
}

func TestIsPalindromeDoubleTrue(t *testing.T) {
	test := "aa"
	want := true
	result := IsPalindrome(test)
	if result != want {
		t.Fatalf(`IsPalindrome("%s") = %t, should be %t, nil`, test, result, want)
	}
}

func TestIsPalindromeDoubleFalse(t *testing.T) {
	test := "ab"
	want := false
	result := IsPalindrome(test)
	if result != want {
		t.Fatalf(`IsPalindrome("%s") = %t, should be %t, nil`, test, result, want)
	}
}

func TestIsPalindromeAccents(t *testing.T) {
	test := "radär"
	want := false
	result := IsPalindrome(test)
	if result != want {
		t.Fatalf(`IsPalindrome("%s") = %t, should be %t, nil`, test, result, want)
	}
}

func TestIsPalindromeCapitals(t *testing.T) {
	test := "Radar"
	want := true
	result := IsPalindrome(test)
	if result != want {
		t.Fatalf(`IsPalindrome("%s") = %t, should be %t, nil`, test, result, want)
	}
}

func TestIsPalindromeSpaces(t *testing.T) {
	test := "  too hot to  hoot  " // multiple leading, in-between, and trailing spaces
	want := true
	result := IsPalindrome(test)
	if result != want {
		t.Fatalf(`IsPalindrome("%s") = %t, should be %t, nil`, test, result, want)
	}
}
