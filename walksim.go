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

func main(){
	roomlist := make(map[string]Room)
	location := "101"
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
				dest, ok := roomlist[location].exits["north"]
				if ok == false {
					fmt.Println("You can't go that way.")
					suppressdesc = true
				} else {
					location = dest
				}
			case uinput == "east" || uinput == "e":
				dest, ok := roomlist[location].exits["east"]
				if ok == false {
					fmt.Println("You can't go that way.")
					suppressdesc = true
				} else {
					location = dest
				}
			case uinput == "south" || uinput == "s":
				dest, ok := roomlist[location].exits["south"]
				if ok == false {
					fmt.Println("You can't go that way.")
					suppressdesc = true
				} else {
					location = dest
				}
			case uinput == "west" || uinput == "w":
				dest, ok := roomlist[location].exits["west"]
				if ok == false {
					fmt.Println("You can't go that way.")
					suppressdesc = true
				} else {
					location = dest
				}
			case uinput == "up" || uinput == "u":
				dest, ok := roomlist[location].exits["up"]
				if ok == false {
					fmt.Println("You can't go that way.")
					suppressdesc = true
				} else {
					location = dest
				}
			case uinput == "down" || uinput == "d":
				dest, ok := roomlist[location].exits["down"]
				if ok == false {
					fmt.Println("You can't go that way.")
					suppressdesc = true
				} else {
					location = dest
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
				for key, _ := range roomlist[location].exits {
					fmt.Printf(" %v", key)
				}
				fmt.Printf(".\n")
				suppressdesc = true
			default:
				fmt.Printf("I don't know what '%v' means.\n", uinput)
		}
	}
	fmt.Println("Goodbye.")
}