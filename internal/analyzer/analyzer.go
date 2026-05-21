package analyzer

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

const longLineLimit = 120

type Metrics struct {
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

func Analyze(sourceCode string) (Metrics, error) {
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, "input.go", sourceCode, parser.ParseComments)
	if err != nil {
		return Metrics{}, err
	}

	var metrics Metrics
	lines := strings.Split(sourceCode, "\n")

	for _, line := range lines {
		if len(line) > longLineLimit {
			metrics.LongLineCount++
		}
	}

	metrics.LineCount = int64(len(lines))

	for _, commentGroup := range file.Comments {
		startLine := fileSet.Position(commentGroup.Pos()).Line
		endLine := fileSet.Position(commentGroup.End()).Line
		metrics.CommentLineCount += int64(endLine - startLine + 1)
	}

	if metrics.LineCount > 0 {
		metrics.CommentRatio = float64(metrics.CommentLineCount) / float64(metrics.LineCount)
	}

	sumOfLengths := 0
	for _, decl := range file.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			metrics.FunctionCount++

			startLine := fileSet.Position(funcDecl.Pos()).Line
			endLine := fileSet.Position(funcDecl.End()).Line
			length := endLine - startLine + 1
			if length > int(metrics.MaxFunctionLength) {
				metrics.MaxFunctionLength = int64(length)
			}

			sumOfLengths += length

			if funcDecl.Body != nil {
				depth := calculateMaxNestingDepth(funcDecl.Body.List, 0)
				if depth > int(metrics.MaxNestingDepth) {
					metrics.MaxNestingDepth = int64(depth)
				}
			}
		}

		if genDecl, ok := decl.(*ast.GenDecl); ok {
			if genDecl.Tok == token.VAR {
				for _, spec := range genDecl.Specs {
					if valueSpec, ok := spec.(*ast.ValueSpec); ok {
						metrics.GlobalVariableCount += int64(len(valueSpec.Names))
					}
				}
			}
		}
	}

	if metrics.FunctionCount > 0 {
		metrics.AverageFunctionLength = float64(sumOfLengths) / float64(metrics.FunctionCount)
	}

	ast.Inspect(file, func(node ast.Node) bool {
		switch node.(type) {
		case *ast.IfStmt, *ast.SwitchStmt, *ast.TypeSwitchStmt, *ast.SelectStmt:
			metrics.ConditionalCount++
		case *ast.ForStmt, *ast.RangeStmt:
			metrics.LoopCount++
		}

		return true
	})

	return metrics, nil
}

func calculateMaxNestingDepth(stmts []ast.Stmt, currentDepth int) int {
	maxDepth := currentDepth
	for _, stmt := range stmts {
		stmtDepth := calculateStmtMaxNestingDepth(stmt, currentDepth)
		maxDepth = max(maxDepth, stmtDepth)
	}

	return maxDepth
}

func calculateStmtMaxNestingDepth(stmt ast.Stmt, currentDepth int) int {
	switch stmt := stmt.(type) {
	case *ast.IfStmt:
		newDepth := currentDepth + 1
		maxDepth := newDepth
		bodyDepth := calculateMaxNestingDepth(stmt.Body.List, newDepth)
		maxDepth = max(maxDepth, bodyDepth)

		if stmt.Else != nil {
			switch elseStmt := stmt.Else.(type) {
			case *ast.BlockStmt:
				elseDepth := calculateMaxNestingDepth(elseStmt.List, newDepth)
				maxDepth = max(maxDepth, elseDepth)
			case *ast.IfStmt:
				elseIfDepth := calculateStmtMaxNestingDepth(elseStmt, currentDepth)
				maxDepth = max(maxDepth, elseIfDepth)
			}
		}
		return maxDepth
	case *ast.ForStmt:
		newDepth := currentDepth + 1
		maxDepth := newDepth
		bodyDepth := calculateMaxNestingDepth(stmt.Body.List, newDepth)
		maxDepth = max(maxDepth, bodyDepth)
		return maxDepth
	case *ast.RangeStmt:
		newDepth := currentDepth + 1
		maxDepth := newDepth
		bodyDepth := calculateMaxNestingDepth(stmt.Body.List, newDepth)
		maxDepth = max(maxDepth, bodyDepth)
		return maxDepth
	case *ast.SwitchStmt:
		newDepth := currentDepth + 1
		maxDepth := newDepth
		for _, switchBody := range stmt.Body.List {
			if caseClause, ok := switchBody.(*ast.CaseClause); ok {
				caseDepth := calculateMaxNestingDepth(caseClause.Body, newDepth)
				maxDepth = max(maxDepth, caseDepth)
			}
		}
		return maxDepth
	case *ast.TypeSwitchStmt:
		newDepth := currentDepth + 1
		maxDepth := newDepth
		for _, switchBody := range stmt.Body.List {
			if caseClause, ok := switchBody.(*ast.CaseClause); ok {
				caseDepth := calculateMaxNestingDepth(caseClause.Body, newDepth)
				maxDepth = max(maxDepth, caseDepth)
			}
		}
		return maxDepth
	case *ast.SelectStmt:
		newDepth := currentDepth + 1
		maxDepth := newDepth
		for _, selectBody := range stmt.Body.List {
			if commClause, ok := selectBody.(*ast.CommClause); ok {
				commDepth := calculateMaxNestingDepth(commClause.Body, newDepth)
				maxDepth = max(maxDepth, commDepth)
			}
		}
		return maxDepth
	}
	return currentDepth
}
