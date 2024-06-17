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
	Public    bool   `json:"public"`
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
		fmt.Fprintf(os.Stderr, "./images is empty, i expect it to have 1 or more image files\n")
		os.Exit(0)
	}

	// show marked files first
	imagePaths = groupImages(imagePaths, preferences)

	currentIndex := 0
	img := canvas.NewImageFromFile(imagePaths[currentIndex])
	img.FillMode = canvas.ImageFillContain
	img.SetMinSize(imageSize)

	public := widget.NewRadioGroup([]string{"pUblic", "pRivate"}, func(string) {})
	updatePublicSelection(preferences, imagePaths[currentIndex], public)

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
			setPreference(preferences, imagePaths[currentIndex], true, public)
			currentIndex = (currentIndex + 1) % len(imagePaths)
		case fyne.KeyR:
			setPreference(preferences, imagePaths[currentIndex], false, public)
			currentIndex = (currentIndex + 1) % len(imagePaths)
		case fyne.KeyD:
			deletePreference(preferences, imagePaths[currentIndex], public)
			currentIndex = (currentIndex + 1) % len(imagePaths)
		case fyne.KeyN:
			currentIndex = (currentIndex + 1) % len(imagePaths)
		case fyne.KeyP:
			currentIndex = (currentIndex - 1 + len(imagePaths)) % len(imagePaths)
		case fyne.KeyT:
			currentIndex = findNextUnmarkedImage(preferences, imagePaths, currentIndex)
		default:
			return
		}

		img.File = imagePaths[currentIndex]
		img.Refresh()

		updatePublicSelection(preferences, imagePaths[currentIndex], public)
	})

	w.ShowAndRun()
}

func updatePublicSelection(preferences map[string]ImagePreference, imagePath string, public *widget.RadioGroup) {
	if preference, ok := preferences[imagePath]; ok {
		if preference.Public {
			public.SetSelected("pUblic")
		} else {
			public.SetSelected("pRivate")
		}
	} else {
		public.SetSelected("")
	}
}

func setPreference(preferences map[string]ImagePreference, imagePath string, isPublic bool, public *widget.RadioGroup) {
	preferences[imagePath] = ImagePreference{ImagePath: imagePath, Public: isPublic}
	if err := savePreferences(preferences, "preferences.json"); err != nil {
		fmt.Fprintf(os.Stderr, "error saving preferences, %v\n", err)
		os.Exit(1)
	}
	if isPublic {
		public.SetSelected("pUblic")
	} else {
		public.SetSelected("pRivate")
	}
}

func deletePreference(preferences map[string]ImagePreference, imagePath string, public *widget.RadioGroup) {
	delete(preferences, imagePath)
	if err := savePreferences(preferences, "preferences.json"); err != nil {
		fmt.Fprintf(os.Stderr, "error saving preferences, %v\n", err)
		os.Exit(1)
	}
	public.SetSelected("")
}

func findNextUnmarkedImage(preferences map[string]ImagePreference, imagePaths []string, currentIndex int) int {
	for i := 0; i < len(imagePaths); i++ {
		if _, ok := preferences[imagePaths[currentIndex]]; !ok {
			break
		}
		currentIndex = (currentIndex + 1) % len(imagePaths)
	}
	return currentIndex
}

func groupImages(imagePaths []string, preferences map[string]ImagePreference) []string {
	var preferredImages, unpreferredImages []string

	for _, imagePath := range imagePaths {
		if _, ok := preferences[imagePath]; ok {
			preferredImages = append(preferredImages, imagePath)
		} else {
			unpreferredImages = append(unpreferredImages, imagePath)
		}
	}

	return append(preferredImages, unpreferredImages...)
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
