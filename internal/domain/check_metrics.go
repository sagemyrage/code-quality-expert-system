package domain

import "time"

type CheckMetrics struct {
	CheckID               int64
	LineCount             int64
	CommentLineCount      int64
	CommentRatio          float64
	FunctionCount         int64
	AverageFunctionLength float64
	MaxFunctionLength     int64
	ConditionalCount      int64
	LoopCount             int64
	MaxNestingDepth       int64
	GlobalVariableCount   int64
	LongLineCount         int64
	CreatedAt             time.Time
}
