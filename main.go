package main

import (
	_ "embed"
	"fmt"
	"os"

	"hashbyte/brute"
	"hashbyte/rainbow"
)

//go:embed tables.bin
var embeddedTables []byte

const usage = `Activation Byte Recovery Tool

Usage:
  %s brute-force <sha1_hash>
  %s rainbow generate <output_file>
  %s rainbow lookup [-v|--verbose] [table_file] <sha1_hash>

Commands:
  brute-force    Search the full 2^32 keyspace (no tables needed)
  rainbow        Use pre-computed rainbow tables for fast lookup

Rainbow Subcommands:
  generate       Build rainbow tables and save to a binary file
  lookup         Look up a hash in pre-generated rainbow tables

Examples:
  %s brute-force 999a6ab85e8d1d...  
  %s rainbow generate tables.bin
  %s rainbow lookup 999a6ab85e8d1d... (uses embedded tables)
  %s rainbow lookup -v my_tables.bin 999a6ab85e8d1d... (uses my_tables.bin)
`

func printUsage() {
	name := os.Args[0]
	fmt.Printf(usage, name, name, name, name, name, name, name)
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "brute-force":
		if len(os.Args) != 3 {
			fmt.Printf("Usage: %s brute-force <sha1_hash>\n", os.Args[0])
			os.Exit(1)
		}
		brute.Run(os.Args[2])

	case "rainbow":
		if len(os.Args) < 3 {
			fmt.Printf("Usage:\n")
			fmt.Printf("  %s rainbow generate <output_file>\n", os.Args[0])
			fmt.Printf("  %s rainbow lookup [-v|--verbose] [table_file] <sha1_hash>\n", os.Args[0])
			os.Exit(1)
		}

		switch os.Args[2] {
		case "generate":
			if len(os.Args) != 4 {
				fmt.Printf("Usage: %s rainbow generate <output_file>\n", os.Args[0])
				os.Exit(1)
			}
			rainbow.Generate(os.Args[3])

		case "lookup":
			verbose := false
			args := os.Args[3:]
			var filteredArgs []string
			for _, arg := range args {
				if arg == "-v" || arg == "--verbose" {
					verbose = true
				} else {
					filteredArgs = append(filteredArgs, arg)
				}
			}

			if len(filteredArgs) < 1 || len(filteredArgs) > 2 {
				fmt.Printf("Usage: %s rainbow lookup [-v|--verbose] [table_file] <sha1_hash>\n", os.Args[0])
				os.Exit(1)
			}

			var fileData []byte
			var hashArg string

			if len(filteredArgs) == 1 {
				// Uses embedded tables
				fileData = embeddedTables
				hashArg = filteredArgs[0]
				if verbose {
					fmt.Println("Using embedded rainbow tables.")
				}
			} else {
				// Uses provided file
				filename := filteredArgs[0]
				hashArg = filteredArgs[1]
				if verbose {
					fmt.Printf("Loading custom rainbow tables from %s...\n", filename)
				}
				var err error
				fileData, err = os.ReadFile(filename)
				if err != nil {
					fmt.Printf("Error reading %s: %v\n", filename, err)
					os.Exit(1)
				}
			}

			if verbose {
				rainbow.RunLookupAndPrint(fileData, hashArg)
			} else {
				rainbow.RunLookup(fileData, hashArg)
			}

		default:
			fmt.Printf("Unknown rainbow subcommand: %s\n", os.Args[2])
			fmt.Printf("Available subcommands: generate, lookup\n")
			os.Exit(1)
		}

	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}
