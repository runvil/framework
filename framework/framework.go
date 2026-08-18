// Package framework is the entry point for the Runvil meta-framework.
//
// Rather than being built as a single, self-contained ecosystem, Runvil
// composes modules sourced across multiple ecosystems and repositories,
// orchestrating them into one cohesive, high-level framework.
package framework

import "fmt"

const (
	// Name is the canonical display name of the Runvil framework.
	Name = "Runvil Framework"
	// Version is the semantic version of the framework.
	Version = "0.3.0"
)

// Banner returns the framework banner identifying the current framework
// version.
func Banner() string {
	return fmt.Sprintf("%s v%s", Name, Version)
}
