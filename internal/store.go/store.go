package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/tryoasnafi/gate/internal/crypto"
)

type Entry struct {
	Label     string    `json:"label"`
	User      string    `json:"user"`
	Host      string    `json:"host"`
	Port      int       `json:"port"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
}

type storage struct {
	Entries map[string]Entry `json:"entries"`
}

var (
	fileName = filepath.Join(os.Getenv("HOME"), ".gate", "store.enc")
	mutex    sync.Mutex
)

func InitStore(password []byte) error {
	mutex.Lock()
	defer mutex.Unlock()

	if _, err := os.Stat(fileName); err == nil {
		return errors.New("store already exists")
	}
	os.MkdirAll(filepath.Dir(fileName), 0700)
	store := storage{Entries: make(map[string]Entry)}
	return save(password, store)
}

func AddEntry(password []byte, entry Entry) error {
	mutex.Lock()
	defer mutex.Unlock()
	store, err := load(password)
	if err != nil {
		return err
	}
	if _, exists := store.Entries[entry.Label]; exists {
		return errors.New("label already exists")
	}
	store.Entries[entry.Label] = entry
	return save(password, store)
}

func GetEntry(password []byte, label string) (Entry, error) {
	store, err := load(password)
	if err != nil {
		return Entry{}, err
	}
	entry, exists := store.Entries[label]
	if !exists {
		return Entry{}, fmt.Errorf("label not found: %s", label)
	}
	return entry, nil
}

func ListEntries(password []byte) ([]Entry, error) {
	store, err := load(password)
	if err != nil {
		return nil, err
	}
	entries := []Entry{}
	for _, e := range store.Entries {
		entries = append(entries, e)
	}
	return entries, nil
}

func DeleteEntry(password []byte, label string) error {
	mutex.Lock()
	defer mutex.Unlock()
	store, err := load(password)
	if err != nil {
		return err
	}
	delete(store.Entries, label)
	return save(password, store)
}

func RotateMasterPassword(oldPwd, newPwd []byte) error {
	mutex.Lock()
	defer mutex.Unlock()
	store, err := load(oldPwd)
	if err != nil {
		return err
	}
	return save(newPwd, store)
}

func ImportEntries(masterPwd, importPwd []byte, filePath string) error {
	mutex.Lock()
	defer mutex.Unlock()

	existing, err := load(masterPwd)
	if err != nil {
		return err
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	plain, err := crypto.Decrypt(importPwd, data)
	if err != nil {
		return err
	}
	var imported storage
	if err := json.Unmarshal(plain, &imported); err != nil {
		return err
	}
	for k, v := range imported.Entries {
		existing.Entries[k] = v
	}
	return save(masterPwd, existing)
}

func TryDecryptOnly(password []byte) error {
	_, err := load(password)
	return err
}

func load(password []byte) (storage, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return storage{}, err
	}
	plain, err := crypto.Decrypt(password, data)
	if err != nil {
		return storage{}, err
	}
	var store storage
	if err := json.Unmarshal(plain, &store); err != nil {
		return storage{}, err
	}
	return store, nil
}

func save(password []byte, s storage) error {
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}
	encrypted, err := crypto.Encrypt(password, data)
	if err != nil {
		return err
	}
	return os.WriteFile(fileName, encrypted, 0600)
}
