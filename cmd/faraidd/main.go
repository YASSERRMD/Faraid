// Command faraidd is the HTTP server entrypoint for Faraid, an Islamic
// inheritance (ilm al-faraid) calculation system.
//
// The deterministic legal engine lives under internal/core and never performs
// I/O or talks to an LLM. This entrypoint is expanded in later phases to load
// configuration, set up structured logging, and serve the HTTP API.
package main

import "fmt"

func main() {
	fmt.Println("faraidd: server entrypoint, not yet wired (see later phases)")
}
