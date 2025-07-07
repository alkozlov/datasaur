package builtin

import (
	"block-flow/internal/blocks"
)

// RegisterBuiltinBlocks registers all built-in blocks with the registry
func RegisterBuiltinBlocks(registry *blocks.Registry) {
	// Input blocks
	registry.Register(&InjectBlockFactory{})

	// Output blocks
	registry.Register(&DebugBlockFactory{})

	// Math blocks
	registry.Register(&AdditionBlockFactory{})
	registry.Register(&SubtractionBlockFactory{})
	registry.Register(&MultiplicationBlockFactory{})
	registry.Register(&DivisionBlockFactory{})
}
