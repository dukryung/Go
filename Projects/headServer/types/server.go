package types

type Server interface {
	Run()
	Close()
}
