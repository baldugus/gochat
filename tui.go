package main

import (
	"github.com/baldugus/gochat/internal/broker"

	"github.com/marcusolsson/tui-go"
)

func StartUi(c *broker.Broker) {
	loginView := NewLoginView()
	chatView := NewChatView()

	ui, err := tui.New(loginView)
	if err != nil {
		panic(err)
	}

	quit := func() { ui.Quit() }

	ui.SetKeybinding("Esc", quit)
	ui.SetKeybinding("Ctrl+c", quit)

	loginView.OnLogin(func(username string) {
		c.SetName(username)
		ui.SetWidget(chatView)
	})

	chatView.OnSubmit(func(msg string) {
		c.SendMessage(msg)
	})

	go func() {
		msgs := c.Incoming()
		for {
			msg := <-msgs
			// we need to make the change via ui update to make sure the ui is repaint correctly
			ui.Update(func() {
				chatView.AddMessage(msg)
			})
		}
	}()

	if err := ui.Run(); err != nil {
		panic(err)
	}
}
