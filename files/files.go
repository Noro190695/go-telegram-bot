package files

import (
	"bot/lib/myError"
	"bot/storage"
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

type Storage struct {
	basePath string
}

const (
	defaultPerm = 0774
)

var ErrNoSavedPages = errors.New("no saved pages")

// func New(basePath string) Storage {

// }

func (s Storage) Save(page *storage.Page) (err error) {
	defer func() {
		err = myError.WrapIfErr("can't save page", err)
	}()

	filePath := filepath.Join(s.basePath, page.UserName)

	if err := os.MkdirAll(filePath, defaultPerm); err != nil {
		return err
	}
	fName, err := fileName(page)
	if err != nil {
		return err
	}

	fPath := filepath.Join(filePath, fName)

	file, err := os.Create(fPath)
	defer func() { _ = file.Close() }()

	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return err
	}
	return nil
}

func (s Storage) PickRandom(userName string) (page *storage.Page, err error) {
	defer func() {
		err = myError.WrapIfErr("can't pick random", err)
	}()
	path := filepath.Join(s.basePath, userName)

	files, err := os.ReadDir(path)

	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, ErrNoSavedPages
	}

	rand.Seed(time.Now().UnixMicro())
	n := rand.Intn(len(files))

	file := files[n]

	return s.decodePage(filepath.Join(path, file.Name()))

}

func (s Storage) Remove(p *storage.Page) error {
	fName, err := fileName(p)
	if err != nil {
		return myError.Wrap("can't remove file", err)
	}
	path := filepath.Join(s.basePath, p.UserName, fName)
	if err := os.Remove(path); err != nil {
		errMsg := fmt.Sprintf("can't remove file: %s", path)
		return myError.Wrap(errMsg, err)
	}
	return nil
}

func (s Storage) IsExists(p *storage.Page) (bool, error) {
	fName, err := fileName(p)
	if err != nil {
		return false, myError.Wrap("file not exist", err)
	}
	path := filepath.Join(s.basePath, p.UserName, fName)
	switch _, err = os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		errMsg := fmt.Sprintf("can't check if file:: %s", path)
		return false, myError.Wrap(errMsg, err)
	}
	return true, nil
}

func (s Storage) decodePage(filePath string) (*storage.Page, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, myError.Wrap("can't decode path", err)
	}

	defer func() { _ = f.Close() }()

	var p storage.Page

	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, myError.Wrap("can't decode page", err)
	}
	return &p, nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
