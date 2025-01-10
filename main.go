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
	
	hasFoundPath, _, _ := bfs()
	if !hasFoundPath {
		fmt.Println("Bfs didn't find any paths")
		return
	}
	saveBeforeInPath()
	removePathLinks()
	fmt.Println(Rooms, len(Rooms))


	hasFoundPath, revNode, toRemove := bfs()
	if !hasFoundPath {
		fmt.Println("Bfs didn't find any paths")
		return
	}
	saveBeforeInPath()
	removePathLinks()
	
	validPaths = [][]string{}
	Rooms = make(map[string]Room)


	// Normal BFS
	ParsingData(string(dataBytes))

	if revNode != "" && toRemove != "" {
		room := Rooms[revNode]
		room.links = removeLink(room.links, toRemove)
		Rooms[revNode] = room

		
		room2 := Rooms[toRemove]
		room2.links = removeLink(room2.links, revNode)
		Rooms[toRemove] = room2
	}
	
	fmt.Println(Rooms, len(Rooms))
	
	hasFoundPath, _, _ = bfs()
	if !hasFoundPath {
		fmt.Println("Bfs didn't find any paths")
		return
	}
	saveBeforeInPath()
	removePathLinks()
	fmt.Println(Rooms, len(Rooms))


	hasFoundPath, revNode, toRemove = bfs()
	if !hasFoundPath {
		fmt.Println("Bfs didn't find any paths")
		return
	}
	saveBeforeInPath()
	removePathLinks()
	


	// hasFoundPath = bfs()
	// if !hasFoundPath {
	// 	fmt.Println("Bfs didn't find any paths")
	// 	return
	// }
	// saveBeforeInPath()
	// removePathLinks()
}

func removePathLinks() {
	lastPath := validPaths[len(validPaths)-1]
	for i := 0; i < len(lastPath)-1; i++ {
		node := Rooms[lastPath[i]]
		node.links = removeLink(node.links, lastPath[i+1])
		Rooms[lastPath[i]] = node
	}
}

func saveBeforeInPath() {
	lastPath := validPaths[len(validPaths)-1]
	for i := 1; i < len(lastPath)-1; i++ {
		room := Rooms[lastPath[i]]
		room.beforeInPath = lastPath[i-1]
		Rooms[lastPath[i]] = room
	}
}

// Removes a link from a vertex
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

func bfs() (bool, string, string) {

	startRoom := Rooms[start]
	startRoom.isChecked = true
	Rooms[start] = startRoom

	alreadyInRevesedPath := false
	firstEncoutredAsReversed := ""
	toRemoveLink := ""

	paths := [][]string{{start}}
	for len(paths) != 0 {
		for i:= 0; i < len(paths); i++ {
			validLinks := 0


			lastInPath := paths[i][len(paths[i])-1]
			fmt.Println(paths)

			if Rooms[lastInPath].beforeInPath != "" && !alreadyInRevesedPath {
				beforeInPath := Rooms[lastInPath].beforeInPath
				fmt.Println(beforeInPath)
				firstEncoutredAsReversed = beforeInPath
				toRemoveLink = lastInPath
				validLinks++

				room := Rooms[beforeInPath]
				room.isChecked = true
				Rooms[beforeInPath] = room
				fmt.Println(paths[i], beforeInPath)
				
				paths[i] = append(paths[i], beforeInPath)
				fmt.Println(paths[i], beforeInPath)

				alreadyInRevesedPath = true

				continue
			}

			for j, link := range Rooms[lastInPath].links {

				if link == end {
					paths[i] = append(paths[i], link)
					validPaths = append(validPaths, paths[i])
					paths = [][]string{{start}}
					resetIsChecked(validPaths)
					fmt.Println("Valid paths :", validPaths)
					return true, firstEncoutredAsReversed, toRemoveLink
				}

				if !Rooms[link].isChecked {
					// fmt.Println(lastInPath, link)
					room := Rooms[link]
					validLinks++
					room.isChecked = true
					Rooms[link] = room
					if validLinks == 1 {
						paths[i] = append(paths[i], link)
					} else {
						var path = make([]string, len(paths[i]))
						copy(path, paths[i])
						paths = append(paths, append(path[:len(path)-1], link))
					}

				} else if j == len(Rooms[lastInPath].links)-1 && validLinks == 0 {
					if i + 1 < len(paths) {
						paths = append(paths[:i], paths[i+1:]...)
					} else {
						paths = paths[:i]			
					}

				}
			}
		}
	}
	return false, "", ""
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
			// if index == allPathRooms[i] {
			// 	break
			// }
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