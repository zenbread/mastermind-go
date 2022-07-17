package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type Guess struct {
	Num   []int
	Color []Color
}

func NewGuess(sz int) *Guess {
	nums := make([]int, sz)
	col := make([]Color, sz)
	for i := 0; i < sz; i++ {
		col[i] = Wrong
	}
	guess := Guess{nums, col}
	return &guess
}

type Color int

const (
	Wrong Color = iota
	White
	Right
)

type stats struct {
	timeTotal      time.Duration
	fastestTime    time.Duration
	longestTime    time.Duration
	played         int
	guesses        int
	firstGuessWins int
}

type gamestat struct {
	guesses       int
	time          time.Duration
	firstGuessWin int
}

func assignColors(guess *Guess, password []int, sz int) {
	marked := make([]bool, sz)
	for i := 0; i < sz; i++ {
		if guess.Num[i] == password[i] {
			guess.Color[i] = Right
			marked[i] = true
		}
	}
	for i := 0; i < 5; i++ {
		for ind, x := range password {
			if !marked[ind] && guess.Color[i] != Right && guess.Num[i] == x {
				marked[ind] = true
				guess.Color[i] = White
				break
			}
		}
	}
}

func getPassword(sz int) []int {
	rand.Seed(time.Now().UnixNano())
	password := make([]int, sz)

	for i := 0; i < sz; i++ {
		password[i] = rand.Intn(9) + 1
	}

	return password
}

func getUserGuess(s *bufio.Scanner, guess *Guess, sz int) {
	errorMSG := "Please enter %d numbers separated by spaces. E.G.: 1 2 3 4 5\n"
	for {
		fmt.Print("Please enter your guess:> ")
		s.Scan()
		values := strings.Split(s.Text(), " ")
		if len(values) != sz {
			fmt.Printf(errorMSG, sz)
			continue
		}

		redo := false
		for i := 0; i < sz; i++ {
			if val, err := strconv.Atoi(values[i]); err != nil {
				fmt.Printf(errorMSG, sz)
				if val == 0 {
					fmt.Printf("Not a number (%v) - Try again\n", val)
				}
				redo = true
				break
			} else {
				guess.Num[i] = val
			}
		}
		if !redo {
			break
		}
	}
}

func checkWin(guess *Guess, sz int) bool {
	for i := 0; i < sz; i++ {
		if guess.Color[i] != Right {
			return false
		}
	}
	return true
}

func printHistory(guesses []*Guess, sz int) {
	for _, guess := range guesses {
		for i := 0; i < sz; i++ {
			switch guess.Color[i] {
			case Right:
				fmt.Print("✅ ")
			case White:
				fmt.Print("➖ ")
			case Wrong:
				fmt.Print("⭕ ")
			}
		}
		fmt.Println("")
		// add the colors to each numbers

		// get rid of the wrong ones
		for i := 0; i < sz; i++ {
			fmt.Printf("%-3d", guess.Num[i])
		}
		fmt.Println("\n------------")
	}
}

func updateStats(gs gamestat, as stats) stats {

	if as.fastestTime == 0 || gs.time < as.fastestTime {
		as.fastestTime = gs.time
	}
	if as.longestTime == 0 || gs.time > as.longestTime {
		as.longestTime = gs.time
	}
	if gs.guesses == 1 {
		as.firstGuessWins++
	}
	as.guesses += gs.guesses
	as.played++
	as.timeTotal += gs.time
	return as
}

func printStats(s stats) {
	fmt.Printf("Time:\n\t\tTotal time played: %s\n\t\tFastest win: %s\n\t\tLongest time: %s\n", s.timeTotal, s.fastestTime, s.longestTime)
	fmt.Printf("Game stats:\n\t\tGames played: %d\n", s.played)
	fmt.Printf("Guessing:\n\t\tAverage guesses: %.2f\n\t\tFirst guess wins: %d\n\t\tTotal guesses: %d\n", float64(s.guesses)/float64(s.played), s.firstGuessWins, s.guesses)
}

func main() {
	size := 5
	password := getPassword(size)
	// get input
	fmt.Println(password)

	// game setup stats and variables
	s := bufio.NewScanner(os.Stdin)
	allstats := stats{}
	gs := gamestat{}
	start := time.Now()
	guesses := make([]*Guess, 0)
	for {
		guess := NewGuess(size)
		getUserGuess(s, guess, size)

		// check if numbers match
		assignColors(guess, password, size)
		guesses = append(guesses, guess)
		printHistory(guesses, size)

		// check win
		if checkWin(guess, size) {
			gs.time = time.Since(start)
			gs.guesses = len(guesses)
			fmt.Printf("You did it!\n\t\tGuesses: %d\n\t\tTime: %s\n\n", gs.guesses, gs.time)
			allstats = updateStats(gs, allstats)
			printStats(allstats)
			fmt.Printf("Play again? (y/N):> ")
			s.Scan()
			if strings.ToLower(s.Text()) != "y" {
				break
			}
			// reset game variables
			gs = gamestat{}
			start = time.Now()
			guesses = make([]*Guess, 0)
			password = getPassword(size)
			fmt.Println(password)

		}
	}
}
