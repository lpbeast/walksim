package main

import (
	"fmt"
	"strings"
	"os"
	"bufio"
	"strconv"
)

type Room struct {
	name string
	desc string
	exits map[string]int
}

func main(){
	var roomlist []Room
	location := 0
	runloop := true
	suppressdesc := false
	// read the rooms file and parse it out into the room list
	roomfile, err := os.Open("rooms.txt")
    if err != nil {
        fmt.Println(err)
		fmt.Println("Exiting.")
		runloop = false
    }
    defer roomfile.Close()

	scanner := bufio.NewScanner(roomfile)
    for scanner.Scan() {
		var newroom Room
        rawroom := strings.Split(scanner.Text(), "|")
		newroom.name = rawroom[0]
		newroom.desc = rawroom[1]
		newroom.exits = make(map[string]int)
		newroom.exits["north"], _ = strconv.Atoi(rawroom[2])
		newroom.exits["east"], _ = strconv.Atoi(rawroom[3])
		newroom.exits["south"], _ = strconv.Atoi(rawroom[4])
		newroom.exits["west"], _ = strconv.Atoi(rawroom[5])
		newroom.exits["up"], _ = strconv.Atoi(rawroom[6])
		newroom.exits["down"], _ = strconv.Atoi(rawroom[7])
		roomlist = append(roomlist, newroom)
    }

	//print out the help entry once to get the player started
	fmt.Println("Available commands:")
	fmt.Println("north, east, south, west, up, down, quit, help, look")
	fmt.Println("You can use the first letter of a command instead of typing the whole thing.")
	
	//main loop - print name and desc of current room, accept input, act on input
	for runloop == true{
		if suppressdesc == false {
			fmt.Println(roomlist[location].name)
			fmt.Println(roomlist[location].desc)
		}
		fmt.Printf(">> ")
		suppressdesc = false
		var uinput string
		fmt.Scanln(&uinput)
		switch{
			case uinput == "north" || uinput == "n":
				if roomlist[location].exits["north"] < 0{
					fmt.Println("You can't go that way.")
					suppressdesc = true
				} else {
					location = roomlist[location].exits["north"]
				}
			case uinput == "east" || uinput == "e":
				if roomlist[location].exits["east"] < 0{
					fmt.Println("You can't go that way.")
					suppressdesc = true
				} else {
					location = roomlist[location].exits["east"]
				}
			case uinput == "south" || uinput == "s":
				if roomlist[location].exits["south"] < 0{
					fmt.Println("You can't go that way.")
					suppressdesc = true
				} else {
					location = roomlist[location].exits["south"]
				}
			case uinput == "west" || uinput == "w":
				if roomlist[location].exits["west"] < 0{
					fmt.Println("You can't go that way.")
					suppressdesc = true
				} else {
					location = roomlist[location].exits["west"]
				}
			case uinput == "up" || uinput == "u":
				if roomlist[location].exits["up"] < 0{
					fmt.Println("You can't go that way.")
					suppressdesc = true
				} else {
					location = roomlist[location].exits["up"]
				}
			case uinput == "down" || uinput == "d":
				if roomlist[location].exits["down"] < 0{
					fmt.Println("You can't go that way.")
					suppressdesc = true
				} else {
					location = roomlist[location].exits["down"]
				}
			case uinput == "quit" || uinput == "q":
				runloop = false
			case uinput == "help" || uinput == "h":
				fmt.Println("Available commands:")
				fmt.Println("north, east, south, west, up, down, quit, help, look")
				fmt.Println("You can use the first letter of a command instead of typing the whole thing.")
				suppressdesc = true
			case uinput == "look" || uinput == "l":
				fmt.Println(roomlist[location].desc)
				fmt.Printf("Exits:")
				if roomlist[location].exits["north"] >= 0 { fmt.Printf(" north") }
				if roomlist[location].exits["east"] >= 0 { fmt.Printf(" east") }
				if roomlist[location].exits["south"] >= 0 { fmt.Printf(" south") }
				if roomlist[location].exits["west"] >= 0 { fmt.Printf(" west") }
				if roomlist[location].exits["up"] >= 0 { fmt.Printf(" up") }
				if roomlist[location].exits["down"] >= 0 { fmt.Printf(" down") }
				fmt.Printf(".\n")
				suppressdesc = true
			default:
				fmt.Printf("I don't know what '%v' means.\n", uinput)
		}
	}
	fmt.Println("Goodbye.")
}