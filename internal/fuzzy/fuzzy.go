package fuzzy

import "errors"

const (
	weightLineCount             = 1
	weightCommentRatio          = 2
	weightFunctionCount         = 1
	weightAverageFunctionLength = 2
	weightMaxFunctionLength     = 2
	weightConditionalCount      = 1
	weightLoopCount             = 1
	weightMaxNestingDepth       = 3
	weightGlobalVariableCount   = 2
	weightLongLineCount         = 1
	totalWeight                 = weightLineCount +
		weightCommentRatio +
		weightFunctionCount +
		weightAverageFunctionLength +
		weightMaxFunctionLength +
		weightConditionalCount +
		weightLoopCount +
		weightMaxNestingDepth +
		weightGlobalVariableCount +
		weightLongLineCount
)

type Evaluation struct {
	Score           float64
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

	var result Evaluation
	s := computeScores(input)
	sumOfScores := s.lineCount*weightLineCount +
		s.commentRatio*weightCommentRatio +
		s.functionCount*weightFunctionCount +
		s.averageFunctionLength*weightAverageFunctionLength +
		s.maxFunctionLength*weightMaxFunctionLength +
		s.conditionalCount*weightConditionalCount +
		s.loopCount*weightLoopCount +
		s.maxNestingDepth*weightMaxNestingDepth +
		s.globalVariableCount*weightGlobalVariableCount +
		s.longLineCount*weightLongLineCount
	result.Score = (sumOfScores / totalWeight) * 100
	result.Level = classifyLevel(result.Score)
	result.Recommendations = buildRecommendations(input, s)

	return result, nil
}
