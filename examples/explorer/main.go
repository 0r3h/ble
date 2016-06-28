package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/currantlabs/ble/gatt"
	"github.com/currantlabs/x/io/bt"
)

var (
	name = flag.String("name", "Gopher", "name of remote peripheral")
	addr = flag.String("addr", "", "address of remote peripheral (MAC on Linux, UUID on OS X)")
	sub  = flag.Duration("sub", 0, "subscribe to notification and indication for a specified period")
)

// matcher returns true if the advertisement matches our search criteria.
type matcher func(a bt.Advertisement) bool

// explorer connects to a remote peripheral and explores its GATT server.
type explorerer struct {
	bt.Central

	match matcher
	ch    chan bt.Advertisement
}

func (e *explorerer) Handle(a bt.Advertisement) {
	if e.match(a) {
		e.StopScanning()
		e.ch <- a
	}
}

func explorer(cln bt.Client) error {
	l := log.New(os.Stdout, "["+cln.Address().String()+"] ", log.Lmicroseconds)

	ss, err := cln.DiscoverServices(nil)
	if err != nil {
		return fmt.Errorf("can't discover services: %s\n", err)
	}
	for _, s := range ss {
		l.Printf("Service: %s %s\n", s.UUID.String(), bt.Name(s.UUID))

		cs, err := cln.DiscoverCharacteristics(nil, s)
		if err != nil {
			return fmt.Errorf("can't discover characteristics: %s\n", err)
		}
		for _, c := range cs {
			l.Printf("  Characteristic: %s, Property: 0x%02X, %s\n", c.UUID, c.Property, bt.Name(c.UUID))
			if (c.Property & bt.CharRead) != 0 {
				b, err := cln.ReadCharacteristic(c)
				if err != nil {
					l.Printf("Failed to read characteristic: %s\n", err)
					continue
				}
				l.Printf("    Value         %x | %q\n", b, b)
			}

			for _, c := range cs {
				ds, err := cln.DiscoverDescriptors(nil, c)
				if err != nil {
					return fmt.Errorf("can't discover descriptors: %s\n", err)
				}
				for _, d := range ds {
					l.Printf("    Descriptor: %s, %s\n", d.UUID, bt.Name(d.UUID))
					b, err := cln.ReadDescriptor(d)
					if err != nil {
						l.Printf("Failed to read descriptor: %s\n", err)
						continue
					}
					l.Printf("    Value         %x | %q\n", b, b)
				}
				if *sub != 0 {
					if (c.Property & bt.CharNotify) != 0 {
						h := func(req []byte) { l.Printf("Notified: %q [ % X ]", string(req), req) }
						cln.Subscribe(c, false, h)
						time.Sleep(*sub)
						cln.Unsubscribe(c, false)
					}
					if (c.Property & bt.CharIndicate) != 0 {
						h := func(req []byte) { l.Printf("Indicated: %q [ % X ]", string(req), req) }
						cln.Subscribe(c, true, h)
						time.Sleep(*sub)
						cln.Unsubscribe(c, true)
					}
				}
			}

		}
		l.Printf("\n")
	}
	return nil
}

func main() {
	flag.Parse()

	// Default to search device with name of Gopher (or specified by user).
	match := func(a bt.Advertisement) bool {
		return strings.ToUpper(a.LocalName()) == strings.ToUpper(*name)
	}

	// If addr is specified, search for addr instead.
	if len(*addr) != 0 {
		match = func(a bt.Advertisement) bool {
			return strings.ToUpper(a.Address().String()) == strings.ToUpper(*addr)
		}
	}

	dev, err := gatt.NewCentral()
	if err != nil {
		log.Fatalf("can't create central: %s", err)
	}
	exp := &explorerer{
		Central: dev,
		ch:      make(chan bt.Advertisement),
		match:   match,
	}

	if err = dev.SetAdvHandler(exp); err != nil {
		log.Fatalf("can't set adv handler: %s", err)
	}

	if err = dev.Scan(false); err != nil {
		log.Fatalf("can't scan: %s", err)
	}

	// Wait for the exploration is done.
	a := <-exp.ch

	// Dial connects to the remote device.
	c, err := exp.Dial(a.Address())
	if err != nil {
		log.Fatalf("can't dial: %s", err)
	}

	// Create and attach a GATT client to the connection.
	cln, err := gatt.NewClient(c)
	if err != nil {
		log.Fatalf("can't create client: %s", err)
	}

	// Start the exploration.
	explorer(cln)

	// Disconnect the connection. (On OS X, this might take a while.)
	log.Printf("Disconnecting... ")
	cln.CancelConnection()
}
