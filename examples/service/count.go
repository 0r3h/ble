package service

import (
	"fmt"
	"log"
	"time"

	"github.com/currantlabs/bt"
	"github.com/currantlabs/bt/gatt"
	"github.com/currantlabs/bt/uuid"
)

// NewCountService ...
func NewCountService() bt.Service {
	n := 0
	s := gatt.NewService(uuid.MustParse("09fc95c0-c111-11e3-9904-0002a5d5c51b"))
	s.AddCharacteristic(uuid.MustParse("11fac9e0-c111-11e3-9246-0002a5d5c51b")).HandleRead(
		bt.ReadHandlerFunc(func(req bt.Request, rsp bt.ResponseWriter) {
			fmt.Fprintf(rsp, "count: Read %d", n)
			log.Printf("count: Read %d", n)
			n++
		}))

	s.AddCharacteristic(uuid.MustParse("16fe0d80-c111-11e3-b8c8-0002a5d5c51b")).HandleWrite(
		bt.WriteHandlerFunc(func(req bt.Request, rsp bt.ResponseWriter) {
			log.Printf("count: Wrote %s", string(req.Data()))
		}))

	s.AddCharacteristic(uuid.MustParse("1c927b50-c116-11e3-8a33-0800200c9a66")).HandleNotify(
		false, bt.NotifyHandlerFunc(func(req bt.Request, n bt.Notifier) {
			cnt := 0
			log.Printf("count: Notification subscribed")
			for {
				select {
				case <-n.Context().Done():
					log.Printf("count: Notification unsubscribed")
					return
				case <-time.After(time.Second):
					log.Printf("count: Notify: %d", cnt)
					if _, err := fmt.Fprintf(n, "Count: %d", cnt); err != nil {
						// Client disconnected prematurely before unsubscription.
						log.Printf("count: Failed to notify : %s", err)
						return
					}
					cnt++
				}
			}
		}))

	s.AddCharacteristic(uuid.MustParse("2da38c61-c116-11e3-8a33-0800200c9a66")).HandleNotify(
		true, bt.NotifyHandlerFunc(func(req bt.Request, n bt.Notifier) {
			cnt := 0
			log.Printf("count: Indication subscribed")
			for {
				select {
				case <-n.Context().Done():
					log.Printf("count: Indication unsubscribed")
					return
				case <-time.After(time.Second):
					log.Printf("count: Indicate: %d", cnt)
					if _, err := fmt.Fprintf(n, "Count: %d", cnt); err != nil {
						// Client disconnected prematurely before unsubscription.
						log.Printf("count: Failed to notify : %s", err)
						return
					}
					cnt++
				}
			}
		}))
	return s
}
