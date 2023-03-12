package plog

type Handler func(logger *Logger, record Record)
