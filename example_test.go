package coro_test

import (
	"fmt"

	"github.com/eihigh/coro"
)

func Example_basic() {
	t := 0

	f := func(y coro.Yield) {
		fmt.Println(t, "Hello")
		y()
		fmt.Println(t, "World")
		y()
		fmt.Println(t, "!")
	}

	co := coro.New(f)
	for co.Next() {
		t++
	}

	// Output:
	// 0 Hello
	// 1 World
	// 2 !
}

func Example() {
	t := 0

	rivalAI := func(y coro.Yield) string {
		fmt.Println(t, "Rival AI is thinking...")
		y()
		fmt.Println(t, "Rival AI is still thinking...")
		y()
		fmt.Println(t, "Rival AI has finished thinking.")
		return "rock"
	}
	playerInput := func(y coro.Yield) string {
		fmt.Println(t, "Waiting for player input...")
		y()
		fmt.Println(t, "Waiting for player input...")
		y()
		fmt.Println(t, "Waiting for player input...")
		y()
		fmt.Println(t, "Received player input.")
		return "paper"
	}
	gameMain := func(y coro.Yield) {
		fmt.Println(t, "Rival AI's turn")
		rivalChoice := rivalAI(y)
		fmt.Println(t, "Rival AI chose:", rivalChoice)
		fmt.Println(t, "Player's turn")
		playerChoice := playerInput(y)
		fmt.Println(t, "Player chose:", playerChoice)
		fmt.Println(t, "Choice:", rivalChoice, "vs", playerChoice)
	}

	co := coro.New(gameMain)
	for co.Next() {
		t++
	}

	// Output:
	// 0 Rival AI's turn
	// 0 Rival AI is thinking...
	// 1 Rival AI is still thinking...
	// 2 Rival AI has finished thinking.
	// 2 Rival AI chose: rock
	// 2 Player's turn
	// 2 Waiting for player input...
	// 3 Waiting for player input...
	// 4 Waiting for player input...
	// 5 Received player input.
	// 5 Player chose: paper
	// 5 Choice: rock vs paper
}
