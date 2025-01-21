package parser

const (
	// states
	start = iota
	roomsList
	end
	roomLinks
	// tokens
	antsNumber
	roomCharacter
	x
	y
	space
	dash
	hashtag = '#'
)

type state struct {
	prevToken, expectedToken, prevState, expectedState, linePosition int
}

type AntFarm struct {
	Number             int
	xyPairs            map[int]struct{}
	Rooms              map[string]*room
	StartRoom, EndRoom string
	state              *state
	currentLine        string
}

type room struct {
	Links map[string]struct{}
	X, Y  int
}

func initFarm() *AntFarm {
	return &AntFarm{
		xyPairs: map[int]struct{}{},
		Rooms:   map[string]*room{},
		state: &state{
			expectedState: antsNumber,
		},
	}
}
