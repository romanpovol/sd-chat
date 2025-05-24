package main

import (
	"flag"
	"log"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/romanpovol/sd-chat/internal/rabbitmq"
)

type ChatGUI struct {
	app        fyne.App
	mainWindow fyne.Window
	client     *rabbitmq.ChatClient

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

	serverAddr := widget.NewEntry()
	serverAddr.SetPlaceHolder("Write server addr")
	initChannel := widget.NewEntry()
	initChannel.SetPlaceHolder("Write init channel")

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Server Address", Widget: serverAddr},
			{Text: "Initial Channel", Widget: initChannel},
		},
		OnSubmit: func() {
			loginWindow.Hide()

			serverAddr := serverAddr.Text
			initChannel := initChannel.Text

			client := rabbitmq.NewClient()
			defer client.Close()

			err := client.Connect(serverAddr)
			if err != nil {
				log.Fatalf("Failed to connect: %v", err)
			}
			client.SwitchChannel(initChannel)

			c.client = client

			c.createMainWindow()

			c.mainWindow.Show()
		},
	}

	loginWindow.SetContent(container.NewVBox(
		widget.NewLabel("Enter connection details:"),
		form,
	))
	loginWindow.Show()
}

func (c *ChatGUI) createMainWindow() {
	c.mainWindow = c.app.NewWindow("RabbitMQ Chat")
	c.mainWindow.Resize(fyne.NewSize(600, 400))

	c.messages = widget.NewMultiLineEntry()
	c.messages.Disable()

	c.input = widget.NewEntry()
	c.input.SetPlaceHolder("Type message here...")
	c.input.OnSubmitted = func(_ string) {
		c.client.SendMessage(c.input.Text)
		c.input.SetText("")
	}

	c.channelEntry = widget.NewEntry()
	c.channelEntry.SetPlaceHolder("Type channel name here...")

	switchBtn := widget.NewButton("Switch Channel", func() {
		if strings.TrimSpace(c.channelEntry.Text) == "" {
			c.status.SetText("Channel name cannot be empty!")
			return
		}
		c.client.SwitchChannel(c.channelEntry.Text)
	})

	sendBtn := widget.NewButton("Send", func() {
		log.Println("SUBMIT " + c.input.Text)
		if c.input.Text != "" {
			c.client.SendMessage(c.input.Text)
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
		c.client.Close()
		c.mainWindow.Close()
	})
}
