package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Room struct {
	links []string
	isChecked bool
	pathIndex int
}

var (
	antsNumber = 0
	Rooms = make(map[string]Room)
	start = ""
	end = ""
)

func main() {
	dataBytes, _ := os.ReadFile("data2.txt")
	ParsingData(string(dataBytes))
	fmt.Println(Rooms)
	FindPaths()
}

func FindPaths() {
	startRoom := Rooms[start]
	startRoom.isChecked = true
	Rooms[start] = startRoom
	paths := [][]string{{start}}
	validPaths := [][]string{}
	validLinks := 0
	needToBreak := false
	for len(paths) != 0 {
		for i:= 0; i < len(paths); i++ {
			currentPath := paths[i]
			lastInPath := paths[i][len(paths[i])-1]

			for j, link := range Rooms[lastInPath].links {
				// fmt.Println("paths", paths)
				if link != end {

					if !Rooms[link].isChecked {
						// fmt.Println("LINKS :", Rooms[lastInPath].links, link, lastInPath)
						
						room := Rooms[link]
						validLinks++
						room.isChecked = true
						Rooms[link] = room
						if validLinks == 1 {
							paths[i] = append(paths[i], link)
						} else {
							paths = append(paths, append(currentPath, link))
						}
					} else if j == len(Rooms[lastInPath].links)-1 && validLinks == 0 {
						if i + 1 < len(paths) {
							paths = append(paths[:i], paths[i+1:]...)
						} else {
							paths = paths[:i]			
						}

					}
				} else {
					paths[i] = append(paths[i], link)
					validPaths = append(validPaths, paths[i])
					fmt.Println(validPaths)
					paths = [][]string{{start}}
					resetIsChecked(validPaths)
					needToBreak = true
					break
				}
			}
			validLinks = 0
			if needToBreak {
				needToBreak = false
				break
			}
		}
	}
	fmt.Println("Valid paths :", validPaths)

}
func resetIsChecked(pathRooms [][]string) {
	allPathRooms := []string{}
	for i := range pathRooms {
		for j := range pathRooms[i] {
			allPathRooms = append(allPathRooms, pathRooms[i][j])
		}
	}

	for index := range Rooms {
		for i := range allPathRooms {
			if index == allPathRooms[i] {
				break
			}
			if i == len(allPathRooms)-1 {
				room := Rooms[index]
				room.isChecked = false
				Rooms[index] = room
			}
		}
	}
}
func ParsingData(str string) any {
	var err error
	split := strings.Split(str, "\n")
	if antsNumber, err = strconv.Atoi(split[0]); err != nil {
		log.Fatal(err)
	}
	for i:= 1; i < len(split); i++ {
		if len(strings.Split(split[i], " ")) == 3 {
			continue
		} else if split[i] == "##start" {
			start = strings.Split(split[i+1], " ")[0]
			i++
		} else if  split[i] == "##end" {
			end = strings.Split(split[i+1], " ")[0]
			i++
		} else if split[i][0] == '#' {
			continue
		} else if len(strings.Split(split[i], "-")) == 2 {
			fillRoom(split[i])
		} else {
			log.Fatal("Errorrr")
		}
	}
	fmt.Println(start, end)
	return nil
}

func fillRoom(str string) error {
	split := strings.Split(str, "-")


	room := Rooms[split[0]]
	room.links = append(room.links, split[1])
	Rooms[split[0]] = room

	room = Rooms[split[1]]
	room.links = append(room.links, split[0])
	Rooms[split[1]] = room
	return nil
}