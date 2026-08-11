package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Bubble Tea apps follow the "Elm architecture": a single Model holds all
// state, Update() reacts to messages (keypresses, etc.) and returns a new
// Model, and View() renders that Model to a string every frame. There is no
// other way for state to change — everything flows through Update.

// mode tracks whether the user is browsing the list (normal) or typing a
// new item (insert). Update() branches on this to decide which keys mean
// what, similar to Vim's modal editing.
type mode int

const (
	modeNormal mode = iota
	modeInsert
)

// item is one row in the task list.
type item struct {
	text string
	done bool
}

// model is THE state for the whole program. Bubble Tea passes it by value
// into Update and View each time, so every state change is expressed as
// "return a modified copy" rather than mutating shared state in place.
type model struct {
	items    []item
	cursor   int  // index of the currently selected item
	mode     mode // modeNormal or modeInsert
	input    textinput.Model // a reusable sub-component from the bubbles library that owns its own text-editing state (cursor position, blink, etc.)
	quitting bool
}

// initialModel builds the starting state. This is the seed value handed to
// tea.NewProgram in main().
func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "New todo..."
	ti.CharLimit = 80

	return model{
		items: []item{
			{text: "Press 'a' to add a todo", done: false},
			{text: "Use j/k or arrows to move", done: false},
			{text: "Press space/enter to toggle done", done: false},
			{text: "Press 'd' to delete, 'q' to quit", done: false},
		},
		input: ti,
	}
}

// Init runs once when the program starts, before the first View render. It
// can return a tea.Cmd to kick off async work (HTTP calls, timers, etc.) —
// this program has nothing to do at startup, so it returns nil.
func (m model) Init() tea.Cmd {
	return nil
}

// Update is the only place state changes. Bubble Tea calls it once per
// incoming tea.Msg (a keypress, a window resize, a command's result, ...)
// and it must return the (possibly new) model plus an optional tea.Cmd —
// a side-effecting function Bubble Tea will run and feed back in as a
// future Msg. Returning tea.Quit is how the program exits.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Dispatch to a different handler depending on which "mode" we're
		// in, so the same key (e.g. Enter) can mean different things.
		if m.mode == modeInsert {
			return m.updateInsert(msg)
		}
		return m.updateNormal(msg)
	}
	return m, nil
}

// updateNormal handles keys while browsing the list.
func (m model) updateNormal(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		m.quitting = true
		// tea.Quit is a built-in tea.Cmd that tells the runtime to stop the
		// program after this Update call finishes and View renders once more.
		return m, tea.Quit

	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}

	case "down", "j":
		if m.cursor < len(m.items)-1 {
			m.cursor++
		}

	case " ", "enter":
		if len(m.items) > 0 {
			m.items[m.cursor].done = !m.items[m.cursor].done
		}

	case "d":
		if len(m.items) > 0 {
			// Remove the item at cursor by slicing around it, then keep the
			// cursor in bounds if we just deleted the last row.
			m.items = append(m.items[:m.cursor], m.items[m.cursor+1:]...)
			if m.cursor >= len(m.items) && m.cursor > 0 {
				m.cursor--
			}
		}

	case "a":
		m.mode = modeInsert
		m.input.Focus()
		// textinput.Blink is the Cmd that drives the cursor's blink
		// animation — without returning it here, the input's cursor
		// would sit still instead of blinking.
		return m, textinput.Blink
	}
	return m, nil
}

// updateInsert handles keys while typing a new task into the text input.
func (m model) updateInsert(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.mode = modeNormal
		m.input.Blur()
		m.input.SetValue("")
		return m, nil

	case "enter":
		if text := m.input.Value(); text != "" {
			m.items = append(m.items, item{text: text})
		}
		m.mode = modeNormal
		m.input.Blur()
		m.input.SetValue("")
		return m, nil
	}

	// Any other key (letters, backspace, arrow keys inside the field, ...)
	// gets forwarded to the textinput sub-component so it can handle its
	// own editing logic. This is the standard pattern for composing bubbles
	// components: delegate the message, capture the returned Cmd, and pass
	// it back up so Bubble Tea still runs it.
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

// Lip Gloss styles are immutable, reusable style definitions — declare them
// once and call .Render(text) wherever you need that look.
var (
	titleStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7D56F4"))
	cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#F25D94"))
	doneStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#5C5C5C")).Strikethrough(true)
	helpStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#5C5C5C"))
	insertStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4"))
)

// View renders the current model to a plain string. Bubble Tea calls this
// after every Update and redraws the terminal with whatever it returns —
// View must never mutate m or have side effects, it just describes what the
// screen should look like right now.
func (m model) View() string {
	if m.quitting {
		return "Bye!\n"
	}

	s := titleStyle.Render("My Todos") + "\n\n"

	for i, it := range m.items {
		cursor := "  "
		if i == m.cursor {
			cursor = cursorStyle.Render("> ")
		}

		checkbox := "[ ]"
		text := it.text
		if it.done {
			checkbox = "[x]"
			text = doneStyle.Render(it.text)
		}

		s += fmt.Sprintf("%s%s %s\n", cursor, checkbox, text)
	}

	s += "\n"

	// Swap in the text input's own View() while adding, otherwise show the
	// normal-mode help text. This is the same delegation pattern as
	// updateInsert: the sub-component knows how to render itself.
	if m.mode == modeInsert {
		s += insertStyle.Render("Add: ") + m.input.View() + "\n"
		s += helpStyle.Render("enter: save  •  esc: cancel")
	} else {
		s += helpStyle.Render("j/k: move  •  space: toggle  •  a: add  •  d: delete  •  q: quit")
	}

	return s
}

func main() {
	// tea.NewProgram wires up the terminal (raw mode, alt screen, input
	// reading) and drives the Init -> Update -> View loop until a Cmd
	// returns tea.Quit or the process is interrupted.
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
