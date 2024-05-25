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
			msg, err := communication.Decode(message)
			t.errorsChannel <- err
			if err != nil {
				continue
			}
			switch msg.From {
			case communication.FROM_SERVER:
				t.handleServerCommands(*msg)
			case communication.FROM_CLIENT:
				t.createNotification(string(msg.Data), "Server Message!")
			}
		}
	}()
	return nil
}

func (t *TUI) handleServerCommands(data communication.Message) {
	if data.About == communication.ABOUT_NICKNAME {
		t.SendMessageToServer("Nickname", 16)
	}
}
