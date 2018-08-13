package main

import "fmt"

// SelectMenuEntry reads a key from the terminal and changes the
// current key map to reflect this selection.
type SelectMenuEntry struct {
	*Context
}

// NewSelectMenuEntry creates a new instance of this command bound to the given context.
func NewSelectMenuEntry(ctx *Context) *SelectMenuEntry {
	return &SelectMenuEntry{
		Context: ctx,
	}
}

// Execute runs this command.
func (cmd *SelectMenuEntry) Execute() error {
	for {
		menu, err := cmd.displayMenu()
		if err != nil {
			return err
		}
		key, err := cmd.Terminal.ReadKey()
		if err != nil {
			return fmt.Errorf("SelectMenuEntry: %s", err)
		}
		binding := cmd.CurrentKeyMap.LookupKey(key)
		if binding.HasChildren() {
			if err := menu.Erase(cmd.Terminal); err != nil {
				return err
			}
			cmd.CurrentKeyMap = binding.Children()
		} else {
			if err := cmd.Terminal.Restore(); err != nil {
				return err
			}
			return binding.Execute()
		}
	}
}

// displayMenu displays a menu for the current keymap
func (cmd *SelectMenuEntry) displayMenu() (*MenuView, error) {
	menuEntries := []*MenuEntry{}
	for _, binding := range cmd.CurrentKeyMap.Bindings() {
		menuEntries = append(menuEntries, NewMenuEntryForKeyBinding(binding))
	}

	menu := NewMenuView(menuEntries)
	return menu, menu.Render(cmd.Terminal)
}