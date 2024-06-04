package monkeylang

import (
	"math/rand"
)

var vowels = []string{"a", "e", "i", "o", "u", "ai", "ei", "oi", "ou", "au"}
var consonants = []string{"b", "c", "d", "f", "g", "h", "j", "k", "l", "m", "n", "p", "qu", "r", "s", "t", "v", "w", "x", "y", "z"}
var consonantClusters = []string{"bl", "br", "ch", "cl", "cr", "dr", "fl", "fr", "gl", "gr", "pl", "pr", "sc", "sh", "sk", "sl", "sm", "sn", "sp", "st", "str", "sw", "th", "tr", "wh", "wr"}
var consonants2 = []string{"b", "c", "d", "f", "g", "h", "j", "k", "l", "m", "n", "p", "qu", "r", "s", "t", "v", "w", "x", "y", "z", "bl", "br", "ch", "cl", "cr", "dr", "fl", "fr", "gl", "gr", "pl", "pr", "sc", "sh", "sk", "sl", "sm", "sn", "sp", "st", "str", "sw", "th", "tr", "wh", "wr"}

var syllablePatterns = []string{"CV", "CVC", "VC"}

func randomElement(choices []string) string {
	return choices[rand.Intn(len(choices))]
}

func generateSyllable() string {
	pattern := randomElement(syllablePatterns)
	syllable := ""
	for _, char := range pattern {
		switch char {
		case 'C':
			syllable += randomElement(consonants2)
		case 'V':
			syllable += randomElement(vowels)
		}
	}
	return syllable
}

func generateWord(syllableCount int) string {
	word := ""
	for i := 0; i < syllableCount; i++ {
		word += generateSyllable()
	}
	return word
}

func GenerateRandomWords(numWords, minSyllables, maxSyllables int) []string {
	words := make([]string, numWords)
	seen := make(map[string]bool)
	i := 0
	for i < numWords {
		syllableCount := rand.Intn(maxSyllables-minSyllables+1) + minSyllables
		tmpWord := generateWord(syllableCount)
		_, saw := seen[tmpWord]
		if saw {
			continue
		}
		words[i] = tmpWord
		i++
		seen[tmpWord] = true
	}
	return words
}

func GenerateRandom_NON_UNIQUE_Words(numWords, minSyllables, maxSyllables int) []string {
	words := make([]string, numWords)
	for i := 0; i < numWords; i++ {
		syllableCount := rand.Intn(maxSyllables-minSyllables+1) + minSyllables
		words[i] = generateWord(syllableCount)
	}
	return words
}

func GenerateSimpleWords(numWords int, unique bool) []string {
	return GenerateRandomWords(numWords, 2, 3)
}

func GenerateStrongWords(numWords int) []string {
	return GenerateRandomWords(numWords, 2, 3)
}
