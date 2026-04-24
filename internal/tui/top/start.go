package top

import (
	"context"
	"sync"
	"testing"

	tea "charm.land/bubbletea/v2"
	teav1 "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/leg100/pug/internal/app"
	"github.com/leg100/pug/internal/resource"
	"github.com/stretchr/testify/require"
)

// modelAdapter adapts our v2 model to the v1 interface for teatest
type modelAdapter struct {
	v2Model model
}

func (a modelAdapter) Init() teav1.Cmd {
	cmd := a.v2Model.Init()
	if cmd == nil {
		return nil
	}
	// Wrap the v2 cmd to return v1 msg
	return func() teav1.Msg {
		msg := cmd()
		return msg
	}
}

func (a modelAdapter) Update(msg teav1.Msg) (teav1.Model, teav1.Cmd) {
	newModel, cmd := a.v2Model.Update(msg)
	a.v2Model = newModel.(model)
	if cmd == nil {
		return a, nil
	}
	// Wrap the v2 cmd to return v1 msg
	return a, func() teav1.Msg {
		msg := cmd()
		return msg
	}
}

func (a modelAdapter) View() string {
	return a.v2Model.View().Content
}

// Start starts the TUI and blocks until the user exits.
func Start(cfg app.Config) error {
	app, err := app.New(cfg)
	if err != nil {
		return err
	}
	defer app.Cleanup()

	m, err := newModel(cfg, app)
	if err != nil {
		return err
	}

	p := tea.NewProgram(m)// Enabling mouse cell motion removes the ability to "blackboard" text
	// with the mouse, which is useful for then copying text into the
	// clipboard. Therefore we've decided to disable it and leave it
	// commented out for posterity.
	//
	// tea.WithMouseCellMotion(),

	ch, unsub := setupSubscriptions(app, cfg)
	defer unsub()

	// Relay events to model in background
	go func() {
		for msg := range ch {
			p.Send(msg)
		}
	}()

	// Blocks until user quits
	_, err = p.Run()
	return err
}

// StartTest starts the TUI and returns a test model for testing purposes.
func StartTest(t *testing.T, cfg app.Config, width, height int) *teatest.TestModel {
	app, err := app.New(cfg)
	if err != nil {
		return nil
	}
	t.Cleanup(app.Cleanup)

	m, err := newModel(cfg, app)
	require.NoError(t, err)

	ch, unsub := setupSubscriptions(app, cfg)
	t.Cleanup(unsub)

	// Wrap our v2 model with an adapter for teatest (which still uses v1 API)
	adapter := modelAdapter{v2Model: m}
	tm := teatest.NewTestModel(t, adapter, teatest.WithInitialTermSize(width, height))

	// Relay events to model in background
	go func() {
		for msg := range ch {
			tm.Send(msg)
		}
	}()

	t.Cleanup(func() {
		tm.Quit()
	})
	return tm
}

func setupSubscriptions(app *app.App, cfg app.Config) (chan tea.Msg, func()) {
	// Relay resource events to TUI. Deliberately set up subscriptions *before*
	// any events are triggered, to ensure the TUI receives all messages.
	ch := make(chan tea.Msg)
	wg := sync.WaitGroup{} // sync closure of subscriptions

	ctx, cancel := context.WithCancel(context.Background())

	{
		sub := app.Logger.Subscribe(ctx)
		wg.Add(1)
		go func() {
			for ev := range sub {
				ch <- ev
			}
			wg.Done()
		}()
	}
	{
		sub := app.Modules.Subscribe(ctx)
		wg.Add(1)
		go func() {
			for ev := range sub {
				ch <- ev
			}
			wg.Done()
		}()
	}
	{
		sub := app.Workspaces.Subscribe(ctx)
		wg.Add(1)
		go func() {
			for ev := range sub {
				ch <- ev
			}
			wg.Done()
		}()

	}
	{
		sub := app.States.Subscribe(ctx)
		wg.Add(1)
		go func() {
			for ev := range sub {
				ch <- ev
			}
			wg.Done()
		}()

	}
	{
		sub := app.Plans.Subscribe(ctx)
		wg.Add(1)
		go func() {
			for ev := range sub {
				ch <- ev
			}
			wg.Done()
		}()

	}
	{
		sub := app.Tasks.TaskBroker.Subscribe(ctx)
		wg.Add(1)
		go func() {
			for ev := range sub {
				ch <- ev
			}
			wg.Done()
		}()

	}
	{
		sub := app.Tasks.GroupBroker.Subscribe(ctx)
		wg.Add(1)
		go func() {
			for ev := range sub {
				ch <- ev
			}
			wg.Done()
		}()
	}
	// Automatically load workspaces whenever modules are loaded.
	{
		sub := app.Modules.Subscribe(ctx)
		go app.Workspaces.LoadWorkspacesUponModuleLoad(sub)
	}
	// Automatically load workspaces whenever init is run and workspaces have
	// not yet been loaded.
	{
		sub := app.Tasks.TaskBroker.Subscribe(ctx)
		go app.Workspaces.LoadWorkspacesUponInit(sub)
	}
	// Whenever a workspace is loaded, pull its state
	{
		sub := app.Workspaces.Subscribe(ctx)
		go func() {
			for event := range sub {
				if event.Type == resource.CreatedEvent {
					_, err := app.States.CreateReloadTask(event.Payload.ID)
					if err != nil {
						app.Logger.Error("loading state after loading workspace", "error", err)
					}
				}
			}
		}()
	}
	// Whenever an apply is successful, pull workspace state
	if !cfg.DisableReloadAfterApply {
		sub := app.Tasks.TaskBroker.Subscribe(ctx)
		go app.Plans.ReloadAfterApply(sub)
	}
	// cleanup function to be invoked when program is terminated.
	return ch, func() {
		cancel()
		// Wait for relays to finish before closing channel, to avoid sends
		// to a closed channel, which would result in a panic.
		wg.Wait()
		close(ch)
	}
}
