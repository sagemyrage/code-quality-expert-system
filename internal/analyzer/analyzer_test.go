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
