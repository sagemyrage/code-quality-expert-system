package fuzzy

import "errors"

type Evaluation struct {
	Score           float64
	Summary         string
	Level           string
	Recommendations []string
}

type Input struct {
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
}

func trapezoid(x, a, b, c, d float64) float64 {
	if x >= b && x <= c {
		return 1
	}
	if x > a && x < b {
		return (x - a) / (b - a)
	}
	if x > c && x < d {
		return (d - x) / (d - c)
	}

	return 0
}

func Evaluate(input Input) (Evaluation, error) {
	if input.LineCount < 0 ||
		input.CommentLineCount < 0 ||
		input.CommentRatio < 0 ||
		input.CommentRatio > 1 ||
		input.FunctionCount < 0 ||
		input.AverageFunctionLength < 0 ||
		input.MaxFunctionLength < 0 ||
		input.ConditionalCount < 0 ||
		input.LoopCount < 0 ||
		input.MaxNestingDepth < 0 ||
		input.GlobalVariableCount < 0 ||
		input.LongLineCount < 0 {
		return Evaluation{}, errors.New("invalid fuzzy input")
	}

	return Evaluation{}, nil
}
