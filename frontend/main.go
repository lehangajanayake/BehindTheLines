package main

import (
	"errors"
	"fmt"
	"image"
	_ "image/png"

	//"io/ioutil"
	"log"
	//"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/lehangajanayake/MissionImposible/frontend/models"
)

//Game the game
type Game struct{
	Player models.Player
	Bullets []*models.Bullet
	BulletImg *ebiten.Image
	Frames int
	Map *models.Map
	Camera *models.Camera
	ScreenWidth, ScreenHeight int
}

//Update updates  the game 
func (g *Game) Update()error{
	g.Player.IdleAnimation.Reset()
	g.Player.WalkingAnimation.Reset()
	
	
	
	if g.Player.IsShooting(){
		g.Player.Shoot()
		return nil
	}

	if ebiten.IsKeyPressed(ebiten.KeyD) {
		if ebiten.IsKeyPressed(ebiten.KeyShift) && !(ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyS)){
			g.Player.Run("F")
		}else{
			g.Player.Walk("F")
		}
	}else if ebiten.IsKeyPressed(ebiten.KeyA) {
		if ebiten.IsKeyPressed(ebiten.KeyShift) && !(ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyS)){
			g.Player.Run("B")
		}else{
			g.Player.Walk("B")
		}	
	}
	if ebiten.IsKeyPressed(ebiten.KeyW){
		if ebiten.IsKeyPressed(ebiten.KeyShift) && !(ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyD)){
			g.Player.Run("U")
		}else{
			g.Player.Walk("U")
		}	
	}else if ebiten.IsKeyPressed(ebiten.KeyS) {
		if ebiten.IsKeyPressed(ebiten.KeyShift) && !(ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyD)){
			g.Player.Run("D")
		}else{
			g.Player.Walk("D")
		}	}
	

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft){
		g.Player.Shoot()
		bullet := new(models.Bullet)
		g.Bullets = append(g.Bullets, bullet.New(g.Player.Coords, g.Player.FacingFront))
	}

	if g.Player.IsIdle(){
		g.Player.Idle()
	}

	if ebiten.IsKeyPressed(ebiten.KeyQ){
		return errors.New("Game Exited by the user")
	}
	for _, obj := range g.Map.Obstacles{
		if g.Player.Collides(image.Rect(int(obj.X), int(obj.Y), int(obj.X + obj.Width), int(obj.Y + obj.Height))){
			if g.Player.LastPos != g.Player.Coords {
				g.Player.Coords = g.Player.LastPos
				return nil
			}
		}
		
	}
	g.Camera.Position = g.Player.Coords
	g.Player.LastPos = g.Player.Coords
	return nil
}
// Draw draws to the screen every update
func (g *Game) Draw(screen *ebiten.Image){
	
	//var colliding bool
	screen.DrawImage(g.Map.Img, g.Map.Op)
	g.Player.Op.GeoM.Reset()
	g.Player.Op.GeoM.Translate(-float64(g.Player.WalkingAnimation.FrameWidth/2), -float64(g.Player.WalkingAnimation.FrameHeight/2)) //,ake the axis of the player in teh middle instead of the upper left conner
	if g.Player.FacingFront{
		g.Player.Op.GeoM.Scale(0.5,0.5)
	}else{
		g.Player.Op.GeoM.Scale(-0.5,0.5)
	}
	g.Player.Op.GeoM.Translate(float64(g.Player.Coords.X),float64(g.Player.Coords.Y))
	
	if g.Player.IsWalking(){
		f := (g.Player.WalkingAnimation.CurrentFrame / 10 ) % g.Player.WalkingAnimation.FrameNum
		x, y := g.Player.WalkingAnimation.FrameWidth*f, g.Player.WalkingAnimation.StartY
		g.Player.Render(screen, g.Player.Img.SubImage(image.Rect(x, y, x + g.Player.WalkingAnimation.FrameWidth, y + g.Player.WalkingAnimation.FrameHeight)).(*ebiten.Image))
	}else if g.Player.IsIdle(){
		f := (g.Player.IdleAnimation.CurrentFrame / 20) % g.Player.IdleAnimation.FrameNum
		x, y := g.Player.IdleAnimation.FrameWidth*f, g.Player.IdleAnimation.StartY
		g.Player.Render(screen, g.Player.Img.SubImage(image.Rect(x, y, x + g.Player.WalkingAnimation.FrameWidth, y + g.Player.WalkingAnimation.FrameHeight)).(*ebiten.Image))
		
	}else if g.Player.IsShooting(){
		f := (g.Player.ShootingAnimation.CurrentFrame / 3) % g.Player.ShootingAnimation.FrameNum
		x, y := g.Player.ShootingAnimation.FrameWidth*f, g.Player.ShootingAnimation.StartY
		g.Player.Render(screen, g.Player.Img.SubImage(image.Rect(x, y, x + g.Player.WalkingAnimation.FrameWidth, y + g.Player.WalkingAnimation.FrameHeight)).(*ebiten.Image))
		
		if g.Player.ShootingAnimation.CurrentFrame == g.Player.ShootingAnimation.FrameNum*3{ //done shooting
			g.Player.ShootingAnimation.Reset()
			g.Player.ShootingAnimation.CurrentFrame = 1
			g.Player.Gun.Shoot()
		}
	}
	if len(g.Bullets) != 0{
		for _, v := range g.Bullets{
			v.Move()
			v.Op.GeoM.Reset()
			v.Op.GeoM.Scale(4,2)
			v.Op.GeoM.Translate(float64(v.Coords.X), float64(v.Coords.Y))
			v.Render(screen, g.BulletImg)
		}
	}
	g.Camera.Render(screen, )
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Bullets Left: %v", g.Player.Gun.Bullets))
}

