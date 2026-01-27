package intro

import (
	_ "embed"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/nfnt/resize"
)

//go:embed animation.gif
var animationData []byte

var (
	// Gradient colors for the loading bar
	barGradient = []string{
		"#292f56", "#2b3960", "#2c446b", "#2d5277", "#2e6282",
		"#2f758e", "#2f8c9b", "#2ea5a7", "#2db5a7", "#2cc2a1",
		"#2bd097", "#2ed989", "#35de7b", "#3de36d", "#44e860",
		"#4cec54", "#60f055", "#7af45d", "#94f766", "#acfa70",
	}
	white = lipgloss.Color("#FFFFFF")
)

var (
	preRenderedFrames []string
	animationDelay    int // in frames (approx)
)

func init() {
	loadAnimation()
}

func loadAnimation() {
	reader := strings.NewReader(string(animationData))
	g, err := gif.DecodeAll(reader)
	if err != nil {
		return
	}

	width := 60 // Fixed width for the intro
	bounds := g.Image[0].Bounds()
	canvas := image.NewRGBA(bounds)

	ratio := float64(bounds.Dy()) / float64(bounds.Dx())
	targetHeight := int(float64(width) * ratio * 2)

	for i, srcImg := range g.Image {
		drawOver(canvas, srcImg)
		resized := resize.Resize(uint(width), uint(targetHeight), canvas, resize.Lanczos3)
		preRenderedFrames = append(preRenderedFrames, renderImageToString(resized))

		if len(g.Disposal) > i && g.Disposal[i] == 2 {
			for y := 0; y < bounds.Dy(); y++ {
				for x := 0; x < bounds.Dx(); x++ {
					canvas.Set(x, y, color.Transparent)
				}
			}
		}
	}
}

func drawOver(dst *image.RGBA, src image.Image) {
	b := src.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			c := src.At(x, y)
			_, _, _, a := c.RGBA()
			if a > 0 {
				dst.Set(x, y, c)
			}
		}
	}
}

func renderImageToString(img image.Image) string {
	var sb strings.Builder
	b := img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y += 2 {
		for x := b.Min.X; x < b.Max.X; x++ {
			cTop := img.At(x, y)
			r1, g1, b1, _ := cTop.RGBA()
			var r2, g2, b2 uint32
			if y+1 < b.Max.Y {
				cBot := img.At(x, y+1)
				r2, g2, b2, _ = cBot.RGBA()
			}
			style := lipgloss.NewStyle().
				Foreground(lipgloss.Color(fmt.Sprintf("#%02x%02x%02x", r1>>8, g1>>8, b1>>8))).
				Background(lipgloss.Color(fmt.Sprintf("#%02x%02x%02x", r2>>8, g2>>8, b2>>8)))
			sb.WriteString(style.Render("▀"))
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

// RenderProgressBar creates a loading bar with the specified gradient
func RenderProgressBar(percent int, width int) string {
	if percent > 100 {
		percent = 100
	}

	filledCells := (percent * width) / 100
	var bar strings.Builder

	for i := 0; i < width; i++ {
		if i < filledCells {
			// Calculate gradient index
			gradIdx := (i * len(barGradient)) / width
			color := barGradient[gradIdx]
			bar.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render("█"))
		} else {
			bar.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#333333")).Render("░"))
		}
	}

	percentText := lipgloss.NewStyle().Foreground(white).Bold(true).Render(fmt.Sprintf(" %d%%", percent))
	return bar.String() + percentText
}

// RenderText renders centered text
func RenderText(text string, width int) string {
	return lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Bold(true).Render(text)
}
