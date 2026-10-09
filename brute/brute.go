package brute

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/kmuchiri/hashbyte/hash"
	"os"
	"runtime"
	"sync"
	"time"
)

// Run performs a parallel brute-force search across the full 2^32 keyspace
// to find the 4-byte activation bytes matching the given SHA-1 hash.
func Run(targetHashHex string) {
	targetHash, err := hex.DecodeString(targetHashHex)
	if err != nil || len(targetHash) != sha1.Size {
		fmt.Println("Error: Invalid SHA1 checksum provided. Must be 40 hex characters (20 bytes).")
		os.Exit(1)
	}

	match, elapsed, found := findHash(targetHash)
	printResult(targetHashHex, match, found, elapsed)
}

func returnHex(match [4]byte) string {
	return fmt.Sprintf("%02X%02X%02X%02X", match[0], match[1], match[2], match[3])
}

func GetBytesOnly(targetHashHex string) string {
	targetHash, err := hex.DecodeString(targetHashHex)
	if err != nil || len(targetHash) != sha1.Size {
		fmt.Println("Error: Invalid SHA1 checksum provided. Must be 40 hex characters (20 bytes).")
		os.Exit(1)
	}

	match, _, found := findHash(targetHash)
	if found {
		return returnHex(match)
	} else {
		return ""
	}
}

func findHash(targetHash []byte) ([4]byte, time.Duration, bool) {
	numWorkers := runtime.NumCPU()
	fmt.Printf("Starting brute force with %d workers...\n", numWorkers)
	fmt.Printf("Target hash: %x\n", targetHash)

	startTime := time.Now()

	resultChan := make(chan [4]byte)
	doneChan := make(chan struct{})

	var wg sync.WaitGroup

	const totalKeys = 4294967296 // 0x100000000
	chunkSize := uint32(totalKeys / int64(numWorkers))

	for i := range numWorkers {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			start := uint32(workerID) * chunkSize
			var maxVal uint32
			if workerID == numWorkers-1 {
				maxVal = 0xFFFFFFFF
			} else {
				maxVal = start + chunkSize - 1
			}

			buf1, buf2, buf3 := hash.NewBuffers()

			for k := start; k <= maxVal; k++ {
				if k&0xFFFF == 0 {
					select {
					case <-doneChan:
						return
					default:
					}
				}

				res := hash.Audible(k, buf1, buf2, buf3)

				if bytes.Equal(res[:], targetHash) {
					var result [4]byte
					result[0] = byte(k >> 24)
					result[1] = byte(k >> 16)
					result[2] = byte(k >> 8)
					result[3] = byte(k)
					resultChan <- result
					return
				}

				if k == maxVal {
					break
				}
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	match, found := <-resultChan
	if found {
		close(doneChan)
	}
	elapsed := time.Since(startTime)

	return match, elapsed, found
}

func printResult(targetHashHex string, match [4]byte, found bool, elapsed time.Duration) {
	if found {
		fmt.Printf("\nMetrics\n")
		fmt.Printf("-------------------------------------------------------\n")
		fmt.Printf("plaintext found:                              1 of 1\n")
		fmt.Printf("total time:                                   %.2f s\n", elapsed.Seconds())
		fmt.Printf("-------------------------------------------------------\n")
		fmt.Printf("result\n")
		fmt.Printf("-------------------------------------------------------\n")
		fmt.Printf("%s...                               hex:%s\n",
			targetHashHex[:8], returnHex(match))
	} else {
		fmt.Printf("\nNo match found.\n")
	}
}
