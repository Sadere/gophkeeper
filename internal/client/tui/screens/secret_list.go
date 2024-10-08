package screens

import (
	"context"
	"fmt"
	"io"

	tuiMsg "github.com/Sadere/gophkeeper/internal/client/tui/msg"
	"github.com/Sadere/gophkeeper/internal/client/tui/style"
	"github.com/Sadere/gophkeeper/pkg/constants"
	"github.com/Sadere/gophkeeper/pkg/model"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"

	tea "github.com/charmbracelet/bubbletea"
)

var errInvalidItem = "invalid item"

type secretListItem struct {
	Preview *model.SecretPreview
}

func (i secretListItem) Title() string {
	return fmt.Sprintf(
		"%s %s %s",
		i.Status(),
		i.Icon(),
		i.Preview.Metadata,
	)
}
func (i secretListItem) Description() string {
	return fmt.Sprintf(
		"Created: %v Updated: %v",
		i.Preview.CreatedAt.Format(constants.TimeFormat),
		i.Preview.UpdatedAt.Format(constants.TimeFormat),
	)
}
func (i secretListItem) FilterValue() string { return i.Preview.Metadata }

func (i secretListItem) Icon() string {
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

func (i secretListItem) Status() string {
	status := ""

	switch i.Preview.Status {
	case model.SecretPreviewNew:
		status = "*** NEW ***"
	case model.SecretPreviewUpdated:
		status = "*** EDITED ***"
	}

	return status
}

// Renders list of secret previews, shows type of secret, status of secret (new, updated, normal)
type SecretListModel struct {
	state    *State
	list     list.Model
	errorMsg string
}

type MyDelegate struct {
	list.DefaultDelegate
}

func (d MyDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	d.DefaultDelegate.Styles.NormalDesc = style.BlurredStyle
	d.DefaultDelegate.Styles.NormalTitle = style.BlurredStyle

	d.DefaultDelegate.Styles.SelectedDesc = style.FocusedStyle
	d.DefaultDelegate.Styles.SelectedTitle = style.FocusedStyle

	item, ok := listItem.(secretListItem)
	if ok {
		switch item.Preview.Status {
		case model.SecretPreviewNew:
			d.DefaultDelegate.Styles.NormalDesc = style.NewSecretStyle
			d.DefaultDelegate.Styles.NormalTitle = style.NewSecretStyle
		case model.SecretPreviewUpdated:
			d.DefaultDelegate.Styles.NormalDesc = style.UpdatedSecretStyle
			d.DefaultDelegate.Styles.NormalTitle = style.UpdatedSecretStyle
		}
	}

	d.DefaultDelegate.Render(w, m, index, listItem)
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
		items = append(items, list.Item(secretListItem{
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
		m.list.SetSize(msg.Width, msg.Height-4)
	case tea.KeyMsg:
		switch msg.String() {
		case "r":
			return m.Reload()
		case "a":
			return m.AddNew()
		case "enter":
			selectedItem, ok := m.list.SelectedItem().(secretListItem)
			if !ok {
				m.errorMsg = errInvalidItem
				return m, nil
			}

			return m.Edit(selectedItem.Preview)
		}
	case tuiMsg.ReloadSecretList:
		return m.Reload()
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m *SecretListModel) Reload() (tea.Model, tea.Cmd) {
	previews, err := m.state.client.LoadPreviews(context.Background())
	if err != nil {
		m.errorMsg = fmt.Sprintf("failed to refresh secrets: %s", err.Error())
		return m, nil
	}

	items := PreviewsToItems(previews)
	return m, m.list.SetItems(items)
}

func (m *SecretListModel) AddNew() (tea.Model, tea.Cmd) {
	chooseTypeModel := NewChooseTypeModel(m.state)
	return NewRootModel(m.state).SwitchScreen(chooseTypeModel)
}

func (m *SecretListModel) Edit(preview *model.SecretPreview) (tea.Model, tea.Cmd) {
	// Set as read
	preview.Status = model.SecretPreviewRead

	switch preview.SType {
	case string(model.CredSecret):
		credModel := NewCredentialModel(m.state, preview.ID)
		return NewRootModel(m.state).SwitchScreen(credModel)
	case string(model.TextSecret):
		textModel := NewTextModel(m.state, preview.ID)
		return NewRootModel(m.state).SwitchScreen(textModel)
	case string(model.CardSecret):
		cardModel := NewCardModel(m.state, preview.ID)
		return NewRootModel(m.state).SwitchScreen(cardModel)
	case string(model.BlobSecret):
		fileModel := NewFileModel(m.state, preview.ID)
		return NewRootModel(m.state).SwitchScreen(fileModel)
	}

	return m, nil
}

func (m *SecretListModel) View() string {
	listView := m.list.View()

	if len(m.errorMsg) > 0 {
		return lipgloss.JoinVertical(lipgloss.Top, listView, style.ErrorStyle.Render(m.errorMsg))
	}

	return listView
}
