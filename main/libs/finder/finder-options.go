package finder

const (
	DefaultLimit = 10
)

type FinderOptions struct {
	Limit int `json:"limit" form:"limit"`
	Skip  int `json:"skip" form:"skip"`
}

func NewFinderOptions() *FinderOptions {
	return &FinderOptions{
		Skip:  0,
		Limit: DefaultLimit,
	}
}
