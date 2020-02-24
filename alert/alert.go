package alert

type Alert interface {
    // Message takes a message as a string and sends it
    // to the configured destination
    Message(message string) error
}
