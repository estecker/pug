package tui

import tea "charm.land/bubbletea/v2"

func ReportError(err error) tea.Cmd {
	return CmdHandler(ErrorMsg(err))
}

type ErrorMsg error
