package main

import (
	"fmt"
	"fractalview/mandelbrot"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type appConfig struct {
	buttons      []buttonConfig
	status       *widget.Label
	imageContent *fyne.Container
	imageConfig  *mandelbrot.ImgConfig
}

type buttonConfig struct {
	title  string
	icon   fyne.Resource
	widget *widget.Button
}

func main() {
	myApp := app.New()
	appConfig := initAppConfig(myApp)

	myWindow := myApp.NewWindow("FractalView")
	buttonBar := container.New(layout.NewHBoxLayout(), appConfig.buttons[0].widget, appConfig.buttons[1].widget,
		appConfig.buttons[2].widget, appConfig.buttons[3].widget, appConfig.buttons[4].widget,
		appConfig.buttons[5].widget, layout.NewSpacer())
	statusBar := container.New(layout.NewHBoxLayout(), appConfig.status)
	myWindow.SetContent(container.New(layout.NewVBoxLayout(), buttonBar, appConfig.imageContent, statusBar))
	myWindow.ShowAndRun()
}

func (appConfig *appConfig) UpdateStatus() {
	appConfig.status.SetText(fmt.Sprintf("X: %f, Y: %f, Zoom: %f", appConfig.imageConfig.CenterX,
		appConfig.imageConfig.CenterY, 1/appConfig.imageConfig.ZoomX))
}

func initAppConfig(app fyne.App) *appConfig {
	iconScheme := "white"
	if app.Settings().ThemeVariant() != 0 {
		iconScheme = "black"
	}

	result := appConfig{}
	result.imageConfig = mandelbrot.NewConfig()
	result.imageContent = container.New(layout.NewMaxLayout(), createImage(result.imageConfig))
	result.buttons = make([]buttonConfig, 6)
	result.buttons[0] = *initButton("left", iconScheme, &result)
	result.buttons[1] = *initButton("right", iconScheme, &result)
	result.buttons[2] = *initButton("up", iconScheme, &result)
	result.buttons[3] = *initButton("down", iconScheme, &result)
	result.buttons[4] = *initButton("plus", iconScheme, &result)
	result.buttons[5] = *initButton("minus", iconScheme, &result)
	result.status = widget.NewLabel("")
	result.UpdateStatus()
	return &result
}

func initButton(name string, iconScheme string, appConfig *appConfig) *buttonConfig {
	config := buttonConfig{}
	config.title = name

	var error error
	config.icon, error = fyne.LoadResourceFromPath(fmt.Sprintf("assets/%s-%s.svg", name, iconScheme))

	buttonText := ""
	if error != nil {
		buttonText = config.title
	}
	config.widget = widget.NewButtonWithIcon(buttonText, config.icon, func() {
		updateContent(appConfig, name)
	})
	return &config
}

func createImage(imgConfig *mandelbrot.ImgConfig) *canvas.Image {
	image := canvas.NewImageFromImage(mandelbrot.CreateImage(imgConfig))
	image.FillMode = canvas.ImageFillContain
	image.SetMinSize(fyne.NewSize(800, 480))
	return image
}

func updateContent(appConfig *appConfig, source string) {
	for _, button := range appConfig.buttons {
		button.widget.Disable()
		button.widget.OnTapped = func() {}
	}

	//log.Println("update", source)
	appConfig.imageConfig.Move(source)
	appConfig.imageConfig.Scale(source)
	appConfig.UpdateStatus()
	appConfig.imageContent.Objects[0] = createImage(appConfig.imageConfig)
	appConfig.imageContent.Refresh()
	time.AfterFunc(1*time.Second, func() {
		for _, button := range appConfig.buttons {
			button.widget.Enable()
			buttonText := button.title
			button.widget.OnTapped = func() {
				updateContent(appConfig, buttonText)
			}
		}
	})
}
