
# Hangman-Web

## Description
**Hangman-Web** est une version améliorée du projet **Hangman-Classic**, développé en GoLang, qui permet de jouer au jeu du pendu via une interface web. Le projet exploite les fonctionnalités de l'ancien projet en transformant les fonctions principales en module réutilisable, tout en ajoutant une interface utilisateur conviviale et moderne pour jouer directement depuis un navigateur, grâce à un serveur web Go.

## Objectifs
- Réutiliser et transformer les fonctionnalités du projet **Hangman-Classic** en module.
- Offrir une expérience de jeu plus accessible grâce à une interface web interactive.
- Appliquer un design esthétique pour améliorer l'ergonomie du jeu.
- Gérer les interactions utilisateur via le protocole HTTP au lieu de la console.

---

## Fonctionnalités
- **Modularité** : Les fonctions principales du projet **Hangman-Classic** (choix des mots, gestion des tentatives, vérification des lettres) sont encapsulées dans un module appelé `Hangmanclassic`, permettant leur réutilisation dans différents projets.
- **Interface web** :
  - Une page design pour afficher le mot à deviner, les lettres incorrectes, et le nombre de tentatives restantes.
  - Un formulaire interactif pour entrer des lettres et vérifier leur présence dans le mot.
- **Gestion des états du jeu** :
  - Affichage du mot masqué avec les lettres révélées.
  - Affichage des lettres incorrectes et des tentatives restantes.
  - Messages spécifiques pour indiquer la victoire ou la défaite.
- **Server-Side Rendering** : Le jeu est géré entièrement côté serveur en Go, avec des modèles HTML.

---

## Fonctionnement
1. **Initialisation** : Le serveur Go initialise le jeu en sélectionnant un mot aléatoire depuis le fichier `words.txt` et en réinitialisant les variables du jeu.
2. **Interactions utilisateur** : L'utilisateur propose des lettres via un formulaire web.
3. **Mécanique de jeu** : 
   - Si la lettre est correcte, elle est révélée dans le mot.
   - Si la lettre est incorrecte, une tentative est décomptée.
4. **Fin de partie** : Le jeu se termine lorsque le mot est entièrement découvert (victoire) ou lorsque toutes les tentatives sont épuisées (défaite).

---

## Architecture
### 1. Modules et fichiers principaux :
- **Hangman-Web** (Dossier principal) :
  - `hangmanweb.go` : Contient le code du serveur web et la logique principale du jeu.
  - `Templates/hangman.tmpl` : Modèle HTML pour l'interface utilisateur.
  - `css/style.css` : Feuille de style pour le design.
- **Hangman-Classic (module)** :
  - `hangman.go` : Contient les fonctions réutilisées (choix des mots, gestion des lettres, vérification des états).
  - `words.txt` : Liste des mots utilisés pour le jeu.

### 2. Serveur Web
Le serveur web est implémenté en Go avec les routes suivantes :
- `/` : Affiche la page principale du jeu.
- `/hangman` : Traite les propositions de lettres.
- `/restart` : Réinitialise le jeu avec un nouveau mot.

---

## Installation
1. Clonez le dépôt :
   ```bash
   git clone https://github.com/votre-utilisateur/hangman-web.git
   cd hangman-web
   ```

2. Assurez-vous que Go est installé sur votre système. Configurez le module :
   ```bash
   go mod tidy
   ```

3. Lancez le serveur :
   ```bash
   go run ./hangmanweb.go
   ```

4. Accédez au jeu via votre navigateur à l'adresse : `http://localhost:8080`.

---

## Lien avec Hangman-Classic
- Le projet **Hangman-Classic** a été transformé en un module Go (`Hangmanclassic`) contenant les fonctionnalités principales :
  - Sélection aléatoire d'un mot.
  - Gestion des lettres révélées et incorrectes.
  - Vérification des conditions de victoire ou de défaite.
- Ces fonctions sont appelées dans le fichier `hangmanweb.go` de **Hangman-Web**, permettant de réutiliser le code existant tout en ajoutant une nouvelle couche fonctionnelle (interface web).

---

## Auteur
- Projet réalisé par Accoows & Willpns dans le cadre de l'évolution du projet **Hangman-Classic**.
