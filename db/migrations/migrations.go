package migrations

import "embed"

// Files embeds Goose migrations for runtime execution.
//go:embed *.sql
var Files embed.FS

