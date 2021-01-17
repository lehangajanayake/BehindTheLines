package models

import (
	//"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
)

//Gun model gun
type Gun struct {
	Sound   audio.Context
	Bullets int
}

//Shoot shoots one bullet at a time
func (g *Gun) Shoot() {
	g.Bullets--
}
