package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type answer struct {
	num   int
	right bool
	in    bool
}

func main() {
	rand.Seed(time.Now().UnixNano())
	passSize := 5

	password := make([]int, passSize, passSize)

	for i := 0; i < passSize; i++ {
		password[i] = rand.Intn(9) + 1
	}

	// get input
	fmt.Println(password)

	s := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Please enter your guess:> ")
		s.Scan()
		input := s.Text()
		var nums [5]answer
		if n, err := fmt.Sscanf(input, "%d %d %d %d %d", &nums[0].num, &nums[1].num, &nums[2].num, &nums[3].num, &nums[4].num); err != nil || n != 5 {
			fmt.Println("Please enter 5 numbers 1-9 separated by spaces. E.G.: 1 2 3 4 5")
			continue
		}

		// check if numbers match
		var marked [5]bool
		for i := 0; i < 5; i++ {
			if nums[i].num == password[i] {
				nums[i].right = true
			} else {
				for ind, x := range password {
					if nums[i].num == x && !marked[ind] {
						marked[ind] = true
						nums[i].in = true
						break
					}
				}
			}

		}
		check := true
		for i := 0; i < 5; i++ {
			if !nums[i].right {
				check = false
				break
			}
		}
		if check {
			fmt.Println("You got it!")
			break
		}
		for i := 0; i < 5; i++ {
			if nums[i].right {
				fmt.Print("Right ")
			} else if nums[i].in {
				fmt.Print("White ")
			} else {
				fmt.Print("Wrong ")
			}
		}
		fmt.Println("")
		// add the colors to each numbers

		// get rid of the wrong ones
		fmt.Println(input)
	}
}
