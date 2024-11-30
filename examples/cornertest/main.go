package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/setanarut/tilecollider"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 400
	screenHeight = 200
)

var (
	PlayerX = 0.
	PlayerY = 0.
	PlayerW = 50.
	PlayerH = 50.

	testVelX = 1.
	testVelY = 1.

	TileMap = [][]uint8{
		{0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0},
		{0, 0, 1, 0, 0},
		{0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0},
	}

	collider = tilecollider.NewCollider(TileMap, 50, 50)
)

type Game struct {
}

func (g *Game) Update() error {
	deltaX, deltaY := collider.Collide(PlayerX, PlayerY, PlayerW, PlayerH, testVelX, testVelY, nil)

	if testVelX != deltaX {
		fmt.Println("X COLLISION")
	}
	if testVelY != deltaY {
		fmt.Println("Y COLLISION")
	}

	PlayerX += deltaX
	PlayerY += deltaY

	if PlayerX > screenWidth || PlayerY > screenHeight {
		PlayerX = 0
		PlayerY = 0
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	// Draw tiles
	for y := 0; y < len(TileMap); y++ {
		for x := 0; x < len(TileMap[y]); x++ {
			if TileMap[y][x] != 0 {
				vector.DrawFilledRect(screen,
					float32(x*collider.TileSize[0]),
					float32(y*collider.TileSize[1]),
					float32(collider.TileSize[0]),
					float32(collider.TileSize[1]),
					color.Gray{Y: 128},
					true)
			}
		}
	}

	// Draw player
	vector.DrawFilledRect(screen,
		float32(PlayerX),
		float32(PlayerY),
		float32(PlayerW),
		float32(PlayerH),
		color.RGBA{47, 36, 254, 255},
		false)

	// Print collisions to the screen
	for i, c := range collider.Collisions {
		ebitenutil.DebugPrintAt(screen,
			fmt.Sprintf(
				"Tile ID: %d, Tile Coords: %v, Collision Normal: %v",
				c.TileID,
				c.TileCoords,
				c.Normal,
			), 20, 20+(i*20))
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	// ebiten.SetTPS(6)

	ebiten.SetWindowSize(screenWidth, screenHeight)
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
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
