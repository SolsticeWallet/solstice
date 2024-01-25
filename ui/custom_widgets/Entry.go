package custom_widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type Entry struct {
	widget.Entry
	readonly bool
}

func NewEntry() *Entry {
	e := &Entry{}
	e.ExtendBaseWidget(e)
	return e
}

func (e *Entry) ReadOnly() bool {
	return e.readonly
}

func (e *Entry) SetReadonly(readonly bool) {
	e.readonly = readonly
}

func (e *Entry) KeyDown(key *fyne.KeyEvent) {
	if e.readonly {
		return
	}
	e.Entry.KeyDown(key)
}

func (e *Entry) KeyUp(key *fyne.KeyEvent) {
	if e.readonly {
		return
	}
	e.Entry.KeyUp(key)
}

func (e *Entry) TypedKey(key *fyne.KeyEvent) {
	if e.readonly {
		return
	}
	e.Entry.TypedKey(key)
}

func (e *Entry) TypedRune(r rune) {
	if e.readonly {
		return
	}
	e.Entry.TypedRune(r)
}
