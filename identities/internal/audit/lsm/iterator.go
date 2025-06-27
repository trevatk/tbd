package lsm

import (
	"bufio"
	"os"
)

// Iterator
type iterator struct {
	lsm            *LSM
	currentIndex   int
	currentSSTable *os.File
	scanner        *bufio.Scanner
}

// HasNext
func (it *iterator) HasNext() bool {
	n := it.scanner.Scan()
	if n && it.currentIndex != len(it.lsm.sstables) {
		return true
	}

	// scanner will not return EOF
	// so if unable to scan then
	// open new scanner with next sstable in slice
	if !n {
		if it.currentIndex+1 > len(it.lsm.sstables) {
			return false
		}

		it.currentIndex += 1
		it.scanner = bufio.NewScanner(it.lsm.sstables[it.currentIndex])
		return true
	}

	return false
}

// Next
func (it *iterator) Next() (key string, value []byte, err error) {
	keyvalue, err := unmarshalKeyValue(it.scanner.Bytes())
	if err != nil {
		return
	}

	err = it.scanner.Err()
	if err != nil {
		return
	}

	key = keyvalue.Key
	value = keyvalue.Value

	return
}
