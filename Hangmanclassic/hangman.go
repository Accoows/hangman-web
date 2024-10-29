package Hangmanclassic

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"unicode"

	"github.com/fatih/color"
)

var maxtentative int = 10
var hangmanStages []string

// Fonction pour lire les lignes d'un fichier et du hangman.
func readFileLines(filename string, isHangmanStages bool) ([]string, error) {
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
func findWord(words []string) (string, []bool) {
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

// Fonction pour afficher le mot
func displayWord(word string, revealed []bool, lastCorrectLetter rune) {
	for i, letter := range word {
		if revealed[i] {
			if letter == lastCorrectLetter { // Vérifie si la lettre est la dernière correcte
				color.Set(color.FgGreen) // Change la couleur de la lettre
				fmt.Printf("%c ", letter)
				color.Unset()
			} else {
				fmt.Printf("%c ", letter)
			}
		} else {
			fmt.Print("_ ")
		}
	}
	fmt.Println()
}

// Fonction pour vérifier si toutes les lettres sont révélées
func allRevealed(revealed []bool) bool {
	for _, r := range revealed {
		if !r {
			return false
		}
	}
	return true
}

func main() {
	var err error
	if hangmanStages, err = readFileLines("hangman.txt", true); err != nil {
		fmt.Println("Erreur lors de la lecture des étapes :", err)
		return
	}

	words, err := readFileLines("words.txt", false)
	if err != nil {
		fmt.Println("Erreur lors de la lecture du fichier :", err)
		return
	}

	word, revealed := findWord(words) //Choix du mot aléatoire
	tentatives := maxtentative
	lettersRevealed := make(map[rune]bool) //Trace des lettres déjà révélées
	var lastCorrectLetter rune

	fmt.Println("Bonne chance, vous avez", maxtentative, "tentatives.")
	displayWord(word, revealed, lastCorrectLetter)

	user := bufio.NewReader(os.Stdin)
	for tentatives > 0 {
		fmt.Print("Choisissez une lettre : ")
		input, _ := user.ReadString('\n')
		letter := strings.TrimSpace(strings.ToUpper(input)) // Enlève les espaces au cas ou

		if len(letter) != 1 {
			fmt.Println("Veuillez saisir une seule lettre.")
			continue
		}

		revealedRune := rune(letter[0])
		if !unicode.IsLetter(revealedRune) {
			fmt.Println("Ce n'est pas une lettre.")
			continue
		}

		if lettersRevealed[revealedRune] {
			fmt.Println("Vous avez déjà donné cette lettre.")
			continue
		}
		lettersRevealed[revealedRune] = true

		wordFound := false
		for i, ltr := range word {
			if unicode.ToUpper(ltr) == revealedRune {
				revealed[i] = true
				wordFound = true
				lastCorrectLetter = ltr
			}
		}

		if wordFound {
			if allRevealed(revealed) {
				fmt.Println("Bravo ! Le mot est :")
				displayWord(word, revealed, lastCorrectLetter)
				return
			}
		} else {
			tentatives-- // Décrémente le nombre de tentatives
			fmt.Printf("Pas présent dans le mot, il vous reste %d tentatives.\n", tentatives)
			if maxtentative-tentatives-1 < len(hangmanStages) {
				fmt.Println(hangmanStages[maxtentative-tentatives-1]) // Affiche l'étape du pendu
			}
		}

		if tentatives == 0 {
			for i := range revealed {
				revealed[i] = true
			}
			fmt.Println("Perdu ! Le mot était :")
			displayWord(word, revealed, lastCorrectLetter)
			return
		}

		displayWord(word, revealed, lastCorrectLetter)
	}
}
