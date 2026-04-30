// Package audit provides structured audit logging for vaultdiff operations.
//
// It records every diff comparison performed — including the secret path,
// version range, number of changes, and the acting user — to either a file
// or any io.Writer in JSON or human-readable text format.
//
// Basic usage:
//
//	logger := audit.NewLogger(os.Stdout, "json")
//	err := logger.Record("secret/myapp/config", 2, 3, changes, os.Getenv("USER"))
//
// For persistent file-backed logging:
//
//	fl, err := audit.NewFileLogger("/var/log/vaultdiff/audit.log", "json")
//	if err != nil { ... }
//	defer fl.Close()
//	err = fl.Record("secret/myapp/config", 2, 3, changes, "alice")
package audit
