package app

import (
	"flag"
	"fmt"
	"test/config"
	server "test/internal"
)

// App defines the application setting
type App struct {
	flagSet *flag.FlagSet
	cfg     *mergedConfig
}

// New instantiates a new App instance. ui must not be a nil.
func New() *App {
	return &App{}
}

// Run the application
func (a *App) Run(args []string) int {
	err := a.run(args)
	if err == nil {
		return 0
	}
	fmt.Println(err)
	return 1
}

func (a *App) run(args []string) error {
	_, err := a.parseFlags(args)

	if err != nil {
		return err
	}

	cfg, err := mergeConfig(a.flagSet)

	if err != nil {
		return err
	}

	a.cfg = cfg

	srv := server.New(a.cfg.Config)
	if err := srv.Start(); err != nil {
		return err
	}

	return nil
}

type mergedConfig struct {
	*config.Config
}

func mergeConfig(fs *flag.FlagSet) (*mergedConfig, error) {
	cfg, err := config.Get(fs)

	if err != nil {
		return nil, err
	}
	return &mergedConfig{
		Config: cfg,
	}, nil
}
