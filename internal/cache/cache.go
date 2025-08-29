package cache

import (
	"sync"
	"time"

	"github.com/nais/kolide-api-proxy/internal/kolide"
)

type Cache struct {
	devices      []kolide.Device
	lastModified time.Time
	lock         sync.RWMutex
}

func New() *Cache {
	return &Cache{
		devices: make([]kolide.Device, 0),
	}
}

func (c *Cache) SetDevices(devices []kolide.Device) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.devices = devices
	c.lastModified = time.Now().UTC()
}

func (c *Cache) GetDevices() []kolide.Device {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.devices
}

func (c *Cache) LastModified() time.Time {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.lastModified
}
