package context

type ContextKeys string

const (
	ActionKey    ContextKeys = "action"
	XCallerIDKey ContextKeys = "x-caller-id"
)

func (key ContextKeys) String() string {
	return string(key)
}
