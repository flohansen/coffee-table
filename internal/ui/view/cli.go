package view

import (
	"fmt"

	"github.com/flohansen/coffee-table/internal/ui/viewmodel"
	"github.com/flohansen/coffee-table/pkg/proto"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type CliView struct {
	app          *tview.Application
	userList     *tview.List
	messageList  *tview.TextView
	messageInput *tview.InputField
	viewModel    *viewmodel.MainViewModel
}

func NewCliView() *CliView {
	view := &CliView{
		app:          tview.NewApplication(),
		messageInput: tview.NewInputField(),
		userList:     tview.NewList(),
		viewModel:    viewmodel.NewMainViewModel(),
	}

	loginView := view.NewLoginView()
	chatView := view.NewChatView()

	view.viewModel.CurrentView.Bind(func(value viewmodel.View) {
		switch value {
		case viewmodel.ViewLogin:
			view.app.SetRoot(loginView, true)
		case viewmodel.ViewChat:
			view.app.SetRoot(chatView, true)
		}
	})

	view.viewModel.Message.Bind(func(value string) {
		view.messageInput.SetText(value)
	})

	view.viewModel.Users.Bind(func(value []*proto.User) {
		view.userList.Clear()
		for _, user := range value {
			view.userList.AddItem(user.Username, "", ' ', nil)
		}
	})

	view.viewModel.CurrentMessage.Bind(func(msg string) {
		if msg == "" {
			return
		}

		view.app.QueueUpdateDraw(func() {
			fmt.Fprint(view.messageList, msg)
			view.messageList.ScrollToEnd()
		})
	})

	return view
}

func (a *CliView) Run() error {
	return a.app.Run()
}

func (a *CliView) NewLoginView() tview.Primitive {
	form := tview.NewForm()
	form.SetBorder(true).
		SetTitle("Login")

	form.AddInputField("Username: ", "", 0, nil, a.viewModel.UpdateUsername).
		AddInputField("Server URL: ", "", 0, nil, a.viewModel.UpdateServerURL).
		AddButton("Connect", a.viewModel.Connect)

	return form
}

func (a *CliView) NewChatView() tview.Primitive {
	a.messageInput.SetLabel("Enter message: ").
		SetLabelColor(tcell.ColorWhite).
		SetDoneFunc(func(key tcell.Key) {
			switch key {
			case tcell.KeyEnter:
				a.viewModel.SendMessage(a.messageInput.GetText())
			}
		})

	a.messageList = tview.NewTextView().
		SetDynamicColors(true)
	a.messageList.SetBorder(true).
		SetTitle("Messages")

	a.userList.SetUseStyleTags(false, false)
	a.userList.SetBorder(true).
		SetTitle("Users")

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexColumn).
			AddItem(a.messageList, 0, 1, false).
			AddItem(a.userList, 20, 0, false), 0, 1, false).
		AddItem(a.messageInput, 3, 0, true)

	return flex
}
