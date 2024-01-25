package custom_widgets

import (
	"unicode"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/solsticewallet/solstice-core/blockchains/ethereum/utils"
)

type Bip39Entry struct {
	widget.Entry
}

func NewBip39Entry() *Bip39Entry {
	e := &Bip39Entry{}
	e.ExtendBaseWidget(e)
	return e
}

func (e *Bip39Entry) KeyDown(key *fyne.KeyEvent) {
	// Override current behaviour so we do not allow selection
	// via shift keys
}

func (e *Bip39Entry) KeyUp(key *fyne.KeyEvent) {
	// Override current behaviour so we do not allow selection
	// via shift keys
}

func (e *Bip39Entry) TypedKey(key *fyne.KeyEvent) {
	switch key.Name {
	case fyne.KeyUp, fyne.KeyDown,
		fyne.KeyPageUp, fyne.KeyPageDown,
		fyne.KeyDelete:
		return
	case fyne.KeyBackspace:
		e.Entry.TypedKey(key)
		e.TypedRune(0)
	default:
		e.Entry.TypedKey(key)
	}
}

func (e *Bip39Entry) TypedRune(r rune) {
	var suggestions []string
	txt := e.Text[:e.CursorColumn]
	if unicode.IsLower(r) {
		txt += string(r)
	}

	suggestions = utils.Bip39SuggestWords(txt)
	if len(suggestions) == 0 {
		return
	}

	if unicode.IsLower(r) {
		e.Entry.TypedRune(r)
	}

	//e.partialWord = txt
	e.Text = suggestions[0]
	e.Refresh()
}

func (e *Bip39Entry) FocusGained() {
	e.Entry.FocusGained()
	e.CursorColumn = len(e.Text)
	e.Refresh()
}

func (e *Bip39Entry) FocusLost() {
	e.Validate()
	e.Entry.FocusLost()
}
