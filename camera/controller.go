package camera

import (
	"porridgo/datatypes"
	"porridgo/window"
)

type Controller struct {
	Camera     *Camera
	Speed      float32
	velocity   datatypes.Vec3f
	lastCursor datatypes.Vec3f
	rVel       datatypes.Vec3f
}

func (c *Controller) ProcessKey(key window.Key, action window.Action, modifier window.ModifierKey) {
	if action == window.Press || action == window.Repeat {
		switch key {
		case window.KeyA:
			c.velocity.X = -c.Speed
		case window.KeyD:
			c.velocity.X = c.Speed
		case window.KeyW:
			c.velocity.Z = c.Speed
		case window.KeyS:
			c.velocity.Z = -c.Speed
		case window.KeyE:
			c.velocity.Y = c.Speed
		case window.KeyQ:
			c.velocity.Y = -c.Speed
		}
	}
	if action == window.Release {
		switch key {
		case window.KeyA:
			c.velocity.X = 0.0
		case window.KeyD:
			c.velocity.X = 0.0
		case window.KeyW:
			c.velocity.Z = 0.0
		case window.KeyS:
			c.velocity.Z = 0.0
		case window.KeyE:
			c.velocity.Y = 0.0
		case window.KeyQ:
			c.velocity.Y = 0.0
		}
	}
}

func (c *Controller) ProcessMouse(x float32, y float32) {
	current := datatypes.NewVec3f(x, y, 0.0)
	c.rVel = c.lastCursor.Sub(current)
	c.lastCursor = current
}

func (c *Controller) UpdateCamera() {
	c.Camera.Position = c.Camera.Position.Add(c.velocity)
	c.Camera.Rotation = c.Camera.Rotation.Add(c.rVel)
}
