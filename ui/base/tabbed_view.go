package base

import (
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/solsticewallet/solstice/i18n"
)

type TabbedView interface {
	View

	// Append adds the provided views as children to the TabbedView
	Append(...View)

	// RemoveActive removes the current active tab from the TabbedView
	RemoveActive()

	// NumTabs returns the number of tabs present in the TabbedView
	NumTabs() int

	// TabIndex returns the index of the requested tab in the TabbedView
	TabIndex(View) int

	// SelectTabByIndex selects the tab in the TabbedView
	SelectTabByIndex(int)

	// SelectTab selects the tab in the TabbedView
	SelectTab(View)

	// OnTabSelected is called when a tab is selected
	OnTabSelected(View)

	// OnTabUnselected is called when a tab is unselected
	OnTabUnselected(View)
}

type AbstractTabbedView struct {
	*AbstractView

	tabviews       []View
	tabviewIndex   map[string]View
	currentTabView View
	tabs           *container.AppTabs
	tabLocation    container.TabLocation
	mu             *sync.Mutex
}

// NewAbstractTabbedView creates a new AbstractTabbedView.
//
// It takes a title string, tabLocation container.TabLocation, and zero or more
// tabs of type View. It returns a pointer to AbstractTabbedView.
func NewAbstractTabbedView(
	title string,
	tabLocation container.TabLocation,
	tabs ...View,
) *AbstractTabbedView {
	return &AbstractTabbedView{
		AbstractView: NewAbstractView(title),
		tabviews:     tabs,
		tabviewIndex: make(map[string]View),
		tabLocation:  tabLocation,
		mu:           &sync.Mutex{},
	}
}

// Initialize initializes the AbstractTabbedView.
//
// It returns a CanvasObject and an error.
func (v *AbstractTabbedView) Initialize() (fyne.CanvasObject, error) {
	if err := v.configureTabs(); err != nil {
		Logger.Error(
			i18n.T("Err.Initialization"),
			i18n.T("Err.Arg.Error"), err)
		return nil, err
	}
	return v.tabs, nil
}

// Append appends the given views to the AbstractTabbedView.
//
// It takes a variadic parameter of views and does not return anything.
func (v *AbstractTabbedView) Append(views ...View) {
	v.mu.Lock()
	defer v.mu.Unlock()

	for _, view := range views {
		if _, exists := v.tabviewIndex[view.ID()]; !exists {
			if canvas, err := view.Initialize(); err == nil {
				v.tabviews = append(v.tabviews, view)
				v.tabviewIndex[view.ID()] = view
				v.tabs.Append(container.NewTabItem(view.Title(), canvas))
			}
		}
	}
}

// RemoveActive removes the active tab from the AbstractTabbedView.
func (v *AbstractTabbedView) RemoveActive() {
	v.mu.Lock()
	defer v.mu.Unlock()

	if v.currentTabView == nil {
		return
	}

	curIdx := v.tabs.SelectedIndex()
	delete(v.tabviewIndex, v.currentTabView.ID())
	v.currentTabView.OnHide()
	v.currentTabView = nil
	v.tabviews = append(v.tabviews[0:curIdx], v.tabviews[curIdx+1:]...)
	v.tabs.RemoveIndex(curIdx)

	if curIdx > 0 {
		curIdx--
	}
	if curIdx < len(v.tabviews) {
		v.tabs.SelectIndex(curIdx)
	}
}

// NumTabs returns the number of tabs in the AbstractTabbedView.
//
// Returns an integer.
func (v AbstractTabbedView) NumTabs() int {
	return len(v.tabviews)
}

// TabIndex returns the index of the specified view in the AbstractTabbedView.
//
// It takes a View as a parameter and returns an integer.
func (v AbstractTabbedView) TabIndex(view View) int {
	for idx, tabview := range v.tabviews {
		if tabview.Equals(view) {
			return idx
		}
	}
	return -1
}

// SelectTabByIndex selects a tab by its index.
//
// idx int
func (v *AbstractTabbedView) SelectTabByIndex(idx int) {
	v.tabs.SelectIndex(idx)
}

// SelectTab selects the tab for the given view.
//
// view View
func (v *AbstractTabbedView) SelectTab(view View) {
	idx := v.TabIndex(view)
	if idx < 0 {
		return
	}
	v.tabs.SelectIndex(idx)
}

// configureTabs configures the tabs for the AbstractTabbedView.
// It initializes the tabs, sets the tab location, and appends the tabs
// to the tab view. It sets the current tab view to the first one.
// It returns an error if there is an error initializing the tab views.
func (v *AbstractTabbedView) configureTabs() error {
	v.mu.Lock()
	defer v.mu.Unlock()

	v.tabs = container.NewAppTabs()
	v.tabs.SetTabLocation(v.tabLocation)
	v.tabs.OnSelected = v.onTabSelected
	v.tabs.OnUnselected = v.onTabUnselected

	for _, tabView := range v.tabviews {
		canvasObj, err := tabView.Initialize()
		if err != nil {
			return err
		}
		v.tabs.Append(container.NewTabItem(tabView.Title(), canvasObj))
		v.tabviewIndex[tabView.ID()] = tabView
	}
	if len(v.tabviews) > 0 {
		v.currentTabView = v.tabviews[0]
	}
	return nil
}

// onTabSelected is a function that handles the tab selection in the
// AbstractTabbedView.
//
// It takes a tab *container.TabItem as a parameter.
func (v *AbstractTabbedView) onTabSelected(tab *container.TabItem) {
	v.mu.Lock()
	defer v.mu.Unlock()

	idx := v.tabs.SelectedIndex()
	if idx > -1 && idx < len(v.tabviews) {
		v.currentTabView = v.tabviews[idx]
		v.currentTabView.OnShow()
	}
}

// onTabUnselected is a function that handles the unselect event of a tab in the
// AbstractTabbedView.
//
// It takes a tab *container.TabItem as a parameter.
func (v *AbstractTabbedView) onTabUnselected(tab *container.TabItem) {
	v.mu.Lock()
	defer v.mu.Unlock()

	idx := v.tabs.SelectedIndex()
	if idx > -1 && idx < len(v.tabviews) {
		if v.tabviews[idx].Equals(v.currentTabView) {
			v.currentTabView = nil
		}
		v.tabviews[idx].OnHide()
	}
}
