package delivery_express

import (
	"strings"

	deliveryDriver "github.com/ijry/lyshop/core/driver/delivery"
)

type expressDriver struct{}

func (d *expressDriver) Name() string { return "express" }

func (d *expressDriver) ValidateShipment(req deliveryDriver.ValidateReq) error {
	if strings.TrimSpace(req.Company) == "" {
		return deliveryDriver.ErrMissingCompany
	}
	if strings.TrimSpace(req.TrackingNo) == "" {
		return deliveryDriver.ErrMissingTrackingNo
	}
	return nil
}

func (d *expressDriver) NeedTrackingSync() bool { return true }
