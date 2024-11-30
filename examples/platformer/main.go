package main

import (
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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

const (
	minSpeed              = 0.07421875
	maxSpeed              = 2.5625
	maxWalkSpeed          = 1.5625
	maxFallSpeed          = 4.5
	maxFallSpeedCap       = 4
	minSlowDownSpeed      = 0.5625
	walkAcceleration      = 0.037109375
	runAcceleration       = 0.0556640625
	walkFriction          = 0.05078125
	skidFriction          = 0.1015625
	stompSpeed            = 4
	stompSpeedCap         = 4
	jumpSpeedNormal       = -4
	jumpSpeedRun          = -4
	jumpSpeedLong         = -5
	longJumpGravityNormal = 0.12
	longJumpGravityRun    = 0.11
	longJumpGravityLong   = 0.15
	gravity               = 0.43
	speedThreshold1       = 1
	speedThreshold2       = 2.3125
)

type PlayerController struct {
	// Constants (replaced magic numbers with named constants)
	MinSpeed         float64
	MaxSpeed         float64
	MaxWalkSpeed     float64
	MaxFallSpeed     float64
	MaxFallSpeedCap  float64
	MinSlowDownSpeed float64
	WalkAcceleration float64
	RunAcceleration  float64
	WalkFriction     float64
	SkidFriction     float64
	StompSpeed       float64
	StompSpeedCap    float64
	JumpSpeed        [3]float64
	LongJumpGravity  [3]float64
	Gravity          float64
	SpeedThresholds  [2]float64
	// states
	IsFacingLeft bool
	IsRunning    bool
	IsJumping    bool
	IsFalling    bool
	IsSkidding   bool
	IsCrouching  bool
	IsOnFloor    bool
	// private
	minSpeedValue       float64
	maxSpeedValue       float64
	accel               float64
	speedThresholdIndex int
}

func NewPlayerController() *PlayerController {

	pc := &PlayerController{
		MinSpeed:         minSpeed,
		MaxSpeed:         maxSpeed,
		MaxWalkSpeed:     maxWalkSpeed,
		MaxFallSpeed:     maxFallSpeed,
		MaxFallSpeedCap:  maxFallSpeedCap,
		MinSlowDownSpeed: minSlowDownSpeed,
		WalkAcceleration: walkAcceleration,
		RunAcceleration:  runAcceleration,
		WalkFriction:     walkFriction,
		SkidFriction:     skidFriction,
		StompSpeed:       stompSpeed,
		StompSpeedCap:    stompSpeedCap,

		JumpSpeed:       [3]float64{jumpSpeedNormal, jumpSpeedRun, jumpSpeedLong},
		LongJumpGravity: [3]float64{longJumpGravityNormal, longJumpGravityRun, longJumpGravityLong},
		Gravity:         gravity,
		SpeedThresholds: [2]float64{speedThreshold1, speedThreshold2},
		IsFacingLeft:    false,
		IsRunning:       false,

		IsJumping:   false,
		IsFalling:   false,
		IsSkidding:  false,
		IsCrouching: false,
		IsOnFloor:   false,

		speedThresholdIndex: 0,
	}

	pc.minSpeedValue = pc.MinSpeed
	pc.maxSpeedValue = pc.MaxSpeed
	pc.accel = pc.WalkAcceleration

	return pc
}

func (pc *PlayerController) SetPhyicsScale(s float64) {
	pc.MinSpeed *= s
	pc.MaxSpeed *= s
	pc.MaxWalkSpeed *= s
	pc.MaxFallSpeed *= s
	pc.MaxFallSpeedCap *= s
	pc.MinSlowDownSpeed *= s
	pc.WalkAcceleration *= s
	pc.RunAcceleration *= s
	pc.WalkFriction *= s
	pc.SkidFriction *= s
	pc.StompSpeed *= s
	pc.StompSpeedCap *= s
	pc.JumpSpeed[0] *= s
	pc.JumpSpeed[1] *= s
	pc.JumpSpeed[2] *= s
	pc.LongJumpGravity[0] *= s
	pc.LongJumpGravity[1] *= s
	pc.LongJumpGravity[2] *= s
	pc.Gravity *= s
	pc.SpeedThresholds[0] *= s
	pc.SpeedThresholds[1] *= s
}

func (pc *PlayerController) ProcessVelocity(vel [2]float64) [2]float64 {
	inputAxisX, inputAxisY := getAxis()

	if pc.IsOnFloor {
		pc.IsRunning = ebiten.IsKeyPressed(ebiten.KeyShift)
		pc.IsCrouching = ebiten.IsKeyPressed(ebiten.KeyDown)
		if pc.IsCrouching && inputAxisX != 0 {
			pc.IsCrouching = false
			inputAxisX = 0.0
		}
	}

	if pc.IsOnFloor {
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			pc.IsJumping = true
			speed := math.Abs(vel[0])
			pc.speedThresholdIndex = 0
			if speed >= pc.SpeedThresholds[1] {
				pc.speedThresholdIndex = 2
			} else if speed >= pc.SpeedThresholds[0] {
				pc.speedThresholdIndex = 1
			}

			vel[1] = pc.JumpSpeed[pc.speedThresholdIndex]

		}
	} else {
		gravityValue := pc.Gravity
		if ebiten.IsKeyPressed(ebiten.KeySpace) && pc.IsJumping && vel[1] < 0 {
			gravityValue = pc.LongJumpGravity[pc.speedThresholdIndex]
		}
		vel[1] += gravityValue
		if vel[1] > pc.MaxFallSpeedCap {
			vel[1] = pc.MaxFallSpeedCap
		}
	}

	// Update states
	if vel[1] > 0 {
		pc.IsJumping = false
		pc.IsFalling = true
	} else if pc.IsOnFloor {
		pc.IsFalling = false
	}

	if inputAxisX != 0 {
		if pc.IsOnFloor {
			if vel[0] != 0 {
				pc.IsFacingLeft = inputAxisX < 0.0
				pc.IsSkidding = vel[0] < 0.0 != pc.IsFacingLeft
			}
			if pc.IsSkidding {
				pc.minSpeedValue = pc.MinSlowDownSpeed
				pc.maxSpeedValue = pc.MaxWalkSpeed
				pc.accel = pc.SkidFriction
			} else if pc.IsRunning {
				pc.minSpeedValue = pc.MinSpeed
				pc.maxSpeedValue = pc.MaxSpeed
				pc.accel = pc.RunAcceleration
			} else {
				pc.minSpeedValue = pc.MinSpeed
				pc.maxSpeedValue = pc.MaxWalkSpeed
				pc.accel = pc.WalkAcceleration
			}
		} else if pc.IsRunning && math.Abs(vel[0]) > pc.MaxWalkSpeed {
			pc.maxSpeedValue = pc.MaxSpeed
		} else {
			pc.maxSpeedValue = pc.MaxWalkSpeed
		}
		targetSpeed := inputAxisX * pc.maxSpeedValue

		// Manually implementing moveToward()
		if vel[0] < targetSpeed {
			vel[0] += pc.accel
			if vel[0] > targetSpeed {
				vel[0] = targetSpeed
			}
		} else if vel[0] > targetSpeed {
			vel[0] -= pc.accel
			if vel[0] < targetSpeed {
				vel[0] = targetSpeed
			}
		}

	} else if pc.IsOnFloor && vel[0] != 0 {
		if !pc.IsSkidding {
			pc.accel = pc.WalkFriction
		}
		if inputAxisY != 0 {
			pc.minSpeedValue = pc.MinSlowDownSpeed
		} else {
			pc.minSpeedValue = pc.MinSpeed
		}
		if math.Abs(vel[0]) < pc.minSpeedValue {
			vel[0] = 0.0
		} else {
			// Manually implementing moveToward() for deceleration
			if vel[0] > 0 {
				vel[0] -= pc.accel
				if vel[0] < 0 {
					vel[0] = 0
				}
			} else {
				vel[0] += pc.accel
				if vel[0] > 0 {
					vel[0] = 0
				}
			}
		}
	}
	if math.Abs(vel[0]) < pc.MinSlowDownSpeed {
		pc.IsSkidding = false
	}

	return vel
}

func getAxis() (axisX, axisY float64) {
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		axisY -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		axisY += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		axisX -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		axisX += 1
	}
	return axisX, axisY
}
