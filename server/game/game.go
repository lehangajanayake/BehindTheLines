package game

import "fmt"

type Game struct{
	Player [3]*Player
	Guards [3]*Player
}

type Player struct{
	Coords struct{X int; Y int}
	FacingFront bool
}

//String returns usefull data abt the player as a string 
func (p *Player) String()string{
	return fmt.Sprintf("X: %v, Y: %v, Front: %v", p.Coords.X, p.Coords.Y, p.FacingFront)
}

//Update updates the player by whole
func (p *Player) Update(player Player){
	p = &player
}

