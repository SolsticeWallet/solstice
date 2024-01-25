package base

import "fyne.io/fyne/v2"

type Shortcut struct {
	key      fyne.KeyName
	modifier fyne.KeyModifier
	name     string
}

// NewShortcut creates a new Shortcut with the given name, key, and optional
// modifiers.
//
// Parameters:
//
//	name string - the name of the shortcut
//	key fyne.KeyName - the key associated with the shortcut
//	modifier ...fyne.KeyModifier - optional modifiers for the shortcut
//
// Return:
//
//	*Shortcut - the newly created Shortcut
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

// Key returns the key name for the Shortcut.
//
// Returns fyne.KeyName.
func (s Shortcut) Key() fyne.KeyName {
	return s.key
}

// Mod returns the key modifier of the Shortcut.
//
// Returns fyne.KeyModifier.
func (s Shortcut) Mod() fyne.KeyModifier {
	return s.modifier
}

// ShortcutName description of the Go function.
//
// Returns a string.
func (s Shortcut) ShortcutName() string {

	return s.name
}
