package internal

type Status int

const (
	StatusQueued     = 1000
	StatusWorking    = 10000
	StatusSuccess    = 20000
	StatusFailed     = 50000
	StatusCancelling = 58000
	StatusCancelled  = 60000
)

func (s Status) String() string {
	switch s {
	case StatusQueued:
		return "queued"
	case StatusWorking:
		return "working"
	case StatusSuccess:
		return "success"
	case StatusFailed:
		return "failed"
	case StatusCancelling:
		return "cancelling"
	case StatusCancelled:
		return "cancelled"
	default:
		panic("internal error: entered unreachable code")
	}
}
