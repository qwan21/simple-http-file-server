package app

import "flag"

// flags defines available command line flags.
type flags struct {
	service struct {
		service   string
		host      string
		port      string
		directory string
		route     string
	}
}

func (a *App) parseFlags(args []string) (*flags, error) {
	f := flag.NewFlagSet("main", flag.ContinueOnError)
	var flags flags

	f.StringVar(&flags.service.service, "service", "http-file-server", "")
	f.StringVar(&flags.service.host, "host", "127.0.0.1", "http-file-server host")
	f.StringVar(&flags.service.port, "port", "8080", "http-file-server port")
	f.StringVar(&flags.service.directory, "directory", ".", "http-file-server root directory")
	f.StringVar(&flags.service.route, "route", "/store/", "http-file-server route")

	err := f.Parse(args)

	if err != nil {
		return nil, err
	}
	a.flagSet = f

	return &flags, err
}
