package SimplyP2P

import (
	"github.com/andlabs/ui"
	"time"
)

type SimplyGui struct {
	window *ui.Window
}

func (n *Node) setupUI() {
	gui := new(SimplyGui)
	gui.window = ui.NewWindow("State Peer" , 300, 200, false)

	// Close current Node at exit
	gui.window.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		_ = n.Close()
		return true
	})
	ui.OnShouldQuit(func() bool {
		gui.window.Destroy()
		_ = n.Close()
		return true
	})

	box := ui.NewVerticalBox()
	gui.window.SetChild(box)

    // Current state label that will be updated
	state := ui.NewLabel("")

	n.state.event = func(actualState bool) {
		if actualState {
			state.SetText("Actual: True")
		} else {
			state.SetText("Actual: False")
		}
	}

	// Function that create a new change state button
	addButton := func(text string, newState bool) *ui.Button {
		button := ui.NewButton(text)
		button.OnClicked(func(btn *ui.Button) {
			currentTime := time.Now()
			if n.state.Update(newState, currentTime) {
				n.peers.Broadcast(ChangeState{
					State: newState,
					time:  currentTime,
				})
			}
		})

		return button
	}

	box.Append(state, false)
	box.Append(addButton("False", false), false)
	box.Append(addButton("True", true), false)

	gui.window.Show()
}

func (n *Node) newGui() {
	go ui.Main(n.setupUI)
}