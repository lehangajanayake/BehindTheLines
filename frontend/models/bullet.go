package models

import (
	"image"
	//"log"
	//"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

//Bullet model
type Bullet struct {
	Img *ebiten.Image
	Op *ebiten.DrawImageOptions
	Coords Coordinates
	Hit, Moving, FacingFront bool

}

//New creates a new bullet
func (b *Bullet) New(coords Coordinates, facingFront bool)*Bullet{
	b.Coords = coords
	b.FacingFront = facingFront
	b.Op = &ebiten.DrawImageOptions{}
	return b
}

//Move the bullet forward
func (b *Bullet) Move(){
	b.Moving = true
	switch b.FacingFront{
	case true:
		b.Coords.X +=15
	case false:
		b.Coords.X -= 15
	}
}
//IsHit returns tru is it coliided woth something
func (b *Bullet) IsHit()bool{
	return b.Hit
}
//Render the bullet
func (b *Bullet) Render(screen, bullet *ebiten.Image)  {
	screen.DrawImage(bullet, b.Op)
}

//Collides telles whether the player is goona collide with te  object
func (b *Bullet) Collides(object image.Rectangle)bool {
	bullet := image.Rect(b.Coords.X , b.Coords.Y , b.Coords.X + 4, b.Coords.Y + 2)
	if bullet.Overlaps(object){
		
		//go log.Println(b.IsHit())
		return true
	}
	return false
}