package lsm

import (
	"fmt"
	"os"
	"strings"
)

// WAL write ahead log
type WAL struct {
	f *os.File
}

func newWAL(f *os.File) *WAL {
	return &WAL{
		f: f,
	}
}

// Size
func (w *WAL) Size() (int64, error) {
	stat, err := w.f.Stat()
	if err != nil {
		return 0, fmt.Errorf("file.Size: %v", err)
	}
	return stat.Size(), nil
}

func (w *WAL) appendEntry(op, key, value string) error {
	entry := fmt.Sprintf("%s|%s|%s\n", strings.ToUpper(op), key, value)
	_, err := w.f.WriteString(entry)
	if err != nil {
		return fmt.Errorf("file.WriteString: %v", err)
	}
	return nil
}

func (w *WAL) flush() error {
	return w.f.Truncate(0)
}
