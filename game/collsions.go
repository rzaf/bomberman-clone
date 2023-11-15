package game

import (
	"sync"

	ray "github.com/gen2brain/raylib-go/raylib"
)

type Collision interface {
	GetCollision() ray.Rectangle
	Enabled() bool
}

type collisionObjects struct {
	objects []Collision
	sync.Mutex
}

var (
	ShowCollsions    bool = false
	CollisionManager collisionObjects
)

func (c *collisionObjects) Reset() {
	t := c.TryLock()
	if t {
		defer c.Unlock()
	}
	c.objects = make([]Collision, 0, 250)
}

func (c *collisionObjects) Add(obj Collision) {
	t := c.TryLock()
	if t {
		defer c.Unlock()
	}
	c.objects = append(c.objects, obj)
}

func (c *collisionObjects) Remove(obj Collision) {
	t := c.TryLock()
	if t {
		defer c.Unlock()
	}
	for index, col := range c.objects {
		if col == obj {
			c.objects[index] = c.objects[len(c.objects)-1]
			c.objects = c.objects[:len(c.objects)-1]
			return
		}
	}
}

func (c *collisionObjects) CheckCollisions(rec ray.Rectangle) []Collision {
	t := c.TryLock()
	if t {
		defer c.Unlock()
	}
	collisions := make([]Collision, 0)
	for _, collision := range c.objects {
		if collision.Enabled() && ray.CheckCollisionRecs(rec, collision.GetCollision()) {
			collisions = append(collisions, collision)
		}
	}
	return collisions
}

func (c *collisionObjects) CheckCollision(rec ray.Rectangle) Collision {
	t := c.TryLock()
	if t {
		defer c.Unlock()
	}
	for _, collision := range c.objects {
		if collision.Enabled() && ray.CheckCollisionRecs(rec, collision.GetCollision()) {
			return collision
		}
	}
	return nil
}
