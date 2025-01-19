package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/setanarut/tilecollider"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 200
	screenHeight = 200
)

var (
	PlayerX = 20.
	PlayerY = 20.
	PlayerW = 50.
	PlayerH = 50.

	TileMap = [][]uint8{
		{0, 0, 0, 0, 0},
		{0, 1, 0, 0, 0},
		{0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0},
	}

	collider = tilecollider.NewCollider(TileMap, 50, 50)
)

func init() {
	// Set the collider to always check for static collisions (no movement)
	collider.StaticCheck = true
}

type Game struct {
}

func (g *Game) Update() error {

	// Teleport player to mouse position
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		PlayerX = float64(x)
		PlayerY = float64(y)
	}
	// Toggle static check
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		collider.StaticCheck = !collider.StaticCheck
	}

	// Collide with zero movement (teleport)
	deltaX, deltaY := collider.Collide(PlayerX, PlayerY, PlayerW, PlayerH, 0, 0, nil)

	if deltaX != 0 || deltaY != 0 {
		collisionsLabel = fmt.Sprintf("Collision! %v", collider.Collisions)
	}

	PlayerX += deltaX
	PlayerY += deltaY

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

	ebitenutil.DebugPrintAt(screen, collisionsLabel, 10, 10)

}

var collisionsLabel string

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
