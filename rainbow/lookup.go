package rainbow

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/kmuchiri/hashbyte/hash"
)

// LookupResult holds the result of a rainbow table lookup.
type LookupResult struct {
	Found           bool
	Table           int
	Step            int
	ActivationBytes uint32
	TimeTaken       time.Duration
}

// Lookup searches for the 4-byte activation bytes that produce the given SHA-1 hash.
func Lookup(fileData []byte, targetHashHex string) (*LookupResult, error) {
	targetHashBytes, err := hex.DecodeString(targetHashHex)
	if err != nil || len(targetHashBytes) != 20 {
		return nil, fmt.Errorf("invalid SHA1 checksum provided. Must be 40 hex characters (20 bytes)")
	}
	var targetHash [20]byte
	copy(targetHash[:], targetHashBytes)

	tableSize := ChainCount * 8
	expectedSize := NumTables * tableSize
	if len(fileData) < expectedSize {
		return nil, fmt.Errorf("invalid table data size: expected at least %d bytes, got %d. Did you generate the tables?", expectedSize, len(fileData))
	}

	var tables [][]Chain
	for t := 0; t < NumTables; t++ {
		tableBytes := fileData[t*tableSize : (t+1)*tableSize]
		chains := make([]Chain, ChainCount)
		for i := 0; i < ChainCount; i++ {
			chains[i].Endpoint = binary.LittleEndian.Uint32(tableBytes[i*8 : i*8+4])
			chains[i].Startpoint = binary.LittleEndian.Uint32(tableBytes[i*8+4 : i*8+8])
		}
		tables = append(tables, chains)
	}

	startTime := time.Now()
	buf1, buf2, buf3 := hash.NewBuffers()

	for t := range NumTables {
		for step := ChainLength - 1; step >= 0; step-- {
			// Compute chain from step to end
			pt := reduce(targetHash, t, step)
			for i := step + 1; i < ChainLength; i++ {
				h := hash.Audible(pt, buf1, buf2, buf3)
				pt = reduce(h, t, i)
			}

			// Binary search the endpoint in the current table
			idx := sort.Search(ChainCount, func(i int) bool {
				return tables[t][i].Endpoint >= pt
			})

			// If matching endpoint found, regenerate the chain from the start
			if idx < ChainCount && tables[t][idx].Endpoint == pt {
				candidatePt := tables[t][idx].Startpoint
				for i := 0; i < step; i++ {
					h := hash.Audible(candidatePt, buf1, buf2, buf3)
					candidatePt = reduce(h, t, i)
				}

				// Verify it is not a collision by checking the final hash
				finalHash := hash.Audible(candidatePt, buf1, buf2, buf3)
				if finalHash == targetHash {
					return &LookupResult{
						Found:           true,
						Table:           t,
						Step:            step,
						ActivationBytes: candidatePt,
						TimeTaken:       time.Since(startTime),
					}, nil
				}
			}
		}
	}
	return &LookupResult{
		Found:     false,
		TimeTaken: time.Since(startTime),
	}, nil
}

// PrintLookupResult prints the formatted output of a LookupResult.
func PrintLookupResult(res *LookupResult) {
	if res.Found {
		fmt.Printf("\nSUCCESS!\n")
		fmt.Printf("Found in Table %d, Step %d\n", res.Table, res.Step)
		fmt.Printf("Activation Bytes: hex:%08x\n", res.ActivationBytes)
		fmt.Printf("Time taken: %v\n", res.TimeTaken)
	} else {
		fmt.Printf("Hash not found in rainbow tables (Time: %v).\n", res.TimeTaken)
	}
}

// HexOnly prints the activation bytes in hex format if found, otherwise prints an empty string.
func HexOnly(res *LookupResult) {
	if res.Found {
		fmt.Printf("%08x\n", res.ActivationBytes)
	} else {
		fmt.Printf("")
	}
}

// RunLookupAndPrint runs the lookup and prints the result.
func RunLookupAndPrint(fileData []byte, targetHashHex string) {
	fmt.Printf("Starting lookup for hash %s...\n", targetHashHex)
	res, err := Lookup(fileData, targetHashHex)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	PrintLookupResult(res)
}

// RunLookup runs the lookup and prints only the activation bytes in hex format if found, otherwise prints an empty string.
func RunLookup(fileData []byte, targetHashHex string) {
	res, err := Lookup(fileData, targetHashHex)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	HexOnly(res)
}

// RunLookupAndReturnHex runs the lookup and returns the activation bytes in hex format if found, otherwise returns an empty string.
func RunLookupAndReturnHex(fileData []byte, targetHashHex string) string {
	res, err := Lookup(fileData, targetHashHex)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	if res.Found {
		return fmt.Sprintf("%08x", res.ActivationBytes)
	} else {
		return ""
	}
}
