package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Item struct {
	name     string
	desc     string
	gettable bool
	hasinv   bool
	inv      map[string]Item
	onuse    func(*Item) string
	users    map[string]func(*Item) string
}

type Room struct {
	name  string
	desc  string
	exits map[string]string
	inv   map[string]Item
}

func printhelp() {
	fmt.Println("Available Commands")
	fmt.Println("These commands can be abbreviated with the first letter of the word:")
	fmt.Println("north, east, south, west, up, down, quit, help, look")
	fmt.Println("These commands must be typed out in full:")
	fmt.Println("get, drop, inv, use, go, put")
	fmt.Println("")
}

func (r Room) String() string {
	var result string
	result += r.desc + "\n" + "Exits:"
	for key, _ := range r.exits {
		result += " " + key
	}
	result += "."
	for key, _ := range r.inv {
		result += "\nYou see " + r.inv[key].name + " here."
	}
	return result
}

func (i Item) String() string {
	return i.desc
}

func (i *Item) use() string {
	return i.onuse(i)
}

func (i *Item) useon(target *Item) string {
	instr, ok := target.users[i.name]
	if ok {
		return instr(i)
	}
	return "You can't use " + i.name + " on " + target.name + ".\n"
}

func contains(sl []string, target string) bool {
	for _, elem := range sl {
		if elem == target {
			return true
		}
	}
	return false
}

