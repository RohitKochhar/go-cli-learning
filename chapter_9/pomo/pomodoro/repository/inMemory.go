package repository

import (
	"fmt"
	"rohitsingh/pomo/pomodoro"
	"sync"
)

// inMemoryRepo type represents the in-memory repository and implements the required methods
// needed for use in the Pomodoro app
type inMemoryRepo struct {
	sync.RWMutex                     // Allows methods to be accessed directly from this type without concurrent accesses
	intervals    []pomodoro.Interval // Stores Interval objects in memory
}

// NewInMemoryRepo instantiates a new inMemoryRepo type with an empty slice of intervals
func NewInMemoryRepo() *inMemoryRepo {
	return &inMemoryRepo{
		intervals: []pomodoro.Interval{},
	}
}

// Create method takes an instance of interval and saves its values to the data store
func (r *inMemoryRepo) Create(i pomodoro.Interval) (int64, error) {
	// Lock the store to prevent concurrent writes
	r.Lock()
	defer r.Unlock()
	// Place the new object at the end of the memory array and get it's ID
	i.ID = int64(len(r.intervals) + 1)
	r.intervals = append(r.intervals, i)
	return i.ID, nil
}

// Update method updates the values of an existing data store entry
func (r *inMemoryRepo) Update(i pomodoro.Interval) error {
	// Lock the store to prevent concurrent writes
	r.Lock()
	defer r.Unlock()
	// Set the value at the intervals location to the newly upated interval
	if i.ID == 0 {
		return fmt.Errorf("%w: %d", pomodoro.ErrInvalidId, i.ID)
	}
	r.intervals[i.ID-1] = i
	return nil
}

// ByID retrieves and returns an item from it's ID
func (r *inMemoryRepo) ByID(id int64) (pomodoro.Interval, error) {
	// Lock the store to prevent concurrent reads
	r.RLock()
	defer r.RUnlock()
	i := pomodoro.Interval{}
	if id == 0 {
		return i, fmt.Errorf("%w: %d", pomodoro.ErrInvalidId, i.ID)
	}
	i = r.intervals[id-1]
	return i, nil
}

// Last retrieves and returns the last interval from the data store
func (r *inMemoryRepo) Last() (pomodoro.Interval, error) {
	// Lock the store to prevent concurrent reads
	r.RLock()
	defer r.RUnlock()
	i := pomodoro.Interval{}
	if len(r.intervals) == 0 {
		return i, pomodoro.ErrNoIntervals
	}
	return r.intervals[len(r.intervals)-1], nil
}

// Breaks method retrieves a given number n of intervals of category break
func (r *inMemoryRepo) Breaks(n int) ([]pomodoro.Interval, error) {
	// Lock the store to prevent concurrent reads
	r.RLock()
	defer r.RUnlock()
	data := []pomodoro.Interval{}
	for k := len(r.intervals) - 1; k >= 0; k-- {
		if r.intervals[k].Category == pomodoro.CategoryPomodoro {
			// Ignore working intervals
			continue
		}
		data = append(data, r.intervals[k])
		if len(data) == n {
			return data, nil
		}
	}
	return data, nil
}
