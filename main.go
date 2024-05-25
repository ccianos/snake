package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"

	"github.com/eiannone/keyboard"
)

const (
	width  = 30
	height = 20
)

type Point struct {
	x, y int
}

var snake []Point
var direction Point
var food Point
var score int

func main() {
	initGame()
	defer keyboard.Close()

	ticker := time.NewTicker(200 * time.Millisecond)
	go func() {
		for {
			<-ticker.C
			update()
			if checkCollision() {
				ticker.Stop()
				fmt.Println("Game over!")
				fmt.Printf("Final score: %d\n", score)
				os.Exit(0)
			}
			render()
		}
	}()

	for {
		processInput()
	}
}

// initGame initializes the game state
func initGame() {
	snake = []Point{{width / 2, height / 2}}
	direction = Point{0, -1}
	spawnFood()
	score = 0

	if err := keyboard.Open(); err != nil {
		panic(err)
	}
}

// render clears the screen and draws the game board, snake, and food
func render() {
	clearScreen()

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if x == 0 || x == width-1 || y == 0 || y == height-1 {
				fmt.Print("#") // Border of the game board
			} else if contains(snake, Point{x, y}) {
				fmt.Print("O") // Snake body
			} else if food.x == x && food.y == y {
				fmt.Print("X") // Food
			} else {
				fmt.Print(" ") // Empty space
			}
		}
		fmt.Println()
	}

	fmt.Printf("Score: %d\n", score)
}

// processInput reads and processes user input to change the snake's direction
func processInput() {
	if char, key, err := keyboard.GetKey(); err == nil {
		if key == keyboard.KeyEsc || char == 'q' {
			os.Exit(0)
		}

		switch key {
		case keyboard.KeyArrowUp:
			if direction != (Point{0, 1}) {
				direction = Point{0, -1}
			}
		case keyboard.KeyArrowDown:
			if direction != (Point{0, -1}) {
				direction = Point{0, 1}
			}
		case keyboard.KeyArrowLeft:
			if direction != (Point{1, 0}) {
				direction = Point{-1, 0}
			}
		case keyboard.KeyArrowRight:
			if direction != (Point{-1, 0}) {
				direction = Point{1, 0}
			}
		}
	}
}

// update moves the snake and handles food consumption
func update() {
	head := snake[0]
	newHead := Point{head.x + direction.x, head.y + direction.y}

	// Check if the snake has eaten the food
	if newHead == food {
		snake = append([]Point{newHead}, snake...)
		spawnFood()
		score++
	} else {
		snake = append([]Point{newHead}, snake[:len(snake)-1]...)
	}
}

// checkCollision checks if the snake has collided with itself or the walls
func checkCollision() bool {
	head := snake[0]

	// Check for collision with walls
	if head.x <= 0 || head.x >= width-1 || head.y <= 0 || head.y >= height-1 {
		return true
	}

	// Check for collision with itself
	for _, segment := range snake[1:] {
		if segment == head {
			return true
		}
	}

	return false
}

// spawnFood places the food at a random position not occupied by the snake
func spawnFood() {
	for {
		food = Point{rand.Intn(width-2) + 1, rand.Intn(height-2) + 1}
		if !contains(snake, food) {
			break
		}
	}
}

// contains checks if a slice of Points contains a specific Point
func contains(snake []Point, point Point) bool {
	for _, p := range snake {
		if p == point {
			return true
		}
	}
	return false
}

// clearScreen clears the terminal screen
func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
