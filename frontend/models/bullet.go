package models

import "github.com/hajimehoshi/ebiten/v2"

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
	return b.Moving
}
//Render the bullet
func (b *Bullet) Render(screen, bullet *ebiten.Image)  {
	screen.DrawImage(bullet, b.Op)
}