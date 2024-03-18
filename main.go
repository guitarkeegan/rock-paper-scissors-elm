package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
)

const (
	Playing  GameState   = "playing"
	Ending   GameState   = "ending"
	Rock     Choice      = "rock"
	Paper    Choice      = "paper"
	Scissors Choice      = "scissors"
	Won      GameMessage = "You won!!"
	Lost     GameMessage = "Sorry, you lost..."
	Drew     GameMessage = "Draw!!"
)

var (
	InputError error = errors.New("Not a valid input. Must be r, p, or s")
)

type GameState string
type Choice string
type GameMessage string
type GameError string

type model struct {
	state    GameState
	userPick Choice
	compPick Choice
	won      GameMessage
	message  GameMessage
	error    error
}

func main() {

	// create initial model
	s := bufio.NewScanner(os.Stdin)
	for {
		m := initialModel()
		View(m)
		var userInput string

		if s.Scan() {
			userInput = s.Text()
		}
		fmt.Printf("got user input: %s", userInput)
		if userInput == "" {
			// something better here
			continue
		}
		m, err := Update(m, userInput)
		if err != nil {
			fmt.Println(m.error)
			// then what?
		}
		View(m)
		if s.Scan() {
			userInput = s.Text()
		}
		if userInput == "no" || userInput == "n" {
			fmt.Println("Thank you for playing!!!")
			os.Exit(0)
		}
	}

}

func initialModel() model {
	return model{
		state:    Playing,
		userPick: "",
		compPick: "",
	}
}

func View(m model) {

	switch m.state {
	case Playing:

		Clear()
		fmt.Print("[r]ock, [p]aper, [s]cissors?\n")
		fmt.Print("> ")

	case Ending:
		Clear()

		if m.error != nil {
			fmt.Println(m.error)
		} else {
			fmt.Printf("User chose %s, and computer chose %s. %s\n", m.userPick, m.compPick, m.message)
		}

		fmt.Print("Would you like to play again?\n")
		fmt.Print("> ")
	default:
		log.Fatal("State was unnaccounted for!")
	}
}

func Update(m model, message string) (model, error) {
	var msg Choice
	switch message {
	case "r":
		msg = Rock
	case "p":
		msg = Paper
	case "s":
		msg = Scissors
	default:
		// error to model
		m.error = InputError
		m.state = Ending
		return m, m.error
	}
	cChoice := getRandomChoice()
	m.compPick = cChoice
	m.userPick = msg
	result := userWon(m.userPick, m.compPick)
	m.won = result
	m.message = result
	m.state = Ending
	return m, nil
}

func userWon(u Choice, c Choice) GameMessage {
	if u == Scissors && c == Paper ||
		u == Paper && c == Rock ||
		u == Rock && c == Scissors {
		return Won
	} else if c == Scissors && u == Paper ||
		c == Paper && u == Rock ||
		c == Rock && u == Scissors {
		return Lost
	} else {
		return Drew
	}
}

// generate random computer choice
func getRandomChoice() Choice {
	choices := [...]Choice{Rock, Paper, Scissors}
	idx := rand.Intn(3)

	return choices[idx]
}

func Clear() {

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	default:
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	cmd.Run()
}
