package pomodoro

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// Category contants represent different categories and states
// for a pomodoro interval
const (
	CategoryPomodoro   = "Pomodoro"   // Represents the working interval
	CategoryShortBreak = "ShortBreak" // Represents the short break between working intervals
	CategoryLongBreak  = "LongBreak"  // Represents the break between a set of pomodoros
)

// State constants are integers representations for an interval
// of a given state
const (
	StateNotStarted = iota // Interval has not been started
	StateRunning           // Interval is currently in progress
	StatePaused            // Interval is active, but paused
	StateDone              // Interval has been completed
	StateCancelled         // Interval was planned but cancelled
)

// Interval type represents a pomodoro interval
type Interval struct {
	ID              int64         // A unique ID used to access an interval instance
	StartTime       time.Time     // The time at which the interval is planned to start
	PlannedDuration time.Duration // The planned duration of an interval in minutes
	ActualDuration  time.Duration // The actual completed duration on an interval in minutes
	Category        string        // Name of the interval type
	State           int           // Integer representation of the intervals state
}

// Repository interface allows for abstracted data sources by creating specific signatures
// for the Interval type's methods.
type Repository interface {
	Create(i Interval) (int64, error) // Creates an interval, returns Interval ID or error
	Update(i Interval) error          // Updates an interval, returns error if update fails
	ByID(id int64) (Interval, error)  // Retrieves an interval by ID, returns Interval object or error
	Last() (Interval, error)          // Retrieves the last interval, returns Interval object or error
	Breaks(n int) ([]Interval, error) // Retrieves intervals of type breaks, returns slice of Interval objects or error
}

// Error values for representing particular errors that could be returned
var (
	ErrNoIntervals        = errors.New("error: No intervals")
	ErrIntervalNotRunning = errors.New("error: Interval not running")
	ErrIntervalCompleted  = errors.New("error: Interval is completed or cancelled")
	ErrInvalidState       = errors.New("error: Invalid State")
	ErrInvalidId          = errors.New("error: Invalid ID")
)

// IntervalConfig type represents the configuration requirues to instantiate an Interval type
type IntervalConfig struct {
	repo               Repository    // The data store repository to use
	PomodoroDuration   time.Duration // The length of working intervals in minutes
	ShortBreakDuration time.Duration // The length of short breaks in minutes
	LongBreakDuration  time.Duration // The length of long breaks in minutes
}

// NewConfig function constructs a new IntervalConfig from values provided by user
func NewConfig(repo Repository, pomodoro, shortBreak, longBreak time.Duration) *IntervalConfig {
	// Create an IntervalConfig object from the provided repo and default values
	c := &IntervalConfig{
		repo:               repo,
		PomodoroDuration:   25 * time.Minute,
		ShortBreakDuration: 5 * time.Minute,
		LongBreakDuration:  15 * time.Minute,
	}
	// Check if non-default values were specified, if they were, set the config values accordingly
	if pomodoro > 0 {
		c.PomodoroDuration = pomodoro
	}
	if shortBreak > 0 {
		c.ShortBreakDuration = shortBreak
	}
	if longBreak > 0 {
		c.LongBreakDuration = longBreak
	}
	return c
}

// nextCategory function takes a reference to the repository as input and returns the next interval category string (or error)
func nextCategory(r Repository) (string, error) {
	li, err := r.Last()
	if err != nil && err == ErrNoIntervals {
		// If this is the first interval, assume it is a pomodoro
		return CategoryPomodoro, nil
	}
	if err != nil {
		return "", err
	}
	if li.Category == CategoryLongBreak || li.Category == CategoryShortBreak {
		// If the last interval was a break, the next category will be a pomodoro
		return CategoryPomodoro, nil
	}
	// Get the slice of breaks that have occured
	lastBreaks, err := r.Breaks(3)
	if err != nil {
		return "", err
	}
	if len(lastBreaks) < 3 {
		// If less than 3 shortbreaks have occured, it is not longbreak time yet
		return CategoryShortBreak, nil
	}
	for _, i := range lastBreaks {
		if i.Category == CategoryLongBreak {
			return CategoryShortBreak, nil
		}
	}
	return CategoryLongBreak, nil
}

// Callback type allows callers of the package to pass callback function to execute during the interval
type Callback func(Interval)

