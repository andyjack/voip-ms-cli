package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"

	input "github.com/tcnksm/go-input"
)

func blockRecent(r GetCallDataRecord) {
	if r.Status == "no_cdr" {
		// printRecents() handled it already
		return
	}
	if r.Status != "success" || len(r.CallDataRecords) == 0 {
		fmt.Printf("Status on call is %s\n", r.Status)
		fmt.Println("Did not see any recents")
		return
	}
	var blockList []string
	seen := make(map[string]struct{}, len(r.CallDataRecords))
	regex := regexp.MustCompile(`<(\d+)>`)
	for _, cdr := range r.CallDataRecords {
		loc := regex.FindStringSubmatchIndex(cdr.CallerID)
		if loc == nil {
			fmt.Printf("help! no number in %s: %s\n", cdr.Date, cdr.CallerID)
			continue
		}
		if len(loc) < 4 {
			fmt.Printf("help! had trouble matching %s: %v\n", cdr.CallerID, loc)
			continue
		}
		number := cdr.CallerID[loc[2]:loc[3]]
		if _, ok := seen[number]; ok {
			continue
		}
		blockList = append(blockList, number)
		seen[number] = struct{}{}
	}
	if len(blockList) == 0 {
		fmt.Println("No numbers to block")
		return
	}
	ui := input.DefaultUI()
	number, err := ui.Select("Pick a number to block", blockList, &input.Options{Loop: true})
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
	block-recent [D]
		- pick a number to block from a list of recent calls. Display calls from
		today to [D] days ago; D defaults to 1
	show-balance
		- show account balance
	show-recent [D]
		- show recent calls from today to [D] days ago; D defaults to 1
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
	case "show-recent", "block-recent":
		var daysAgo = 1
		if len(rest) >= 1 {
			parsed, err := strconv.Atoi(rest[0])
			if err != nil {
				log.Fatal(err)
			}
			daysAgo = parsed
		}
		dateFrom := time.Now().AddDate(0, 0, -1*daysAgo)
		r := getRecent(dateFrom)
		fmt.Println("Calls since", dateFrom.Format("2006-Jan-02"))
		printRecent(r)
		if cmd == "block-recent" {
			blockRecent(r)
		}
	default:
		usage()
	}
}
