package printer

import "errors"

var (
	ErrHandlerStarted = errors.New("unable to start handler again")
)
