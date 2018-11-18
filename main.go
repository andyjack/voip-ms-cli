package main

import (
	"fmt"
	"log"
	"os"
	"regexp"

	input "github.com/tcnksm/go-input"
)

func showRecents() {
	r := getRecents()
	printRecents(r)
}

func blockRecent() {
	r := getRecents()
	printRecents(r)
	if r.Status != "success" || len(r.CallDataRecords) == 0 {
		fmt.Printf("Status on call is %s\n", r.Status)
		fmt.Println("Did not see any recents")
		return
	}
	var blockList []string
	regex := regexp.MustCompile(`<(\d+)>`)
	for _, cdr := range r.CallDataRecords {
		loc := regex.FindStringSubmatchIndex(cdr.CallerID)
		if loc == nil {
			fmt.Printf("help! no number in %s: %s\n", cdr.Date, cdr.CallerID)
			break
		}
		if len(loc) < 4 {
			fmt.Printf("help! had trouble matching %s: %v\n", cdr.CallerID, loc)
			break
		}
		blockList = append(blockList, cdr.CallerID[loc[2]:loc[3]])
	}
	if len(blockList) == 0 {
		fmt.Println("No numbers to block")
		return
	}
	ui := input.DefaultUI()
	number, err := ui.Select("Pick a number to block", blockList, &input.Options{ Loop: true })
	if err != nil {
		log.Fatal(err)
	}
	note, err := ui.Ask("Input a note?", &input.Options{})
	if err != nil {
		log.Fatal(err)
	}
	blockNumber(&number, &note)
}

func usage() {
	fmt.Println(`
Specify a command:
	block-number number [note]
	  - add a caller ID filter for the provided number, with optional note
	block-recent
	  - pick a number to block from a list of recent calls
	show-balance
	  - show account balance
	show-recents
	  - show recent calls (yesterday and today)
	`)
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		usage()
		return
	}

	cmd := args[0]
	rest := args[1:]

	switch cmd {
	case "show-balance":
		printBalance()
		break
	case "block-number":
		if len(rest) == 0 {
			fmt.Println("block-number needs a number and an optional note")
			return
		}
		var note = ""
		if len(rest) > 1 {
			note = rest[1]
		}
		number := rest[0]
		blockNumber(&number, &note)
		break
	case "show-recents":
		showRecents()
		break
	case "block-recent":
		blockRecent()
		break
	default:
		usage()
	}
}
