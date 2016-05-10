package gatt

import "github.com/currantlabs/bt/uuid"

// DefaultMTU defines the default MTU of ATT protocol excluding ATT header.
const DefaultMTU = 23 - 3

// MaxMTU is maximum of ATT_MTU, which is 512 bytes of value length excluding ATT header.
// The maximum length of an attribute value shall be 512 octets [Vol 3, Part F, 3.2.9]
const MaxMTU = 512

var (
	attrGAPUUID  = uuid.UUID16(0x1800)
	attrGATTUUID = uuid.UUID16(0x1801)

	attrPrimaryServiceUUID   = uuid.UUID16(0x2800)
	attrSecondaryServiceUUID = uuid.UUID16(0x2801)
	attrIncludeUUID          = uuid.UUID16(0x2802)
	attrCharacteristicUUID   = uuid.UUID16(0x2803)

	attrClientCharacteristicConfigUUID = uuid.UUID16(0x2902)
	attrServerCharacteristicConfigUUID = uuid.UUID16(0x2903)

	attrDeviceNameUUID        = uuid.UUID16(0x2A00)
	attrAppearanceUUID        = uuid.UUID16(0x2A01)
	attrPeripheralPrivacyUUID = uuid.UUID16(0x2A02)
	attrReconnectionAddrUUID  = uuid.UUID16(0x2A03)
	attrPeferredParamsUUID    = uuid.UUID16(0x2A04)
	attrServiceChangedUUID    = uuid.UUID16(0x2A05)
)

const (
	cccNotify   = 0x0001
	cccIndicate = 0x0002
)
