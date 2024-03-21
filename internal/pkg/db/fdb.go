package db

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

//FileDB use file as database,test only
//implement Database interface
type FileDB struct {
	dir string
}

func NewFileDB() *FileDB {
	return &FileDB{}
}

func (d *FileDB) Open(dir string) error {
	d.dir = dir
	return nil
}

func (d *FileDB) Close() error {
	return nil
}

func makeFileName(dir string, key string) string {
	return path.Join(dir, fmt.Sprintf("%s.json", key))
}

func (d *FileDB) Read(key string) ([]byte, error) {
	file := makeFileName(d.dir, key)
	return ioutil.ReadFile(file)
}

func (d *FileDB) Write(key string, value []byte) error {
	file := makeFileName(d.dir, key)
	//ensure directory exist
	err := os.MkdirAll(filepath.Dir(d.dir), os.ModePerm)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, value, 0644)
}

func (d *FileDB) Delete(key string) error {
	file := makeFileName(d.dir, key)
	return os.Remove(file)
}

func (d *FileDB) ReadAll() (map[string][]byte, error) {
	dir, err := ioutil.ReadDir(d.dir)
	if err != nil {
		return nil, err
	}
	entries := make(map[string][]byte)
	for _, fi := range dir {
		if fi.IsDir() {
			continue
		}
		file := path.Join(d.dir, fi.Name())
		b, err := ioutil.ReadFile(file)
		if err != nil {
			continue
		}
		name := strings.TrimSuffix(fi.Name(), path.Ext(fi.Name()))
		entries[name] = b
	}
	return entries, nil
}
