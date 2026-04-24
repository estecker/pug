package app

import (
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/exp/teatest"
)

func TestQuit(t *testing.T) {
	t.Parallel()

	tm := setup(t, "./testdata/module_list")

	tm.Send(tea.KeyPressMsg{
		Code: 'c',
		Mod:  tea.ModCtrl,
	})

	waitFor(t, tm, func(s string) bool {
		return strings.Contains(s, "Quit pug? (y/N): ")
	})

	tm.Type("y")

	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
}
