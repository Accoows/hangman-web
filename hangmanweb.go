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
	}{
		WordDisplay:  wordDisplay,
		AttemptsLeft: attemptsLeft,
	}
	renderTemplate(w, data)
}

func Hangman(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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
	}

	if attemptsLeft <= 0 || Hangmanclassic.AllRevealed(revealed) {
		initializeHangman() // Réinitialiser le jeu si toutes les lettres sont trouvées ou si les tentatives sont épuisées
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func main() {
	initializeHangman()

	http.HandleFunc("/", Home)
	http.HandleFunc("/hangman", Hangman)

	port := "8080"
	fmt.Printf("Serveur démarré sur http://localhost:%s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Échec du démarrage du serveur : %v", err)
	}
}
