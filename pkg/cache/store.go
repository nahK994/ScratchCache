package cache

import (
	"strconv"
	"time"
)

func InitCache() *Cache {
	cache := &Cache{
		info: make(map[string]Data),
	}
	go cache.activeExpiration()
	return cache
}

func (c *Cache) GET(key string) interface{} {
	c.mu.RLock()         // Acquire read lock
	defer c.mu.RUnlock() // Release read lock
	return c.info[key].val
}

func (c *Cache) SET(key string, value interface{}) {
	c.mu.Lock()         // Acquire write lock
	defer c.mu.Unlock() // Release write lock
	str, ok_str := value.(string)
	if ok_str {
		num, err := strconv.Atoi(str)
		if err == nil {
			c.info[key] = Data{
				val: num,
			}
		} else {
			c.info[key] = Data{
				val: str,
			}
		}
	} else {
		c.info[key] = Data{
			val: value,
		}
	}
}

func (c *Cache) EXISTS(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	_, exists := c.info[key]
	return exists
}

func (c *Cache) DEL(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.info, key)
}

func (c *Cache) INCR(key string) int {
	c.mu.Lock()
	defer c.mu.Unlock()

	val, _ := c.info[key].val.(int)
	c.info[key] = Data{
		val: val + 1,
	}
	return val + 1
}

func (c *Cache) DECR(key string) int {
	c.mu.Lock()
	defer c.mu.Unlock()

	val, _ := c.info[key].val.(int)
	c.info[key] = Data{
		val: val - 1,
	}
	return val - 1
}

func (c *Cache) LPUSH(key string, values []string) {
	c.mu.Lock()         // Acquire write lock
	defer c.mu.Unlock() // Release write lock

	oldData, _ := c.info[key].val.([]string)
	vals := make([]string, len(values)+len(oldData))
	for i := 0; i < len(values); i++ {
		vals[i] = values[len(values)-1-i]
	}
	copy(vals[len(values):], oldData)
	oldData = nil
	c.info[key] = Data{
		val: vals,
	}
}

func (c *Cache) RPUSH(key string, values []string) {
	c.mu.Lock()         // Acquire write lock
	defer c.mu.Unlock() // Release write lock

	oldData, _ := c.info[key].val.([]string)
	vals := make([]string, len(values)+len(oldData))
	copy(vals[0:], oldData)
	copy(vals[len(oldData):], values)
	oldData = nil
	c.info[key] = Data{
		val: vals,
	}
}

func processIdx(vals []string, idx int) int {
	if idx > len(vals) {
		idx = len(vals) - 1
	} else if idx < 0 {
		if -1*len(vals) > idx {
			idx = 0
		} else {
			idx = len(vals) + idx
		}
	}

	return idx
}

func (c *Cache) LRANGE(key string, startIdx, endIdx int) []string {
	c.mu.RLock()         // Acquire write lock
	defer c.mu.RUnlock() // Release write lock

	vals, _ := c.info[key].val.([]string)
	startIdx = processIdx(vals, startIdx)
	endIdx = processIdx(vals, endIdx)
	if len(vals) > 0 && startIdx <= endIdx {
		return vals[startIdx : endIdx+1]
	}
	return nil
}

func (c *Cache) LPOP(key string) {
	c.mu.Lock()         // Acquire write lock
	defer c.mu.Unlock() // Release write lock

	vals, _ := c.info[key].val.([]string)

	newVals := make([]string, len(vals)-1)
	copy(newVals, vals[1:])

	c.info[key] = Data{
		val: newVals,
	}
	vals = nil
}

func (c *Cache) RPOP(key string) {
	c.mu.Lock()         // Acquire write lock
	defer c.mu.Unlock() // Release write lock

	vals, _ := c.info[key].val.([]string)

	newVals := make([]string, len(vals)-1)
	copy(newVals, vals[0:len(vals)-1])

	c.info[key] = Data{
		val: newVals,
	}
	vals = nil
}

func (c *Cache) FLUSHALL() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for k := range c.info {
		delete(c.info, k)
	}
}

func (c *Cache) EXPIRE(key string, ttl int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	val := c.info[key]

	val.expiryTime = time.Now().Add(time.Second * time.Duration(ttl))
	c.info[key] = val
}

func (c *Cache) activeExpiration() {
	for {
		time.Sleep(60 * time.Second) // Check every second
		c.mu.Lock()
		for key, item := range c.info {
			if time.Now().After(item.expiryTime) {
				delete(c.info, key)
			}
		}
		c.mu.Unlock()
	}
}
