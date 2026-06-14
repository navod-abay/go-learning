package engine

import (
	"os"
	"testing"
)

func TestMandelbrotSetZero(t *testing.T) {
	tests := []struct {
		name     string
		num      complex128
		max_iter int
		expected bool
	}{
		{"zero", complex(0.0, 0.0), 1000, true},
		{"one", complex(1.0, 0.0), 1000, false},
		{"-one", complex(-1.0, 0.0), 1000, true},
		{"i", complex(0.0, 1.0), 1000, true},
		{"-i", complex(0.0, -1.0), 1000, true},
		{"-0.5", complex(-0.5, 0), 1000, true},
		{"-0.5 + 0.8i", complex(-0.5, 0.8), 1000, false},
		{"-0.25 + 0.8i", complex(-0.25, 0.8), 1000, false},
		{"0.25,0.8i", complex(-.25, 0.8), 1000, false},
		{"0.25 + 0.65i", complex(0.25, 0.65), 1000, false},
		{"0.25 + 0.55i", complex(0.25, 0.65), 1000, true},
		{"0.28 + 0.51i", complex(0.28, 0.51), 1000, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := checkMandelbrotSetInclusion(tt.num, tt.max_iter)
			if actual != tt.expected {
				t.Errorf("checkMandelbrotSet: Expencted %v, got %v", tt.expected, actual)
			}
		})
	}
}
func TestWriteToCSV(t *testing.T) {
	testArray := [][]Pixel{
		{{complex(-1, 1), false}, {complex(0, 1), true}, {complex(1, 1), false}},
		{{complex(-1, 0), true}, {complex(0, 0), true}, {complex(1, 0), false}},
		{{complex(-1, -1), false}, {complex(0, -1), true}, {complex(1, -1), false}},
	}
	WriteToCSV(testArray)
	file, err := os.ReadFile("output.csv")
	if err != nil {
		t.Errorf("Couldn't open file. %v", err)
	}
	csvString := string(file)
	if csvString != "0,255,0,\n255,255,0,\n0,255,0,\n" {
		t.Error("Written file is wrong")
		t.Errorf("Got: %v", csvString)
	}
}
