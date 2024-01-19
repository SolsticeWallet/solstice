package base

import "fyne.io/fyne/v2"

type Shortcut struct {
	key      fyne.KeyName
	modifier fyne.KeyModifier
	name     string
}

func NewShortcut(
	name string,
	key fyne.KeyName,
	modifier ...fyne.KeyModifier,
) *Shortcut {
	sc := &Shortcut{
		key:  key,
		name: name,
	}
	for _, mod := range modifier {
		sc.modifier |= mod
	}
	return sc
}

func (s Shortcut) Key() fyne.KeyName {
	return s.key
}

func (s Shortcut) Mod() fyne.KeyModifier {
	return s.modifier
}

func (s Shortcut) ShortcutName() string {
	return s.name
}
