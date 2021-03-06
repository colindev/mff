package main

import "fmt"

type environments struct {
	path       string `env:"-"`
	Debug      bool   `env:"DEBUG"`
	DSN        string `env:"DSN"`
	AdminAddr  string `env:"ADMIN_ADDR"`
	PublicAddr string `env:"PUBLIC_ADDR"`
	PublicUI   string `env:"PUBLIC_UI"`
}

func (env *environments) String() string {
	return fmt.Sprintf(`
	path: %s
	Debug: %v
	DSN: %s
	AdminAddr: %s
	PublicAddr: %s
	PublicUI: %s
	`,
		env.path,
		env.Debug,
		env.DSN,
		env.AdminAddr,
		env.PublicAddr,
		env.PublicUI)
}
