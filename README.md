# coro
Package coro makes Go coroutines simple and useful.

## Usage

```go
var t = 0 // time

// Coroutine takes coro.Yield as an argument
func f(y coro.Yield) {
	fmt.Println(t, "Hello")
	// Calling y() suspends the coroutine and returns control
	y()
	// t is now 1 (incremented by the caller)
	fmt.Println(t, "World")
	y()
	fmt.Println(t, "!")
}

func main() {
	// Create coroutine
	co := coro.New(f)
	// Advance coroutine
	for co.Next() {
		t++ // Increment time
	}

	// Output:
	// 0 Hello
	// 1 World
	// 2 !
}
```

### More Examples

```go
var t = 0

// Rival AI thinking simulation
func rivalAI(y coro.Yield) string {
	fmt.Println(t, "Rival AI is thinking...")
	y() // Yield control

	fmt.Println(t, "Rival AI is still thinking...")
	y() // Yield again

	fmt.Println(t, "Rival AI has finished thinking.")
	return "rock" // Return the AI's choice
}

// Player input simulation
func playerInput(y coro.Yield) string {
	fmt.Println(t, "Waiting for player input...")
	y()

	fmt.Println(t, "Waiting for player input...")
	y()

	fmt.Println(t, "Waiting for player input...")
	y()

	fmt.Println(t, "Received player input.")
	return "paper" // Return the player's choice
}

// Game logic coroutine
func gameMain(y coro.Yield) {
	fmt.Println(t, "Rival AI's turn")
	rivalChoice := rivalAI(y) // Internally yields multiple times
	fmt.Println(t, "Rival AI chose:", rivalChoice)

	fmt.Println(t, "Player's turn")
	playerChoice := playerInput(y) // Internally yields multiple times
	fmt.Println(t, "Player chose:", playerChoice)

	fmt.Println(t, "Choice:", rivalChoice, "vs", playerChoice)
}

func main() {
	// Run game coroutine
	co := coro.New(gameMain)
	for co.Next() {
		t++ // Increment time
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
```