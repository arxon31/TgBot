package telegram

import (
	"TgBot/lib/e"
	"TgBot/storage"
	"context"
	"errors"
	"log"
	"net/url"
	"strings"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(text string, chatID int, userName string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s'", text, userName)

	if isAddCmd(text) {
		return p.savePage(chatID, text, userName)
	}
	// add page: http://...
	// rnd page: /rnd
	// help: /help
	// start: /start: hi + help

	switch text {
	case RndCmd:
		return p.sendRandom(chatID, userName)
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendHello(chatID)
	default:
		return p.tg.SendMessage(chatID, msgUnknown)

	}
}

func (p *Processor) savePage(chatID int, pageURL string, userName string) (err error) {
	defer func() { err = e.Wrap("can't do command save page", err) }()

	page := &storage.Page{
		URL:      pageURL,
		Username: userName,
	}

	isExists, err := p.storage.IsExists(context.Background(), page)
	if err != nil {
		return err
	}
	if isExists {
		return p.tg.SendMessage(chatID, msgAlreadyExist)
	}

	if err := p.storage.Save(context.Background(), page); err != nil {
		return err
	}
	if err := p.tg.SendMessage(chatID, msgSaved); err != nil {
		return err
	}

	return nil
}

func (p *Processor) sendRandom(chatID int, userName string) (err error) {
	defer func() { err = e.Wrap("can't do command send random", err) }()

	page, err := p.storage.PickRandom(context.Background(), userName)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}

	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.tg.SendMessage(chatID, msgNoSavedPages)
	}

	if err := p.tg.SendMessage(chatID, page.URL); err != nil {
		return err
	}

	return p.storage.Remove(context.Background(), page)
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}
func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
