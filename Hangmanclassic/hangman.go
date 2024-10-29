package Hangmanclassic

import (
	"bufio"
	"math/rand"
	"os"
	"strings"
	"unicode"
)

var Maxtentative int = 10
var HangmanStages []string

// Fonction pour lire les lignes d'un fichier et du hangman.
func ReadFileLines(filename string, isHangmanStages bool) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)

	if isHangmanStages {
		var stage strings.Builder // Builder pour accumuler les lignes
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" && stage.Len() > 0 {
				lines = append(lines, stage.String())
				stage.Reset()
			} else {
				stage.WriteString(line + "\n")
			}
		}
		if stage.Len() > 0 {
			lines = append(lines, stage.String())
		}
	} else {
		for scanner.Scan() {
			lines = append(lines, strings.ToUpper(scanner.Text()))
		}
	}

	return lines, scanner.Err()
}

// Choisir un mot et révéler quelques lettres
func FindWord(words []string) (string, []bool) {
	word := words[rand.Intn(len(words))]
	revealed := make([]bool, len(word))

	indicesRevealed := make(map[int]bool) // Map pour stocker les lettres à révéler
	lettersCount := len(word)/2 - 1

	for len(indicesRevealed) < lettersCount {
		randomIndex := rand.Intn(len(word))
		char := unicode.ToUpper(rune(word[randomIndex]))
		for i, c := range word {
			if unicode.ToUpper(c) == char {
				indicesRevealed[i] = true
			}
		}
	}
	for i := range word {
		if indicesRevealed[i] {
			revealed[i] = true
		}
	}
	return word, revealed
}

// Fonction pour afficher le mot en tant que chaîne
func DisplayWord(word string, revealed []bool, lastCorrectLetter rune) string {
	var display strings.Builder
	for i, letter := range word {
		if revealed[i] {
			if letter == lastCorrectLetter {
				display.WriteRune(letter) // Ajoute la lettre révélée
			} else {
				display.WriteRune(letter)
			}
		} else {
			display.WriteRune('_') // Affiche "_" pour les lettres masquées
		}
		display.WriteRune(' ')
	}
	return display.String()
}

// Fonction pour vérifier si toutes les lettres sont révélées
func AllRevealed(revealed []bool) bool {
	for _, r := range revealed {
		if !r {
			return false
		}
	}
	return true
}
