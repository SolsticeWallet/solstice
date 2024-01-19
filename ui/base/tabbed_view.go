package base

import (
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
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

func (v *AbstractTabbedView) Initialize() (fyne.CanvasObject, error) {
	if err := v.configureTabs(); err != nil {
		Logger.Error("initialization failed", "error", err)
		return nil, err
	}
	return v.tabs, nil
}

func (v *AbstractTabbedView) Append(views ...View) {
	v.mu.Lock()
	defer v.mu.Unlock()

	for _, view := range views {
		if _, ok := v.tabviewIndex[view.ID()]; !ok {
			if canvasObj, err := view.Initialize(); err == nil {
				v.tabviews = append(v.tabviews, view)
				v.tabviewIndex[view.ID()] = view
				v.tabs.Append(container.NewTabItem(view.Title(), canvasObj))
			} else {
				Logger.Error(
					"append failed",
					"error", err,
					"view", view.Title())
			}
		}
	}
}

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
		curIdx = curIdx - 1
	}
	if curIdx < len(v.tabviews) {
		v.tabs.SelectIndex(curIdx)
	}
}

func (v AbstractTabbedView) NumTabs() int {
	return len(v.tabviews)
}

func (v AbstractTabbedView) TabIndex(view View) int {
	for idx, tabview := range v.tabviews {
		if tabview.Equals(view) {
			return idx
		}
	}
	return -1
}

func (v *AbstractTabbedView) SelectTabByIndex(idx int) {
	v.tabs.SelectIndex(idx)
}

func (v *AbstractTabbedView) SelectTab(view View) {
	idx := v.TabIndex(view)
	if idx < 0 {
		return
	}
	v.tabs.SelectIndex(idx)
}

func (v *AbstractTabbedView) configureTabs() error {
	v.mu.Lock()
	defer v.mu.Unlock()

	v.tabs = container.NewAppTabs()
	v.tabs.SetTabLocation(v.tabLocation)
	v.tabs.OnSelected = v.onTabSelected
	v.tabs.OnUnselected = v.onTabUnselected

	for _, tv := range v.tabviews {
		canvasObj, err := tv.Initialize()
		if err != nil {
			return err
		}
		v.tabs.Append(container.NewTabItem(tv.Title(), canvasObj))
		v.tabviewIndex[tv.ID()] = tv
	}
	if len(v.tabviews) > 0 {
		v.currentTabView = v.tabviews[0]
	}
	return nil
}

func (v *AbstractTabbedView) onTabSelected(tab *container.TabItem) {
	v.mu.Lock()
	defer v.mu.Unlock()

	idx := v.tabs.SelectedIndex()
	if idx > -1 && idx < len(v.tabviews) {
		v.currentTabView = v.tabviews[idx]
		v.currentTabView.OnShow()
	}
}

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
