package payment

import "strings"

// registry holds all registered PaymentProvider instances keyed by canonical name.
var registry = map[string]PaymentProvider{}

// Register adds a PaymentProvider to the global registry.
// It is typically called during ServiceContext initialisation.
func Register(p PaymentProvider) {
	if p != nil {
		registry[strings.ToLower(p.Name())] = p
	}
}

// Get returns the PaymentProvider for the given name (case-insensitive).
func Get(name string) (PaymentProvider, bool) {
	p, ok := registry[strings.ToLower(name)]
	return p, ok
}

// List returns the names of all registered providers.
func List() []string {
	names := make([]string, 0, len(registry))
	for k := range registry {
		names = append(names, k)
	}
	return names
}
