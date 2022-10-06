package app

import "github.com/pkg/errors"

// ErrNotFound used when an entity is not found.
var ErrNotFound = errors.New("not found")

// ErrChannelDisplayNameInvalid is used when a channel name is too long.
var ErrChannelDisplayNameInvalid = errors.New("channel name is invalid or too long")

// ErrWikiDocNotActive occurs when trying to run a command on a wikiDoc run that has ended.
var ErrWikiDocNotActive = errors.New("already ended")

// ErrWikiDocActive occurs when trying to run a command on a wikiDoc run that is active.
var ErrWikiDocActive = errors.New("already active")

// ErrMalformedWikiDoc occurs when a wikiDoc run is not valid.
var ErrMalformedWikiDoc = errors.New("malformed")

// ErrDuplicateEntry occurs when failing to insert because the entry already existed.
var ErrDuplicateEntry = errors.New("duplicate entry")
