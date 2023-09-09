package utils

type Marker string

const (
	NONE    Marker = "X"
	CLIENT  Marker = "C"
	SERVER  Marker = "S"
	CONTROL Marker = "M" // for Metadata
	SEP     string = ";"
)
