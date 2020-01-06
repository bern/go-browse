package tests

import (
	"github.com/fogleman/gg"
)

// DrawSomething paints a simple circle onto a gg context
func DrawSomething() {
	// Create a new context (canvas)
	dc := gg.NewContext(1000, 1000)

	// Draw a rectangle to take over the canvas
	dc.DrawRectangle(0, 0, 1000, 1000)

	// Set to white
	dc.SetRGB(1, 1, 1)
	dc.Fill()

	// Draw a centered circle with a rad of 400
	dc.DrawCircle(500, 500, 400)

	// Set to black
	dc.SetRGB(0, 0, 0)
	dc.Fill()

	// Output as a png
	dc.SavePNG("out.png")
}
