package main

import (
	"image"
	_ "image/png"
	"log"
	"strconv"
	"sync/atomic"

	//"math"

	"github.com/hajimehoshi/ebiten/v2"
	//"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	//"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/lehangajanayake/MissionImposible/frontend/models"
)

//Game the game
type Game struct{
	Player models.Player
	inited bool
	Frames int
}

//Update updates  the game 
func (g *Game) Update()error{
	g.Player.IdleAnimation.CurrentFrame ++
	g.Player.WalkingAnimation.Animate = false
	if ebiten.IsKeyPressed(ebiten.KeyW){
		g.Player.WalkingAnimation.Animate = true
		g.Player.WalkingAnimation.CurrentFrame ++
		g.Player.Coords.Y --
	}else if ebiten.IsKeyPressed(ebiten.KeyS){
		g.Player.WalkingAnimation.Animate = true
		g.Player.WalkingAnimation.CurrentFrame ++
		g.Player.Coords.Y ++ 
		
	}

	if ebiten.IsKeyPressed(ebiten.KeyA){
		g.Player.WalkingAnimation.Animate = true
		g.Player.Coords.X --
		g.Player.FacingFront = false
	}else if ebiten.IsKeyPressed(ebiten.KeyD){
		g.Player.WalkingAnimation.Animate = true
		g.Player.WalkingAnimation.CurrentFrame ++
		g.Player.Coords.X ++
		g.Player.FacingFront = true
		
	
	}

	return nil
}
// Draw draws to the screen every update
func (g *Game) Draw(screen *ebiten.Image){
	g.Player.Op.GeoM.Reset()
	g.Player.Op.GeoM.Translate(-float64(g.Player.WalkingAnimation.FrameWidth/2), -float64(g.Player.WalkingAnimation.FrameHeight/2)) //,ake the axiz of the player in teh midlle instead of the upper left conner
	if g.Player.FacingFront{
		g.Player.Op.GeoM.Scale(1.5,1.5)
	}else{
		g.Player.Op.GeoM.Scale(-1.5,1.5)
	}
	g.Player.Op.GeoM.Translate(float64(g.Player.Coords.X),float64(g.Player.Coords.Y))
	
	if g.Player.WalkingAnimation.Animate {
		f := (g.Player.WalkingAnimation.CurrentFrame / 20) % len(g.Player.WalkingAnimation.AnimationArray)
		x, y := g.Player.WalkingAnimation.AnimationArray[f].X, g.Player.WalkingAnimation.AnimationArray[f].Y
		screen.DrawImage(g.Player.Img.SubImage(image.Rect(x, y, x + g.Player.WalkingAnimation.FrameWidth, y + g.Player.WalkingAnimation.FrameHeight)).(*ebiten.Image), g.Player.Op)
		g.Player.WalkingAnimation.CurrentFrame ++
	}else{
		f := (g.Player.IdleAnimation.CurrentFrame / 30) % len(g.Player.IdleAnimation.AnimationArray)
		x, y := g.Player.IdleAnimation.AnimationArray[f].X, g.Player.IdleAnimation.AnimationArray[f].Y
		screen.DrawImage(g.Player.Img.SubImage(image.Rect(x, y, x + g.Player.IdleAnimation.FrameWidth, y + g.Player.IdleAnimation.FrameHeight)).(*ebiten.Image), g.Player.Op)
		g.Player.IdleAnimation.CurrentFrame ++
	}
	ebitenutil.DebugPrint(screen, strconv.Itoa(int(ebiten.CurrentFPS())))
}

// Layout lays the screen
func (g *Game) Layout(outsideWidth, outsideHeight int)(int, int){
	return 640, 480
}

// func (g *Game) init(){
// 	defer func ()  {
// 		g.inited = true
// 	}()
// 	//var err error
// 	img,_, err := image.Decode(bytes.NewReader(images.Runner_png))
// 	if err != nil {
// 		panic("Cannot load the assets")
// 	}
// 	g.Player.Img = ebiten.NewImageFromImage(img)
// 	if err != nil{
// 		log.Fatal("Error loading the sprite: ", err.Error())
// 	}
// 	g.Player = models.Player{
// 		Coords : models.Coordinates{
// 			X: 0,
// 			Y: 0,
// 		},
// 		Op: &ebiten.DrawImageOptions{},
// 	}
// 	g.inited = true
// }

func main(){
	img,_, err := ebitenutil.NewImageFromFile("assets/runner.png")
	if err != nil {
		log.Fatalf("Cannot load the assets err : %w", err)
	}
	g := Game{
		Player: models.Player{
			Img: img,
			Op: &ebiten.DrawImageOptions{},
			WalkingAnimation: models.Animation{
				AnimationArray: []models.Coordinates{{X:0, Y:32},{X:32, Y:32},{X:64, Y:32},{X:96, Y:32},{X:128, Y:32},{X:160, Y:32},{X:192, Y:32},{X:224, Y:32}},
				FrameNum: 8,
				FrameWidth: 32,
				FrameHeight: 32,
				Animate: false,
			},
			IdleAnimation: models.Animation{
				AnimationArray: []models.Coordinates{{X:0, Y:0},{X:32, Y:0},{X:64, Y:0},{X:96, Y:0},{X:128, Y:0}},
				FrameNum: 5,
				FrameHeight: 32,
				FrameWidth: 32,
				Animate: false,

			},
			FacingFront: true,
		},
		Frames: 0,
	}
	if err := ebiten.RunGame(&g); err != nil || ebiten.IsKeyPressed(ebiten.KeyQ){
		log.Fatal(err)
	}
}