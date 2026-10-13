package esepunittests

import "testing"

// Pass/Fail grading is selected when the calculator is created:
//
//	NewGradeCalculator()                    letter grades (unchanged default)
//	NewGradeCalculatorWithMode(LetterGrade) letter grades
//	NewGradeCalculatorWithMode(PassFail)    "Pass" or "Fail"
//
// In PassFail mode GetFinalGrade returns "Pass" for a numerical grade of 70
// or above (a C or better) and "Fail" for anything below 70.

// finalGradeWithMode builds a calculator in the given mode from the given
// scores and returns its final grade.
func finalGradeWithMode(mode GradingMode, scores []score) string {
	gradeCalculator := NewGradeCalculatorWithMode(mode)

	for _, s := range scores {
		gradeCalculator.AddGrade(s.gradeType.String(), s.grade, s.gradeType)
	}

	return gradeCalculator.GetFinalGrade()
}

func TestPassFailGrades(t *testing.T) {
	tests := []struct {
		name     string
		scores   []score
		expected string
	}{
		{"an A passes", oneOfEach(95, 95, 95), "Pass"},
		{"a B passes", oneOfEach(85, 85, 85), "Pass"},
		{"a C passes", oneOfEach(75, 75, 75), "Pass"},
		{"a D fails", oneOfEach(65, 65, 65), "Fail"},
		{"an F fails", oneOfEach(30, 30, 30), "Fail"},
		{"no grades fails", nil, "Fail"},
		{"a missing category can fail strong grades: 50+0+15 = 65", []score{{100, Assignment}, {100, Essay}}, "Fail"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := finalGradeWithMode(PassFail, tc.scores)

			if tc.expected != actual {
				t.Errorf("Expected GetFinalGrade to return '%s'; got '%s' instead", tc.expected, actual)
			}
		})
	}
}

func TestPassFailBoundary(t *testing.T) {
	tests := []struct {
		name       string
		assignment int
		exam       int
		essay      int
		expected   string
	}{
		{"71 passes", 71, 71, 71, "Pass"},
		{"70 passes", 70, 70, 70, "Pass"},
		{"69 fails", 69, 69, 69, "Fail"},
		{"a weighted grade of exactly 70 passes: 30.5+33.95+5.55", 61, 97, 37, "Pass"},
		{"a weighted grade just under 70 fails: 35+24.15+10.5 = 69.65", 70, 69, 70, "Fail"},
		{"100 passes", 100, 100, 100, "Pass"},
		{"0 fails", 0, 0, 0, "Fail"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := finalGradeWithMode(PassFail, oneOfEach(tc.assignment, tc.exam, tc.essay))

			if tc.expected != actual {
				t.Errorf("Expected assignment %d, exam %d, essay %d to return '%s'; got '%s' instead",
					tc.assignment, tc.exam, tc.essay, tc.expected, actual)
			}
		})
	}
}

func TestLetterGradeModeIsUnchanged(t *testing.T) {
	tests := []struct {
		score    int
		expected string
	}{
		{95, "A"},
		{90, "A"},
		{89, "B"},
		{80, "B"},
		{79, "C"},
		{70, "C"},
		{69, "D"},
		{60, "D"},
		{59, "F"},
	}

	for _, tc := range tests {
		scores := oneOfEach(tc.score, tc.score, tc.score)

		// the original constructor still produces letter grades
		actual := finalGradeFor(scores)

		if tc.expected != actual {
			t.Errorf("Expected NewGradeCalculator with a score of %d to return '%s'; got '%s' instead", tc.score, tc.expected, actual)
		}

		// selecting LetterGrade explicitly behaves the same way
		actual = finalGradeWithMode(LetterGrade, scores)

		if tc.expected != actual {
			t.Errorf("Expected LetterGrade mode with a score of %d to return '%s'; got '%s' instead", tc.score, tc.expected, actual)
		}
	}
}

func TestUnknownGradingModeUsesLetterGrades(t *testing.T) {
	expected_value := "C"

	actual_value := finalGradeWithMode(GradingMode(99), oneOfEach(75, 75, 75))

	if expected_value != actual_value {
		t.Errorf("Expected an unknown grading mode to fall back to letter grades and return '%s'; got '%s' instead", expected_value, actual_value)
	}
}
