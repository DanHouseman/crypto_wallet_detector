package main

import (
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/dialog"
    "fyne.io/fyne/v2/theme"
    "fyne.io/fyne/v2/widget"
    "log"
    "os"
    "io"
    "sync"
    "go_bitcoin_detector/detector"
)

// CustomTheme extends the default theme to modify the scrollbar size
type CustomTheme struct {
    fyne.Theme
}

// ScrollBarSize overrides the default scrollbar width
func (c CustomTheme) ScrollBarSize() float32 {
    return 20 // Set the width of the scrollbar here
}

func main() {
    a := app.New()

    // Set the custom theme to increase the scrollbar width
    a.Settings().SetTheme(&CustomTheme{theme.DefaultTheme()})

    w := a.NewWindow("File Scanner")

    // Multi-line text entry for log output
    logOutput := widget.NewMultiLineEntry()
    logOutput.SetText("Logs will appear here...")
    logOutput.Wrapping = fyne.TextWrapWord // Set text wrapping
    logOutput.Disable()                     // Disable user editing

    // Make logOutput scrollable with a wider scrollbar
    scrollableLog := container.NewVScroll(logOutput)
    scrollableLog.SetMinSize(fyne.NewSize(500, 150)) // Set minimum size for the log box

    // File selection dialog
    startButton := widget.NewButton("Select File and Start Scan", func() {
        fileDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
            if err != nil || reader == nil {
                logOutput.SetText("No file selected or an error occurred.")
                return // Handle file dialog cancel or error
            }

            filePath := reader.URI().Path()
            logOutput.SetText("Selected file: " + filePath)
            reader.Close()

            // Start scanning the selected file in a separate goroutine
            go startScan(filePath, logOutput)
        }, w)
        fileDialog.Show()
    })

    stopButton := widget.NewButton("Stop Scan", func() {
        log.Println("Stopping scan...") // Add logic to stop scan
    })

    pauseButton := widget.NewButton("Pause Scan", func() {
        log.Println("Pausing scan...") // Add logic to pause scan
    })

    // Declare themeToggle before defining its functionality
    themeToggle := widget.NewButton("", nil) // Initialize with an empty label and nil function

    // Now define the behavior of themeToggle after it has been declared
    themeToggle.SetText("Dark Mode") // Start with Dark Mode as default text
    themeToggle.OnTapped = func() {
        if a.Settings().Theme() == theme.DarkTheme() {
            a.Settings().SetTheme(theme.LightTheme())
            themeToggle.SetText("Dark Mode")
        } else {
            a.Settings().SetTheme(theme.DarkTheme())
            themeToggle.SetText("Light Mode")
        }
    }

    // Create a container for the control buttons
    controlButtons := container.NewHBox(startButton, stopButton, pauseButton, themeToggle)

    // Use a Border layout to position the buttons at the top and the scrollable log at the bottom
    content := container.NewBorder(controlButtons, scrollableLog, nil, nil)

    w.SetContent(content)

    w.Resize(fyne.NewSize(500, 400))
    w.ShowAndRun()
}

// startScan runs the scan in a separate goroutine to keep the UI responsive
func startScan(filePath string, logOutput *widget.Entry) {
    // Create a wait group to handle goroutines properly
    var wg sync.WaitGroup
    wg.Add(1)

    // Run the scan asynchronously
    go func() {
        defer wg.Done() // Mark goroutine as done when finished

        // Open log file
        logFile, err := os.Create("scan.log")
        if err != nil {
            log.Fatalf("Error creating log file: %v", err)
        }
        defer logFile.Close()

        // MultiWriter to log to both file and console
        multiWriter := io.MultiWriter(os.Stdout, logFile)

        // Set log output to multi-writer
        log.SetOutput(multiWriter)

        logOutput.SetText("Starting scan on: " + filePath)

        onDetection := func(d detector.Detection) {
            log.Printf("Detection: %s", d.Description)
            logOutput.SetText(logOutput.Text + "\nDetection: " + d.Description)
            logOutput.Refresh() // Refresh the log output after adding new text
        }

        onProgress := func(progress detector.ProgressInfo) {
            log.Printf("Scanning target: %s, Progress: %d/%d", progress.CurrentTarget, progress.ScannedBytes, progress.TotalBytes)
            logOutput.SetText(logOutput.Text + "\nScanning target: " + progress.CurrentTarget + " Progress: " +
                string(progress.ScannedBytes) + "/" + string(progress.TotalBytes))
            logOutput.Refresh() // Refresh the log output after adding new text
        }

        // Start scanning the selected file
        err = detector.Scan(0, filePath, onDetection, onProgress)
        if err != nil {
            log.Printf("Scan error: %v", err)
            logOutput.SetText(logOutput.Text + "\nScan error: " + err.Error())
            logOutput.Refresh() // Refresh the log output after adding new text
        }
    }()
    
    // Wait for the scan to complete
    wg.Wait()
}
