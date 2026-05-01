package postgres

import domainerrors "github.com/casatorino/backend/internal/domain/errors"

func domainrepositoriesErrNotFound() error {
	return domainerrors.ErrNotFound
}
