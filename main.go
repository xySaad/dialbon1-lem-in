package main

import (
	"fmt"
	"lemin/parser"
	"os"
	"sort"
)

type Room struct {
	links        []string
	isChecked    bool
	beforeInPath string
}

var (
	antsNumber    int
	Rooms         = make(map[string]Room)
	start         string
	end           string
	validPaths    [][]string
	allValidPaths [][][]string
)

func main() {
	// dataBytes, _ := os.ReadFile(os.Args[1])
	// data = string(dataBytes)
	af, err := parser.ParseFile(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	ConvertStruct(af)
	FindValidPaths()
	allValidPaths = append(allValidPaths, validPaths)
	fmt.Println("len final paths :", len(validPaths))

	// Sorting strings by length
	sort.Slice(validPaths, func(i, j int) bool {
		return len(validPaths[i]) < len(validPaths[j])
	})

	fmt.Println(allValidPaths)
	shortestPathIndex := 0
	lessTurns, antsOrdred := orderAnts(0)

	for i := 1; i < len(allValidPaths); i++ {
		turns, ants := orderAnts(i)
		if turns < lessTurns {
			antsOrdred = ants
			lessTurns = turns
			shortestPathIndex = i
		}
	}
	// fmt.Println(allValidPaths)
	fmt.Println("Final : OrdredAnts", antsOrdred, " - turns :", lessTurns, allValidPaths[shortestPathIndex])
}

func orderAnts(indexValidPaths int) (int, []int) {

	foundPath := allValidPaths[indexValidPaths]
	antsN := antsNumber
	// Order Ants
	ants := make([]int, len(foundPath))
	indexShortestPath := 0
	shortestPathLen := 0

	for antsN > 0 {
		shortestPathLen = len(foundPath[indexShortestPath]) + ants[indexShortestPath] // Should set to make int64 value because it is max value that can be reeturned by len()

		if len(foundPath) > 1 && indexShortestPath+1 < len(foundPath) && len(foundPath[indexShortestPath+1])+ants[indexShortestPath+1] < shortestPathLen {
			shortestPathLen = len(foundPath[indexShortestPath+1]) + ants[indexShortestPath+1]
			indexShortestPath++
		}

		ants[indexShortestPath]++
		antsN--

		if indexShortestPath == len(validPaths)-1 {
			indexShortestPath = 0
		}
	}

	fmt.Println("OrdredAnts", ants, " - turns :", shortestPathLen-1, foundPath)
	return shortestPathLen - 1, ants
}

func FindValidPaths() {

	linksToRemove := make(map[string][]string)

	for {
		// time.Sleep(time.Second)

		hasFoundAny := bfs()

		if !hasFoundAny {
			break
		}
		saveBeforeInPath()

		isBackTracking, revNode, toRemove := checkIfBackTrackingPath()

		if isBackTracking {

			linksToRemove[revNode] = append(linksToRemove[revNode], toRemove)
			allValidPaths = append(allValidPaths, validPaths[:len(validPaths)-1])

			validPaths = [][]string{}

			Rooms = make(map[string]Room)

			af, err := parser.ParseFile(os.Args[1])
			if err != nil {
				fmt.Println(err)
				return
			}
			ConvertStruct(af)

			for rev, links := range linksToRemove {
				for _, toRm := range links {
					room := Rooms[rev]
					room.links = removeLink(room.links, toRm)
					Rooms[rev] = room

					room2 := Rooms[toRm]
					room2.links = removeLink(room2.links, rev)
					Rooms[toRm] = room2
				}
			}
		} else if antsNumber == len(validPaths) {
			return
		} else {
			removePathsLinks()
		}
	}

	if len(validPaths) == 0 {
		fmt.Println("ERROR: invalid data format 0")
		os.Exit(0)
	}
}

func checkIfBackTrackingPath() (bool, string, string) {
	lastPath := validPaths[len(validPaths)-1]
	pathRooms := validPaths[:len(validPaths)-1]
	links := make(map[string]string)
	// get all path links reversed
	for i := len(pathRooms) - 1; i >= 0; i-- {

		for j := len(pathRooms[i]) - 2; j >= 1; j-- {
			links[pathRooms[i][j]] = pathRooms[i][j-1]
		}
	}

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
		for i := 0; i < len(path)-1; i++ {
			node := Rooms[path[i]]
			node.links = removeLink(node.links, path[i+1])
			Rooms[path[i]] = node
		}
	}
}

func saveBeforeInPath() {
	lastPath := validPaths[len(validPaths)-1]
	for i := 1; i < len(lastPath)-1; i++ { // see if the link to the end should be removed
		room := Rooms[lastPath[i]]
		room.beforeInPath = lastPath[i-1]
		Rooms[lastPath[i]] = room
	}
}

// Removes a link from a vertex
func removeLink(links []string, conflictRoom string) []string {
	for i := 0; i < len(links); i++ {
		if links[i] == conflictRoom {
			if i+1 < len(links) {
				links = append(links[:i], links[i+1:]...)
			} else {
				links = links[:i]
			}
		}
	}
	return links
}

func bfs() bool {

	startRoom := Rooms[start]
	startRoom.isChecked = true
	Rooms[start] = startRoom

	alreadyInRevesedPath := false

	paths := [][]string{{start}}
	for len(paths) != 0 {
		if len(Rooms[start].links) == 0 {
			return false
		}

		for i := 0; i < len(paths); i++ {
			validLinks := 0

			lastInPath := paths[i][len(paths[i])-1]

			if Rooms[lastInPath].beforeInPath != "" && !alreadyInRevesedPath {
				beforeInPath := Rooms[lastInPath].beforeInPath
				validLinks++

				room := Rooms[lastInPath]
				room.isChecked = false
				Rooms[lastInPath] = room

				paths[i] = append(paths[i], beforeInPath)

				alreadyInRevesedPath = true

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
					resetIsChecked()
					return true
				}

				if !Rooms[link].isChecked {
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
					if i+1 < len(paths) {
						paths = append(paths[:i], paths[i+1:]...)
					} else {
						paths = paths[:i]
					}

				}
			}
		}
	}
	return false
}

func resetIsChecked() {

	for index := range Rooms {
		room := Rooms[index]
		room.isChecked = false
		Rooms[index] = room
	}
}