// tick function controls the timer for each intervals execution.
// This function takes as input an instance of context that indicates a cancellation,
// the id of the interval to control, and an instance of the configuration, as well as three callbacck functions,
// once to execute at the start, one at the end, and one periodically
// This function returns an error if the tick could not execute successfully
func tick(ctx context.Context, id int64, config *IntervalConfig, start, periodic, end Callback) error {
	// Ticker contains a channel that sends the current time over the channel on each tick
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	// Get the Interval object by ID
	i, err := config.repo.ByID(id)
	if err != nil {
		return err
	}
	// expire triggers after the difference in planned duration and actual duration
	expire := time.After(i.PlannedDuration - i.ActualDuration)
	// Start the interval and loop infinitely until expire is triggered
	start(i)
	for {
		select {
		// Occurs on a change on the ticker channel (occurs after a tick)
		case <-ticker.C:
			// Get the interval again to get changed attributes
			i, err := config.repo.ByID(id)
			if err != nil {
				return err
			}
			if i.State == StatePaused {
				return nil
			}
			// By adding a second to i.ActualDuration, we are effectively counting
			// down on the expire trigger
			i.ActualDuration += time.Second
			// Update the actual duration value of the current interval
			if err := config.repo.Update(i); err != nil {
				return err
			}
			// Execute the periodic callback function
			periodic(i)
		// Occurs when the current interval has counted down
		case <-expire:
			i, err := config.repo.ByID(id)
			if err != nil {
				return err
			}
			// Complete the interval by calling the end callback and updating the object
			i.State = StateDone
			end(i)
			return config.repo.Update(i)
		// Occurs when the current context has been cancelled or timed out
		case <-ctx.Done():
			i, err := config.repo.ByID(id)
			if err != nil {
				return err
			}
			// Update and return the interval object
			i.State = StateCancelled
			return config.repo.Update(i)
		}
	}
}

// newInterval takes an instance of the config and returns an interval instance
func newInterval(config *IntervalConfig) (Interval, error) {
	i := Interval{}
	category, err := nextCategory(config.repo)
	if err != nil {
		return i, err
	}
	i.Category = category
	// Set the duration of the interval depending on it's category
	switch category {
	case CategoryPomodoro:
		i.PlannedDuration = config.PomodoroDuration
	case CategoryShortBreak:
		i.PlannedDuration = config.ShortBreakDuration
	case CategoryLongBreak:
		i.PlannedDuration = config.LongBreakDuration
	}
	if i.ID, err = config.repo.Create(i); err != nil {
		return i, err
	}
	return i, nil
}

// GetInterval takes an instance of intervalconfig and returns an instance of the interval type of an error
// This function attempts to retrieve the last interval from the repo, returning it if it is active or
// returning an error if there is an issue. If the last interval is inactive or unavailable, this function
// returns a new interval using the previously defined function newInterval()
func GetInterval(config *IntervalConfig) (Interval, error) {
	i := Interval{}
	var err error
	// Get the last interval from the given interal
	i, err = config.repo.Last()
	if err != nil && err != ErrNoIntervals {
		return i, err
	}
	if err == nil && i.State != StateCancelled && i.State != StateDone {
		return i, nil
	}
	return newInterval(config)
}

// Start method is used by callers to start the interval timer.
// This function checks the state of the current interval and sets the options accordingly,
// finally calling the tick function to time the interval
func (i Interval) Start(ctx context.Context, config *IntervalConfig, start, periodic, end Callback) error {
	switch i.State {
	case StateRunning:
		// If the state is already running, do nothing
		return nil
	case StateNotStarted:
		// If the state hasn't started, start it
		i.StartTime = time.Now()
		// Fallthrough to the StatePaused case
		fallthrough
	case StatePaused:
		// If the state is paused, resume it
		i.State = StateRunning
		if err := config.repo.Update(i); err != nil {
			return err
		}
		return tick(ctx, i.ID, config, start, periodic, end)
	case StateCancelled, StateDone:
		return fmt.Errorf("%w: Cannot start", ErrIntervalCompleted)
	default:
		return fmt.Errorf("%w: %d", ErrInvalidState, i.State)
	}
}

// Pause method is used by callers to pause a running interval
func (i Interval) Pause(config *IntervalConfig) error {
	if i.State != StateRunning {
		// Cannot pause an interval that is not running
		return ErrIntervalNotRunning
	}
	i.State = StatePaused
	return config.repo.Update(i)
}
