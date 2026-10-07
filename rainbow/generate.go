package rainbow

import (
	"encoding/binary"
	"fmt"
	"hashbyte/hash"
	"math/rand"
	"os"
	"runtime"
	"sort"
	"sync"
	"time"
)

// Table Parameters
// 4 tables * 1,048,576 chains * 8 bytes = 33.5 MB total size
const (
	NumTables   = 4
	ChainCount  = 1048576 // 2^20 chains per table
	ChainLength = 4096    // 2^12 links per chain
)

// Chain stores the start and end points of a rainbow chain.
type Chain struct {
	Endpoint   uint32
	Startpoint uint32
}

// Sort interfaces for binary searching the endpoints
type ByEndpoint []Chain

func (a ByEndpoint) Len() int           { return len(a) }
func (a ByEndpoint) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByEndpoint) Less(i, j int) bool { return a[i].Endpoint < a[j].Endpoint }

// reduce maps a 20-byte hash back to a 4-byte plaintext.
func reduce(hash [20]byte, tableID, step int) uint32 {
	hVal := binary.LittleEndian.Uint32(hash[0:4])
	return hVal + uint32(step) + uint32(tableID*10000)
}

// Generate builds rainbow tables and writes them to the given file.
func Generate(filename string) {
	fmt.Printf("Generating %d Rainbow Tables (%d chains of length %d each)...\n", NumTables, ChainCount, ChainLength)
	fmt.Printf("Estimated output size: %.2f MB\n", float64(NumTables*ChainCount*8)/(1024*1024))

	startTime := time.Now()

	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	for t := range NumTables {
		tableStartTime := time.Now()
		chains := make([]Chain, ChainCount)

		var wg sync.WaitGroup
		numWorkers := runtime.NumCPU()
		chunkSize := ChainCount / numWorkers

		for w := range numWorkers {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()
				startIdx := workerID * chunkSize
				endIdx := startIdx + chunkSize
				if workerID == numWorkers-1 {
					endIdx = ChainCount
				}

				buf1, buf2, buf3 := hash.NewBuffers()

				// Use workerID and time as seed to avoid duplicate starts
				r := rand.New(rand.NewSource(time.Now().UnixNano() + int64(workerID)))

				for i := startIdx; i < endIdx; i++ {
					startPt := r.Uint32()
					pt := startPt

					for step := 0; step < ChainLength; step++ {
						h := hash.Audible(pt, buf1, buf2, buf3)
						pt = reduce(h, t, step)
					}

					chains[i] = Chain{Endpoint: pt, Startpoint: startPt}
				}
			}(w)
		}
		wg.Wait()

		sort.Sort(ByEndpoint(chains))

		// Write to disk
		binary.Write(file, binary.LittleEndian, chains)
		fmt.Printf("Table %d/%d generated in %v\n", t+1, NumTables, time.Since(tableStartTime))
	}

	fmt.Printf("All tables generated successfully in %v. Saved to %s\n", time.Since(startTime), filename)
}
