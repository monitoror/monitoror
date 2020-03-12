package stripe

type (
	Repository interface {
		GetCount(afterTimestamp string) int
	}
)
