package detector

import (
    "bytes"
    "fmt"
    "io"
    "os"
)

// Block represents a chunk of data being processed
type Block struct {
    offset   int64
    data     []byte
    location string
    source   ScanTarget
}

// Detection represents the result of detecting a wallet or file
type Detection struct {
    Description string
}

// ProgressInfo keeps track of the current scan status
type ProgressInfo struct {
    CurrentTarget    string
    ScannedBytes     int64
    TotalBytes       int64
    UnscannedTargets int
}

// ScanTarget defines the interface for targets to be scanned
type ScanTarget interface {
    Describe() string
    StartOffset() int64
    Size() (int64, error)
    Open() (io.ReadCloser, error)
}

// FileScanTarget implements the ScanTarget interface for file-based scanning
type FileScanTarget struct {
    path        string
    startOffset int64
}

func (t *FileScanTarget) Describe() string {
    return t.path
}

func (t *FileScanTarget) StartOffset() int64 {
    return t.startOffset
}

func (t *FileScanTarget) Size() (int64, error) {
    info, err := os.Stat(t.path)
    if err != nil {
        return 0, err
    }
    return info.Size(), nil
}

func (t *FileScanTarget) Open() (io.ReadCloser, error) {
    return os.Open(t.path)
}

// Scan performs the entire scanning operation for detecting wallet traces
func Scan(startOffset int64, path string, onDetection func(Detection), onProgress func(ProgressInfo)) error {
    scanTargets := make(chan ScanTarget, 1024)
    emptyBlocks := make(chan *Block, 30)
    detectionQueue := make(chan *Block, 30)
    signals := make(chan error, 1)

    // Initialize block pool
    for i := 0; i < 20; i++ {
        emptyBlocks <- &Block{
            data: make([]byte, 4*1024),
        }
    }

    // Close channels when complete
    onComplete := func() {
        close(emptyBlocks)
        close(detectionQueue)
        signals <- io.EOF
    }

    // Start scan operations
    go scanBlocks(scanTargets, emptyBlocks, detectionQueue, onProgress)
    go detectWallets(detectionQueue, emptyBlocks, onDetection, onComplete)

    // Publish the initial target
    scanTargets <- &FileScanTarget{
        path:        path,
        startOffset: startOffset,
    }

    // Wait for the signal to finish
    signal := <-signals
    return signal
}

// scanBlocks processes each scan target, splitting the data into blocks and sending it through the pipeline
func scanBlocks(targets chan ScanTarget, emptyBlocks chan *Block, out chan *Block, onProgress func(ProgressInfo)) {
    for target := range targets {
        totalBytes, err := target.Size()
        if err != nil {
            fmt.Printf("[scan] Error getting size: %s\n", err)
            continue
        }

        reader, err := target.Open()
        if err != nil {
            fmt.Printf("[scan] Error opening target: %s", err)
            continue
        }
        defer reader.Close()

        offset := target.StartOffset()
        for block := range emptyBlocks {
            n, err := reader.Read(block.data)
            if err == io.EOF {
                if n == 0 {
                    break
                }
            } else if err != nil {
                fmt.Printf("[scan] Error reading: %s", err)
                break
            }

            block.offset = offset
            block.source = target
            offset += int64(n)
            onProgress(ProgressInfo{
                CurrentTarget: target.Describe(),
                ScannedBytes:  offset,
                TotalBytes:    totalBytes,
            })

            out <- block
        }
    }
}

// detectWallets scans blocks for specific wallet-related data patterns
func detectWallets(in chan *Block, out chan *Block, onDetection func(Detection), onComplete func()) {
    walletSignatures := [][]byte{
        []byte("orderposnext"),
        []byte("wallet.dat"),
    }

    for block := range in {
        for _, signature := range walletSignatures {
            if bytes.Contains(block.data, signature) {
                onDetection(Detection{
                    Description: fmt.Sprintf("Found '%s' in block at offset %d", signature, block.offset),
                })
            }
        }
        out <- block
    }
    onComplete()
}
