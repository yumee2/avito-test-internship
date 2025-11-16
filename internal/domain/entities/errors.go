package entities

import "errors"

var (
	ErrNotFound     = errors.New("not found error")
	ErrDuplicate    = errors.New("duplicate error")
	ErrPrMergerd    = errors.New("pr is already merged")
	ErrNoCandidates = errors.New("no canditates for this pr")
	ErrNotAssigned  = errors.New("reviewer is not assigned to this PR")
)
