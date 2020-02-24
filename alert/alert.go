package alert

type Alert interface {
    Message(message string) error
}
