package parser

func (af *AntFarm) parseRoomLinks() error {
	for i, char := range af.currentLine {
		if char == '-' {
			if i == len(af.currentLine)-1 {
				return af.parsingError("invalid link format", 0)
			}
			if af.currentLine[:i] == af.currentLine[i+1:] {
				return af.parsingError("can't link a room with itself", 0)
			}
			from, exists := af.Rooms[af.currentLine[:i]]
			if !exists {
				return af.parsingError("can't link unexisted room", 0)
			}
			to, exists := af.Rooms[af.currentLine[i+1:]]
			if !exists {
				return af.parsingError("can't link unexisted room", 0)
			}
			_, ok := from.Links[af.currentLine[i+1:]]
			if ok {
				return af.parsingError("duplicated Links", 0)
			}

			from.Links[af.currentLine[i+1:]] = struct{}{}
			to.Links[af.currentLine[:i]] = struct{}{}
			af.state.prevToken = dash
			break
		}
	}
	if af.state.prevToken != dash {
		return af.parsingError("invalid link", 0)
	}
	return nil
}
