package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// function to count the frequency of each word in a given text
func countWordFrequency(text string) map[string]int {
	wordFrequency := make(map[string]int)
	words := strings.Fields(text)

	for _, word := range words {
		word = removePunctuation(word)
		word = strings.ToLower(word)
		if word != "" {
			wordFrequency[word]++
		}
	}
	return wordFrequency
}

// function to remove punctuation from a word
func removePunctuation(word string) string {
	var cleanedWord strings.Builder
	for _, char := range word {
		if strings.ContainsRune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", char) {
			cleanedWord.WriteRune(char)
		}
	}
	return cleanedWord.String()
}

func main() {
	fmt.Println("Enter text:")
	reader := bufio.NewReader(os.Stdin)

	text, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}
	wordFrequency := countWordFrequency(text)

	for word, count := range wordFrequency {
		fmt.Printf("%s: %d\n", word, count)
	}
}
