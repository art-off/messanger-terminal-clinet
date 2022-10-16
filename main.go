package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"os"
	"terminal-client/message"
	"terminal-client/socket_manager"
)

var messages []*mes.TextMessage

type TableData struct {
	tview.TableContentReadOnly
}

func (d *TableData) GetCell(row, column int) *tview.TableCell {
	m := messages[row]
	return tview.NewTableCell(fmt.Sprintf("[red]%s: [green]%s", m.Sender, m.Text))
}

func (d *TableData) GetRowCount() int {
	return len(messages)
}

func (d *TableData) GetColumnCount() int {
	return 1
}

// Tview
var app = tview.NewApplication()
var flex = tview.NewFlex()
var messageTable = tview.NewTable()
var input = tview.NewInputField()

var sm *socket_manager.SocketManager

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Please provide username and room")
		return
	}
	username := os.Args[1]
	room := os.Args[2]

	data := &TableData{}
	messageTable.
		SetBorders(false).
		SetSelectable(false, false).
		SetContent(data)

	flex.SetDirection(tview.FlexRow).
		AddItem(messageTable, 0, 100, false).
		AddItem(input, 0, 1, false)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch key := event.Key(); key {
		case tcell.KeyEnter:
			enteredText := input.GetText()
			if enteredText != "" {
				sm.SendMessage(enteredText)
				input.SetText("")
			}
		}
		return event
	})

	input.Focus(func(p tview.Primitive) {})

	sm = &socket_manager.SocketManager{
		Host:             "localhost:8080",
		Path:             "/",
		OnNewTextMessage: displayNewMessage,
	}

	go sm.ListenAndRegisterUser(username, room)

	if err := app.SetRoot(flex, true).EnableMouse(false).Run(); err != nil {
		panic(err)
	}
}

func displayNewMessage(m mes.TextMessage) {
	app.QueueUpdateDraw(func() {
		messages = append(messages, &m)
		messageTable.ScrollToEnd()
	})
}
