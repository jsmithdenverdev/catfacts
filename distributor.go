package catfacts

type Distributor interface {
	Distribute(to, message string) error
}
