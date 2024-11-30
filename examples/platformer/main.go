package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/setanarut/kamera/v2"
	"github.com/setanarut/tilecollider"
)

func main() {
	// ebiten.SetTPS(10)
	ebiten.SetWindowSize(512, 512)
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}

var TileMap = [][]uint8{
	{1, 0, 1, 0, 1, 1, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 1},
	{1, 1, 0, 0, 0, 1, 0, 1},
	{0, 0, 0, 1, 0, 1, 0, 1},
	{0, 0, 0, 0, 0, 1, 0, 1},
	{1, 0, 1, 1, 1, 1, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 1},
	{1, 1, 1, 1, 1, 1, 1, 1}}

func init() {
	Controller.SetPhyicsScale(2.2)
	cam.LerpEnabled = true
}

var infoText string

func Translate(box *[4]float64, x, y float64) {
	box[0] += x
	box[1] += y
}

var (
	Offset   = [2]int{0, 0}
	GridSize = [2]int{8, 8}
	TileSize = [2]int{64, 64}
	Box      = [4]float64{70, 70, 24, 32}
	Vel      = [2]float64{0, 4}
	cam      = kamera.NewCamera(Box[0], Box[1], 512, 512)
)

var Controller = NewPlayerController()
var collider = tilecollider.NewCollider(TileMap, TileSize[0], TileSize[1])

func (g *Game) Update() error {
	if Vel[1] < 0 {
		Controller.IsOnFloor = false
	}
	Vel = Controller.ProcessVelocity(Vel)
	dx, dy := collider.Collide(
		Box[0],
		Box[1],
		Box[2],
		Box[3],
		Vel[0],
		Vel[1],
		func(ci []tilecollider.CollisionInfo[uint8], f1, f2 float64) {

			for _, v := range ci {
				if v.Normal[1] == -1 {
					Controller.IsOnFloor = true
				}
				if v.Normal[1] == 1 {
					Controller.IsJumping = false
					Vel[1] = 0
				}
				if v.Normal[0] == 1 || v.Normal[0] == -1 {
					Vel[0] = 0
				}
			}
		},
	)
	Translate(&Box, dx, dy)
	cam.LookAt(Box[0], Box[1])
	return nil
}

func (g *Game) Layout(w, h int) (int, int) {
	return 512, 512
}

type Game struct{}

func (g *Game) Draw(s *ebiten.Image) {

	for y, row := range TileMap {
		for x, value := range row {
			if value != 0 {
				px, py := float64(x*TileSize[0]), float64(y*TileSize[1])
				geom := &ebiten.GeoM{}
				cam.ApplyCameraTransform(geom)
				px, py = geom.Apply(px, py)
				vector.DrawFilledRect(
					s,
					float32(px),
					float32(py),
					float32(TileSize[0]),
					float32(TileSize[1]),
					color.Gray{127},
					false,
				)
			}
		}
	}

	// draw collided tiles
	for _, col := range collider.Collisions {

		px, py := float64(col.TileCoords[0]*collider.TileSize[0]), float64(col.TileCoords[1]*collider.TileSize[1])
		geom := &ebiten.GeoM{}
		cam.ApplyCameraTransform(geom)
		px, py = geom.Apply(px, py)

		vector.DrawFilledRect(
			s,
			float32(px),
			float32(py),
			float32(collider.TileSize[0]),
			float32(collider.TileSize[1]),
			color.RGBA{R: 255, G: 255, B: 0, A: 30},
			false,
		)
	}

	// draw player
	x, y := Box[0], Box[1]
	geom := &ebiten.GeoM{}
	cam.ApplyCameraTransform(geom)
	x, y = geom.Apply(x, y)
	vector.DrawFilledRect(
		s,
		float32(x),
		float32(y),
		float32(Box[2]),
		float32(Box[3]),
		color.Gray{180},
		false,
	)

	ebitenutil.DebugPrint(s, infoText)
}

func Axis() (axisX, axisY float64) {
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		axisY -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		axisY += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		axisX -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		axisX += 1
	}
	return axisX, axisY
}
