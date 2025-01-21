package parser

func (af *AntFarm) parseRoomList() error {
	spaceIndex := 0
	for i, char := range af.currentLine {
		switch af.state.expectedToken {
		case roomCharacter:
			if af.state.prevToken == space {
				return af.parsingError("missing y coordinates", i)
			}
			if char == ' ' {
				spaceIndex = i
				_, exist := af.Rooms[af.currentLine[:i]]
				if exist {
					return af.parsingError("duplicated room", 0)
				}
				af.Rooms[af.currentLine[:i]] = &room{
					Links: map[string]struct{}{},
				}
				if af.state.prevState == start {
					af.StartRoom = af.currentLine[:i]
				}
				if af.state.prevState == end {
					af.EndRoom = af.currentLine[:i]
				}
				af.state.prevToken = space
				af.state.expectedToken = x
				continue
			}
			if char == '-' && i != 0 {
				if af.StartRoom == "" {
					return af.parsingError("no start room provided", 0)
				}
				if af.EndRoom == "" {
					return af.parsingError("no end room provided", 0)
				}

				af.state.prevState = roomsList
				af.state.expectedState = roomLinks
				af.parseRoomLinks()
				return nil
			}
		case x:
			if char == ' ' {
				if af.state.prevToken == space {
					return af.parsingError("invalid space after", i)
				}
				af.state.prevToken = space
				af.state.expectedToken = y

				continue
			}
			if char >= '0' && char <= '9' {
				af.Rooms[af.currentLine[:spaceIndex]].X = af.Rooms[af.currentLine[:spaceIndex]].X*10 + int(char-'0')
				af.state.prevToken = x
			} else {
				return af.parsingError("invalid x value", i)
			}
		case y:
			if char == ' ' {
				if af.state.prevToken == space || i < len(af.currentLine) {
					return af.parsingError("invalid space after", i)
				}
				af.state.prevToken = space
				af.state.expectedToken = roomCharacter
				continue
			}
			if char >= '0' && char <= '9' {
				af.Rooms[af.currentLine[:spaceIndex]].Y = af.Rooms[af.currentLine[:spaceIndex]].Y*10 + int(char-'0')
				af.state.prevToken = y
				af.state.expectedToken = roomCharacter
			} else {
				return af.parsingError("invalid y value", i)
			}
		}
	}
	if af.state.prevToken == x || af.state.expectedToken == y {
		return af.parsingError("missing y value", 0)
	}

	if af.state.prevToken == y {
		room, alo := af.Rooms[af.currentLine[:spaceIndex]]
		if !alo {
			return af.parsingError("invalid format", 0)
		}
		if room.X < room.Y {
			uniquePair := room.Y*room.Y + room.X
			_, exists := af.xyPairs[uniquePair]
			if exists {
				return af.parsingError("duplicated coordinates", 0)
			}
			af.xyPairs[uniquePair] = struct{}{}
		} else {
			uniquePair := room.X*room.X + room.X + room.Y
			_, exists := af.xyPairs[uniquePair]
			if exists {
				return af.parsingError("duplicated coordinates", 0)
			}
			af.xyPairs[uniquePair] = struct{}{}
		}
	}
	if af.currentLine[:spaceIndex] == "" && af.state.expectedToken == roomCharacter {
		return af.parsingError("invalid format", 0)
	}
	return nil
}
