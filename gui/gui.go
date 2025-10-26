package gui

import (
	"ChaCha20/chacha20"
	"ChaCha20/internal"
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"math"
	"strconv"
	"strings"
)

type App struct {
	app    fyne.App
	window fyne.Window

	dataFileLabel        *widget.Label
	dataFileButton       *widget.Button
	selectedDataFilePath *string

	keyInputText   *widget.Entry
	generateKeyBtn *widget.Button
	key            *internal.Key

	nonceInputText   *widget.Entry
	generateNonceBtn *widget.Button
	nonce            *internal.Nonce

	counterInputText *widget.Entry
	counter          *uint32

	loader     *widget.ProgressBarInfinite
	encryptBtn *widget.Button
}

func (mApp *App) onKeyInputChange(val string) error {
	key, err := chacha20.GetKeyFromStr(val)
	if err != nil {
		mApp.key = nil
		mApp.updateEncryptBtnState()
		return err
	}

	mApp.updateEncryptBtnState()
	mApp.key = key
	return nil
}

func (mApp *App) onNonceInputChange(val string) error {
	nonce, err := chacha20.GetNonceFromString(val)
	if err != nil {
		mApp.nonce = nil
		mApp.updateEncryptBtnState()
		return err
	}

	mApp.updateEncryptBtnState()
	mApp.nonce = nonce
	return nil
}

func (mApp *App) onCounterInputChange(val string) error {
	num, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		mApp.counter = nil
		mApp.updateEncryptBtnState()
		return err
	}

	if num < 0 || num >= math.MaxUint32 {
		mApp.counter = nil
		mApp.updateEncryptBtnState()
		return fmt.Errorf("number must be uint32")
	}

	mApp.updateEncryptBtnState()
	uint32val := uint32(num)
	mApp.counter = &uint32val
	return nil
}

func (mApp *App) updateEncryptBtnState() {
	if mApp.key == nil || mApp.nonce == nil || mApp.counter == nil || mApp.selectedDataFilePath == nil {
		mApp.encryptBtn.Disable()
		return
	}

	mApp.encryptBtn.Enable()
}

func (mApp *App) setSelectedDataFilePath(val *string) {
	mApp.selectedDataFilePath = val
	mApp.updateEncryptBtnState()

	if mApp.selectedDataFilePath == nil {
		mApp.dataFileLabel.SetText("No file selected")
		return
	}

	mApp.dataFileButton.SetText(*mApp.selectedDataFilePath)
}

func (mApp *App) createGenerateBtns() {
	mApp.generateNonceBtn = widget.NewButton("Generate", func() {
		nonceStr := chacha20.GenerateNonce()
		mApp.nonceInputText.SetText(nonceStr)
	})

	mApp.generateKeyBtn = widget.NewButton("Generate", func() {
		keyStr := chacha20.GenerateKey()
		mApp.keyInputText.SetText(keyStr)
	})
}

func (mApp *App) setLoading(finished bool) {
	if finished {
		mApp.encryptBtn.Enable()
		mApp.keyInputText.Enable()
		mApp.nonceInputText.Enable()
		mApp.counterInputText.Enable()
		mApp.dataFileButton.Enable()
		mApp.loader.Hide()
		return
	}

	mApp.loader.Show()
	mApp.encryptBtn.Disable()
	mApp.keyInputText.Disable()
	mApp.nonceInputText.Disable()
	mApp.counterInputText.Disable()
	mApp.dataFileButton.Disable()
}

func (mApp *App) Init() {
	mApp.app = app.NewWithID("ChaCha20")
	mApp.window = mApp.app.NewWindow("ChaCha20")

	mApp.keyInputText = widget.NewEntry()
	mApp.keyInputText.Validator = mApp.onKeyInputChange
	mApp.keyInputText.AlwaysShowValidationError = true

	mApp.nonceInputText = widget.NewEntry()
	mApp.nonceInputText.Validator = mApp.onNonceInputChange
	mApp.nonceInputText.AlwaysShowValidationError = true

	mApp.counterInputText = widget.NewEntry()
	mApp.counterInputText.SetText("0")
	initialCounter := uint32(0)
	mApp.counter = &initialCounter
	mApp.counterInputText.Validator = mApp.onCounterInputChange
	mApp.counterInputText.AlwaysShowValidationError = true

	mApp.dataFileLabel = widget.NewLabel("No file selected")
	mApp.selectedDataFilePath = nil
	mApp.dataFileButton = widget.NewButton("Select data file", func() {
		d := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, mApp.window)
				mApp.setSelectedDataFilePath(nil)
				return
			}
			if reader == nil {
				mApp.setSelectedDataFilePath(nil)
				dialog.ShowError(errors.New("error occured selecting file"), mApp.window)
				return
			}

			path := reader.URI().Path()
			mApp.setSelectedDataFilePath(&path)

			defer reader.Close()
		}, mApp.window)
		d.Show()
	})

	mApp.loader = widget.NewProgressBarInfinite()
	mApp.loader.Hide()
	mApp.encryptBtn = widget.NewButton("Run", func() {
		encryptedSuffix := ".enc"
		decrypting := strings.HasSuffix(*mApp.selectedDataFilePath, encryptedSuffix)

		var outFile string
		if decrypting {
			outFile = strings.TrimSuffix(*mApp.selectedDataFilePath, encryptedSuffix)
		} else {
			outFile = *mApp.selectedDataFilePath + encryptedSuffix
		}

		mApp.setLoading(false)
		go func() {
			err := chacha20.EncryptFile(*mApp.key, *mApp.counter, *mApp.nonce, *mApp.selectedDataFilePath, outFile)
			mApp.setLoading(true)
			if err != nil {
				dialog.ShowError(err, mApp.window)
				return
			}
		}()
	})
	mApp.encryptBtn.Disable()

	mApp.createGenerateBtns()
}

func (mApp *App) ShowAndRun() {

	mainContent := container.NewVBox(
		&widget.Form{
			Items: []*widget.FormItem{
				{
					Text: "Key",
					Widget: container.NewVBox(
						mApp.keyInputText,
						mApp.generateKeyBtn,
					),
				},
				{
					Text: "Nonce",
					Widget: container.NewVBox(
						mApp.nonceInputText,
						mApp.generateNonceBtn,
					),
				},
				{Text: "Counter", Widget: mApp.counterInputText},
			},
		},
		mApp.dataFileLabel,
		mApp.dataFileButton,
		mApp.encryptBtn,
	)

	mApp.window.SetContent(container.NewStack(
		mainContent,
		mApp.loader,
	))

	mApp.window.ShowAndRun()
}

func ShowGui() {
	myApp := App{}
	myApp.Init()
	myApp.ShowAndRun()
}
