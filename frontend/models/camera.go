package models

import "github.com/hajimehoshi/ebiten/v2"

//Camera is the model camera
type Camera struct{
	Op *ebiten.GeoM
	Position Coordinates
}

//Move moves the camera
func (c *Camera) Move(pos Coordinates){
	c.Position = pos
}

//Reneder reders the camera
func (c *Camera) Render(world, screen *ebiten.Image){
	c.Op.Reset()
	c.Op.Translate()
	screen.DrawImage(world, c.Op)
}