package fuzzy

type scores struct {
	lineCount             float64
	commentRatio          float64
	functionCount         float64
	averageFunctionLength float64
	maxFunctionLength     float64
	conditionalCount      float64
	loopCount             float64
	maxNestingDepth       float64
	globalVariableCount   float64
	longLineCount         float64
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

func evaluateLineCount(x float64) float64 {
	return trapezoid(x, -1, 0, 200, 500)
}

func evaluateCommentRatio(x float64) float64 {
	return trapezoid(x, 0.05, 0.1, 0.25, 0.4)
}

func evaluateFunctionCount(x float64) float64 {
	return trapezoid(x, 1, 2, 10, 20)
}

func evaluateAverageFunctionLength(x float64) float64 {
	return trapezoid(x, -1, 0, 20, 50)
}

func evaluateMaxFunctionLength(x float64) float64 {
	return trapezoid(x, -1, 0, 30, 80)
}

func evaluateConditionalCount(x float64) float64 {
	return trapezoid(x, -1, 0, 5, 15)
}

func evaluateLoopCount(x float64) float64 {
	return trapezoid(x, -1, 0, 3, 8)
}

func evaluateMaxNestingDepth(x float64) float64 {
	return trapezoid(x, -1, 0, 2, 5)
}

func evaluateGlobalVariableCount(x float64) float64 {
	return trapezoid(x, -1, 0, 2, 6)
}

func evaluateLongLineCount(x float64) float64 {
	return trapezoid(x, -1, 0, 2, 8)
}

func computeScores(input Input) scores {
	var scores scores
	scores.lineCount = evaluateLineCount(float64(input.LineCount))
	scores.commentRatio = evaluateCommentRatio(input.CommentRatio)
	scores.functionCount = evaluateFunctionCount(float64(input.FunctionCount))
	scores.averageFunctionLength = evaluateAverageFunctionLength(input.AverageFunctionLength)
	scores.maxFunctionLength = evaluateMaxFunctionLength(float64(input.MaxFunctionLength))
	scores.conditionalCount = evaluateConditionalCount(float64(input.ConditionalCount))
	scores.loopCount = evaluateLoopCount(float64(input.LoopCount))
	scores.maxNestingDepth = evaluateMaxNestingDepth(float64(input.MaxNestingDepth))
	scores.globalVariableCount = evaluateGlobalVariableCount(float64(input.GlobalVariableCount))
	scores.longLineCount = evaluateLongLineCount(float64(input.LongLineCount))

	return scores
}
