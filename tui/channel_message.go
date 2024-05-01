package tui

import (
	"git.qowevisa.me/Qowevisa/gotell/communication"
	"git.qowevisa.me/Qowevisa/gotell/errors"
)

// Basically every X_channel.go file launches some sort of channel

func (t *TUI) launchMessageChannel() error {
	if t.messageChannel == nil {
		return errors.WrapErr("t.messageChannel", errors.NOT_INIT)
	}
	go func() {
		for message := range t.messageChannel {
			t.createNotification(string(message), "Message!")
			msg, err := communication.Decode(message)
			t.errorsChannel <- err
			if err != nil {
				continue
			}
			switch msg.Type {
			case communication.SERVER_COMMAND:
				t.handleServerCommands(msg.Data)
			}
		}
	}()
	return nil
}

func (t *TUI) handleServerCommands(data []byte) {
	if len(data) == 1 {
		if data[0] == communication.NICKNAME {
			t.SendMessageToServer("Nickname", 16)
		}
	}
}
