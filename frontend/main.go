package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"sync"

	//"io/ioutil"
	"log"
	//"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/lafriks/go-tiled"
	"github.com/lehangajanayake/MissionImposible/frontend/models"
	"github.com/lehangajanayake/MissionImposible/frontend/network"
	"github.com/lehangajanayake/MissionImposible/frontend/ray"
)

//Game the game
type Game struct{
	Player *models.Player
	Client *network.Client
	Bullets []*models.Bullet
	BulletImg *ebiten.Image
	Frames int
	RayObjects []ray.Object
	ShadowImg *ebiten.Image
	TriangleImg *ebiten.Image
	Map *models.Map
	Camera *models.Camera
	ScreenWidth, ScreenHeight int
}

//Update updates  the game 
func (g *Game) Update()error{
	g.Player.LastAnimation = g.Player.IdleAnimation.Name
	g.Camera.Zoom = 1
	g.Camera.Move(models.Coordinates{X: g.ScreenWidth /2, Y:g.ScreenHeight/2})
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
	collide := make(chan bool, len(g.Map.TransparentObstacles) + len(g.Map.RayObjects))
	var wg sync.WaitGroup
	for _, obj := range g.Map.TransparentObstacles{
		wg.Add(1)
		go func(obj *tiled.Object, wg *sync.WaitGroup) {
			defer wg.Done()
			if g.Player.Collides(image.Rect(int(obj.X), int(obj.Y), int(obj.X + obj.Width), int(obj.Y + obj.Height))){
				collide <- true
			}	
		}(obj, &wg)
		
	}
	for _, obj := range g.Map.RayObjects{
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
	wg.Wait()
	close(collide)
	
	for v := range collide{
		if v {
			g.Player.Coords = g.Player.LastPos
			break
		}
	}
		
	
	
	if g.Player.LastPos != g.Player.Coords{
		coords := &network.Coordinates{X: g.Player.Coords.X, Y: g.Player.Coords.Y}
		g.Client.UpdatePlayerCoordsWrite <- coords.String()
		g.Player.LastPos = g.Player.Coords
	}

	switch g.Player.LastAnimation{
	case g.Player.IdleAnimation.Name:
		if !g.Player.IdleAnimation.Animate{
			g.Player.LastAnimation = g.Player.IdleAnimation.Name
			g.Client.UpdatePlayerAnimationWrite <- g.Player.IdleAnimation.Name
		}
	case g.Player.WalkingAnimation.Name:
		if !g.Player.WalkingAnimation.Animate{
			g.Player.LastAnimation = g.Player.WalkingAnimation.Name
			g.Client.UpdatePlayerAnimationWrite <- g.Player.WalkingAnimation.Name
		}
	case g.Player.ShootingAnimation.Name:
			if !g.Player.ShootingAnimation.Animate{
				g.Player.LastAnimation = g.Player.ShootingAnimation.Name
				g.Client.UpdatePlayerAnimationWrite <- g.Player.ShootingAnimation.Name
			}
	default:
		println("hello")
	}
	
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
		for _, obj := range g.Map.BlindSpots{
		if g.Player.Collides(image.Rect(int(obj.X), int(obj.Y), int(obj.X + obj.Width), int(obj.Y + obj.Height))){
			g.Camera.Zoom = 0.7
			g.Camera.Move(g.Player.Coords)
		}	
	}
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
	for _, v := range g.Client.Players{ 
		v.CurrentFrame++
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-float64(g.Player.WalkingAnimation.FrameWidth/2), -float64(g.Player.WalkingAnimation.FrameHeight/2)) //,ake the axis of the player in teh middle instead of the upper left conner
		if v.FacingFront{
			op.GeoM.Scale(0.5,0.5)
		}else{
			op.GeoM.Scale(-0.5,0.5)
		}
		op.GeoM.Translate(float64(v.Coords.X), float64(v.Coords.Y))
		switch v.Animation{
		case g.Player.IdleAnimation.Name:
			f := (v.CurrentFrame / 20) % g.Player.IdleAnimation.FrameNum
			x, y := g.Player.IdleAnimation.FrameWidth*f, g.Player.IdleAnimation.StartY
			g.Camera.View.DrawImage(g.Player.Img.SubImage(image.Rect(x, y, x + g.Player.IdleAnimation.FrameWidth, y + g.Player.IdleAnimation.FrameHeight)).(*ebiten.Image), op)
		case g.Player.WalkingAnimation.Name:
			f := (v.CurrentFrame / 10) % g.Player.WalkingAnimation.FrameNum
			x, y := g.Player.WalkingAnimation.FrameWidth*f, g.Player.WalkingAnimation.StartY
			g.Camera.View.DrawImage(g.Player.Img.SubImage(image.Rect(x, y, x + g.Player.WalkingAnimation.FrameWidth, y + g.Player.WalkingAnimation.FrameHeight)).(*ebiten.Image), op)
		case g.Player.ShootingAnimation.Name:
			f := (v.CurrentFrame / 3) % g.Player.ShootingAnimation.FrameNum
			x, y := g.Player.ShootingAnimation.FrameWidth*f, g.Player.ShootingAnimation.StartY
			g.Camera.View.DrawImage(g.Player.Img.SubImage(image.Rect(x, y, x + g.Player.ShootingAnimation.FrameWidth, y + g.Player.ShootingAnimation.FrameHeight)).(*ebiten.Image), op)
		}
		//g.Camera.View.DrawImage(g.Player.Img.SubImage(image.Rect(0, 0, g.Player.WalkingAnimation.FrameWidth,g.Player.WalkingAnimation.FrameHeight)).(*ebiten.Image), op)
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
	
	//render the trees
	g.Camera.View.DrawImage(g.Map.Trees, g.Map.Op)
	//render teh shadow
	g.ShadowImg.Fill(color.Black)
	//log.Println(g.RayObjects)
	rays := ray.Cast(float64(g.Player.Coords.X), float64(g.Player.Coords.Y), g.RayObjects)

	// Subtract ray triangles from shadow
	opt := &ebiten.DrawTrianglesOptions{}
	opt.Address = ebiten.AddressRepeat
	opt.CompositeMode = ebiten.CompositeModeSourceOut
	for i, line := range rays {
		nextLine := rays[(i+1)%len(rays)]

		// Draw triangle of area between rays
		v := ray.Vertices(float64(g.Player.Coords.X), float64(g.Player.Coords.Y), nextLine.X2, nextLine.Y2, line.X2, line.Y2)
		g.ShadowImg.DrawTriangles(v, []uint16{0, 1, 2}, g.TriangleImg, opt)
	}		
	op := &ebiten.DrawImageOptions{}
	op.ColorM.Scale(1, 1, 1, 0.7)
	g.Camera.View.DrawImage(g.ShadowImg, op)
	//render using the camera
	g.Camera.Render(screen)
	//clear the camera 
	g.Camera.View.Clear()
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Bullets Left: %v , Current TPS: %0.2f, Current FPS: %0.2f", g.Player.Gun.Bullets , ebiten.CurrentTPS(), ebiten.CurrentFPS()))
	//ebitenutil.DebugPrint(screen, g.Player.String())
}

// Layout lays the screen
func (g *Game) Layout(outsideWidth, outsideHeight int)(int, int){
	return ebiten.ScreenSizeInFullscreen()
}



func main(){
	//runtime.GOMAXPROCS(1)
	log.SetFlags(log.Ltime | log.Lshortfile)
	player,_, err := ebitenutil.NewImageFromFile("assets/hero_spritesheet.png")
	if err != nil {
		log.Fatalf("Cannot load the assets err : %v", err)
	}
	bullet,_, err := ebitenutil.NewImageFromFile("assets/bullet.png")
	if err != nil {
		log.Fatalf("Cannot load the assets err : %v", err)
	}

	w, h := ebiten.ScreenSizeInFullscreen()
	ebiten.SetWindowSize(w,h)
	g := Game{
		ScreenWidth: w,
		ScreenHeight: h,
		Player: &models.Player{
			Img: player,
			Op: &ebiten.DrawImageOptions{},
			Coords: models.Coordinates{
				X: w/2,
				Y: h/2,
			},
			WalkingAnimation: models.Animation{
				Name: "Walking",
				StartX: 0,
				StartY: 100,
				FrameNum: 6	,
				FrameWidth: 80,
				FrameHeight: 80,
				Animate: false,
			},
			IdleAnimation: models.Animation{
				Name: "Idle",
				StartX: 0,
				StartY: 0,
				FrameNum: 1,
				FrameHeight: 80,
				FrameWidth: 80,
				Animate: false,

			},
			ShootingAnimation: models.Animation{
				Name: "Shooting",
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
	g.Map.LoadTransparentObstacles()
	g.Map.LoadBlindSpots()
	g.Map.LoadRayObjects()
	g.Camera.View = ebiten.NewImage(g.Map.World.Size())
	g.ShadowImg = ebiten.NewImage(g.Map.World.Size())
	g.TriangleImg = ebiten.NewImage(g.Map.World.Size())
	g.TriangleImg.Fill(color.White)
	g.Client, err = network.Connect("192.168.1.7", "8080")
	if err != nil {
		log.Fatal("Error connecting to the server", err)
	}
	g.Client.Run()
	//g.Player.Coords = models.Coordinates{X: g.Player.Network.X, Y: g.Player.Network.Y}
	
	//converting the tiled object layer to ray objects
	for _, v := range g.Map.RayObjects{
		g.RayObjects = append(g.RayObjects, ray.Object{Edges: ray.Rect(v.X, v.Y,  v.Width, v.Height)})
	}
	println(len(g.RayObjects), len(g.Map.RayObjects))
	if err := ebiten.RunGame(&g); err != nil{
		log.Fatal(err)
	}
	
}
