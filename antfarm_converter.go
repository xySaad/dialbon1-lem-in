package main

import "lemin/parser"

func ConvertStruct(af *parser.AntFarm) {
	start = af.StartRoom
	end = af.EndRoom
	antsNumber = af.Number
	for roomName, room := range af.Rooms {
		Rooms[roomName] = Room{
			links: Map2Slice(room.Links),
		}
	}
}

func Map2Slice(links map[string]struct{}) (result []string) {
	for link := range links {
		result = append(result, link)
	}
	return
}
