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
}

var (
	antsNumber = 0
	Rooms = make(map[string]Room)
	start = ""
	end = ""
	validPaths = [][]string{}
)

func main() {
	dataBytes, _ := os.ReadFile("data2.txt")
	
	// Normal BFS
	ParsingData(string(dataBytes))
	fmt.Println(Rooms, len(Rooms))
	FindPaths()
	
	// Reversed BFS


}

func FindPaths() {
	startRoom := Rooms[start]
	startRoom.isChecked = true
	Rooms[start] = startRoom
	paths := [][]string{{start}}
	validLinks := 0
	for len(paths) != 0 {
		for i:= 0; i < len(paths); i++ {
			currentPath := paths[i]
			lastInPath := paths[i][len(paths[i])-1]

			linksLoop: for j, link := range Rooms[lastInPath].links {

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
					// fmt.Println(validPaths)
					paths = [][]string{{start}}
					resetIsChecked(validPaths)
					break linksLoop
				}
			}
			validLinks = 0
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
		} else if strings.HasPrefix(split[i], "#") {
			continue
		} else if len(strings.Split(split[i], "-")) == 2 {
			fillRoomData(split[i])
		} else {
			log.Fatal("Errorrr")
		}
	}
	fmt.Println(start, end)
	return nil
}

func fillRoomData(str string) error {
	split := strings.Split(str, "-")

	if (split[0] == start && split[1] == end) || (split[0] == end && split[1] == start) {
		validPaths = [][]string{{start, end}}
		return nil
	}

	room := Rooms[split[0]]
	room.links = append(room.links, split[1])
	Rooms[split[0]] = room

	room = Rooms[split[1]]
	room.links = append(room.links, split[0])
	Rooms[split[1]] = room
	return nil
}