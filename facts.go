package catfacts

type FactRetriever interface {
	Retrieve() (string, error)
}
