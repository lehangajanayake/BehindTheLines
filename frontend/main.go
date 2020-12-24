package main

import (
	"errors"
	"fmt"
	"image"
	_ "image/png"
	"sync"

	//"io/ioutil"
	"log"
	//"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/lafriks/go-tiled"
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
		}	
	}
	collide := make(chan bool)
	var wg sync.WaitGroup
	for _, obj := range g.Map.Obstacles{
		wg.Add(1)
		go func(obj *tiled.Object, wg *sync.WaitGroup) {
			defer wg.Done()
			if g.Player.Collides(image.Rect(int(obj.X), int(obj.Y), int(obj.X + obj.Width), int(obj.Y + obj.Height))){
				collide <- true
			}	
			for i, v :=  range g.Bullets{
				if v.Collides(image.Rect(int(obj.X), int(obj.Y), int(obj.X + obj.Width), int(obj.Y + obj.Height))){
					g.Bullets = append(g.Bullets[:i], g.Bullets[i+1:]...)
				}
			}
		}(obj, &wg)
		
		
	
		
	}
	done := make(chan bool)
	
	go func(done chan bool) {
		for v := range collide{
			if v {
				g.Player.Coords = g.Player.LastPos
			}
		}
		done <- true
	}(done)
	
	wg.Wait()
	close(collide)
	<-done
	
	
	g.Player.LastPos = g.Player.Coords
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
	if g.Player.Coords.X > g.ScreenWidth /2 && g.Player.Coords.Y > g.ScreenHeight /2 {
		g.Camera.Move(g.Player.Coords)
	}else if g.Player.Coords.X > g.ScreenWidth /2 {
		g.Camera.Move(models.Coordinates{X:g.Player.Coords.X, Y:g.ScreenHeight /2})
	}else if g.Player.Coords.Y > g.ScreenHeight/2{
		g.Camera.Move(models.Coordinates{ X: g.ScreenWidth/2, Y: g.Player.Coords.Y})
	}
	
	return nil
}
// Draw draws to the screen every update
func (g *Game) Draw(screen *ebiten.Image){
	g.Camera.View.DrawImage(g.Map.World, g.Map.Op)
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
		g.Player.Render(g.Camera.View, g.Player.Img.SubImage(image.Rect(x, y, x + g.Player.WalkingAnimation.FrameWidth, y + g.Player.WalkingAnimation.FrameHeight)).(*ebiten.Image))
	}else if g.Player.IsIdle(){
		f := (g.Player.IdleAnimation.CurrentFrame / 20) % g.Player.IdleAnimation.FrameNum
		x, y := g.Player.IdleAnimation.FrameWidth*f, g.Player.IdleAnimation.StartY
		g.Player.Render(g.Camera.View, g.Player.Img.SubImage(image.Rect(x, y, x + g.Player.WalkingAnimation.FrameWidth, y + g.Player.WalkingAnimation.FrameHeight)).(*ebiten.Image))
		
	}else if g.Player.IsShooting(){
		f := (g.Player.ShootingAnimation.CurrentFrame / 3) % g.Player.ShootingAnimation.FrameNum
		x, y := g.Player.ShootingAnimation.FrameWidth*f, g.Player.ShootingAnimation.StartY
		g.Player.Render(g.Camera.View, g.Player.Img.SubImage(image.Rect(x, y, x + g.Player.WalkingAnimation.FrameWidth, y + g.Player.WalkingAnimation.FrameHeight)).(*ebiten.Image))
		
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
			v.Render(g.Camera.View, g.BulletImg)
			
		}
	}
	g.Camera.View.DrawImage(g.Map.Trees, g.Map.Op)
	g.Camera.Render(screen)
	g.Camera.View.Clear()
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Bullets Left: %v , Current TPS: %0.2f, Current FPS: %0.2f", g.Player.Gun.Bullets , ebiten.CurrentTPS(), ebiten.CurrentFPS()))
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
		Map: &models.Map{
			Op: &ebiten.DrawImageOptions{},
		},
		Camera: &models.Camera{
			Position : models.Coordinates{
				X: w/2,
				Y: h/2,
			},
			Op: &ebiten.GeoM{},
		},
	}
	err = g.Map.LoadMap("assets/Map/Map1.tmx")
	if err != nil {
		log.Fatalf("Err loading the map : %v", err)
	}
	g.Map.LoadObstacles(0)
	g.Camera.View = ebiten.NewImage(g.Map.World.Size())
	if err := ebiten.RunGame(&g); err != nil{
		log.Fatal(err)
	}
	
}
