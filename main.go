package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	errMsg error
)

var choices = []string{"buy milk", "buy butter", "water the plants"}

type keyMap struct {
	Up    key.Binding
	Down  key.Binding
	Space key.Binding
	Esc   key.Binding
	Help  key.Binding
	Quit  key.Binding
	Enter key.Binding
}

type model struct {
	textInput textinput.Model
	err       error
	choice    string
	cursor    int
	selected  map[int]struct{}
	keys      keyMap
	help      help.Model
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Esc, k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Esc, k.Space},
		{k.Enter, k.Help, k.Quit},
	}
}

var keys = keyMap{
	Up:    key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "move up")),
	Down:  key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "move down")),
	Esc:   key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "unfocus input")),
	Space: key.NewBinding(key.WithKeys("space"), key.WithHelp("space", "select")),
	Help:  key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "toggle help")),
	Quit:  key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
	Enter: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "add a new todo")),
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Add a new Todo"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return model{
		textInput: ti,
		err:       nil,
		choice:    " ",
		selected:  make(map[int]struct{}),
		keys:      keys,
		help:      help.New(),
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// If we set a width on the help menu it can it can gracefully truncate
		// its view as needed.
		m.help.Width = msg.Width

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Help):
			if !m.textInput.Focused() {
				m.help.ShowAll = !m.help.ShowAll
			}

		case key.Matches(msg, m.keys.Quit):
			if !m.textInput.Focused() {
				return m, tea.Quit
			}

		case key.Matches(msg, m.keys.Esc):
			m.textInput.Blur()
			return m, nil

		}

		switch msg.String() {
		// handle add new todo
		case "enter":
			if m.textInput.Focused() {

				choice := string(m.textInput.Value())
				choices = append(choices, choice)
				m.textInput.Reset()
				return m, nil
			} else {
				m.textInput.Focus()
			}

		// handle up & down
		case "up", "k":
			if !m.textInput.Focused() {
				if m.cursor > 0 {
					m.cursor--
				}
			}
		case "down", "j":
			if !m.textInput.Focused() {
				if m.cursor < len(choices)-1 {
					m.cursor++
				}
			}

		// handle selection
		case " ":
			if !m.textInput.Focused() {
				_, ok := m.selected[m.cursor]
				if ok {
					delete(m.selected, m.cursor)
				} else {
					m.selected[m.cursor] = struct{}{}
				}
			}
		// handle q-uit
		case "q":
			if !m.textInput.Focused() {
				return m, tea.Quit
			}
		}

		// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd

}

func (m model) View() string {
	s := "Your Todos\n\n"

	for i, choice := range choices {
		cursor := "  "
		if m.cursor == i {
			cursor = "👉"
		}
		checked := "  "
		if _, ok := m.selected[i]; ok {
			checked = "❌"
		}
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}
	helpView := m.help.View(m.keys)
	s += fmt.Sprintf("%s\n\n---\n%s\n", m.textInput.View(), helpView)
	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
