package app

import (
	"github.com/mum4k/termdash/align"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/terminalapi"
)

// newGrid defines a new grid layout
// by taking pointers to a buttonSet and widgets, as well as an instance of Terminal
// it returns a pointer to a container or an error if one has occured
func newGrid(b *buttonSet, w *widgets, t terminalapi.Terminal) (*container.Container, error) {
	// Termdash uses a grid.Builder type to build grid layouts
	builder := grid.New()
	// Add the first row by using the builder's add method
	builder.Add(
		grid.RowHeightPerc(30,
			grid.ColWidthPercWithOpts(30,
				[]container.Option{
					container.Border(linestyle.Light),
					container.BorderTitle("Press Q to Quit"),
				},
				grid.RowHeightPerc(80,
					grid.Widget(w.donTimer),
				),
				grid.RowHeightPercWithOpts(20,
					[]container.Option{
						container.AlignHorizontal(align.HorizontalCenter),
					},
					grid.Widget(w.txtTimer,
						container.AlignHorizontal(align.HorizontalCenter),
						container.AlignVertical(align.VerticalMiddle),
						container.PaddingLeftPercent(49),
					),
				),
			),
			grid.ColWidthPerc(70,
				grid.RowHeightPerc(80,
					grid.Widget(w.disType, container.Border(linestyle.Light)),
				),
				grid.RowHeightPerc(20,
					grid.Widget(w.txtInfo, container.Border(linestyle.Light)),
				),
			),
		),
	)
	// Add the second row
	builder.Add(
		grid.RowHeightPerc(10,
			grid.ColWidthPerc(50,
				grid.Widget(b.btStart),
			),
			grid.ColWidthPerc(50,
				grid.Widget(b.btPause),
			),
		),
	)
	// Add a placeholder for the third line using the remaining 60% of the screen space
	builder.Add(
		grid.RowHeightPerc(60),
	)
	// Use the builder.Build method to build the layout and create the container options required to instantiate a container
	gridOpts, err := builder.Build()
	if err != nil {
		return nil, err
	}
	// Use the generated container options to instantiate the container using the method container.New()
	c, err := container.New(t, gridOpts...)
	if err != nil {
		return nil, err
	}
	// Return the newly created container and a nil error
	return c, nil
}
