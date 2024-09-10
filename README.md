
# Bitcoin Wallet and Compressed File Detection System

## Overview

An application written in Go that scans files or devices for evidence of Bitcoin wallets and compressed files (e.g., Gzip, Zip). It breaks down large files into manageable blocks, scans for wallet-related keys or specific byte sequences, and processes compressed files recursively.

The application leverages **goroutines** to ensure that the scan runs asynchronously, allowing the UI to remain responsive. The log output is displayed in a scrollable window, and the UI provides basic controls to start, stop, and pause the scan.

## Features

- **Bitcoin Wallet Detection**: Scans files and devices for specific byte patterns related to Bitcoin wallets.
- **Compressed File Detection**: Supports scanning of compressed files (Gzip, Zip) and recursively scans their contents.
- **Asynchronous Scanning**: Uses Go's goroutines to keep the UI responsive while the scan runs in the background.
- **Log Output**: Displays real-time log information, including scan progress and detection results, in a scrollable text window.

## How to Use

1. **Select a File**: Click on the "Select File and Start Scan" button to choose the file you wish to scan.
2. **View Progress**: The progress of the scan will be shown in the log window, with detection results and current progress.
3. **Stop and Pause**: The Stop and Pause buttons are placeholders (not yet implemented) for future functionality.

## Running the Application

To run the application, ensure you have Go installed (version 1.19 or higher) and the Fyne GUI library installed. Use the following command to run the program:

```bash
go run -tags "gles2" main.go
```

Make sure the Go environment is set up correctly and the necessary dependencies are installed.

## TODOS

- **Implement stopping of scan**: Add functionality to allow the user to stop the scan mid-process.
- **Implement pausing of scan**: Add the ability to pause and resume the scan.
- **Add ability to scan an entire drive**: Extend the program to allow scanning of entire drives or block devices.
- **Fix the toggling of light/dark mode**: Ensure the theme toggle button properly switches between light and dark mode.
- **Fix some layout issues**: Improve the layout of the control buttons and log window for better visual consistency and user experience.

