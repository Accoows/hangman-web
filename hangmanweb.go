package main

import (
	"fmt"
	"hangman_web/Hangmanclassic"
	"html/template"
	"log"
	"net/http"
	"strings"
	"unicode"
)

var word string
var revealed []bool
var attemptsLeft int
var lettersRevealed map[rune]bool
var gameOver bool
var errorMessage string
var incorrectLetters []string

func initializeHangman() {
	words, err := Hangmanclassic.ReadFileLines("words.txt", false)
	if err != nil {
		fmt.Println("Erreur lors de la lecture du fichier :", err)
		return
	}
	word, revealed = Hangmanclassic.FindWord(words)
	attemptsLeft = Hangmanclassic.Maxtentative
	lettersRevealed = make(map[rune]bool)
	gameOver = false
	errorMessage = ""
	incorrectLetters = []string{}
}

func renderTemplate(w http.ResponseWriter, data interface{}) {
	t, err := template.ParseFiles("Templates/hangman.tmpl")
	if err != nil {
		http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
		return
	}
	t.Execute(w, data)
}

func Home(w http.ResponseWriter, r *http.Request) {
	wordDisplay := Hangmanclassic.DisplayWord(word, revealed, ' ')
	data := struct {
		WordDisplay    string
		AttemptsLeft   int
		GameOver       bool
		ErrorMessage   string
		GuessedLetters string
	}{
		WordDisplay:    wordDisplay,
		AttemptsLeft:   attemptsLeft,
		GameOver:       gameOver,
		ErrorMessage:   errorMessage,
		GuessedLetters: strings.Join(incorrectLetters, ", "),
	}
	renderTemplate(w, data)
}

func Hangman(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if gameOver {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	letter := r.FormValue("letter")
	if letter == "" || len(letter) != 1 || !unicode.IsLetter(rune(letter[0])) {
		errorMessage = "Veuillez entrer une lettre valide !"
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	revealedRune := rune(strings.ToUpper(letter)[0])
	if lettersRevealed[revealedRune] {
		errorMessage = "Vous avez déjà proposé cette lettre !"
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	lettersRevealed[revealedRune] = true

	found := false
	for i, ltr := range word {
		if ltr == revealedRune {
			revealed[i] = true
			found = true
		}
	}

	if !found {
		attemptsLeft--
		incorrectLetters = append(incorrectLetters, string(revealedRune))
		errorMessage = "La lettre n'est pas dans le mot !"
	} else {
		errorMessage = ""
	}

	if attemptsLeft <= 0 || Hangmanclassic.AllRevealed(revealed) {
		gameOver = true
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func Restart(w http.ResponseWriter, r *http.Request) {
	initializeHangman()
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func main() {
	initializeHangman()

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.HandleFunc("/", Home)
	http.HandleFunc("/hangman", Hangman)
	http.HandleFunc("/restart", Restart)

	port := "8080"
	fmt.Printf("Serveur démarré sur http://localhost:%s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Échec du démarrage du serveur : %v", err)
	}
}
