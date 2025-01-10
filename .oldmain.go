package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Room struct {
	links []string
	isChecked bool
	isSurballeAlgo bool
	beforeInPath string
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
	foundNewPath := false

	for len(paths) != 0 {
		for i:= 0; i < len(paths); i++ {
			currentPath := paths[i]
			lastInPath := paths[i][len(paths[i])-1]

			linksLoop: for j, link := range Rooms[lastInPath].links {
				time.Sleep(time.Millisecond * 100)
				fmt.Println(paths, "-------")

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
					} else if Rooms[link].isSurballeAlgo && Rooms[link].beforeInPath != start {
						validLinks++
						if validLinks == 1 {
							paths[i] = append(paths[i], link)
						} else {
							paths = append(paths, append(paths[i], link))
						}
						
						fmt.Println(Rooms[Rooms[link].beforeInPath], Rooms[link].beforeInPath, "hereeee before in path")
						
					} else if j == len(Rooms[lastInPath].links)-1 && validLinks == 0 {
						if i + 1 < len(paths) {
							paths = append(paths[:i], paths[i+1:]...)
						} else {
							paths = paths[:i]			
						}

					}
				} else if !Rooms[lastInPath].isSurballeAlgo {
					paths[i] = append(paths[i], link)
					validPaths = append(validPaths, paths[i])
					fmt.Println(validPaths)
					paths = [][]string{{start}}
					resetIsChecked(validPaths)
					foundNewPath = true
					break linksLoop
				}
			}
			if foundNewPath {
				fmt.Println("backtracking used")
				foundNewPath = false
				if len(validPaths) > 1 {
					fmt.Println("backtracking used")
					BTLastValidPath()
					// resetIsSurballe()
				}
				AsignSurballeValues()
			}
			
			validLinks = 0
		}
	}
	fmt.Println("Valid paths :", validPaths)

}

func AsignSurballeValues() {
	// Disable last valid path isSurballe
	path := validPaths[len(validPaths)-1]
	for roomIndex := 1; roomIndex < len(path)-1; roomIndex++ {
		room := Rooms[path[roomIndex]]
		room.isSurballeAlgo = true
		room.beforeInPath = path[roomIndex-1]
		Rooms[path[roomIndex]] = room
	}
}


// func resetIsSurballe() {
// 	// Disable last valid path isSurballe
// 	path := validPaths[len(validPaths)-2]
// 	for roomIndex := range path {
// 		room := Rooms[path[roomIndex]]
// 		room.isSurballeAlgo = false
// 		room.beforeInPath = ""
// 		Rooms[path[roomIndex]] = room
// 	}
// }

// Backtracking last valid path
func BTLastValidPath() {
	lastValidPath := validPaths[len(validPaths)-1]
	for i := len(lastValidPath)-1; i > 0; i-- {
		room := Rooms[lastValidPath[i]]
		if room.isSurballeAlgo && i != 1 {
			room2 := Rooms[lastValidPath[i-1]]
			room2.links = removeLink(room2.links, lastValidPath[i-1])
			// fmt.Println("validpathssss", validPaths)
			fmt.Println(room2, lastValidPath[i-1], "roooooooooom")
			validPaths = validPaths[:len(validPaths)-2]
			return
		}
	}
}

// Removes a conflicted link found by backtracking
func removeLink(links []string, conflictRoom string) []string {
	for i := 0; i < len(links); i++ {
		if links[i] == conflictRoom {
			if i + 1 < len(links) {
				links = append(links[:i], links[i+1:]...)
			} else {
				links = links[:i]	
			}
		}
	}
	return links
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