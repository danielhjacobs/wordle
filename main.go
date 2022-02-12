package main

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gookit/color"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

func onlyCapitalLetters(s string) bool {
	for _, r := range s {
		if r < 'A' || r > 'Z' {
			return false
		}
	}
	return true
}

func scanFileLineByLine(db *sql.DB, filename string, doSomethingWithScannedData func(db *sql.DB, scanner *bufio.Scanner)) {
	file, err := os.Open(filename)
	fatalError(err)
	scanner := bufio.NewScanner(file)
	doSomethingWithScannedData(db, scanner)
	defer file.Close()
}

func createTable(db *sql.DB) {
	createWordsTableSQL := `CREATE TABLE five_letter_words (
		word_id INTEGER PRIMARY KEY,
		word TEXT NOT NULL,
		linux_word_list BOOLEAN NOT NULL DEFAULT TRUE,
		wordle_word_list BOOLEAN NOT NULL DEFAULT FALSE,
		wordle_guess_list BOOLEAN NOT NULL DEFAULT FALSE
	);`
	log.Println("Create five_letter_words table")
	statement, err := db.Prepare(createWordsTableSQL)
	fatalError(err)
	statement.Exec()
	log.Println("five_letter_words table created")
}

func addToSqliteDatabase(db *sql.DB, word string) {
	if !wordExists(db, word) {
		log.Println("Inserting: " + word)
		insertWordSQL := `INSERT INTO five_letter_words(word) VALUES (?)`
		statement, err := db.Prepare(insertWordSQL)
		fatalError(err)
		_, err = statement.Exec(word)
		fatalError(err)
	}
}

func fileFound(name string) bool {
	_, err := os.Stat(name)
	if err == nil {
		return true
	}
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return false
}

func createSqlite3Db() {
	log.Println("Creating words.db")
	file, err := os.Create("words.db")
	fatalError(err)
	file.Close()
	log.Println("words.db created")
}

func wordExists(db *sql.DB, word string) bool {
	row := db.QueryRow("select word from five_letter_words where word= ?", word)
	temp := ""
	err := row.Scan(&temp)
	if err != nil && err != sql.ErrNoRows {
		fatalError(err)
	}
	if temp != "" {
		return true
	}
	return false
}

func updateDatabase(db *sql.DB, word string) {
	if !wordExists(db, word) {
		log.Println(word + " is not in database")
		log.Println("Inserting: " + word)
		insertWordSQL := `INSERT INTO five_letter_words(word, linux_word_list, wordle_word_list, wordle_guess_list) VALUES (?, ?, ?, ?)`
		statement, err := db.Prepare(insertWordSQL)
		fatalError(err)
		_, err = statement.Exec(word, 0, 1, 1)
		fatalError(err)
	} else {
		//log.Println("Updating: " + word)
		updateWordSQL := `UPDATE five_letter_words SET wordle_word_list = 1, wordle_guess_list = 1 WHERE word = ?`
		statement, err := db.Prepare(updateWordSQL)
		fatalError(err)
		_, err = statement.Exec(word)
		fatalError(err)
	}
}

func updateMoreDatabase(db *sql.DB, word string) {
	if !wordExists(db, word) {
		log.Println(word + " is not in database")
		log.Println("Inserting: " + word)
		insertWordSQL := `INSERT INTO five_letter_words(word, linux_word_list, wordle_guess_list) VALUES (?, ?, ?)`
		statement, err := db.Prepare(insertWordSQL)
		fatalError(err)
		_, err = statement.Exec(word, 0, 1)
		fatalError(err)
	} else {
		//log.Println("Updating: " + word)
		updateWordSQL := `UPDATE five_letter_words SET wordle_guess_list = 1 WHERE word = ?`
		statement, err := db.Prepare(updateWordSQL)
		fatalError(err)
		_, err = statement.Exec(word)
		fatalError(err)
	}
}

func fatalError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func addScannedWords(db *sql.DB, scanner *bufio.Scanner) {
	for scanner.Scan() {
		word := scanner.Text()
		if len(word) == 5 {
			word = strings.ToUpper(word)
			if onlyCapitalLetters(word) {
				addToSqliteDatabase(db, word)
			}
		}
	}
	err := scanner.Err()
	fatalError(err)
}

func updateScannedWords(db *sql.DB, scanner *bufio.Scanner) {
	for scanner.Scan() {
		wordle_word := scanner.Text()
		wordle_word = strings.ToUpper(wordle_word)
		updateDatabase(db, wordle_word)
	}
	err := scanner.Err()
	fatalError(err)
}

