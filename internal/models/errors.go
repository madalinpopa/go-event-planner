package models

import "errors"

// ErrNoRecord indicates that no matching record was found in the database or data source.
var ErrNoRecord = errors.New("models: no matching record found")
