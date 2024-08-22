package screens

import (
	"context"
	"fmt"
	"io"

	"github.com/Sadere/gophkeeper/internal/client/tui/style"
	"github.com/Sadere/gophkeeper/pkg/constants"
	"github.com/Sadere/gophkeeper/pkg/model"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"

	tea "github.com/charmbracelet/bubbletea"
)

type item struct {
	Preview *model.SecretPreview
}

func (i item) Title() string {

	return fmt.Sprintf(
		"%s %s Created: %v Updated: %v",
		i.Status(),
		i.Icon(),
		i.Preview.CreatedAt.Format(constants.TimeFormat),
		i.Preview.UpdatedAt.Format(constants.TimeFormat),
	)
}
func (i item) Description() string { return i.Preview.Metadata }
func (i item) FilterValue() string { return i.Preview.Metadata }

func (i item) Icon() string {
	icon := "📃"

	switch i.Preview.SType {
	case string(model.CredSecret):
		icon = "🔑"
	case string(model.TextSecret):
		icon = "📒"
	case string(model.BlobSecret):
		icon = "📁"
	case string(model.CardSecret):
		icon = "💳"
	}

	return icon
}

func (i item) Status() string {
	status := ""

	switch i.Preview.Status {
	case model.SecretPreviewNew:
		status = "*** NEW ***"
	case model.SecretPreviewUpdated:
		status = "*** EDITED ***"
	}

	return status
}

type SecretListModel struct {
	state    *State
	list     list.Model
	errorMsg string
}

type MyDelegate struct {
	list.DefaultDelegate
}

func (d MyDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	d.DefaultDelegate.Styles.DimmedDesc = style.BlurredStyle
	d.DefaultDelegate.Styles.DimmedTitle = style.BlurredStyle

	d.DefaultDelegate.Styles.SelectedDesc = style.FocusedStyle
	d.DefaultDelegate.Styles.SelectedTitle = style.FocusedStyle

	d.DefaultDelegate.Render(w, m, index, item)
}

func NewSecretListModel(state *State) *SecretListModel {
	m := SecretListModel{
		state: state,
		list:  list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0),
	}

	// Load previews
	previews, err := state.client.LoadPreviews(context.Background())
	if err != nil {
		m.errorMsg = fmt.Sprintf("failed to retrieve secret list: %s", err.Error())
		return &m
	}

	items := PreviewsToItems(previews)

	// set up list
	delegate := MyDelegate{
		list.NewDefaultDelegate(),
	}

	m.list = list.New(items, delegate, 0, 0)
	m.list.Title = "My secrets list"
	m.list.DisableQuitKeybindings()
	m.list.SetShowFilter(false)
	m.list.SetShowHelp(false)

	return &m
}

func PreviewsToItems(previews model.SecretPreviews) []list.Item {
	var items []list.Item
	for _, preview := range previews {
		items = append(items, list.Item(item{
			Preview: preview,
		}))
	}

	return items
}

func (m SecretListModel) Init() tea.Cmd {
	return nil
}

func (m *SecretListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch msg.String() {
		case "r":
			previews, err := m.state.client.LoadPreviews(context.Background())
			if err != nil {
				m.errorMsg = fmt.Sprintf("failed to refresh secrets: %s", err.Error())
				return m, nil
			}

			items := PreviewsToItems(previews)
			return m, m.list.SetItems(items)
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m *SecretListModel) View() string {
	listView := m.list.View()

	if len(m.errorMsg) > 0 {
		return lipgloss.JoinVertical(lipgloss.Top, listView, style.ErrorStyle.Render(m.errorMsg))
	}

	return listView
}
