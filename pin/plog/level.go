package plog

type Level int8

const (
	LevelDbg      Level = 0
	LevelInfo     Level = 1
	LevelWarn     Level = 2
	LevelErr      Level = 3
	LevelDisabled Level = 4
)
