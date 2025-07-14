package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// function to clean the string by removing non-alphabetic characters and converting to lowercase
func cleanString(s string) string { 
	cleaned := ""
	for _, char := range s {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') {
			cleaned += string(char)
		}
	}
	return strings.ToLower(cleaned)
}
// function to check if a string is a palindrome
func isPalindrome(s string) bool {
	s = cleanString(s)
	left, right := 0, len(s)-1
	for left < right {
		if s[left] != s[right] {
			return false
		}
		left++
		right--
	}
	return true
}


// main function to test the palindrome check
func main() {
	var input string
	fmt.Println("Enter a string to check if it's a palindrome:")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return	
	}
	input = strings.TrimSpace(input) // Remove any leading/trailing whitespace
	if isPalindrome(input) {
		fmt.Printf("The string '%s' is a palindrome.", input)
	} else {
		fmt.Printf("The string '%s' is not a palindrome.", input)
	}
}