package view

import (
	"fmt"

	"github.com/flohansen/coffee-table/internal/ui/viewmodel"
	"github.com/flohansen/coffee-table/pkg/proto"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type CliView struct {
	app       *tview.Application
	viewModel *viewmodel.MainViewModel
}

func NewCliView() *CliView {
	view := &CliView{
		app:       tview.NewApplication(),
		viewModel: viewmodel.NewMainViewModel(),
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

	return view
}

func (a *CliView) Run() error {
	return a.app.Run()
}

func (a *CliView) NewLoginView() tview.Primitive {
	errorView := tview.NewTextView()
	a.viewModel.Error.Bind(func(err string) {
		if err == "" {
			errorView.Clear()
			return
		}

		errorView.SetText(err)
	})

	form := tview.NewForm().
		AddInputField("Username: ", "", 0, nil, a.viewModel.UpdateUsername).
		AddInputField("Server URL: ", "", 0, nil, a.viewModel.UpdateServerURL).
		AddCheckbox("Secure", false, a.viewModel.Secure.Set).
		AddButton("Connect", a.viewModel.Connect)

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(errorView, 2, 3, false).
		AddItem(form, 0, 1, true)
	flex.SetBorder(true).SetTitle("Login")

	return flex
}

func (a *CliView) NewChatView() tview.Primitive {
	messageInput := tview.NewInputField()
	messageInput.SetLabel("Enter message: ").
		SetLabelColor(tcell.ColorWhite).
		SetDoneFunc(func(key tcell.Key) {
			switch key {
			case tcell.KeyEnter:
				a.viewModel.SendMessage(messageInput.GetText())
			}
		})
	a.viewModel.Message.Bind(func(value string) {
		messageInput.SetText(value)
	})

	messageList := tview.NewTextView().
		SetDynamicColors(true)
	messageList.SetBorder(true).
		SetTitle("Messages")
	a.viewModel.CurrentMessage.Bind(func(msg string) {
		if msg == "" {
			return
		}

		a.app.QueueUpdateDraw(func() {
			fmt.Fprint(messageList, msg)
			messageList.ScrollToEnd()
		})
	})

	userList := tview.NewList().SetUseStyleTags(false, false)
	userList.SetBorder(true).
		SetTitle("Users")
	a.viewModel.Users.Bind(func(value []*proto.User) {
		userList.Clear()
		for _, user := range value {
			userList.AddItem(user.Username, "", ' ', nil)
		}
	})

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexColumn).
			AddItem(messageList, 0, 1, false).
			AddItem(userList, 20, 0, false), 0, 1, false).
		AddItem(messageInput, 3, 0, true)

	return flex
}
