package db

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	errNullEntry = errors.New("entry not found")
)

type Entry struct {
	entries map[string]*Document
	//entry cache
	entryCache map[string][]byte
	index      *Index
	//index cache
	indexCache []byte
	mutex      sync.RWMutex
	db         *LevelDB
}

func Md5Bytes(b []byte) string {
	sum := md5.Sum(b)
	return hex.EncodeToString(sum[:])
}

func NewEntry() *Entry {
	return &Entry{entries: make(map[string]*Document, 500),
		entryCache: make(map[string][]byte, 500),
		index:      NewIndex(),
		indexCache: make([]byte, 0, 1024*500),
		db:         nil}
}

func (e *Entry) Close() error {
	if e.db != nil {
		return e.db.Close()
	}
	return nil
}

//LoadEntry call once only on application startup
func (e *Entry) LoadEntry(file string) error {
	db := NewLevelDB()
	err := db.Open(file)
	if err != nil {
		return err
	}
	c, err := db.ReadAll()
	if err != nil {
		return err
	}
	for k, v := range c {
		doc := NewDocument()
		err = doc.Load(v)
		if err != nil {
			logrus.Errorf("Failed to load %s", k)
			continue
		}
		e.index.Add(k, Md5Bytes(v))
		e.entries[k] = doc
		//cache entry
		e.entryCache[k], _ = doc.ToJson()
	}
	e.db = db
	//cache index
	if !e.index.IsEmpty() {
		e.indexCache, _ = e.index.ToJson()
	}
	return nil
}

func (e *Entry) RemoveEntry(entry string) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	return e.deleteEntry(entry)
}

func (e *Entry) AddItem(entry string, name string, value []byte) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	c, ok := e.entries[entry]
	if !ok {
		doc := NewDocument()
		err := doc.AddItem(name, value)
		if err != nil {
			return err
		}
		e.entries[entry] = doc
	} else {
		err := c.AddItem(name, value)
		if err != nil {
			return err
		}
	}
	return e.commit(entry)
}

func (e *Entry) RemoveItem(entry string, name string, id string) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	c, ok := e.entries[entry]
	if !ok {
		return nil
	}
	//remove from document
	c.RemoveItem(name, id)
	//document is empty, delete it!
	if c.IsEmpty() {
		return e.deleteEntry(entry)
	}
	//commit changed
	return e.commit(entry)
}

func (e *Entry) ToJson(entry string) ([]byte, error) {
	e.mutex.RLock()
	defer e.mutex.RUnlock()

	//try load from cache
	cache, ok := e.entryCache[entry]
	if ok {
		return cache, nil
	}

	c, ok := e.entries[entry]
	if !ok {
		return nil, errNullEntry
	}
	//serialize to json
	j, err := c.ToJson()
	if err != nil {
		return nil, err
	}
	e.entryCache[entry] = j
	return j, nil
}

func (e *Entry) IndexToJson() ([]byte, error) {
	e.mutex.RLock()
	defer e.mutex.RUnlock()

	//load from cache
	if len(e.indexCache) > 0 {
		return e.indexCache, nil
	}
	j, err := e.index.ToJson()
	if err != nil {
		return nil, err
	}
	e.indexCache = j
	return j, nil
}

func (e *Entry) commit(entry string) error {
	c, ok := e.entries[entry]
	if !ok {
		return nil
	}
	data, err := c.ToJson()
	if err != nil {
		return err
	}
	err = e.db.Write(entry, data)
	if err != nil {
		return err
	}
	e.index.Add(entry, Md5Bytes(data))
	e.entryCache[entry] = data
	e.indexCache, _ = e.index.ToJson()
	return nil
}

func (e *Entry) deleteEntry(entry string) error {
	//delete database entry first
	err := e.db.Delete(entry)
	if err != nil {
		return err
	}
	delete(e.entries, entry)
	delete(e.entryCache, entry)
	e.index.Remove(entry)
	e.indexCache = e.indexCache[:0]
	return nil
}
