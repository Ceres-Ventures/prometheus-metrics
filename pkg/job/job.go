package job

type (
	Job struct {
		ID            uint64 `json:"id"`
		IsRepeatable  bool
		RepeatCount   int
		RepeatCounter int
		F             func() error
	}
	JobExecutable func() error
)
