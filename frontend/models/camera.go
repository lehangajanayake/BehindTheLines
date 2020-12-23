package models

import (
	"github.com/hajimehoshi/ebiten/v2"
)

//Camera is the model camera
type Camera struct{
	Op *ebiten.GeoM
	Position Coordinates
	View *ebiten.Image
}

//Move moves the camera
func (c *Camera) Move(pos Coordinates){
	c.Position = pos
}

//Render reders the camera
func (c *Camera) Render(screen *ebiten.Image){
	c.Op.Reset()
	w, h := ebiten.ScreenSizeInFullscreen()
	c.Op.Translate(-float64(w/2), -float64(h/2))
	c.Op.Translate(float64(c.Position.X), float64(c.Position.Y))
	if c.Op.IsInvertible(){
		c.Op.Invert()
	}
	screen.DrawImage(c.View, &ebiten.DrawImageOptions{
		GeoM: *c.Op,
	})
}