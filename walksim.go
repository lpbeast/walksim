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
	printhelp()
	
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
			case uinput == "n":
				uinput = "north"
			case uinput == "e":
				uinput = "east"
			case uinput == "s":
				uinput = "south"
			case uinput == "w":
				uinput = "west"
			case uinput == "u":
				uinput = "up"
			case uinput == "d":
				uinput = "down"
			case uinput == "q":
				uinput = "quit"
			case uinput == "h":
				uinput = "help"
			case uinput == "l":
				uinput = "look"
			default:
		}
		switch{
			case uinput == "quit":
				runloop = false
			case uinput == "help":
				printhelp()
				suppressdesc = true
			case uinput == "look":
				fmt.Println(roomlist[location].desc)
				fmt.Printf("Exits:")
				for key, _ := range roomlist[location].exits {
					fmt.Printf(" %v", key)
				}
				fmt.Printf(".\n")
				suppressdesc = true
			default:
				dest, ok := roomlist[location].exits[uinput]
				if ok == false {
					fmt.Printf("You can't do '%v' here.\n", uinput)
					suppressdesc = true
				} else {
					location = dest
				}
		}
	}
	fmt.Println("Goodbye.")
}