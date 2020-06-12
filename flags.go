package cmd

// Flags is used to define flags with and without arguments. It’s meant to be used through Cmd and
// Group.
type Flags struct {
}

// Flag defines a flag without a value.
func (f *Flags) Flag(spec string, p *bool, usage string) {
}

// String defines a flag with a string value.
func (f *Flags) String(spec string, p *string, name, usage string) {
}

// Int defines a flag with an integer value.
func (f *Flags) Int(spec string, p *int, name, usage string) {
}
