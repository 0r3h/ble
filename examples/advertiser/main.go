package main

import (
	"fmt"
	"log"

	"github.com/currantlabs/bt/adv"
	"github.com/currantlabs/bt/hci"
)

func main() {
	dev := new(hci.HCI)
	if err := dev.Init(-1); err != nil {
		log.Fatalf("can't open HCI device: %s", err)
	}

	// Craft a simple advertising data packet.
	ad := adv.Packet(nil)
	ad = ad.AppendFlags(adv.FlagGeneralDiscoverable | adv.FlagLEOnly)
	ad = ad.AppendCompleteName("Gopher")

	// Set advertising data
	if err := dev.SetAdvertisement(ad, nil); err != nil {
		log.Fatalf("can't set advertisement: %s", err)
	}

	// Set advertising parameter
	if err := dev.SetAdvParams(hci.AdvParams{
		AdvertisingIntervalMin:  0x0020,    // 0x0020 - 0x4000; N * 0.625 msec
		AdvertisingIntervalMax:  0x0020,    // 0x0020 - 0x4000; N * 0.625 msec
		AdvertisingType:         0x00,      // 00: ADV_IND, 0x01: DIRECT(HIGH), 0x02: SCAN, 0x03: NONCONN, 0x04: DIRECT(LOW)
		OwnAddressType:          0x00,      // 0x00: public, 0x01: random
		DirectAddressType:       0x00,      // 0x00: public, 0x01: random
		DirectAddress:           [6]byte{}, // Public or Random Address of the Device to be connected
		AdvertisingChannelMap:   0x7,       // 0x07 0x01: ch37, 0x2: ch38, 0x4: ch39
		AdvertisingFilterPolicy: 0x00,
	}); err != nil {
		log.Fatalf("can't set advertising parameters: %s", err)
	}

	// Start advertising
	if err := dev.Advertise(); err != nil {
		log.Fatalf("can't advertise: %s", err)
	}

	fmt.Printf("Advertising...\n")
	select {} // Prevent program from exiting
}
