package main

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/streadway/amqp"
)

const exchangeName = "chat_exchange"

type ChatGUI struct {
	app            fyne.App
	mainWindow     fyne.Window
	conn           *amqp.Connection
	ch             *amqp.Channel
	currentQueue   string
	currentChannel string
	consumerTag    string

	messages     *widget.Entry
	input        *widget.Entry
	channelEntry *widget.Entry
	status       *widget.Label
}

func main() {
	chatApp := &ChatGUI{app: app.New()}
	chatApp.showLoginWindow()
	chatApp.app.Run()
}

func (c *ChatGUI) showLoginWindow() {
	loginWindow := c.app.NewWindow("Connect to Chat")
	loginWindow.Resize(fyne.NewSize(400, 200))

	serverEntry := widget.NewEntry()
	serverEntry.SetPlaceHolder("127.0.0.1:5672")
	channelEntry := widget.NewEntry()
	channelEntry.SetPlaceHolder("general")

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Server Address", Widget: serverEntry},
			{Text: "Initial Channel", Widget: channelEntry},
		},
		OnSubmit: func() {
			loginWindow.Hide()
			c.createMainWindow()
		},
	}

	loginWindow.SetContent(container.NewVBox(
		widget.NewLabel("Enter connection details:"),
		form,
	))
	loginWindow.Show()
}

func (c *ChatGUI) sendMessage() {

}
func (c *ChatGUI) createMainWindow() {
	c.mainWindow = c.app.NewWindow("RabbitMQ Chat")
	c.mainWindow.Resize(fyne.NewSize(600, 400))

	c.messages = widget.NewMultiLineEntry()
	c.messages.Disable()

	c.input = widget.NewEntry()
	c.input.SetPlaceHolder("Type message here...")
	c.input.OnSubmitted = func(_ string) {
		c.sendMessage()
		c.input.SetText("")
	}

	c.channelEntry = widget.NewEntry()
	c.channelEntry.SetPlaceHolder("Type channel name here...")

	switchBtn := widget.NewButton("Switch Channel", func() {
		if strings.TrimSpace(c.channelEntry.Text) == "" {
			c.status.SetText("Channel name cannot be empty!")
			return
		}
		c.switchChannel(c.channelEntry.Text)
	})

	sendBtn := widget.NewButton("Send", func() {
		if c.input.Text != "" {
			c.sendMessage()
			c.input.SetText("")
		}
	})

	c.status = widget.NewLabel("Disconnected")

	topBar := container.NewBorder(
		nil,
		nil,
		widget.NewLabel("Current Channel:"),
		container.NewHBox(switchBtn, c.status),
		c.channelEntry,
	)

	bottom := container.NewBorder(
		nil,
		nil,
		nil,
		sendBtn,
		c.input,
	)

	content := container.NewBorder(
		topBar,
		bottom,
		nil,
		nil,
		container.NewScroll(c.messages),
	)

	c.mainWindow.SetContent(content)
	c.mainWindow.SetCloseIntercept(func() {
		c.disconnect()
		c.mainWindow.Close()
	})

	c.mainWindow.Show()
}

func (c *ChatGUI) switchChannel(text string) {

}

func (c *ChatGUI) disconnect() {

}
