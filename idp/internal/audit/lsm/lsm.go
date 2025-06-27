package lsm

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/protobuf/proto"

	pb "github.com/structx/tbd/lib/protocol/lsm/v1"
)

const (
	// neverExpire = -1
	defaultCurrentIndex = 0
	zeroLength          = 0

	openFileErrMsg = "os.OpenFile: %w"

	flushToDiskSize = 10000
)

// Store key value store
type Store interface {
	Get(string) ([]byte, error)
	Put(string, []byte, map[string]string, int64) error
	Iterator() Iterator
}

// Iterator key value store iterator
type Iterator interface{}

// LSM log structured merge tree
type LSM struct {
	indice   []*os.File
	sstDir   string
	sstables []*os.File
	memtable *memtable
	// wal         *WAL
	flushToDisk int64
}

// interface compliance
var _ Store = (*LSM)(nil)

// New return new lsm implementation
func New(dir string) (*LSM, error) {
	filePath := filepath.Clean(dir)
	lsm := &LSM{
		indice:      make([]*os.File, zeroLength),
		sstables:    make([]*os.File, zeroLength),
		memtable:    newMemTable(),
		sstDir:      filePath,
		flushToDisk: flushToDiskSize,
	}

	entries, err := os.ReadDir(filePath)
	if err != nil {
		return nil, fmt.Errorf("os.ReadDir: %w", err)
	}

	if len(entries) == zeroLength {
		// create required files

		sstFilePath := filepath.Join(filePath, "sstable_0.data")
		f, err := os.OpenFile(sstFilePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
		if err != nil {
			return nil, fmt.Errorf(openFileErrMsg, err)
		}
		lsm.sstables = append(lsm.sstables, f)

		idxFilePath := filepath.Join(filePath, "index_0.log")
		f, err = os.OpenFile(idxFilePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
		if err != nil {
			return nil, fmt.Errorf(openFileErrMsg, err)
		}
		lsm.indice = append(lsm.indice, f)

		walFilePath := filepath.Join(filePath, "wal.log")
		_, err = os.OpenFile(walFilePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
		if err != nil {
			return nil, fmt.Errorf(openFileErrMsg, err)
		}
		// lsm.wal = newWAL(f)

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
				return nil, fmt.Errorf(openFileErrMsg, err)
			}
			lsm.sstables = append(lsm.sstables, f)
		} else if strings.HasPrefix("index_", entry.Name()) {
			filePath := filepath.Join(filePath, entry.Name())
			f, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
			if err != nil {
				return nil, fmt.Errorf(openFileErrMsg, err)
			}
			lsm.indice = append(lsm.indice, f)
		}
	}

	return lsm, nil
}

// Iterator returns new iterator
func (l *LSM) Iterator() Iterator {
	s := bufio.NewScanner(l.sstables[defaultCurrentIndex])

	return &iterator{
		lsm:            l,
		currentSSTable: l.sstables[defaultCurrentIndex],
		currentIndex:   defaultCurrentIndex,
		scanner:        s,
	}
}

// Put set record
func (l *LSM) Put(key string, value []byte, indice map[string]string, ttl int64) error {
	entry := &pb.KeyValue{
		Key:   key,
		Value: value,
		Ttl:   ttl,
	}

	b64bytes, err := marshalKeyValue(entry)
	if err != nil {
		return fmt.Errorf("marshalKeyValue: %w", err)
	}

	if int64(l.memtable.size) >= l.flushToDisk {
		// compaction
		err = l.compact()
		if err != nil {
			return fmt.Errorf("lsm.compact: %w", err)
		}
		l.memtable.flush()
	}

	l.memtable.put(key, b64bytes)

	// if ttl > neverExpire {
	// TODO
	// add to timeout channel
	// }

	return nil
}

// Get record by key
func (l *LSM) Get(key string) ([]byte, error) {
	vbytes, err := l.memtable.get(key)
	if err != nil && !errors.Is(ErrNotFound, err) {
		return nil, fmt.Errorf("memtable.get: %w", err)
	}

	if vbytes != nil {
		keyvalue, err := unmarshalKeyValue(vbytes)
		if err != nil {
			return nil, fmt.Errorf("unmarshalKeyValue: %w", err)
		}

		return keyvalue.Value, nil
	}

	for _, sstable := range l.sstables {
		s := bufio.NewScanner(sstable)

		for s.Scan() {
			keyvalue, err := unmarshalKeyValue(s.Bytes())
			if err != nil {
				return nil, fmt.Errorf("unmarshalKeyValue: %w", err)
			}

			if strings.Compare(key, keyvalue.Key) == equal {
				return keyvalue.Value, nil
			}
		}
	}

	return nil, ErrNotFound
}

// Snapshot getter snapshot of records
func (l *LSM) Snapshot() ([]byte, error) {
	return nil, nil
}

// ExpirationCh getter expiration channel
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
		return fmt.Errorf("os.OpenFile: %w", err)
	}

	it := l.memtable.newMemTableIterator()

	for it.hasNext() {
		n := it.next()
		_, err = f.Write(n.payload)
		if err != nil {
			return fmt.Errorf("file.Write: %w", err)
		}
	}

	err = f.Sync()
	if err != nil {
		return fmt.Errorf("file.Sync: %w", err)
	}

	l.sstables = append(l.sstables, f)

	return nil
}

func marshalKeyValue(keyvalue *pb.KeyValue) ([]byte, error) {
	pbbytes, err := proto.Marshal(keyvalue)
	if err != nil {
		return nil, fmt.Errorf("proto.Marshal: %w", err)
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
		return nil, fmt.Errorf("proto.Unmarshal: %w", err)
	}

	return &keyvalue, nil
}
