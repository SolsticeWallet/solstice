package base

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/solsticewallet/solstice/i18n"
)

type WizardPanePrecondition func() bool

type WizardPane interface {
	View

	SetParentWindow(fyne.Window)
	GetParentWindow() fyne.Window

	Previous() WizardPane
	Next() WizardPane
	IsFinalPane() bool

	// SetState assigns the wizard state to the current page
	SetState(any) error

	CanTransitionTo(any) bool

	// ResetState undoes any changes this pane introduced into the
	// state. This is required to be able to move back in the wizard.
	ResetState()

	// IsValid should return true if the pane updated the state
	// in such a way that a next state or the end of the wizard
	// can be run.
	IsValid() bool

	AddNextPane(WizardPane, WizardPanePrecondition)
	SetPreviousPane(WizardPane)

	OnChanged(func())
	NotifyOnChanged()

	OnBeforeNext()
	OnBeforePrevious()
}

type WizardView interface {
	View

	AddStep(WizardPane, WizardPane, WizardPanePrecondition)

	NextStep()
	PreviousStep()

	GetState() any
}

type AbstractWizardView struct {
	*AbstractView

	parentWindow fyne.Window

	layout      *fyne.Container
	btnPrevious *widget.Button
	btnNext     *widget.Button
	btnCancel   *widget.Button
	btnFinish   *widget.Button

	rootPane      WizardPane
	currentPane   WizardPane
	currentCanvas fyne.CanvasObject

	getState          func() any
	onCancelCallback  func()
	onConfirmCallback func()
}

func NewAbstractWizardView(
	parentWindow fyne.Window,
	title string,
	rootPane WizardPane,
	getState func() any,
	onCancelCallback func(),
	onConfirmCallback func(),
) *AbstractWizardView {
	v := &AbstractWizardView{
		parentWindow:      parentWindow,
		AbstractView:      NewAbstractView(title),
		rootPane:          rootPane,
		currentPane:       rootPane,
		getState:          getState,
		onCancelCallback:  onCancelCallback,
		onConfirmCallback: onConfirmCallback,
	}
	v.rootPane.OnChanged(v.onPaneChangedCallback)
	return v
}

func (v *AbstractWizardView) Initialize() (fyne.CanvasObject, error) {
	if err := v.configurePanes(); err != nil {
		Logger.Error(
			i18n.T("Err.Initialization"),
			i18n.T("Err.Arg.Error"), err)
		return nil, err
	}
	return v.layout, nil
}

func (v *AbstractWizardView) AddStep(parent WizardPane, pane WizardPane, precondition WizardPanePrecondition) {
	linkWizardPanes(parent, pane, precondition)
	pane.OnChanged(v.onPaneChangedCallback)
	pane.SetParentWindow(v.parentWindow)
}

func (v *AbstractWizardView) NextStep() {
	var err error
	if v.currentPane.IsFinalPane() {
		return
	}
	if !v.currentPane.IsValid() {
		return
	}

	v.currentPane.OnBeforeNext()

	next := v.currentPane.Next()
	if next == nil {
		return
	}
	next.SetPreviousPane(v.currentPane)

	v.currentPane.OnHide()
	v.currentPane = next
	if err = v.currentPane.SetState(v.getState()); err != nil {
		Logger.Error(
			i18n.T("Err.NextWizardPane"),
			i18n.T("Err.Arg.Error"), err)
		return
	}

	v.layout.Remove(v.currentCanvas)
	if v.currentCanvas, err = v.currentPane.Initialize(); err != nil {
		Logger.Error(
			i18n.T("Err.NextWizardPane"),
			i18n.T("Err.Arg.Error"), err)
		return
	}
	v.layout.Add(v.currentCanvas)
	v.currentPane.OnShow()
	v.setButtonStates()
}

func (v *AbstractWizardView) PreviousStep() {
	var err error

	v.currentPane.OnBeforePrevious()

	prev := v.currentPane.Previous()
	if prev == nil {
		return
	}

	v.currentPane.ResetState()
	v.currentPane.OnHide()
	v.currentPane.SetPreviousPane(nil)
	v.currentPane = prev
	if err = v.currentPane.SetState(v.getState()); err != nil {
		Logger.Error(
			i18n.T("Err.PreviousWizardPane"),
			i18n.T("Err.Arg.Error"), err)
		return
	}

	v.layout.Remove(v.currentCanvas)
	if v.currentCanvas, err = v.currentPane.Initialize(); err != nil {
		Logger.Error(
			i18n.T("Err.PreviousWizardPane"),
			i18n.T("Err.Arg.Error"), err)
		return
	}
	v.layout.Add(v.currentCanvas)
	v.currentPane.OnShow()
	v.setButtonStates()
}

