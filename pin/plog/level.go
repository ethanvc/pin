package plog

type Level int8

const (
	LevelDbg      Level = 0
	LevelInfo     Level = 1
	LevelWarn     Level = 2
	LevelErr      Level = 3
	LevelDisabled Level = 4
)

func (l Level) String() string {
	switch l {
	case LevelDbg:
		return "PDBG"
	case LevelInfo:
		return "PINFO"
	case LevelWarn:
		return "PWARN"
	case LevelErr:
		return "PERR"
	default:
		return "PUnknown"
	}
}
