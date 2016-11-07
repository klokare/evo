package web

// Client provides composite interface of clients for each type
type Client interface {
	StudyClient
	ExperimentClient
	TrialClient
	IterationClient
}