func (v *AbstractWizardView) OnShow() {
	v.currentPane.SetState(v.getState())
	v.currentPane.OnShow()
}

func (v *AbstractWizardView) OnHide() {
	v.currentPane.OnHide()
}

func (v *AbstractWizardView) configurePanes() (err error) {
	v.btnPrevious = widget.NewButtonWithIcon(
		i18n.T("Lbl.Previous"),
		theme.NavigateBackIcon(),
		func() {
			v.PreviousStep()
		})
	v.btnNext = widget.NewButtonWithIcon(
		i18n.T("Lbl.Next"),
		theme.NavigateNextIcon(),
		func() {
			v.NextStep()
		})

	v.btnCancel = widget.NewButtonWithIcon(
		i18n.T("Lbl.Cancel"),
		theme.CancelIcon(),
		func() {
			v.onCancelCallback()
		})
	v.btnFinish = widget.NewButtonWithIcon(
		i18n.T("Lbl.Finish"),
		theme.ConfirmIcon(),
		func() {
			v.onConfirmCallback()
		})

	btnLayout := container.NewBorder(
		nil, nil,
		container.NewHBox(v.btnPrevious, v.btnNext),
		container.NewHBox(v.btnCancel, v.btnFinish),
		nil,
	)

	if v.currentCanvas, err = v.currentPane.Initialize(); err != nil {
		return
	}

	v.layout = container.NewBorder(
		nil, btnLayout,
		nil, nil,
		v.currentCanvas,
	)

	v.setButtonStates()
	return
}

func (v *AbstractWizardView) setButtonStates() {
	v.btnPrevious.Enable()
	if v.currentPane.Previous() == nil {
		v.btnPrevious.Disable()
	}

	v.btnNext.Disable()
	if !v.currentPane.IsFinalPane() && v.currentPane.IsValid() {
		v.btnNext.Enable()
	}

	v.btnFinish.Disable()
	if v.currentPane.IsFinalPane() && v.currentPane.IsValid() {
		v.btnFinish.Enable()
	}
}

func (v *AbstractWizardView) onPaneChangedCallback() {
	v.setButtonStates()
}

type nextWizardPane struct {
	pane         WizardPane
	precondition WizardPanePrecondition
}

type AbstractWizardPane struct {
	*AbstractView

	parentWindow fyne.Window

	previous          WizardPane
	next              []nextWizardPane
	onChangedCallback func()
}

func NewAbstractWizardPane(
	title string,
) *AbstractWizardPane {
	pane := &AbstractWizardPane{
		AbstractView: NewAbstractView(title),
		previous:     nil,
		next:         []nextWizardPane{},
	}
	return pane
}

func (p *AbstractWizardPane) SetParentWindow(window fyne.Window) {
	p.parentWindow = window
}

func (p AbstractWizardPane) GetParentWindow() fyne.Window {
	return p.parentWindow
}

func (p AbstractWizardPane) Previous() WizardPane {
	return p.previous
}

func (p AbstractWizardPane) Next() WizardPane {
	for _, np := range p.next {
		if np.precondition() {
			return np.pane
		}
	}
	return nil
}

func (p AbstractWizardPane) IsFinalPane() bool {
	return len(p.next) == 0
}

func (p *AbstractWizardPane) AddNextPane(
	nextPane WizardPane,
	precondition WizardPanePrecondition,
) {
	p.next = append(p.next, nextWizardPane{
		pane:         nextPane,
		precondition: precondition,
	})
}

func (p *AbstractWizardPane) SetPreviousPane(previousPane WizardPane) {
	p.previous = previousPane
}

func (p *AbstractWizardPane) OnChanged(callback func()) {
	p.onChangedCallback = callback
}

func (p AbstractWizardPane) NotifyOnChanged() {
	if p.onChangedCallback != nil {
		p.onChangedCallback()
	}
}

func linkWizardPanes(
	parent WizardPane,
	next WizardPane,
	precondition WizardPanePrecondition,
) {
	parent.AddNextPane(next, precondition)
}
