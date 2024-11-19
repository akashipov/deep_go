package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Node struct {
	Left   *Node
	Right  *Node
	Parent *Node
	Key    int
	Value  int
}

type OrderedMap struct {
	root *Node
	size int
}

func NewOrderedMap() OrderedMap {
	return OrderedMap{
		root: nil,
		size: 0,
	}
}

func (m *OrderedMap) Insert(key, value int) {
	var current *Node

	current = m.root
	for {
		if current == nil {
			m.root = &Node{Key: key, Value: value}
			m.size++
			return
		}
		if current.Key == key {
			current.Value = value
			return
		} else if current.Key < key {
			if current.Right == nil {
				current.Right = &Node{Parent: current, Key: key, Value: value}
				m.size++
				return
			}
			current = current.Right
		} else {
			if current.Left == nil {
				current.Left = &Node{Parent: current, Key: key, Value: value}
				m.size++
				return
			}
			current = current.Left
		}
	}
}

func findMax(n *Node) *Node {
	for n.Right != nil {
		n = n.Right
	}
	return n
}

func swap(n *Node, other *Node) {
	n.Value, other.Value = other.Value, n.Value
	n.Key, other.Key = other.Key, n.Key
}

func (m *OrderedMap) Erase(key int) {
	current := m.root
	for {
		if current == nil {
			return
		} else if current.Key == key {
			defer func() {
				m.size--
			}()
			if current.Left == nil {
				if current == m.root {
					m.root = current.Right
					current.Right.Parent = nil
					return
				}
				if current.Parent.Right == current {
					current.Parent.Right = current.Right
					return
				}
				current.Parent.Left = current.Right
				return
			}
			n := findMax(current.Left)
			swap(current, n)
			if n == current.Left {
				current.Left = n.Left
				return
			}
			n.Parent.Right = n.Left
			return
		} else if current.Key < key {
			current = current.Right
			continue
		}
		current = current.Left
	}
}

func (m *OrderedMap) Contains(key int) bool {
	var current *Node

	current = m.root
	for {
		if current == nil {
			return false
		}
		if current.Key == key {
			return true
		} else if current.Key < key {
			if current.Right == nil {
				return false
			}
			current = current.Right
		} else {
			if current.Left == nil {
				return false
			}
			current = current.Left
		}
	}
}

func (m *OrderedMap) Size() int {
	return m.size
}

func iter(action func(int, int), n *Node) {
	if n != nil {
		iter(action, n.Left)
		action(n.Key, n.Value)
		iter(action, n.Right)
	}
}

func (m *OrderedMap) ForEach(action func(int, int)) {
	iter(action, m.root)
}

func TestCircularQueue(t *testing.T) {
	data := NewOrderedMap()
	assert.Zero(t, data.Size())

	data.Insert(10, 10)
	data.Insert(5, 5)
	data.Insert(15, 15)
	data.Insert(2, 2)
	data.Insert(4, 4)
	data.Insert(12, 12)
	data.Insert(14, 14)

	assert.Equal(t, 7, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(3))
	assert.False(t, data.Contains(13))

	var keys []int
	expectedKeys := []int{2, 4, 5, 10, 12, 14, 15}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Erase(15)
	data.Erase(14)
	data.Erase(2)

	assert.Equal(t, 4, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(2))
	assert.False(t, data.Contains(14))

	keys = nil
	expectedKeys = []int{4, 5, 10, 12}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}
