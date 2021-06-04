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

    // State
	state := ui.NewLabel("")

	n.state.event = func(actualState bool) {
		if actualState {
			state.SetText("Actual: True")
		} else {
			state.SetText("Actual: False")
		}
	}

	falseBtn := ui.NewButton("False")
	falseBtn.OnClicked(func(btn *ui.Button) {
		currentTime := time.Now()
		n.state.Update(false, currentTime)
		n.peers.Broadcast(ChangeState{
			State: false,
			time:  currentTime,
		})
	})

	trueBtn := ui.NewButton("True")
	trueBtn.OnClicked(func(btn *ui.Button) {
		currentTime := time.Now()
		n.state.Update(true, currentTime)
		n.peers.Broadcast(ChangeState{
			State: true,
			time:  currentTime,
		})
	})

	box.Append(state, false)
	box.Append(falseBtn, false)
	box.Append(trueBtn, false)

	gui.window.Show()
}

func (n *Node) newGui() {
	go ui.Main(n.setupUI)
}