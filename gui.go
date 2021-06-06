package SimplyP2P

import (
	"github.com/andlabs/ui"
	"time"
)

// SimplyGui is a simple graphic interface that show
// the current state and some buttons to change its value.
type SimplyGui struct {
	window *ui.Window
}

// setupUI sets the gui elements and show then.
func (n *Node) setupUI() {
	gui := new(SimplyGui)
	gui.window = ui.NewWindow("State Peer", 300, 200, false)

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
	state := ui.NewLabel("")

	// Set state change handler (change GUI label)
	n.state.SetEvent(func(actualState bool) {
		if actualState {
			state.SetText("Actual: True")
		} else {
			state.SetText("Actual: False")
		}
	})

	// Function that create a new change state button
	addChangeStateButton := func(text string, newState bool) *ui.Button {
		button := ui.NewButton(text)
		button.OnClicked(func(*ui.Button) {
			currentTime := time.Now()
			if n.state.Update(newState, currentTime) {
				n.peers.Broadcast(ChangeState{
					State: newState,
					Time:  currentTime,
				})
			}
		})
		return button
	}

	// Add graphic elements to box
	box.Append(state, false)
	box.Append(addChangeStateButton("False", false), false)
	box.Append(addChangeStateButton("True", true), false)

	// Add box to the window and show the window
	gui.window.SetChild(box)
	gui.window.Show()
}

// newGui creates a new GUI linked to current node
// and starts it in a new goroutine.
func (n *Node) newGui() {
	go ui.Main(n.setupUI)
}
