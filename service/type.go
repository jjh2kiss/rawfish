package service

type Type int

const (
	SERVICETYPE_NORMAL = iota
	SERVICETYPE_RAW
)

func (self Type) IsNormalType() bool {
	if self == SERVICETYPE_NORMAL {
		return true
	}
	return false
}

func (self Type) IsRawType() bool {
	if self == SERVICETYPE_RAW {
		return true
	}
	return false
}
