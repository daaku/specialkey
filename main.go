// specialkey shows a translucent popup grid of special characters. Picking one
// copies it to the clipboard (via wl-copy) and types it into the previously
// focused window (via wtype).
package main

import (
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

const columns = 6

type key struct {
	char  string
	label string
}

var keys = []key{
	{"²", "squared"},
	{"³", "cubed"},
	{"¹", "sup 1"},
	{"½", "half"},
	{"¼", "quarter"},
	{"¾", "3 quarters"},
	{"°", "degree"},
	{"±", "plus-minus"},
	{"×", "times"},
	{"÷", "divide"},
	{"≠", "not equal"},
	{"≈", "approx"},
	{"≤", "less-eq"},
	{"≥", "greater-eq"},
	{"√", "sqrt"},
	{"∞", "infinity"},
	{"π", "pi"},
	{"µ", "micro"},
	{"→", "right arrow"},
	{"←", "left arrow"},
	{"€", "euro"},
	{"£", "pound"},
	{"©", "copyright"},
	{"—", "em dash"},
}

const css = `
window.specialkey {
	background-color: rgba(0, 0, 0, 0.9);
	border-radius: 16px;
}
.specialkey grid {
	padding: 0 20px 20px 20px;
}
.specialkey button {
	background: rgba(255, 255, 255, 0.06);
	border: 2px solid transparent;
	border-radius: 10px;
	box-shadow: none;
	color: white;
	min-width: 96px;
	min-height: 96px;
	padding: 4px;
}
.specialkey button:hover {
	background: rgba(255, 255, 255, 0.14);
}
.specialkey button:focus,
.specialkey button:focus-visible {
	background: rgba(255, 255, 255, 0.18);
	border-color: rgba(120, 170, 255, 0.9);
	outline: none;
}
.specialkey button.close {
	background: transparent;
	color: rgba(255, 255, 255, 0.6);
	min-width: 28px;
	min-height: 28px;
	padding: 0;
	border-radius: 14px;
}
.specialkey button.close:hover {
	background: rgba(255, 255, 255, 0.14);
	color: white;
}
.specialkey .char {
	font-size: 40px;
}
.specialkey .name {
	font-size: 10px;
	color: rgba(255, 255, 255, 0.55);
}
`

func main() {
	app := gtk.NewApplication("org.daaku.specialkey", gio.ApplicationNonUnique)

	var picked string
	app.ConnectActivate(func() {
		activate(app, func(s string) { picked = s })
	})

	if code := app.Run(os.Args); code != 0 {
		os.Exit(code)
	}
	if picked == "" {
		return
	}

	// wl-copy forks a background process that keeps serving the selection
	// after we exit, which GTK's own clipboard would not.
	if err := exec.Command("wl-copy", "--", picked).Run(); err != nil {
		log.Printf("wl-copy: %v", err)
	}

	// Our window is gone by now; give the compositor a moment to hand focus
	// back to the previous window before typing into it.
	time.Sleep(150 * time.Millisecond)
	if err := exec.Command("wtype", "--", picked).Run(); err != nil {
		log.Printf("wtype: %v", err)
	}
}

func activate(app *gtk.Application, onPick func(string)) {
	provider := gtk.NewCSSProvider()
	provider.LoadFromString(css)
	gtk.StyleContextAddProviderForDisplay(
		gdk.DisplayGetDefault(), provider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)

	win := gtk.NewApplicationWindow(app)
	win.SetTitle("specialkey")
	win.SetDecorated(false)
	win.SetResizable(false)
	win.SetModal(true)
	win.AddCSSClass("specialkey")

	grid := gtk.NewGrid()
	grid.SetRowSpacing(10)
	grid.SetColumnSpacing(10)
	grid.SetRowHomogeneous(true)
	grid.SetColumnHomogeneous(true)

	pick := func(s string) {
		onPick(s)
		win.SetVisible(false)
		app.Quit()
	}

	var first *gtk.Button
	for i, k := range keys {
		char := gtk.NewLabel(k.char)
		char.AddCSSClass("char")
		name := gtk.NewLabel(k.label)
		name.AddCSSClass("name")

		box := gtk.NewBox(gtk.OrientationVertical, 2)
		box.SetVAlign(gtk.AlignCenter)
		box.Append(char)
		box.Append(name)

		btn := gtk.NewButton()
		btn.SetChild(box)
		btn.SetTooltipText(k.label)
		btn.ConnectClicked(func() { pick(k.char) })
		grid.Attach(btn, i%columns, i/columns, 1, 1)
		if first == nil {
			first = btn
		}
	}

	// Not focusable, so the arrow keys stay within the grid; Escape closes
	// from the keyboard.
	closeBtn := gtk.NewButtonFromIconName("window-close-symbolic")
	closeBtn.AddCSSClass("close")
	closeBtn.SetTooltipText("Close")
	closeBtn.SetFocusable(false)
	closeBtn.SetHAlign(gtk.AlignEnd)
	closeBtn.SetMarginTop(8)
	closeBtn.SetMarginEnd(8)
	closeBtn.ConnectClicked(app.Quit)

	content := gtk.NewBox(gtk.OrientationVertical, 0)
	content.Append(closeBtn)
	content.Append(grid)
	win.SetChild(content)

	keyCtl := gtk.NewEventControllerKey()
	keyCtl.ConnectKeyPressed(func(keyval, _ uint, _ gdk.ModifierType) bool {
		if keyval == gdk.KEY_Escape {
			app.Quit()
			return true
		}
		return false
	})
	win.AddController(keyCtl)

	win.SetVisible(true)
	glib.IdleAdd(func() { first.GrabFocus() })
}
