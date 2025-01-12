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
	dataBytes, _ := os.ReadFile(os.Args[1])
	ParsingData(string(dataBytes))
	fmt.Println(Rooms, len(Rooms))

	

	

	linksToRemove := make(map[string]string)

	for {
		// time.Sleep(time.Second)

		hasFoundAny, _, _ := bfs()

		if !hasFoundAny {
			fmt.Println("Bfs didn't find any paths")
			break
		}
		saveBeforeInPath()
		

		isBackTracking, revNode, toRemove := checkIfBackTrackingPath()
		if isBackTracking {
		
		
			linksToRemove[revNode] = toRemove

		
			fmt.Println(revNode, toRemove, "revnode---------------------------")
			validPaths = [][]string{}
			Rooms = make(map[string]Room)
	
			ParsingData(string(dataBytes))

			for rev, toRm := range linksToRemove {
				room := Rooms[rev]
				room.links = removeLink(room.links, toRm)
				Rooms[rev] = room
		
				
				room2 := Rooms[toRm]
				room2.links = removeLink(room2.links, rev)
				Rooms[toRm] = room2
			}
		}
		removePathsLinks()
	}

	fmt.Println("finnal paths :", validPaths)
}

func checkIfBackTrackingPath() (bool, string, string) {
	lastPath := validPaths[len(validPaths)-1]
	pathRooms := validPaths[:len(validPaths)-1]
	links := make(map[string]string)

	// get all path links reversed
	for i := len(pathRooms)-1; i >= 0; i-- {
		fmt.Println("okk")

		fmt.Println(pathRooms)

		for j := len(pathRooms[i])-2; j >= 1; j-- {
			fmt.Println("ok")
			links[pathRooms[i][j]] = pathRooms[i][j-1]
		}
	}
	fmt.Println("------------------------------------")
	fmt.Println(links)
	fmt.Println("------------------------------------")
	for i := 1; i < len(lastPath)-1; i++ {
		if links[lastPath[i]] == lastPath[i+1] {
			return true, lastPath[i], lastPath[i+1]
		}
	}

	return false, "", ""
}

func removePathsLinks() {
	for index := 0; index < len(validPaths); index++ {
		path := validPaths[index]
		fmt.Println(path)
		for i := 0; i < len(path)-1; i++ {
			node := Rooms[path[i]]
			node.links = removeLink(node.links, path[i+1])
			Rooms[path[i]] = node
		}
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
	fmt.Println(Rooms)


	startRoom := Rooms[start]
	startRoom.isChecked = true
	Rooms[start] = startRoom

	alreadyInRevesedPath := false
	
	firstEncoutredAsReversed := ""
	toRemoveLink := ""


	paths := [][]string{{start}}
	for len(paths) != 0 {
		if len(Rooms[start].links) == 0 {
			return false, "", ""
		}
		
		for i:= 0; i < len(paths); i++ {
			fmt.Println(paths)
			validLinks := 0


			lastInPath := paths[i][len(paths[i])-1]

			if Rooms[lastInPath].beforeInPath != "" && !alreadyInRevesedPath {
				beforeInPath := Rooms[lastInPath].beforeInPath
				firstEncoutredAsReversed = beforeInPath
				toRemoveLink = lastInPath
				validLinks++

				room := Rooms[lastInPath]
				room.isChecked = false
				Rooms[lastInPath] = room
				
				paths[i] = append(paths[i], beforeInPath)

				alreadyInRevesedPath = true
				fmt.Println(Rooms)

				continue
			}

			for j, link := range Rooms[lastInPath].links {

				if link == end {

					if validLinks == 0 {
						paths[i] = append(paths[i], link)
					} else {
						paths = append(paths, append(paths[i][:len(paths[i])-1], link))
					}
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
	return false, firstEncoutredAsReversed, toRemoveLink
}

func resetIsChecked(pathRooms [][]string) {
	// allPathRooms := []string{}
	// for i := range pathRooms {
	// 	for j := range pathRooms[i] {
	// 		allPathRooms = append(allPathRooms, pathRooms[i][j])
	// 	}
	// }

	for index := range Rooms {
		room := Rooms[index]
		room.isChecked = false
		Rooms[index] = room
	}

	// for index := range Rooms {
	// 	for i := range allPathRooms {
	// 		// if index == allPathRooms[i] {
	// 		// 	break
	// 		// }
	// 		if i == len(allPathRooms)-1 {
	// 			room := Rooms[index]
	// 			room.isChecked = false
	// 			Rooms[index] = room
	// 		}
	// 	}
	// }
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