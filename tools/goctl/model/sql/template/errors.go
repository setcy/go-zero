package template

// Error defines an error template
const Error = `package {{.pkg}}

import "github.com/hduhelp/go-zero/core/stores/sqlx"

var ErrNotFound = sqlx.ErrNotFound
`
