package main

import (
	"fmt"
	"hangman_web/Hangmanclassic"
	"html/template"
	"log"
	"net/http"
	"strings"
)

var word string
var revealed []bool
var attemptsLeft int
var lettersRevealed map[rune]bool
var gameOver bool
var errorMessage string

func initializeHangman() {
	// Initialisation du jeu
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
		WordDisplay  string
		AttemptsLeft int
		GameOver     bool
		ErrorMessage string
	}{
		WordDisplay:  wordDisplay,
		AttemptsLeft: attemptsLeft,
		GameOver:     gameOver,
		ErrorMessage: errorMessage,
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
	if letter == "" || len(letter) != 1 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	revealedRune := rune(strings.ToUpper(letter)[0])
	if lettersRevealed[revealedRune] {
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
