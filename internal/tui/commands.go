package tui

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	tea "charm.land/bubbletea/v2"
	"github.com/atotto/clipboard"
)

// NavigateTo sends an instruction to navigate to a page with the given model
// kind, and optionally parent resource.
func NavigateTo(kind Kind, opts ...NavigateOption) tea.Cmd {
	return CmdHandler(NewNavigationMsg(kind, opts...))
}

func ReportInfo(msg string) tea.Cmd {
	return CmdHandler(InfoMsg(msg))
}

func CopyToClipboard(content string) tea.Cmd {
	return func() tea.Msg {
		if err := clipboard.WriteAll(content); err != nil {
			return ErrorMsg(fmt.Errorf("copying to clipboard: %w", err))
		}
		return InfoMsg("Content copied to clipboard")
	}
}

func OpenEditor(path string) tea.Cmd {
	// TODO: check for side effects of exec blocking the tui - do
	// messages get queued up?
	editor, ok := os.LookupEnv("EDITOR")
	if !ok {
		return ReportError(errors.New("cannot open editor: environment variable EDITOR not set"))
	}
	cmd := exec.Command(editor, path)
	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		if err != nil {
			return ReportError(fmt.Errorf("opening %s in editor: %w", path, err))()
		}
		return nil
	})
}
