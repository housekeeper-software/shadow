package db

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
)

var (
	errDBNotOpened = errors.New("database not opened")
)

//LevelDB use leveldb as database
//implement Database interface
type LevelDB struct {
	db *leveldb.DB
}

func NewLevelDB() *LevelDB {
	return &LevelDB{
		db: nil,
	}
}

func (d *LevelDB) Open(dir string) error {
	db, err := leveldb.OpenFile(dir, nil)
	if _, corrupted := err.(*errors.ErrCorrupted); corrupted {
		db, err = leveldb.RecoverFile(dir, nil)
	}
	if err != nil {
		return err
	}
	d.db = db
	return nil
}

func (d *LevelDB) Close() error {
	if d.db == nil {
		return nil
	}
	return d.db.Close()
}

func (d *LevelDB) Read(key string) ([]byte, error) {
	if d.db == nil {
		return nil, errDBNotOpened
	}
	return d.db.Get([]byte(key), nil)
}

func (d *LevelDB) Write(key string, value []byte) error {
	if d.db == nil {
		return errDBNotOpened
	}
	return d.db.Put([]byte(key), value, nil)
}

func (d *LevelDB) Delete(key string) error {
	if d.db == nil {
		return errDBNotOpened
	}
	return d.db.Delete([]byte(key), nil)
}

func (d *LevelDB) ReadAll() (map[string][]byte, error) {
	if d.db == nil {
		return nil, errDBNotOpened
	}
	entries := make(map[string][]byte)
	iter := d.db.NewIterator(nil, nil)
	defer iter.Release()
	for iter.Next() {
		entries[string(iter.Key())] = iter.Value()
	}
	return entries, nil
}
