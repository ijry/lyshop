package delivery_local

import (
	"strings"

	deliveryDriver "github.com/ijry/lyshop/core/driver/delivery"
)

type localDriver struct{}

func (d *localDriver) Name() string { return "local" }

func (d *localDriver) ValidateShipment(req deliveryDriver.ValidateReq) error {
	if strings.TrimSpace(req.RiderName) == "" {
		return deliveryDriver.ErrMissingRiderName
	}
	if strings.TrimSpace(req.RiderPhone) == "" {
		return deliveryDriver.ErrMissingRiderPhone
	}
	return nil
}

func (d *localDriver) NeedTrackingSync() bool { return false }