func updateMoreScannedWords(db *sql.DB, scanner *bufio.Scanner) {
	for scanner.Scan() {
		wordle_word := scanner.Text()
		wordle_word = strings.ToUpper(wordle_word)
		updateMoreDatabase(db, wordle_word)
	}
	err := scanner.Err()
	fatalError(err)
}

func getAllWords(db *sql.DB) []string {
	var wordList []string
	row, err := db.Query("SELECT word FROM five_letter_words ORDER BY word")
	fatalError(err)
	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		var word string
		row.Scan(&word)
		wordList = append(wordList, word)
	}
	return wordList
}

func getWordleGuessWords(db *sql.DB) []string {
	var wordList []string
	row, err := db.Query("SELECT word FROM five_letter_words WHERE wordle_guess_list=1 ORDER BY word")
	fatalError(err)
	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		var word string
		row.Scan(&word)
		wordList = append(wordList, word)
	}
	return wordList
}

func getLinuxWords(db *sql.DB) []string {
	var wordList []string
	row, err := db.Query("SELECT word FROM five_letter_words WHERE linux_word_list=1 ORDER BY word")
	fatalError(err)
	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		var word string
		row.Scan(&word)
		wordList = append(wordList, word)
	}
	return wordList
}

func getWordleWords(db *sql.DB) []string {
	var wordList []string
	row, err := db.Query("SELECT word FROM five_letter_words WHERE wordle_word_list=1 ORDER BY word")
	fatalError(err)
	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		var word string
		row.Scan(&word)
		wordList = append(wordList, word)
	}
	return wordList
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// I'm not 100% certain about this function working correctly, but I think it works well enough at least
func getColors(targetWord, guess string) [5]rune {
	var colors [5]rune
	var matchesFoundInGuess int
	var targetLetter byte
	var letterCopiesInTargetWord int
	for i := 0; i < 5; i++ {
		targetLetter = targetWord[i]
		matchesFoundInGuess = 0
		letterCopiesInTargetWord = strings.Count(targetWord, string(targetLetter))
		for j := 0; j < 5; j++ {
			if targetLetter == guess[j] {
				matchesFoundInGuess = matchesFoundInGuess + 1
				if matchesFoundInGuess > letterCopiesInTargetWord {
					break
				}
				if i == j {
					colors[j] = 'G'
				} else {
					colors[j] = 'Y'
				}

			}
		}
	}
	return colors
}

func main() {
	var sqliteDb *sql.DB
	if !fileFound("./words.db") {
		createSqlite3Db()
		sqliteDb, _ = sql.Open("sqlite3", "./words.db")
		defer sqliteDb.Close()
		createTable(sqliteDb)
		scanFileLineByLine(sqliteDb, "./linux_word_list.txt", addScannedWords)
		scanFileLineByLine(sqliteDb, "./wordle-answers-alphabetical.txt", updateScannedWords)
		scanFileLineByLine(sqliteDb, "./wordle-allowed-guesses.txt", updateMoreScannedWords)
	} else {
		sqliteDb, _ = sql.Open("sqlite3", "./words.db")
	}
	// Limit to only wordle's own answers
	// wordList := getWordleWords(sqliteDb)
	// Allow all wordle words as answers
	wordList := getWordleGuessWords(sqliteDb)
	// Allow all wordle words as guesses
	fullWordList := getWordleGuessWords(sqliteDb)
	// Allow any word as a guess (words include all wordle guesses, all linux words, and all wordle answers)
	// fullWordList := getAllWords(sqliteDb)
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	targetWord := wordList[r1.Intn(len(wordList))]
	green := color.FgGreen.Render
	yellow := color.FgYellow.Render
	var guess string
	var colors [6][5]rune
	var guess_valid bool
	for guess_number := 0; guess_number < 6; guess_number++ {
		fmt.Println("Enter guess " + string('0'+guess_number+1))
		guess_valid = false
		for !guess_valid {
			fmt.Scanln(&guess)
			if len(guess) != 5 {
				fmt.Println("Guess must contain exactly five letters")
				continue
			}
			guess = strings.ToUpper(guess)
			if !onlyCapitalLetters(guess) {
				fmt.Println("Guess must only contain letters [A-Za-z]")
				continue
			}
			if !stringInSlice(guess, fullWordList) {
				fmt.Println("Guess not in our dictionary")
				continue
			}
			guess_valid = true
		}
		if guess == targetWord {
			fmt.Println("That was the target word")
			break
		}
		colors[guess_number] = getColors(targetWord, guess)
		for i := 0; i < 5; i++ {
			if colors[guess_number][i] == 'G' {
				print(green(string(guess[i])))
			} else if colors[guess_number][i] == 'Y' {
				print(yellow(string(guess[i])))
			} else {
				print(string(guess[i]))
			}
		}
		print("\n")
	}
	print("Target Word: " + targetWord + "\n")
}
