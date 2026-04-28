package example // want `bad_interface\.go: missing //go:generate mockgen directive`

type BadService interface {
	Do() error
}
