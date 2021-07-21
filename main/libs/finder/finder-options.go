package finder

const (
	DefaultLimit = 10
)

type FinderOptions struct {
	Limit int `form:"limit"`
	Skip  int `form:"skip"`
}

func NewFinderOptions() *FinderOptions {
	return &FinderOptions{
		Skip:  0,
		Limit: DefaultLimit,
	}
}
