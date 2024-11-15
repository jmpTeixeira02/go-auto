package notifier

import (
	"log"
	"os"
)

type terminalNotifier struct {
	logger *log.Logger
}

func newTerminalNotifier() Notifier {
	return terminalNotifier{logger: log.New(os.Stdout, "", 0)}
}

func (t terminalNotifier) SendMessage(msg string) error {
	t.logger.Println(msg)
	return nil
}
