package models

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"

)

//Gun model gun
type Gun struct {
	Sound audio.Context
	Bullets int
}

//Bullet model
type Bullet struct {
	Img *ebiten.Image
	Op *ebiten.DrawImageOptions
	Coords Coordinates
	Hit, Moving, FacingFront bool

}

//Shoot shoots one bullet at a time
func (g *Gun) Shoot(){
	g.Bullets--
}
//Move the bullet forward
func (b *Bullet) Move(){
	b.Moving = true
	switch b.FacingFront{
	case true:
		b.Coords.X +=10
	case false:
		b.Coords.X -= 10
	}
}
//IsHit returns tru is it coliided woth something
func (b *Bullet) IsHit()bool{
	return b.Moving
}
//Render the bullet
func (b *Bullet) Render(screen, bullet *ebiten.Image)  {
	screen.DrawImage(bullet, b.Op)
}