package main

import (
	"flag"
	"log"

	"github.com/mewmew/psi/emu"
	"github.com/pkg/errors"
)

func main() {
	// Parse command line arguments.
	var (
		biosPath string
	)
	// TODO: Use Japanese BIOS SCPH5500.BIN.
	flag.StringVar(&biosPath, "bios", "SCPH1001.BIN", "PlayStation 1 BIOS path")
	flag.Parse()
	// Start emulator.
	if err := psi(biosPath); err != nil {
		log.Fatalf("%+v", err)
	}
}

// psi initiates the emulator by loading the given BIOS.
func psi(biosPath string) error {
	bios, err := emu.LoadBIOS(biosPath)
	if err != nil {
		return errors.WithStack(err)
	}
	hwregs := emu.NewHWRegs()
	cache := emu.NewCacheControl()
	ram := emu.NewRAM()
	cpu := emu.NewCPU(bios, hwregs, cache, ram)
	for {
		cpu.Step()
		//fmt.Println(cpu)
	}
	return nil
}