func main() {
	roomlist := make(map[string]Room)
	invlist := make(map[string]Item)
	playerinv := make(map[string]Item)
	startloc := "101"
	runloop := true
	suppressdesc := true
	verbs := []string{"go", "north", "east", "south", "west", "up", "down", "help", "look", "quit", "inv", "get", "drop", "use", "put"}
	//north, east south, west, up, and down alias to "go north", etc. drop aliases to "put ITEM here".
	//look by itself becomes "look here"
	prepos := []string{"in", "on", "to", "with"}
	var notacommand bool

	// read the rooms file and parse it out into the room list

	roomfile, err := os.Open("rooms.txt")
	if err != nil {
		fmt.Println(err)
		fmt.Println("Exiting.")
		runloop = false
	}
	defer roomfile.Close()

	itemfile, err := os.Open("items.txt")
	if err != nil {
		fmt.Println(err)
		fmt.Println("Exiting.")
		runloop = false
	}
	defer itemfile.Close()

	iscanner := bufio.NewScanner(itemfile)
	for iscanner.Scan() {
		var newitem Item
		rawitem := strings.Split(iscanner.Text(), "|")
		newitem.name = rawitem[0]
		newitem.desc = rawitem[1]
		newitem.gettable = (rawitem[2] == "T")
		newitem.hasinv = (rawitem[3] == "T")
		newitem.onuse = func(i *Item) string {
			return "You use " + i.name + ".\n"
		}
		parsedusers := strings.Fields(rawitem[5])
		for _, key := range parsedusers {
			newitem.users[key] = func(i *Item) string {
				return "You use " + key + " on " + i.name + ".\n"
			}
		}
		invlist[newitem.name] = newitem
	}

	rscanner := bufio.NewScanner(roomfile)
	for rscanner.Scan() {
		var newroom Room
		rawroom := strings.Split(rscanner.Text(), "|")
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
		rawitems := rawroom[4]
		itemlist := strings.Split(rawitems, ",")
		newroom.inv = make(map[string]Item)
		for _, itemname := range itemlist {
			itm, ok := invlist[itemname]
			if ok {
				newroom.inv[itm.name] = itm
			}
		}
		roomlist[roomid] = newroom
	}

	here := roomlist[startloc]

	//print out the help entry once to get the player started
	printhelp()
	fmt.Println(here)
	uiscanner := bufio.NewScanner(os.Stdin)

	//main loop - print name and desc of current room, accept input, act on input
	for runloop == true {
		var verb, object, prep, recip string
		if suppressdesc == false {
			fmt.Println(here)
		}
		fmt.Printf(">> ")
		suppressdesc = true
		notacommand = false
		uiscanner.Scan()
		uinput := uiscanner.Text()
		for len(uinput) == 0 {
			fmt.Printf(">> ")
			uiscanner.Scan()
			uinput = uiscanner.Text()
		}
		uinputp := strings.Fields(uinput)
		verb = uinputp[0]
		//there has to be a better way to do this but i don't yet know what it is.
		if len(uinputp) > 1 {
			object = uinputp[1]
		}
		if len(uinputp) > 2 {
			if contains(prepos, uinputp[2]) {
				prep = uinputp[2]
			} else {
				recip = uinputp[2]
			}
		}
		if len(uinputp) > 3 {
			recip = uinputp[3]
		}
		fmt.Println(verb, object, prep, recip)

		//this part seems like there should be a better way to do it. I need to distinguish between recognised commands that may not
		//be applicable to the situation, and things that aren't commands, so that I can give an appropriate error message.
		switch {
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
		case contains(verbs, uinputp[0]):
		default:
			notacommand = true
		}

		switch {
		case uinputp[0] == "quit":
			runloop = false

		case uinputp[0] == "help":
			printhelp()

		case uinputp[0] == "look":
			if len(uinputp) == 1 {
				fmt.Println(here)
			} else {
				target, ok := here.inv[uinputp[1]]
				if ok {
					fmt.Println(target)
				} else {
					target, ok := playerinv[uinputp[1]]
					if ok {
						fmt.Println(target)
					} else {
						fmt.Println("You don't see that here.")
					}
				}
			}

		case uinputp[0] == "get":
			if len(uinputp) == 1 {
				fmt.Println("Get what?")
			} else {
				target, ok := here.inv[uinputp[1]]
				if ok {
					playerinv[target.name] = target
					delete(here.inv, target.name)
					fmt.Printf("You get the %s.\n", uinputp[1])
				} else {
					fmt.Printf("You can't get ye %s.\n", uinputp[1])
				}
			}

		case uinputp[0] == "drop":
			if len(uinputp) == 1 {
				fmt.Println("Drop what?")
			} else {
				target, ok := playerinv[uinputp[1]]
				if ok {
					here.inv[target.name] = target
					delete(playerinv, target.name)
					fmt.Printf("You drop the %s.\n", uinputp[1])
				} else {
					fmt.Printf("You don't have %s.\n", uinputp[1])
				}
			}

		case uinputp[0] == "inv":
			fmt.Println("You are carrying:")
			if len(playerinv) == 0 {
				fmt.Println("Nothing.")
			} else {
				for _, itm := range playerinv {
					fmt.Println(itm.name)
				}
			}

		case uinputp[0] == "use":
			switch len(uinputp) {
			case 1:
				fmt.Println("Use what?")
			case 2:
				instr, ok := playerinv[uinputp[1]]
				if ok {
					fmt.Printf(instr.use())
				} else {
					fmt.Printf("You don't have %s.\n", uinputp[1])
				}
			case 3:
				instr, ok := playerinv[uinputp[1]]
				if ok {
					target, ok := playerinv[uinputp[2]]
					if ok {
						fmt.Printf(instr.useon(&target))
					} else {
						target, ok := here.inv[uinputp[2]]
						if ok {
							fmt.Printf(instr.useon(&target))
						} else {
							fmt.Printf("You don't see %s here.\n", uinputp[2])
						}
					}
				} else {
					fmt.Printf("You don't have %s.\n", uinputp[1])
				}
			default:
				fmt.Println("Something went wrong.")
			}
		default:
			dest, ok := here.exits[uinputp[0]]
			if ok == false {
				if notacommand == true {
					fmt.Printf("I don't understand what '%v' means here.\n", uinputp[0])
				} else {
					fmt.Printf("You can't go %v from here.\n", uinputp[0])
				}
			} else {
				here = roomlist[dest]
				suppressdesc = false
			}
		}
	}
	fmt.Println("Goodbye.")
}
