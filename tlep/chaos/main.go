package chaos

import (
	"crypto/rand"
	"encoding/binary"
	"log"
)

// Define parameters for the Lorenz attractor.
const sigma = 10
const rho = 28
const beta = 8.0 / 3.0

// The function that represents the Lorenz system of equations.
func lorenz(x, y, z float64) (float64, float64, float64) {
	dx := sigma * (y - x)
	dy := x*(rho-z) - y
	dz := x*y - beta*z
	return dx, dy, dz
}

// The Runge-Kutta method (4th order) for a single step.
func (c *ChaosSystem) step() (float64, float64, float64) {
	x := c.X
	y := c.Y
	z := c.Z
	if c.Debug {
		log.Printf("In stepping Old vals : x = %f ; y = %f ; z = %f\n", x, y, z)
	}
	h := c.H
	k1x, k1y, k1z := lorenz(x, y, z)
	k2x, k2y, k2z := lorenz(x+h*k1x/2, y+h*k1y/2, z+h*k1z/2)
	k3x, k3y, k3z := lorenz(x+h*k2x/2, y+h*k2y/2, z+h*k2z/2)
	k4x, k4y, k4z := lorenz(x+h*k3x, y+h*k3y, z+h*k3z)

	// Combine the slopes and update the positions.
	newX := x + h*(k1x+2*k2x+2*k3x+k4x)/6
	newY := y + h*(k1y+2*k2y+2*k3y+k4y)/6
	newZ := z + h*(k1z+2*k2z+2*k3z+k4z)/6
	if c.Debug {
		log.Printf("In stepping New vals : x = %f ; y = %f ; z = %f\n", newX, newY, newZ)
	}
	c.X = newX
	c.Y = newY
	c.Z = newZ

	return newX, newY, newZ
}

type ChaosSystem struct {
	Version int
	X       float64
	Y       float64
	Z       float64
	H       float64
	Debug   bool
	XMult   float64
	YMult   float64
	ZMult   float64
}

func CreateNewChaosSystem() *ChaosSystem {
	var _c ChaosSystem
	_c.Version = LATEST_VERSION
	return &_c
}

func (c *ChaosSystem) Init(x, y, z, h float64) {
	c.X = x
	c.Y = y
	c.Z = z
	c.H = h
	c.Debug = false
	c.Version = LATEST_VERSION
}

func randomFloat64() (float64, error) {
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		return 0, err
	}
	return float64(binary.LittleEndian.Uint64(b[:])) / (1 << 64), nil
}

func (c *ChaosSystem) InitRandom() error {
	x, err := randomFloat64()
	if err != nil {
		return err
	}
	y, err := randomFloat64()
	if err != nil {
		return err
	}
	z, err := randomFloat64()
	if err != nil {
		return err
	}

	c.Init(x, y, z, 0.1)
	return nil
}

func (c *ChaosSystem) Burn(h uint) {
	var i uint
	for i = 0; i < h; i++ {
		c.step()
	}
}

func (c *ChaosSystem) SetMults(xM, yM, zM float64) {
	c.XMult = xM
	c.YMult = yM
	c.ZMult = zM
}

func (c *ChaosSystem) ToggleDebug() {
	c.Debug = !c.Debug
}

func (c *ChaosSystem) getByteVal() uint8 {
	var xM, yM, zM float64
	if c.XMult != 0 {
		xM = c.XMult
	} else {
		xM = 13
	}
	if c.YMult != 0 {
		yM = c.YMult
	} else {
		yM = 19
	}
	if c.ZMult != 0 {
		zM = c.ZMult
	} else {
		zM = 11
	}
	hash := func(val float64, mult float64) uint8 {
		scaledVal := val * mult * mult
		// Convert to an int and take the bits we're interested in
		intVal := int(scaledVal * scaledVal * scaledVal)
		// XOR the parts of the integer to spread out the impact
		return uint8((intVal>>8)^intVal) & 0xFF
	}
	xVal := hash(c.X, xM)
	yVal := hash(c.Y, yM)
	zVal := hash(c.Z, zM)
	return xVal ^ yVal ^ zVal
}

func (c *ChaosSystem) GetPassword(lenght uint) []byte {
	var i uint
	var ret []byte
	for i = 0; i < lenght; i++ {
		c.step()
		ret = append(ret, c.getByteVal())
	}
	return ret
}
