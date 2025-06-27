package lsm

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/fx"
	"google.golang.org/protobuf/proto"

	pb "github.com/structx/tbd/lib/protocol/lsm/v1"
	"github.com/structx/tbd/lib/setup"
)

type Store interface {
	Get(string) ([]byte, error)
	Put(string, []byte, map[string]string, int64) error
	Iterator() Iterator
}

type Iterator interface{}

// LSM
type LSM struct {
	indice      []*os.File
	sstDir      string
	sstables    []*os.File
	memtable    *memtable
	wal         *WAL
	flushToDisk int64
}

var _ Store = (*LSM)(nil)

type Params struct {
	fx.In

	Config *setup.KeyValue
}

type Result struct {
	fx.Out

	LSM *LSM
}

// Module
var Module = fx.Module("lsm", fx.Provide(func(p Params) (Result, error) {
	lsm, err := New(p.Config.Dir)
	if err != nil {
		return Result{}, fmt.Errorf("failed to create lsm: %w", err)
	}
	return Result{
		LSM: lsm,
	}, nil
}))

func New(dir string) (*LSM, error) {
	filePath := filepath.Clean(dir)
	lsm := &LSM{
		indice:      make([]*os.File, 0),
		sstables:    make([]*os.File, 0),
		memtable:    newMemTable(),
		sstDir:      filePath,
		flushToDisk: 10000,
	}

	entries, err := os.ReadDir(filePath)
	if err != nil {
		return nil, fmt.Errorf("os.ReadDir: %v", err)
	}

	if len(entries) == 0 {

		// create required files

		sstFilePath := filepath.Join(filePath, "sstable_0.data")
		f, err := os.OpenFile(sstFilePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
		if err != nil {
			return nil, fmt.Errorf("os.OpenFile: %v", err)
		}
		lsm.sstables = append(lsm.sstables, f)

		idxFilePath := filepath.Join(filePath, "index_0.log")
		f, err = os.OpenFile(idxFilePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
		if err != nil {
			return nil, fmt.Errorf("os.OpenFile: %v", err)
		}
		lsm.indice = append(lsm.indice, f)

		walFilePath := filepath.Join(filePath, "wal.log")
		f, err = os.OpenFile(walFilePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
		if err != nil {
			return nil, fmt.Errorf("os.OpenFile: %v", err)
		}
		lsm.wal = newWAL(f)

		return lsm, nil
	}

	for _, entry := range entries {

		if entry.IsDir() {
			continue
		}

		if strings.HasPrefix("sstable_", entry.Name()) {

			filePath := filepath.Join(filePath, entry.Name())
			f, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
			if err != nil {
				return nil, fmt.Errorf("failed to open %s file %v", filePath, err)
			}
			lsm.sstables = append(lsm.sstables, f)

		} else if strings.HasPrefix("index_", entry.Name()) {

			filePath := filepath.Join(filePath, entry.Name())
			f, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
			if err != nil {
				return nil, fmt.Errorf("failed to open %s file %v", filePath, err)
			}
			lsm.indice = append(lsm.indice, f)

		}
	}

	return lsm, nil
}

// Iterator
func (l *LSM) Iterator() Iterator {
	s := bufio.NewScanner(l.sstables[0])

	return &iterator{
		lsm:            l,
		currentSSTable: l.sstables[0],
		currentIndex:   0,
		scanner:        s,
	}
}

// Put
func (l *LSM) Put(key string, value []byte, indice map[string]string, ttl int64) error {
	entry := &pb.KeyValue{
		Key:   key,
		Value: value,
		Ttl:   ttl,
	}

	b64bytes, err := marshalKeyValue(entry)
	if err != nil {
		return fmt.Errorf("marshalKeyValue: %v", err)
	}

	if int64(l.memtable.size) >= l.flushToDisk {
		// compaction
		err = l.compact()
		if err != nil {
			return fmt.Errorf("lsm.compact: %v", err)
		}
		l.memtable.flush()
	}

	l.memtable.put(key, b64bytes)

	if ttl > 0 {
		// add to timeout channel
	}

	return nil
}

// Get
func (l *LSM) Get(key string) ([]byte, error) {
	vbytes, err := l.memtable.get(key)
	if err != nil && errors.Is(ErrNotFound, err) {
		// do nothing
	} else if err != nil {
		return nil, fmt.Errorf("memtable.get: %v", err)
	}

	if vbytes != nil {

		keyvalue, err := unmarshalKeyValue(vbytes)
		if err != nil {
			return nil, fmt.Errorf("unmarshalKeyValue: %v", err)
		}

		return keyvalue.Value, nil
	}

	for _, sstable := range l.sstables {

		s := bufio.NewScanner(sstable)

		for s.Scan() {

			keyvalue, err := unmarshalKeyValue(s.Bytes())
			if err != nil {
				return nil, fmt.Errorf("unmarshalKeyValue: %v", err)
			}

			if strings.Compare(key, keyvalue.Key) == 0 {
				return keyvalue.Value, nil
			}
		}
	}

	return nil, ErrNotFound
}

// Snapshot
func (l *LSM) Snapshot() ([]byte, error) {
	return nil, nil
}

func (l *LSM) ExpirationCh() (chan interface{}, error) {
	ch := make(chan interface{})
	go func() {
		defer close(ch)
	}()
	return ch, nil
}

// TODO:
// filter existing records only use the latest value
func (l *LSM) compact() error {
	fp := filepath.Join(l.sstDir, fmt.Sprintf("sstable_%d.data", len(l.sstables)))
	f, err := os.OpenFile(fp, os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModePerm)
	if err != nil {
		return fmt.Errorf("os.OpenFile: %v", err)
	}

	it := l.memtable.newMemTableIterator()

	for it.hasNext() {

		n := it.next()
		_, err = f.Write(n.payload)
		if err != nil {
			return fmt.Errorf("file.Write: %v", err)
		}
	}

	err = f.Sync()
	if err != nil {
		return fmt.Errorf("file.Sync: %v", err)
	}

	l.sstables = append(l.sstables, f)

	return nil
}

func marshalKeyValue(keyvalue *pb.KeyValue) ([]byte, error) {
	pbbytes, err := proto.Marshal(keyvalue)
	if err != nil {
		return nil, fmt.Errorf("proto.Marshal: %v", err)
	}

	// l := base64.StdEncoding.EncodedLen(len(pbbytes))
	// buf := make([]byte, l)

	// base64.StdEncoding.Encode(buf, pbbytes)

	return pbbytes, nil
}

func unmarshalKeyValue(in []byte) (*pb.KeyValue, error) {
	// l := base64.StdEncoding.DecodedLen(len(in))
	// pbbytes := make([]byte, l)
	// n, err := base64.StdEncoding.Decode(pbbytes, in)
	// if err != nil {
	// 	return nil, fmt.Errorf("base64.Decode: %v", err)
	// }

	var keyvalue pb.KeyValue
	err := proto.Unmarshal(in, &keyvalue)
	if err != nil {
		return nil, fmt.Errorf("proto.Unmarshal: %v", err)
	}

	return &keyvalue, nil
}
