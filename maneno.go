package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"sort"
	"strings"
	"time"
)

const WORD_LENGTH = 5
const MAX_GUESSES = 6

func get_filled_color(color string) [WORD_LENGTH]string {
	color_vector := [WORD_LENGTH]string{}
	for i := range color_vector {
		color_vector[i] = color
	}

	return color_vector
}

func display_word(word string, color_vector [WORD_LENGTH]string) {
	for i, c := range word {
		switch color_vector[i] {
		case "Green":
			fmt.Print("\033[42m\033[1;30m")
		case "Yellow":
			fmt.Print("\033[43m\033[1;30m")
		case "Grey":
			fmt.Print("\033[40m\033[1;37m")
		}

		fmt.Printf(" %c ", c)
		fmt.Print("\033[m\033[m")
	}

	fmt.Println()
}

func main() {
	fmt.Println(`

█▀▀ ▄▀█ █▀▄▀█ █▀▀   █░░ ▄▀█   █▀▄▀█ ▄▀█ █▄░█ █▀▀ █▄░█ █▀█
█▄█ █▀█ █░▀░█ ██▄   █▄▄ █▀█   █░▀░█ █▀█ █░▀█ ██▄ █░▀█ █▄█
`)
	fmt.Println("Kisia neno la kiswahili lenye herufi 5.\n")
	rand.Seed(time.Now().Unix())

	body, err := ioutil.ReadFile("maneno.txt")
	if err != nil {
		log.Fatalln(err)
	}
	words := strings.Split(string(body), "\n")

	game_words := []string{}

	for _, word := range words {
		if len(word) == WORD_LENGTH {
			game_words = append(game_words, strings.ToUpper(word))
		}
	}

	sort.Strings(game_words)

	selected_word := game_words[rand.Intn(len(game_words))]

	reader := bufio.NewReader(os.Stdin)
	guesses := []map[string][WORD_LENGTH]string{}
	var guess_count int
	for guess_count = 0; guess_count < MAX_GUESSES; guess_count++ {
		fmt.Printf("Kisio la (%d/%d): ", guess_count+1, MAX_GUESSES)
		guess_word, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalln(err)
		}

		guess_word = strings.ToUpper(guess_word[:len(guess_word)-1])

		if guess_word == selected_word {
			fmt.Println("Umepatia!!")
			color_vector := get_filled_color("Green")

			guesses = append(guesses, map[string][WORD_LENGTH]string{guess_word: color_vector})

			for _, guess := range guesses {
				for guess_word, color_vector := range guess {
					display_word(guess_word, color_vector)
				}
			}

			break
		} else {
			i := sort.SearchStrings(game_words, guess_word)

			if i < len(game_words) && game_words[i] == guess_word {
				color_vector := get_filled_color("Grey")

				yellow_lock := [WORD_LENGTH]bool{}

				for j, guess_letter := range guess_word {
					for k, letter := range selected_word {
						if guess_letter == letter && j == k {
							color_vector[j] = "Green"

							yellow_lock[k] = true
							break
						}
					}
				}

				for j, guess_letter := range guess_word {
					for k, letter := range selected_word {
						if guess_letter == letter && color_vector[j] != "Green" && yellow_lock[k] == false {
							color_vector[j] = "Yellow"
							yellow_lock[k] = true
						}
					}
				}

				guesses = append(guesses, map[string][WORD_LENGTH]string{guess_word: color_vector})
				display_word(guess_word, color_vector)
			} else {
				guess_count--
				fmt.Printf("Ingiza neno lenye herufi %d\n", WORD_LENGTH)
			}
		}
	}

	if guess_count == MAX_GUESSES {
		color_vector := get_filled_color("Green")
		fmt.Print("Neno sahihi ilikua: ")
		display_word(selected_word, color_vector)
		fmt.Println("Kalale tu ujaribu kesho")
	}
}
