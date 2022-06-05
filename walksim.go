package main

import (
	"fmt"
	"strings"
	"os"
	"bufio"
)

type Room struct {
	name string
	desc string
	exits map[string]string
}

func printhelp() {
	fmt.Println("Available Commands")
	fmt.Println("These commands can be abbreviated with the first letter of the word:")
	fmt.Println("north, east, south, west, up, down, quit, help, look")
	fmt.Println("These commands must be typed out in full:")
	fmt.Println("put, drop, open, close")
	fmt.Println("")
}

func main(){
	roomlist := make(map[string]Room)
	location := "101"
	runloop := true
	suppressdesc := false
	var notacommand bool
	
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
		roomid := rawroom[0]
		newroom.name = rawroom[1]
		newroom.desc = rawroom[2]
		rawexits := rawroom[3]
		newroom.exits = make(map[string]string)
		exitlist := strings.Split(rawexits, ",")
		for _, exitrow := range exitlist {
			exitrowparsed := strings.Fields(exitrow)
			dir := exitrowparsed[0]
			dest := exitrowparsed[1]
			newroom.exits[dir] = dest
		}
		roomlist[roomid] = newroom
    }

	//print out the help entry once to get the player started
	printhelp()
	
	//main loop - print name and desc of current room, accept input, act on input
	for runloop == true{
		if suppressdesc == false {
			fmt.Println(roomlist[location].name)
			fmt.Println(roomlist[location].desc)
		}
		fmt.Printf(">> ")
		suppressdesc = false
		notacommand = false
		var uinput string
		fmt.Scanln(&uinput)
		uinputp := strings.Fields(uinput)
		switch{
			case uinputp[0] == "n":
				uinputp[0] = "north"
			case uinputp[0] == "e":
				uinputp[0] = "east"
			case uinputp[0] == "s":
				uinputp[0] = "south"
			case uinputp[0] == "w":
				uinputp[0] = "west"
			case uinputp[0] == "u":
				uinputp[0] = "up"
			case uinputp[0] == "d":
				uinputp[0] = "down"
			case uinputp[0] == "q":
				uinputp[0] = "quit"
			case uinputp[0] == "h":
				uinputp[0] = "help"
			case uinputp[0] == "l":
				uinputp[0] = "look"
			case uinputp[0] == "north" || uinputp[0] == "south" || uinputp[0] == "east" || uinputp[0] == "west":
			case uinputp[0] == "up" || uinputp[0] == "down" || uinputp[0] == "help" || uinputp[0] == "quit" || uinputp[0] == "look":
			case uinputp[0] == "get" || uinputp[0] == "drop" || uinputp[0] == "open" || uinputp[0] == "close":
			default:
				notacommand = true
		}
		switch{
			case uinputp[0] == "quit":
				runloop = false
			case uinputp[0] == "help":
				printhelp()
				suppressdesc = true
			case uinputp[0] == "look":
				fmt.Println(roomlist[location].desc)
				fmt.Printf("Exits:")
				for key, _ := range roomlist[location].exits {
					fmt.Printf(" %v", key)
				}
				fmt.Printf(".\n")
				suppressdesc = true
			default:
				dest, ok := roomlist[location].exits[uinputp[0]]
				if ok == false {
					if notacommand == true {
						fmt.Printf("I don't understand what '%v' means here.\n", uinputp[0])
					} else {
						fmt.Printf("You can't go %v from here.\n", uinputp[0])
					}
					suppressdesc = true
				} else {
					location = dest
				}
		}
	}
	fmt.Println("Goodbye.")
}