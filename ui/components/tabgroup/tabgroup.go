package tabgroup

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/remogatto/sugarfoam/ui"
)

type Option func(*TabGroup)

type KeyMap struct {
	TabNext key.Binding
	TabPrev key.Binding
}

type Styles struct {
	FocusedBorder         lipgloss.Style
	BlurredBorder         lipgloss.Style
	Navbar                lipgloss.Style
	NavbarTitleUnselected lipgloss.Style
	NavbarTitleSelected   lipgloss.Style
}

type TabItem struct {
	ui.Groupable

	Title  string
	Active bool
}

type TabGroup struct {
	KeyMap KeyMap

	items         []*TabItem
	currItemIndex int

	width, height int
	focused       bool
	styles        *Styles
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		TabNext: key.NewBinding(
			key.WithKeys("alt+right"),
			key.WithHelp("alt+right", "Next tab"),
		),
		TabPrev: key.NewBinding(
			key.WithKeys("alt+left"),
			key.WithHelp("alt+left", "Prev tab")),
	}
}

func New(opts ...Option) *TabGroup {
	tg := &TabGroup{
		items: make([]*TabItem, 0),
	}

	tg.KeyMap = DefaultKeyMap()
	tg.styles = DefaultStyles()

	return tg
}

func (tg *TabGroup) Items() []*TabItem {
	return tg.items
}

func (tg *TabGroup) AddItem(item *TabItem) *TabGroup {
	tg.items = append(tg.items, item)

	return tg
}

func DefaultStyles() *Styles {
	return &Styles{
		FocusedBorder: lipgloss.NewStyle().Margin(2, 0, 0, 0), /*Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("5")),*/
		BlurredBorder: lipgloss.NewStyle().Margin(2, 0, 0, 0), /*Border(lipgloss.RoundedBorder()),*/
		Navbar:        lipgloss.NewStyle().Padding(0, 1, 0, 1),
		NavbarTitleUnselected: lipgloss.NewStyle().
			Background(lipgloss.Color("#373B41")).
			Foreground(lipgloss.Color("240")).
			Padding(0, 2, 0, 2),
		NavbarTitleSelected: lipgloss.NewStyle().
			Background(lipgloss.Color("5")).
			Foreground(lipgloss.Color("#ffffff")).
			Padding(0, 2, 0, 2),
	}
}

func (tg *TabGroup) Init() tea.Cmd {
	var cmds []tea.Cmd

	for _, item := range tg.items {
		cmds = append(cmds, item.Init())
	}

	return tea.Batch(cmds...)
}

func (tg *TabGroup) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, tg.KeyMap.TabNext):
			tg.nextTab()
		case key.Matches(msg, tg.KeyMap.TabPrev):
			tg.prevTab()
		}
	}

	cmds = append(cmds, tg.updateTabItems(msg)...)

	return tg, tea.Batch(cmds...)

}

func (tg *TabGroup) Focus() tea.Cmd {
	tg.focused = true
	return nil
}

func (tg *TabGroup) Blur() {
	tg.focused = false
}

func (tg *TabGroup) SetSize(width int, height int) {
	tg.width = width - tg.styles.FocusedBorder.GetHorizontalFrameSize()
	tg.height = height

	tg.styles.FocusedBorder = tg.styles.FocusedBorder.Width(tg.width)
	tg.styles.BlurredBorder = tg.styles.BlurredBorder.Width(tg.width)

	tg.styles.Navbar = tg.styles.Navbar.Width(tg.width)

	for _, item := range tg.items {
		item.SetSize(tg.width, tg.height)
	}
}

func (tg *TabGroup) Width() int  { return tg.width }
func (tg *TabGroup) Height() int { return tg.height }

func (tg *TabGroup) View() string {
	var navbar string

	for i, item := range tg.Items() {
		tabTitle := tg.styles.NavbarTitleUnselected.Render(item.Title)
		if tg.currItemIndex == i {
			tabTitle = tg.styles.NavbarTitleSelected.Render(item.Title)
		}
		navbar += fmt.Sprintf("%s • ", tabTitle)
	}

	navbar = tg.styles.Navbar.Render(strings.TrimRight(navbar, "• "))

	if len(tg.items) > 0 {
		return tg.styles.FocusedBorder.Render(lipgloss.JoinVertical(lipgloss.Top, navbar, tg.items[tg.currItemIndex].View()))
	}

	return tg.styles.FocusedBorder.Render(navbar)

}

func (tg *TabGroup) nextTab() {
	tg.currItemIndex = (tg.currItemIndex + 1) % len(tg.items)
}

func (tg *TabGroup) prevTab() {
	tg.currItemIndex = (tg.currItemIndex - 1) % len(tg.items)

	if tg.currItemIndex < 0 {
		tg.currItemIndex = len(tg.items) - 1
	}
}

func (tg *TabGroup) updateTabItems(msg tea.Msg) []tea.Cmd {
	cmds := make([]tea.Cmd, 0)

	for _, item := range tg.items {
		_, cmd := item.Update(msg)
		cmds = append(cmds, cmd)
	}

	return cmds
}
