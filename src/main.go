package main

import (
	"fmt"
	"fractalview/mandelbrot"
	"image/jpeg"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type appConfig struct {
	exploreButtons []buttonConfig
	downloadButton buttonConfig
	status         *widget.Label
	imageContent   *fyne.Container
	imageConfig    *mandelbrot.ImgConfig
	window         fyne.Window
}

type buttonConfig struct {
	title  string
	icon   fyne.Resource
	widget *widget.Button
}

func main() {
	myApp := app.New()
	appConfig := initAppConfig(myApp)

	appConfig.window = myApp.NewWindow("FractalView")
	buttonBar := container.New(layout.NewHBoxLayout(), appConfig.exploreButtons[0].widget, appConfig.exploreButtons[1].widget,
		appConfig.exploreButtons[2].widget, appConfig.exploreButtons[3].widget, appConfig.exploreButtons[4].widget,
		appConfig.exploreButtons[5].widget, layout.NewSpacer(), appConfig.downloadButton.widget)
	statusBar := container.New(layout.NewHBoxLayout(), appConfig.status, layout.NewSpacer())
	appConfig.window.SetContent(container.New(layout.NewVBoxLayout(), buttonBar, appConfig.imageContent, statusBar))
	appConfig.window.ShowAndRun()
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
	result.exploreButtons = make([]buttonConfig, 6)
	result.exploreButtons[0] = *initButton("left", iconScheme, &result)
	result.exploreButtons[1] = *initButton("right", iconScheme, &result)
	result.exploreButtons[2] = *initButton("up", iconScheme, &result)
	result.exploreButtons[3] = *initButton("down", iconScheme, &result)
	result.exploreButtons[4] = *initButton("plus", iconScheme, &result)
	result.exploreButtons[5] = *initButton("minus", iconScheme, &result)
	result.status = widget.NewLabel("")
	result.downloadButton = buttonConfig{
		title: "Download",
		widget: widget.NewButtonWithIcon("Download", nil, func() {
			downloadImage(&result)
		}),
	}
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
	for _, button := range appConfig.exploreButtons {
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
		for _, button := range appConfig.exploreButtons {
			button.widget.Enable()
			buttonText := button.title
			button.widget.OnTapped = func() {
				updateContent(appConfig, buttonText)
			}
		}
	})
}

func downloadImage(appConfig *appConfig) {
	fileDialog := dialog.NewFileSave(appConfig.dialogCallback, appConfig.window)
	fileDialog.SetFileName("mandelbrot.jpg")
	fileDialog.Show()
}

func (config *appConfig) dialogCallback(closer fyne.URIWriteCloser, err error) {
	if closer.URI() != nil {
		img := mandelbrot.CreateImage(config.imageConfig)
		toimg, _ := os.Create(closer.URI().Path())
		defer toimg.Close()
		jpeg.Encode(toimg, img, &jpeg.Options{Quality: jpeg.DefaultQuality})
		fmt.Println("File saved at", closer.URI().Path())
	}
}
