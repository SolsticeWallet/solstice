package base

import (
	"fyne.io/fyne/v2"
	"github.com/google/uuid"
)

// View represents the base interview for views
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

	// OnShow is called when the view becomes visible on the screen
	OnShow()

	// OnHide is called when the view becomes invisible on the screen
	OnHide()
}

type AbstractView struct {
	id    string
	title string
}

// NewAbstractView creates a new AbstractView with the given title.
//
// It takes a title string as a parameter and returns a pointer to an AbstractView.
func NewAbstractView(
	title string,
) *AbstractView {
	return &AbstractView{
		id:    uuid.NewString(),
		title: title,
	}
}

// ID returns the ID of the AbstractView.
//
// Returns a string.
func (v AbstractView) ID() string {
	return v.id
}

// SetID sets the ID of the AbstractView.
//
// id string
func (v *AbstractView) SetID(id string) {
	v.id = id
}

// Equals checks if the AbstractTabbedView is equal to the given View.
// It takes a View as a parameter and returns a boolean.
func (v AbstractView) Equals(view View) bool {
	return v.id == view.ID()
}

// Title returns the title of the AbstractView.
//
// Returns a string.
func (v AbstractView) Title() string {
	return v.title
}
