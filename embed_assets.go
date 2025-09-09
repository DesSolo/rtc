/*
	This file is intended to be embedded into the user interface binary file.
*/

package assets

import "embed"

//go:embed frontend/ui/dist/*
var FS embed.FS // nolint:revive
