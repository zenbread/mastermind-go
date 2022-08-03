package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
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

type Stats struct {
	TimeTotal      time.Duration
	FastestTime    time.Duration
	LongestTime    time.Duration
	Streak         int
	Played         int
	Quit           int
	Guesses        int
	FirstGuessWins int
}

type gamestat struct {
	time    time.Duration
	quit    int
	guesses int
}

func assignColors(guess *Guess, password []int, sz int) {
	marked := make([]bool, sz)
	for i := 0; i < sz; i++ {
		if guess.Num[i] == password[i] {
			guess.Color[i] = Right
			marked[i] = true
		}
	}
	for i := 0; i < sz; i++ {
		for ind, x := range password {
			if !marked[ind] && guess.Color[i] != Right && guess.Num[i] == x {
				marked[ind] = true
				guess.Color[i] = White
				break
			}
		}
	}
}

func getPassword(o Options) []int {
	rand.Seed(time.Now().UnixNano())
	password := make([]int, o.size)

	for i := 0; i < o.size; i++ {
		password[i] = rand.Intn(o.max-o.min) + o.min
	}

	return password
}

func getUserGuess(s *bufio.Scanner, guess *Guess, sz int) {
	errorMSG := "Please enter %d numbers separated by spaces. E.G.: 0 1 2 3 4<enter>\n"
	for {
		fmt.Print("Please enter your guess:> ")
		s.Scan()
		userText := strings.ToLower(s.Text())
		if userText == "quit" || userText == "q" {
			guess.Num[0] = -1
			break
		}
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

func printHistory(o Options, guesses []*Guess) {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%%%ds ", o.spaces)
	for _, guess := range guesses {
		for i := 0; i < o.size; i++ {
			switch guess.Color[i] {
			case Right:
				fmt.Printf(sb.String(), "!")
			case White:
				fmt.Printf(sb.String(), "^")
			case Wrong:
				fmt.Printf(sb.String(), "x")
			}
		}
		fmt.Println("")

		// get rid of the wrong ones
		for i := 0; i < o.size; i++ {
			fmt.Printf(fmt.Sprintf("%%%dd", o.spaces+1), guess.Num[i])
		}
		fmt.Println("\n------------")
	}
}

func updateStats(gs gamestat, as Stats) Stats {

	if as.FastestTime == 0 || gs.time < as.FastestTime {
		as.FastestTime = gs.time
	}
	if as.LongestTime == 0 || gs.time > as.LongestTime {
		as.LongestTime = gs.time
	}
	if gs.guesses == 1 {
		as.FirstGuessWins++
	}
	if gs.quit == 1 {
		as.Quit++
		as.Streak = 0
	} else {
		as.Streak++
	}
	as.Guesses += gs.guesses
	as.Played++
	as.TimeTotal += gs.time
	return as
}

func printStats(s Stats) {
	fmt.Printf("Time:\n\t\tTotal time played: %s\n\t\tFastest win: %s\n\t\tLongest time: %s\n",
		s.TimeTotal,
		s.FastestTime,
		s.LongestTime,
	)
	fmt.Printf("Game stats:\n\t\tGames played: %d\n\t\tTimes Quit: %d\n\t\tLongest Streak: %d\n",
		s.Played,
		s.Quit,
		s.Streak,
	)
	fmt.Printf("Guessing:\n\t\tAverage guesses: %.2f\n\t\tFirst guess wins: %d\n\t\tTotal guesses: %d\n",
		float64(s.Guesses)/float64(s.Played),
		s.FirstGuessWins,
		s.Guesses,
	)
}

// True is yes, False is no
func userYN(s *bufio.Scanner) bool {
	fmt.Print("Enter y/N:> ")
	s.Scan()
	output := s.Text()
	return strings.ToLower(output) == "y"
}

type Options struct {
	digits int
	spaces int
	min    int
	max    int
	size   int
}

func gameplayOptions(s *bufio.Scanner) Options {
	options := Options{
		digits: 1,
		spaces: 3,
		min:    0,
		max:    9,
		size:   5,
	}

	fmt.Print("Double Digits (Up to 25)? ")
	if userYN(s) {
		options.digits = -2
		options.spaces = 4
		options.max = 25
	}
	fmt.Print("Four numbers? ")
	if userYN(s) {
		options.size = 4
	}
	return options
}

func writeData(f *os.File, s Stats) error {
	enc := gob.NewEncoder(f)
	f.Truncate(0)
	f.Seek(0, 0)
	return enc.Encode(s)
}

func main() {
	// game setup stats and variables
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Unable to access file system.")
	}
	path := filepath.Join(dirname, "mastermind.bin")
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		fmt.Print(err)
		log.Fatal("Unable to access file system.")
	}
	defer file.Close()
	dec := gob.NewDecoder(file)
	var allstats Stats
	if err = dec.Decode(&allstats); err != nil {
		allstats = Stats{}
	}
	s := bufio.NewScanner(os.Stdin)
	options := gameplayOptions(s)
	password := getPassword(options)
	// fmt.Println(password)
	gs := gamestat{}
	start := time.Now()
	guesses := make([]*Guess, 0)
	for {
		guess := NewGuess(options.size)
		getUserGuess(s, guess, options.size)

		// check for quitting
		if guess.Num[0] == -1 {
			gs.time = time.Since(start)
			gs.guesses = len(guesses)
			gs.quit = 1
			allstats = updateStats(gs, allstats)
			printStats(allstats)
			writeData(file, allstats)
			fmt.Println("Thanks for playing.")
			break
		}
		// check if numbers match
		assignColors(guess, password, options.size)
		guesses = append(guesses, guess)
		printHistory(options, guesses)

		// check win
		if checkWin(guess, options.size) {
			gs.time = time.Since(start)
			gs.guesses = len(guesses)
			gs.quit = 0
			fmt.Printf("You did it!\n\t\tGuesses: %d\n\t\tTime: %s\n\n", gs.guesses, gs.time)
			allstats = updateStats(gs, allstats)
			printStats(allstats)
			fmt.Print("Play again? ")
			if !userYN(s) {
				writeData(file, allstats)
				break
			}
			// reset game variables
			gs = gamestat{}
			start = time.Now()
			guesses = make([]*Guess, 0)
			password = getPassword(options)
			fmt.Println(password)
		}
	}
}
