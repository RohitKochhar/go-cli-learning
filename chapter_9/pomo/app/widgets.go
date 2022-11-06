package app

import (
	"context"

	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/widgets/donut"
	"github.com/mum4k/termdash/widgets/segmentdisplay"
	"github.com/mum4k/termdash/widgets/text"
)

// widgets represents a collection of widgets pointing to the four main status widgets in our application
type widgets struct {
	donTimer       *donut.Donut                   // `Donut` widget in the Timer section
	disType        *segmentdisplay.SegmentDisplay // Display type for the SegmentDisplay widget in the type section
	txtInfo        *text.Text                     // `Text` widget in the Info section
	txtTimer       *text.Text                     // `Text` widget in the Timer section
	updateDonTimer chan []int                     // Channel to update donut timer with new values
	updateTxtInfo  chan string                    // Channel to update the info text section
	updateTxtTimer chan string                    // Channel to update the timer text section
	updateTxtType  chan string                    // Channel to update the display type
}

// update method updates widgets with new data
// function takes update values for all widgets as well as a redrawCh var to specify whether the screen needs to be redrawn
func (w *widgets) update(timer []int, txtType, txtInfo, txtTimer string, redrawCh chan<- bool) {
	if txtInfo != "" {
		w.updateTxtInfo <- txtInfo
	}
	if txtType != "" {
		w.updateTxtType <- txtType
	}
	if txtTimer != "" {
		w.updateTxtTimer <- txtTimer
	}
	if len(timer) > 0 {
		w.updateDonTimer <- timer
	}
	redrawCh <- true
}

// newWidgets method initializes the widget type
// This function calls other functions to instantiate each widget
func newWidgets(ctx context.Context, errorCh chan<- error) (*widgets, error) {
	w := &widgets{}
	var err error
	// Define channels for each update channel
	w.updateDonTimer = make(chan []int)
	w.updateTxtType = make(chan string)
	w.updateTxtInfo = make(chan string)
	w.updateTxtTimer = make(chan string)
	w.donTimer, err = newDonut(ctx, w.updateDonTimer, errorCh)
	if err != nil {
		return nil, err
	}
	w.disType, err = newSegmentDisplay(ctx, w.updateTxtType, errorCh)
	if err != nil {
		return nil, err
	}
	w.txtInfo, err = newText(ctx, w.updateTxtInfo, errorCh)
	if err != nil {
		return nil, err
	}
	w.txtTimer, err = newText(ctx, w.updateTxtTimer, errorCh)
	if err != nil {
		return nil, err
	}
	return w, nil
}

// newText initializes a new text widget
func newText(ctx context.Context, updateText <-chan string, errorCh chan<- error) (*text.Text, error) {
	// Instantiate a new text widget and check for errors
	txt, err := text.New()
	if err != nil {
		return nil, err
	}
	// Use an anonymous function to allow txt widget to be updated concurrently from outside this function
	go func() {
		for {
			select {
			// If we were given text to update, reset the current value and update
			case t := <-updateText:
				txt.Reset()
				errorCh <- txt.Write(t)
			// If the context is complete, return to close the widget
			case <-ctx.Done():
				return
			}
		}
	}()
	return txt, nil
}

// newText initializes a new donut widget
func newDonut(ctx context.Context, donUpdater <-chan []int, errorCh chan<- error) (*donut.Donut, error) {
	// Instantiate a new donut widget and check for errors
	don, err := donut.New(
		donut.Clockwise(),
		donut.CellOpts(cell.FgColor(cell.ColorBlue)),
	)
	if err != nil {
		return nil, err
	}
	// Use an anonymous goroutine to instantiate the widget and keep it open to be accessed later
	go func() {
		for {
			select {
			case d := <-donUpdater:
				if d[0] <= d[1] {
					errorCh <- don.Absolute(d[0], d[1])
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	return don, nil
}

// newSegmentDisplay initializes a segmentDisplay widget
func newSegmentDisplay(ctx context.Context, updateText <-chan string, errorCh chan<- error) (*segmentdisplay.SegmentDisplay, error) {
	sd, err := segmentdisplay.New()
	if err != nil {
		return nil, err
	}
	// Use an anonymous goroutine to instantiate the widget and keep it open to be accessed later
	go func() {
		for {
			select {
			case t := <-updateText:
				if t == "" {
					t = " "
				}
				errorCh <- sd.Write([]*segmentdisplay.TextChunk{
					segmentdisplay.NewChunk(t),
				})
			case <-ctx.Done():
				return
			}
		}
	}()
	return sd, nil
}
