package finder

const (
	DefaultLimit = 10
)

type FinderOptions struct {
	Limit int `json:"limit"`
	Skip  int `json:"skip"`
}

func NewFinderOptions() *FinderOptions {
	return &FinderOptions{
		Skip:  0,
		Limit: DefaultLimit,
	}
}
