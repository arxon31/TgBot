package telegram

const msgHelp = `I can save and keep your pages. Also I can offer you this pages to read
	
	In order to save the page, just send me a link to it.

	In order to get a random page from your list, send me command /rnd.
	
	Caution! After /rnd command returned page will be removed from your list`

const msgHello = "Hi there! \n\n" + msgHelp

const (
	msgUnknown      = "Unknown command"
	msgNoSavedPages = "You have no saved pages"
	msgSaved        = "Saved!"
	msgAlreadyExist = "You have already have this page in your list"
)
