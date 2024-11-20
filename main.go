package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Notepad struct {
	capacity int
	notes    []string
}

type Arg struct {
	cmd      string
	data     string
	position string
}

func main() {
	var capacity int
	fmt.Println("Enter the maximum number of notes:")
	fmt.Scan(&capacity)

	notepad := NewNotepad(capacity)
	scanner := bufio.NewScanner(os.Stdin)

outer:
	for {
		fmt.Print("Enter a command and data: ")

		arg, err := getCmd(scanner)

		if err != nil {
			fmt.Printf("[Error] %s\n", err)
			continue
		}

		switch arg.cmd {
		case "create":
			if err := notepad.addNote(arg.data); err != nil {
				fmt.Printf("[Error] %s\n", err)
			} else {
				fmt.Println("[OK] The note was successfully created")
			}
		case "update":
			if err := notepad.updateNote(arg.position, arg.data); err != nil {
				fmt.Printf("[Error] %s\n", err)
			} else {
				fmt.Printf("[OK] The note at position %s was successfully updated\n", arg.position)
			}
		case "list":
			if len(notepad.notes) == 0 {
				fmt.Println("[Info] Notepad is empty")
			} else {
				for i, note := range notepad.notes {
					fmt.Printf("[Info] %d: %s\n", i+1, note)
				}
			}
		case "delete":
			if err := notepad.deleteNote(arg.position); err != nil {
				fmt.Printf("[Error] %s\n", err)
			} else {
				fmt.Printf("[OK] The note at position %s was successfully deleted\n", arg.position)
			}
		case "clear":
			notepad.clear()
			fmt.Println("[OK] All notes were successfully deleted")
		case "exit":
			break outer
		default:
			fmt.Println("[Error] Unknown command")
		}
	}

	fmt.Println("[Info] Bye!")
}

func getCmd(scanner *bufio.Scanner) (Arg, error) {

	arg := Arg{}

	if scanner.Scan() {
		input := strings.Split(scanner.Text(), " ")

		arg.cmd = strings.ToLower(input[0])

		switch arg.cmd {
		case "create":
			if len(input) > 1 {
				arg.data = strings.Join(input[1:], " ")
			}
		case "delete", "update":
			if len(input) > 1 {
				arg.position = input[1]
			}
			if len(input) > 2 {
				arg.data = strings.Join(input[2:], " ")
			}
		}

		return arg, nil
	}

	return arg, scanner.Err()
}

func NewNotepad(capacity int) *Notepad {
	return &Notepad{capacity: capacity, notes: make([]string, 0, capacity)}
}

func (n *Notepad) addNote(note string) error {
	if len(n.notes) == n.capacity {
		return errors.New("Notepad is full")
	} else if len(note) == 0 {
		return errors.New("Missing note argument")
	} else {
		n.notes = append(n.notes, note)
		return nil
	}
}

func (n *Notepad) updateNote(position, note string) error {
	if position == "" {
		return errors.New("Missing position argument")
	}

	if note == "" {
		return errors.New("Missing note argument")
	}

	pos, err := strconv.Atoi(position)
	if err != nil {
		return fmt.Errorf("Invalid position: %s", position)
	}
	pos--

	if pos < 0 || pos >= n.capacity {
		return fmt.Errorf("Position %d is out of the boundaries [1, %d]", pos+1, n.capacity)
	}
	if pos >= len(n.notes) {
		return errors.New("There is nothing to update")
	}

	n.notes[pos] = note

	return nil
}

func (n *Notepad) deleteNote(position string) error {
	if position == "" {
		return errors.New("Missing position argument")
	}
	pos, err := strconv.Atoi(position)
	if err != nil {
		return fmt.Errorf("Invalid position: %s", position)
	}
	pos--

	if pos < 0 || pos >= n.capacity {
		return fmt.Errorf("Position %d is out of the boundaries [1, %d]", pos+1, n.capacity)
	}
	if pos >= len(n.notes) {
		return errors.New("There is nothing to delete")
	}

	switch pos {
	case 0:
		n.notes = n.notes[1:]
	case len(n.notes) - 1:
		n.notes = n.notes[:pos]
	default:
		notes := n.notes[:pos]
		notes = append(notes, n.notes[pos+1:]...)
		n.notes = notes
	}
	return nil
}

func (n *Notepad) clear() {
	n.notes = make([]string, 0, n.capacity)
}
