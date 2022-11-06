package app

import (
	"context"
	"image"
	"rohitsingh/pomo/pomodoro"
	"time"

	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/terminal/tcell"
	"github.com/mum4k/termdash/terminal/terminalapi"
)

// App type is used by callers to instantiate and control the interface
// this type includes private fields that are required to control, redraw and resize the interface
type App struct {
	ctx        context.Context
	controller *termdash.Controller
	redrawCh   chan bool
	errorCh    chan error
	term       *tcell.Terminal
	size       image.Point
}

// New instantiates a new app, and creates its required widgets, buttons and grid,
// combining them all in a new instance of termdash.Controller
func New(config *pomodoro.IntervalConfig) (*App, error) {
	// Define a new cancellation context that is used when the application closes
	ctx, cancel := context.WithCancel(context.Background())
	// Define a quitter function that maps to the quit key to exit on user request
	quitter := func(k *terminalapi.Keyboard) {
		if k.Key == 'q' || k.Key == 'Q' {
			cancel()
		}
	}
	// Define two channels to control the async redrawing
	redrawCh := make(chan bool)
	errorCh := make(chan error)
	// Instantiate the widgets and buttons
	w, err := newWidgets(ctx, errorCh)
	if err != nil {
		return nil, err
	}
	b, err := newButtonSet(ctx, config, w, redrawCh, errorCh)
	if err != nil {
		return nil, err
	}
	// Define a new instance of Terminal to use as the backend for the application
	term, err := tcell.New()
	if err != nil {
		return nil, err
	}
	// Define a new container using the newGrid function
	c, err := newGrid(b, w, term)
	if err != nil {
		return nil, err
	}
	// Use NewController function to create a controller for our application
	controller, err := termdash.NewController(term, c, termdash.KeyboardSubscriber(quitter))
	if err != nil {
		return nil, err
	}
	// Return a pointer of the app with the defined instances
	return &App{
		ctx:        ctx,
		controller: controller,
		redrawCh:   redrawCh,
		errorCh:    errorCh,
		term:       term,
	}, nil
}

// resize method resizes the interface if needed
// this method is run periodically to verify whether we need to do anything
func (a *App) resize() error {
	if a.size.Eq(a.term.Size()) {
		return nil
	}
	a.size = a.term.Size()
	if err := a.term.Clear(); err != nil {
		return err
	}
	return a.controller.Redraw()
}

// Run method runs and controls the application
func (a *App) Run() error {
	// Specify cleanup processses
	defer a.term.Close()
	defer a.controller.Close()
	ticker := time.NewTicker(2 * time.Second) // 2-second ticker to check if a resize is needed
	defer ticker.Stop()
	// Run the application logic depending on data arriving on one of four channels
	for {
		select {
		// Redraw the application using the controller's built-in redraw method
		case <-a.redrawCh:
			if err := a.controller.Redraw(); err != nil {
				return err
			}
		// Return the error received by the channel that causes the program to exit
		case err := <-a.errorCh:
			if err != nil {
				return err
			}
		// Data in the ctx channel indicates the main context was cancelled by user, exit gracefully
		case <-a.ctx.Done():
			return nil
		// This implies the ticker has expired, so we check for a resize.
		case <-ticker.C:
			if err := a.resize(); err != nil {
				return err
			}
		}

	}
}
