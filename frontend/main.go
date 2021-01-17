package main

import (
	"image/color"
	_ "image/png"

	//"io/ioutil"
	"log"
	//"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/lehangajanayake/MissionImposible/frontend/models"
	"github.com/lehangajanayake/MissionImposible/frontend/network"
	"github.com/lehangajanayake/MissionImposible/frontend/ray"
)

//Game the game
type Game struct {
	Player                    *models.Player
	Client                    *network.Client
	Bullets                   []*models.Bullet
	BulletImg                 *ebiten.Image
	Frames                    int
	RayObjects                []ray.Object
	ShadowImg                 *ebiten.Image
	TriangleImg               *ebiten.Image
	Map                       *models.Map
	Camera                    *models.Camera
	ScreenWidth, ScreenHeight int
}



// Layout lays the screen
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ebiten.ScreenSizeInFullscreen()
}

func main() {
	//runtime.GOMAXPROCS(1)
	log.SetFlags(log.Ltime | log.Lshortfile)
	player, _, err := ebitenutil.NewImageFromFile("assets/hero_spritesheet.png")
	if err != nil {
		log.Fatalf("Cannot load the assets err : %v", err)
	}
	bullet, _, err := ebitenutil.NewImageFromFile("assets/bullet.png")
	if err != nil {
		log.Fatalf("Cannot load the assets err : %v", err)
	}

	w, h := ebiten.ScreenSizeInFullscreen()
	ebiten.SetWindowSize(w, h)
	g := Game{
		ScreenWidth:  w,
		ScreenHeight: h,
		Player: &models.Player{
			Img: player,
			Op:  &ebiten.DrawImageOptions{},
			Coords: models.Coordinates{
				X: w / 2,
				Y: h / 2,
			},
			WalkingAnimation: models.Animation{
				Name:        "Walking",
				StartX:      0,
				StartY:      100,
				FrameNum:    6,
				FrameWidth:  80,
				FrameHeight: 80,
				Animate:     false,
			},
			IdleAnimation: models.Animation{
				Name:        "Idle",
				StartX:      0,
				StartY:      0,
				FrameNum:    1,
				FrameHeight: 80,
				FrameWidth:  80,
				Animate:     false,
			},
			ShootingAnimation: models.Animation{
				Name:        "Shooting",
				StartX:      0,
				StartY:      0,
				FrameNum:    8,
				FrameHeight: 80,
				FrameWidth:  80,
				Animate:     false,
			},
			Gun: models.Gun{
				Bullets: 60,
			},
			FacingFront: true,
		},

		Frames:    0,
		BulletImg: bullet,
		Map: &models.Map{
			Op: &ebiten.DrawImageOptions{},
		},
		Camera: &models.Camera{
			Position: models.Coordinates{
				X: w / 2,
				Y: h / 2,
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
	for _, v := range g.Map.RayObjects {
		g.RayObjects = append(g.RayObjects, ray.Object{Edges: ray.Rect(v.X, v.Y, v.Width, v.Height)})
	}
	println(len(g.RayObjects), len(g.Map.RayObjects))
	if err := ebiten.RunGame(&g); err != nil {
		log.Fatal(err)
	}

}
