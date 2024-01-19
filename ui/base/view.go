package base

import (
	"fyne.io/fyne/v2"
	"github.com/google/uuid"
)

// IView represents the base interview for views
type View interface {
	// ID Uniquely identifiers the view within a context.
	ID() string

	// SetID assigns an id to the view
	SetID(string)

	// Equals checks if 2 views are identical by comparing their ID
	Equals(View) bool

	// Title returns the title of the view
	// The title can be used as a dialog title, tab title, etc.
	Title() string

	// Initialize initializes the view's state and visual representation
	Initialize() (fyne.CanvasObject, error)

	// Refresh is used to signal to the view that it should refres it's
	// internal state and potentially it's visuals.
	Refresh() error

	// OnShow is called when the view becomes visible on the screen
	OnShow()

	// OnHide is called when the view becomes invisible on the screen
	OnHide()
}

type AbstractView struct {
	id    string
	title string
}

func NewAbstractView(
	title string,
) *AbstractView {
	return &AbstractView{
		id:    uuid.NewString(),
		title: title,
	}
}

func (v AbstractView) ID() string {
	return v.id
}

func (v *AbstractView) SetID(id string) {
	v.id = id
}

func (v AbstractTabbedView) Equals(view View) bool {
	return v.id == view.ID()
}

func (v AbstractView) Title() string {
	return v.title
}
