package app

import (
	"context"
	"fmt"
	"rohitsingh/pomo/pomodoro"

	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/widgets/button"
)

// buttonSet represents the start and pause buttons
type buttonSet struct {
	btStart *button.Button // Start interval button
	btPause *button.Button // Pause interval button
}

// newButtonSet instantiates a buttonSet object,
// takes a context to carry cancellation signals, a config instance to call Pomodoro functions,
// a pointer to the widgets and redrawCh and errorCh channels
func newButtonSet(ctx context.Context, config *pomodoro.IntervalConfig, w *widgets, redrawCh chan<- bool, errorCh chan<- error) (*buttonSet, error) {
	// Button actions must be lightweight and blocking since they can be pressed many times,
	// We need a lightweight function to start and interval and assign it to a variable
	startInterval := func() {
		// Try to get the interval using pomodoro package
		i, err := pomodoro.GetInterval(config)
		errorCh <- err
		// Define the callbacks that will be used in the interval
		start := func(i pomodoro.Interval) {
			message := "Take a break"
			if i.Category == pomodoro.CategoryPomodoro {
				message = "Focus on your task"
			}
			w.update([]int{}, i.Category, message, "", redrawCh)
		}
		end := func(pomodoro.Interval) {
			w.update([]int{}, "", "Nothing running...", "", redrawCh)
		}
		periodic := func(i pomodoro.Interval) {
			w.update(
				[]int{int(i.ActualDuration), int(i.PlannedDuration)},
				"", "",
				fmt.Sprint(i.PlannedDuration-i.ActualDuration),
				redrawCh,
			)
		}
		// Attempt to start the interval by calling the interval's start message and send any errors to errorCh for handling
		errorCh <- i.Start(ctx, config, start, periodic, end)
	}
	// Similarly, we need a lightweight function to puase an interval and assign it to a variable
	pauseInterval := func() {
		i, err := pomodoro.GetInterval(config)
		if err != nil {
			errorCh <- err
			return
		}
		if err := i.Pause(config); err != nil {
			if err == pomodoro.ErrIntervalNotRunning {
				return
			}
			errorCh <- err
			return
		}
		w.update([]int{}, "", "Paused... press start to continue", "", redrawCh)
	}
	// With the two actions defined, we must instantiate buttons and attach the functions to them
	btStart, err := button.New("(s)tart", func() error {
		go startInterval()
		return nil
	},
		button.GlobalKey('s'),      // Global key that users use to activate the button
		button.WidthFor("(p)ause"), // Width for the button as if it was displaying (p)ause to ensure all buttons are the same size
		button.Height(2),           // Setting the button height to 2 cells
	)
	if err != nil {
		return nil, err
	}
	btPause, err := button.New("(p)ause", func() error {
		go pauseInterval()
		return nil
	},
		button.FillColor(cell.ColorNumber(220)),
		button.GlobalKey('p'),
		button.Height(2),
	)
	if err != nil {
		return nil, err
	}
	return &buttonSet{btStart, btPause}, nil
}
