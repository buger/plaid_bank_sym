package main

import (
	"encoding/json"
	"io/ioutil"
	_ "log"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

type FileDB struct {
	sync.Mutex
	root string
}

func OpenDB(root string) *FileDB {
	root, _ = filepath.Abs(root)

	if root[len(root)-1] != '/' {
		root += "/"
	}

	return &FileDB{root: root}
}

func (db *FileDB) Write(path string, data []byte) error {
	db.Lock()
	defer db.Unlock()

	dir, _ := filepath.Split(db.root + path)
	os.MkdirAll(dir, 0644)

	return ioutil.WriteFile(db.root+path, data, 0644)
}

func (db *FileDB) WriteJSON(path string, v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	db.Write(path, b)
	return nil
}

func (db *FileDB) WriteString(path string, data string) error {
	return db.Write(path, []byte(data))
}

func (db *FileDB) Read(path string) ([]byte, error) {
	db.Lock()
	defer db.Unlock()

	return ioutil.ReadFile(db.root + path)
}

func (db *FileDB) ReadJSON(path string, v interface{}) error {
	data, err := db.Read(path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, v)
	if err != nil {
		return err
	}

	return nil
}

func (db *FileDB) ReadString(path string) (string, error) {
	data, err := db.Read(path)
	return string(data), err
}

func (db *FileDB) Remove(path string) error {
	db.Lock()
	defer db.Unlock()

	return os.RemoveAll(db.root + path)
}

func (db *FileDB) RemoveAll() error {
	return db.Remove("")
}

type ByLastModified []os.FileInfo

func (a ByLastModified) Len() int           { return len(a) }
func (a ByLastModified) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByLastModified) Less(i, j int) bool { return a[j].ModTime().Before(a[i].ModTime()) }

func (db *FileDB) AllKeys(path string) ([]string, error) {
	var files []os.FileInfo
	content, err := ioutil.ReadDir(db.root + path)

	if err != nil {
		return nil, err
	}

	for _, c := range content {
		if !c.IsDir() {
			files = append(files, c)
		}
	}

	sort.Sort(ByLastModified(files))

	var keys []string
	for _, f := range files {
		keys = append(keys, f.Name())
	}

	return keys, err
}

func (db *FileDB) LastModified(path string) (time.Time, error) {
	info, err := os.Stat(db.root + path)

	if err != nil {
		return time.Time{}, err
	}

	return info.ModTime(), nil
}

func (db *FileDB) IfModifiedSince(path string, cacheTime time.Duration) bool {
	mod, err := db.LastModified(path)

	if err != nil {
		return false
	}

	if !mod.IsZero() && time.Since(mod) < cacheTime {
		return true
	}

	return false
}

func (db *FileDB) Exists(path string) bool {
	db.Lock()
	defer db.Unlock()

	_, err := os.Stat(db.root + path)

	return err == nil
}

func (db *FileDB) CreateLink(path, linkto string) (err error) {
	db.Lock()
	defer db.Unlock()

	dir, _ := filepath.Split(db.root + linkto)
	os.MkdirAll(dir, 0644)

	err = os.Symlink(db.root+path, db.root+linkto)

	return
}

var tMux sync.Mutex

func (db *FileDB) Transaction(cb func(*FileDB)) {
	tMux.Lock()
	defer tMux.Unlock()

	cb(db)
}
