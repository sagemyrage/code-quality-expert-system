package analyzer

import (
	"testing"
)

func TestAnalyzeMaxNestingDepth(t *testing.T) {
	tests := []struct {
		name       string
		sourceCode string
		wantDepth  int64
	}{
		{
			name: "nested if",
			sourceCode: `
package main
			
func main() {
	if true {
		if true {
		}
	}
}
			`,
			wantDepth: 2,
		},
		{
			name: "else if chain",
			sourceCode: `
package main

func main() {
    if true {
    } else if false {
    } else if true {
    }
}
			`,
			wantDepth: 1,
		},
		{
			name: "if inside else block",
			sourceCode: `
package main

func main() {
    if true {
    } else {
    	if false {
    	}
	}
}
			`,
			wantDepth: 2,
		},
		{
			name: "for inside if",
			sourceCode: `
package main

func main() {
    if true {
		for i := 0; i < 10; i++ {
		}
    } 
}
			`,
			wantDepth: 2,
		},
		{
			name: "switch with if inside case",
			sourceCode: `
package main

func main() {
    switch 1 {
	case 1:
		if true {
		}
	}
}
			`,
			wantDepth: 2,
		},
		{
			name: "select with if inside case",
			sourceCode: `
package main

func main() {
	ch := make(chan int)
    			
	select {
	case <-ch:
		if true {
		}
	default:
	}
}
			`,
			wantDepth: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metrics, err := Analyze(tt.sourceCode)
			if err != nil {
				t.Fatalf("Analyze() error = %v", err)
			}

			if metrics.MaxNestingDepth != tt.wantDepth {
				t.Errorf("MaxNestingDepth = %d, want = %d", metrics.MaxNestingDepth, tt.wantDepth)
			}
		})
	}
}

func TestAnalyzeBasicMetrics(t *testing.T) {
	test := struct {
		name                    string
		sourceCode              string
		wantFunctionCount       int64
		wantConditionalCount    int64
		wantLoopCount           int64
		wantGlobalVariableCount int64
		wantLongLineCount       int64
		wantCommentLineCount    int64
	}{
		name: "basic metrics",
		sourceCode: `
package main

import "fmt"

// Глобальные переменные
var global1 = 10
var global2 = "hello"

// Функция main
func main() {
	var input int
	fmt.Scan(&input)

	switch input {
	case 1:
		for i := 0; i < 100; i++ {
		}
	case 2:
		s := []int{1,2,3,4,5}
		for k, v := range s {
		}
	default:
		if input > global1 {
			panic("паника")
		}
	}
}

// Функция anotherMain и длинная строка
func anotherMainAndVeryyyyyyyyyyyyyLoooooooooooooooooooooooooongStriiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiinggggggggggggggggggggggg() {
	if true {
			fmt.Println(global2)
	}
}
		`,
		wantFunctionCount:       2,
		wantConditionalCount:    3,
		wantLoopCount:           2,
		wantGlobalVariableCount: 2,
		wantLongLineCount:       1,
		wantCommentLineCount:    3,
	}

	metrics, err := Analyze(test.sourceCode)
	if err != nil {
		t.Fatalf("Analyze() error: %v", err)
	}

	if metrics.FunctionCount != test.wantFunctionCount {
		t.Errorf("FunctionCount = %d, want = %d", metrics.FunctionCount, test.wantFunctionCount)
	}
	if metrics.ConditionalCount != test.wantConditionalCount {
		t.Errorf("ConditionalCount = %d, want = %d", metrics.ConditionalCount, test.wantConditionalCount)
	}
	if metrics.LoopCount != test.wantLoopCount {
		t.Errorf("LoopCount = %d, want = %d", metrics.LoopCount, test.wantLoopCount)
	}
	if metrics.GlobalVariableCount != test.wantGlobalVariableCount {
		t.Errorf("GlobalVariableCount = %d, want = %d", metrics.GlobalVariableCount, test.wantGlobalVariableCount)
	}
	if metrics.LongLineCount != test.wantLongLineCount {
		t.Errorf("LongLineCount = %d, want = %d", metrics.LongLineCount, test.wantLongLineCount)
	}
	if metrics.CommentLineCount != test.wantCommentLineCount {
		t.Errorf("CommentLineCount = %d, want = %d", metrics.CommentLineCount, test.wantCommentLineCount)
	}
}
