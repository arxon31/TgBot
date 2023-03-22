package storage

import (
	"TgBot/lib/e"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
)

type Storage interface {
	Save(p *Page) error
	PickRandom(Username string) (*Page, error)
	Remove(p *Page) error
	IsExists(p *Page) (bool, error)
}

type Page struct {
	URL      string
	Username string
}

var ErrNoSavedPages = errors.New("no saved pages in this directory")

func (p *Page) Hash() (string, error) {
	h := sha1.New()
	if _, err := io.WriteString(h, p.Username); err != nil {
		return "", e.Wrap("can't calculate hash", err)
	}
	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", e.Wrap("can't calculate hash", err)
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
