package lsm

import (
	"errors"
	"strings"
)

type node struct {
	key     string
	payload []byte
	next    *node
}

type memtable struct {
	head *node
	size int
}

type memtableIterator struct {
	currentNode *node
}

func newMemTable() *memtable {
	return &memtable{
		head: nil,
		size: 0,
	}
}

func (mt *memtable) newMemTableIterator() *memtableIterator {
	return &memtableIterator{
		currentNode: mt.head,
	}
}

func (m *memtable) put(key string, payload []byte) {
	defer func() { m.size += len(payload) }()
	if m.head == nil {
		m.head = &node{
			key:     key,
			payload: payload,
			next:    nil,
		}

		return
	}

	n := m.head
	for {
		if n.next != nil {
			n = n.next
			continue
		}
		break
	}

	n.next = &node{
		key:     key,
		payload: payload,
		next:    nil,
	}
}

func (m *memtable) get(key string) ([]byte, error) {
	if m.head == nil {
		return nil, errors.New("head is nil")
	}

	n := m.head

	for {

		if strings.Compare(key, n.key) == 0 {
			return n.payload, nil
		}

		if n.next == nil {
			break
		}

		n = n.next
	}

	return nil, ErrNotFound
}

func (m *memtable) flush() {
	defer func() { m.size = 0 }()
	n := m.head
	for n != nil {
		t := n.next
		n = nil
		n = t
	}
	m.head = nil
}

func (it *memtableIterator) hasNext() bool {
	return it.currentNode != nil && it.currentNode.next != nil
}

func (it *memtableIterator) next() node {
	n := it.currentNode
	it.currentNode = n.next
	return *n
}
