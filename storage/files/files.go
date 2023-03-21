package files

import (
	"TgBot/lib/e"
	"TgBot/storage"
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

const defaultPermission = 0774

type Storage struct {
	basePath string
}

func New(basePath string) Storage {
	return Storage{
		basePath: basePath,
	}
}

func (s Storage) Save(page *storage.Page) (err error) {
	defer func() { err = e.Wrap("can't save page", err) }()

	filePath := filepath.Join(s.basePath, page.Username)

	if err := os.MkdirAll(filePath, defaultPermission); err != nil {
		return err
	}

	fName, err := fileName(page)
	if err != nil {
		return err
	}

	filePath = filepath.Join(filePath, fName)

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer func() { _ = file.Close() }()

	gob.NewEncoder(file).Encode(page)
	return nil
}
func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
func (s Storage) PickRandom(userName string) (page *storage.Page, err error) {
	defer func() { err = e.Wrap("can't pick random page", err) }()

	filePath := filepath.Join(s.basePath, userName)

	files, err := os.ReadDir(filePath)
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, errors.New("no saved pages in this directory")
	}
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(files))

	file := files[n]

	return s.decodePage(filepath.Join(filePath, file.Name()))

}

func (s Storage) Remove(p *storage.Page) error {
	fileName, err := fileName(p)
	if err != nil {
		return e.Wrap("can't remove file", err)
	}
	filePath := filepath.Join(s.basePath, p.Username, fileName)

	if err := os.Remove(filePath); err != nil {
		msg := fmt.Sprintf("can't remove file from dir %s", filePath)
		return e.Wrap(msg, err)
	}
	return nil
}

func (s Storage) IsExists(p *storage.Page) (bool, error) {
	fileName, err := fileName(p)
	if err != nil {
		return false, e.Wrap("can't check if file is exist", err)
	}
	filePath := filepath.Join(s.basePath, p.Username, fileName)

	switch _, err = os.Stat(filePath); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		msg := fmt.Sprintf("can't check if file is exist in dir %s", filePath)
		return false, e.Wrap(msg, err)
	}
	return true, nil
}

func (s Storage) decodePage(filePath string) (*storage.Page, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, e.Wrap("can't decode file", err)
	}
	defer func() { _ = f.Close() }()

	var p storage.Page
	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, e.Wrap("can't decode file", err)
	}
	return &p, nil
}
