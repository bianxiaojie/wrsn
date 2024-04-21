package state

type WRSNEntityState int64

const (
	None WRSNEntityState = iota
	Charging
	Charged
	SwitchingBattery
	Moving
	Sensing
)
