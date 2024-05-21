package types

type Coords struct {
	Lat float64
	Lon float64
}
type OTUStates struct {
	OTUFirstState *OTU
	OTULastState  *OTU
	Distance      float64
	Cash          float64
}
type OTU struct {
	OTUID  string
	Coords Coords
}
