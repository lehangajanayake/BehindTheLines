package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

//Game tha game
type Game struct{}

//Update updates  the game 
func (g *Game) Update()error{
	return nil
}
// Draw draws to the screen every update
func (g *Game) Draw(screen *ebiten.Image){

}

// Layout lays the screen
func (g *Game) Layout(outsideWidth, outsideHeight int)(int, int){
	return 640, 480
}

func init(){
}

func main(){
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}