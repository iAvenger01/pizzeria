package errors

import "errors"

var (
	ErrToGetCourier               = errors.New("failed to get courier")
	ErrToParseCourierInStruct     = errors.New("failed to parse courier")
	ErrToGetAllCouriers           = errors.New("failed to get couriers")
	ErrToParseAllCouriersInStruct = errors.New("failed to parse couriers in struct")
	ErrToCreateEmployee           = errors.New("failed to create employee")
	ErrToCreateCourier            = errors.New("failed to create courier")
)
