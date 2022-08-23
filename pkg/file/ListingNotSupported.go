package file

import "developer.zopsmart.com/go/gofr/pkg/errors"

const ErrListingNotSupported = errors.Error(`Listing not supported for provided file store.` +
	` Please set a valid value of FILE_STORE:{LOCAL or SFTP}`)
