package ports

type RootController interface {
	BeforeStart() bool
}

type VersionController interface {
	PrintVersion()
}

type MonitorController interface {
	Static() error
	Dynamic() error
}