// Layout lays the screen
func (g *Game) Layout(outsideWidth, outsideHeight int)(int, int){
	return ebiten.ScreenSizeInFullscreen()
}



func main(){
	player,_, err := ebitenutil.NewImageFromFile("assets/hero_spritesheet.png")
	if err != nil {
		log.Fatalf("Cannot load the assets err : %v", err)
	}
	bullet,_, err := ebitenutil.NewImageFromFile("assets/bullet.png")
	if err != nil {
		log.Fatalf("Cannot load the assets err : %v", err)
	}

	ebiten.SetFullscreen(true)
	w, h := ebiten.ScreenSizeInFullscreen()
	g := Game{
		ScreenWidth: w,
		ScreenHeight: h,
		Player: models.Player{
			Img: player,
			Op: &ebiten.DrawImageOptions{},
			Coords: models.Coordinates{
				X: w/2,
				Y: h/2,
			},
			WalkingAnimation: models.Animation{
				StartX: 0,
				StartY: 100,
				FrameNum: 6	,
				FrameWidth: 80,
				FrameHeight: 80,
				Animate: false,
			},
			IdleAnimation: models.Animation{
				StartX: 0,
				StartY: 0,
				FrameNum: 1,
				FrameHeight: 80,
				FrameWidth: 80,
				Animate: false,

			},
			ShootingAnimation: models.Animation{
				StartX: 0,
				StartY: 0,
				FrameNum: 8,
				FrameHeight: 80,
				FrameWidth: 80,
				Animate: false,
			},
			Gun: models.Gun{
				Bullets: 60,
			},
			FacingFront: true,
		},
		Frames: 0,
		BulletImg: bullet,
		Map: new(models.Map),
		Camera: &models.Camera{
			Position : models.Coordinates{
				X: w/2,
				Y: h/2,
			},
		},
	}
	err = g.Map.LoadMap("assets/grass.tmx")
	if err != nil {
		log.Fatalf("Err loading the map : %v", err)
	}
	g.Map.LoadObstacles(0)
	if err := ebiten.RunGame(&g); err != nil{
		log.Fatal(err)
	}
	
}

//Iscolliding
func IsColliding(object1, object2 image.Rectangle) bool{
	if object1.Min.X < object2.Max.X &&
	object1.Max.X > object2.Min.X &&
	object1.Min.Y < object2.Max.Y &&
	object1.Max.Y > object2.Min.Y{
		return true
	}
	return false
}