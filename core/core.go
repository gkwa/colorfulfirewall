package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type ImagePreference struct {
	ImagePath string `json:"imagePath"`
	Liked     bool   `json:"liked"`
}

func Run() {
	a := app.New()
	w := a.NewWindow("Image Feedback")

	imageSize := fyne.NewSize(3*2*96, 3*2*96)
	windowSize := fyne.NewSize(imageSize.Width, imageSize.Height+50)

	w.Resize(windowSize)
	w.SetFixedSize(true)

	preferences, _ := loadPreferences("preferences.json")

	imagePaths, _ := filepath.Glob("images/*.png")

	if len(imagePaths) == 0 {
		fmt.Fprintf(os.Stderr, "./images is empty, i expect it to have 1 or morg image files\n")
		os.Exit(0)
	}

	currentIndex := 0
	img := canvas.NewImageFromFile(imagePaths[currentIndex])
	img.FillMode = canvas.ImageFillContain
	img.SetMinSize(imageSize)

	public := widget.NewRadioGroup([]string{"pUblic", "pRivate"}, func(string) {})

	if preference, ok := preferences[imagePaths[currentIndex]]; ok {
		if preference.Liked {
			public.SetSelected("pUblic")
		} else {
			public.SetSelected("pRivate")
		}
	}

	content := container.NewVBox(
		img,
		container.NewVBox(
			widget.NewLabel("Is this image public or private?"),
			public,
		),
	)

	w.SetContent(container.NewCenter(content))

	w.Canvas().SetOnTypedKey(func(ev *fyne.KeyEvent) {
		switch ev.Name {
		case fyne.KeyU:
			public.SetSelected("pUblic")
			preferences[imagePaths[currentIndex]] = ImagePreference{ImagePath: imagePaths[currentIndex], Liked: true}
			savePreferences(preferences, "preferences.json")
			currentIndex = (currentIndex + 1) % len(imagePaths)
		case fyne.KeyR:
			public.SetSelected("pRivate")
			preferences[imagePaths[currentIndex]] = ImagePreference{ImagePath: imagePaths[currentIndex], Liked: false}
			savePreferences(preferences, "preferences.json")
			currentIndex = (currentIndex + 1) % len(imagePaths)
		case fyne.KeyD:
			public.SetSelected("")
			delete(preferences, imagePaths[currentIndex])
			savePreferences(preferences, "preferences.json")
			currentIndex = (currentIndex + 1) % len(imagePaths)
		case fyne.KeyN:
			currentIndex = (currentIndex + 1) % len(imagePaths)
		case fyne.KeyP:
			currentIndex = (currentIndex - 1 + len(imagePaths)) % len(imagePaths)
		default:
			return
		}

		img.File = imagePaths[currentIndex]
		img.Refresh()

		if preference, ok := preferences[imagePaths[currentIndex]]; ok {
			if preference.Liked {
				public.SetSelected("pUblic")
			} else {
				public.SetSelected("pRivate")
			}
		} else {
			public.SetSelected("")
		}
	})

	w.ShowAndRun()
}

func loadPreferences(filename string) (map[string]ImagePreference, error) {
	preferences := make(map[string]ImagePreference)

	data, err := os.ReadFile(filename)
	if err != nil {
		return preferences, nil
	}

	err = json.Unmarshal(data, &preferences)
	if err != nil {
		return nil, fmt.Errorf("failed to parse preferences: %v", err)
	}

	return preferences, nil
}

func savePreferences(preferences map[string]ImagePreference, filename string) error {
	data, err := json.MarshalIndent(preferences, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal preferences: %v", err)
	}

	err = os.WriteFile(filename, data, 0o644)
	if err != nil {
		return fmt.Errorf("failed to save preferences: %v", err)
	}

	return nil
}
