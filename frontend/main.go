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
	GuardAnimation            map[string]*models.Animation
	NinjaAnimation            map[string]*models.Animation
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
	log.SetFlags(log.Ltime | log.Lshortfile)
	player, _, err := ebitenutil.NewImageFromFile("assets/hero_spritesheet.png")
	if err != nil {
		log.Fatalf("Cannot load the assets err : %v", err)
	}
	bullet, _, err := ebitenutil.NewImageFromFile("assets/bullet.png")
	if err != nil {
		log.Fatalf("Cannot load the assets err : %v", err)
	}
	guardWalking, _, err := ebitenutil.NewImageFromFile("assets/Dungeon/sprtitesheets/Guard/spr_Walk_strip.png")
	if err != nil {
		log.Fatalf("Cannot load the assets err : %v", err)
	}
	guardIdle, _, err := ebitenutil.NewImageFromFile("assets/Dungeon/sprtitesheets/Guard/spr_Idle_strip.png")
	if err != nil {
		log.Fatalf("Cannot load the assets err : %v", err)
	}
	guardAttack, _, err := ebitenutil.NewImageFromFile("assets/Dungeon/sprtitesheets/Guard/spr_SpinAttack_strip.png")
	if err != nil {
		log.Fatalf("Cannot load the assets err : %v", err)
	}
	guardDie, _, err := ebitenutil.NewImageFromFile("assets/Dungeon/sprtitesheets/Guard/spr_Death_strip.png")
	if err != nil {
		log.Fatalf("Cannot load the assets err : %v", err)
	}
	ninjaWalking, _, err := ebitenutil.NewImageFromFile("assets/Dungeon/sprtitesheets/Ninja/spr_ArcherRun_strip_NoBkg.png")
	if err != nil {
		log.Fatalf("Cannot load the assets err : %v", err)
	}
	ninjaMelee, _, err := ebitenutil.NewImageFromFile("assets/Dungeon/sprtitesheets/Ninja/spr_ArcherMelee_strip_NoBkg.png")
	if err != nil {
		log.Fatalf("Cannot load the assets err : %v", err)
	}
	ninjaIdle, _, err := ebitenutil.NewImageFromFile("assets/Dungeon/sprtitesheets/Ninja/spr_ArcherIdle_strip_NoBkg.png")
	if err != nil {
		log.Fatalf("Cannot load the assets err : %v", err)
	}
	ninjaDie, _, err := ebitenutil.NewImageFromFile("assets/Dungeon/sprtitesheets/Ninja/spr_ArcherDeath_strip_NoBkg.png")
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
			Gun: models.Gun{
				Bullets: 60,
			},
			FacingFront: true,
		},
		GuardAnimation: map[string]*models.Animation{
			"Walking": {
				Img:          guardWalking,
				StartX:       8,
				StartY:       0,
				FrameNum:     8,
				FrameHeight:  94,
				FrameWidth:   170,
				CurrentFrame: 0,
			},
			"Attacking": {
				Img:          guardAttack,
				StartX:       8,
				StartY:       0,
				FrameNum:     30,
				FrameHeight:  94,
				FrameWidth:   170,
				CurrentFrame: 0,
			},
			"Idle": {
				Img:          guardIdle,
				StartX:       8,
				StartY:       0,
				FrameNum:     16,
				FrameHeight:  94,
				FrameWidth:   170,
				CurrentFrame: 0,
			},
			"Die": {
				Img:          guardDie,
				StartX:       8,
				StartY:       0,
				FrameNum:     8,
				FrameHeight:  94,
				FrameWidth:   170,
				CurrentFrame: 0,
			},
		},
		NinjaAnimation: map[string]*models.Animation{
			"Walking": {
				Img:          ninjaWalking,
				StartX:       0,
				StartY:       0,
				FrameNum:     8,
				FrameHeight:  128,
				FrameWidth:   128,
				CurrentFrame: 0,
			},
			"Idle": {
				Img:          ninjaIdle,
				StartX:       0,
				StartY:       0,
				FrameNum:     8,
				FrameHeight:  128,
				FrameWidth:   128,
				CurrentFrame: 0,
			},
			"Attacking": {
				Img:          ninjaMelee,
				StartX:       0,
				StartY:       0,
				FrameNum:     8,
				FrameHeight:  128,
				FrameWidth:   128,
				CurrentFrame: 0,
			},
			"Death": {
				Img:          ninjaDie,
				StartX:       0,
				StartY:       0,
				FrameNum:     24,
				FrameHeight:  128,
				FrameWidth:   128,
				CurrentFrame: 0,
			},
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
	g.Client, err = network.Connect("localhost", "8080")
	if err != nil {
		log.Fatal("Error connecting to the server", err)
	}
	err = g.Client.Run(g.Player)
	if g.Player.Guard {
		g.Player.Animations = g.GuardAnimation
	} else {
		g.Player.Animations = g.NinjaAnimation
	}
	if err != nil {
		log.Fatal("Error loadign the player from the server, ", err)
	}
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
