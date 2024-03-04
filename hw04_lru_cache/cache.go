package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   Key
	value interface{}
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	v, ok := c.items[key]
	if ok {
		c.queue.MoveToFront(v)
		item := v.Value.(cacheItem)
		return item.value, ok
	}
	return nil, ok
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	v, ok := c.items[key]
	if ok {
		c.queue.Remove(v)
		cItem := cacheItem{
			key:   key,
			value: value,
		}
		item := c.queue.PushFront(cItem)
		cItem.value = value
		c.items[key] = item
	} else {
		item := c.queue.PushFront(cacheItem{
			key:   key,
			value: value,
		})
		c.items[key] = item
	}

	if c.queue.Len() > c.capacity {
		item := c.queue.Back()
		c.queue.Remove(item)
		cItem := item.Value.(cacheItem)
		delete(c.items, cItem.key)
	}
	return ok
}

func (c *lruCache) Clear() {
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
