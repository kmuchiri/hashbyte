package rainbow

import _ "embed"

// EmbeddedTables contains the pre-computed rainbow tables (tables.bin) 
// embedded directly into the binary.
//go:embed tables.bin
var EmbeddedTables []byte
