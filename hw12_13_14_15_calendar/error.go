package calendar

import "errors"

// ErrNotFound не найдено.
var ErrNotFound = errors.New("not found")

// ErrDateBusy данное время уже занято другим событием.
var ErrDateBusy = errors.New("that date is busy")
