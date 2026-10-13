package esepunittests

import "testing"

func TestGetGradeA(t *testing.T) {
	expected_value := "A"

	gradeCalculator := NewGradeCalculator()

	gradeCalculator.AddGrade("open source assignment", 100, Assignment)
	gradeCalculator.AddGrade("exam 1", 100, Exam)
	gradeCalculator.AddGrade("essay on ai ethics", 100, Essay)

	actual_value := gradeCalculator.GetFinalGrade()

	if expected_value != actual_value {
		t.Errorf("Expected GetGrade to return '%s'; got '%s' instead", expected_value, actual_value)
	}
}

func TestGetGradeB(t *testing.T) {
	expected_value := "B"

	gradeCalculator := NewGradeCalculator()

	gradeCalculator.AddGrade("open source assignment", 80, Assignment)
	gradeCalculator.AddGrade("exam 1", 81, Exam)
	gradeCalculator.AddGrade("essay on ai ethics", 85, Essay)

	actual_value := gradeCalculator.GetFinalGrade()

	if expected_value != actual_value {
		t.Errorf("Expected GetGrade to return '%s'; got '%s' instead", expected_value, actual_value)
	}
}

func TestGetGradeF(t *testing.T) {
	expected_value := "F"

	gradeCalculator := NewGradeCalculator()

	gradeCalculator.AddGrade("open source assignment", 50, Assignment)
	gradeCalculator.AddGrade("exam 1", 55, Exam)
	gradeCalculator.AddGrade("essay on ai ethics", 45, Essay)

	actual_value := gradeCalculator.GetFinalGrade()

	if expected_value != actual_value {
		t.Errorf("Expected GetGrade to return '%s'; got '%s' instead", expected_value, actual_value)
	}
}

// score is one AddGrade call in a table-driven test.
type score struct {
	grade     int
	gradeType GradeType
}

// finalGradeFor builds a calculator from the given scores and returns its
// final grade.
func finalGradeFor(scores []score) string {
	gradeCalculator := NewGradeCalculator()

	for _, s := range scores {
		gradeCalculator.AddGrade(s.gradeType.String(), s.grade, s.gradeType)
	}

	return gradeCalculator.GetFinalGrade()
}

// oneOfEach returns one score per category.
func oneOfEach(assignment, exam, essay int) []score {
	return []score{
		{assignment, Assignment},
		{exam, Exam},
		{essay, Essay},
	}
}

func TestGetFinalGradeLetterGrades(t *testing.T) {
	tests := []struct {
		name     string
		scores   []score
		expected string
	}{
		{"A", oneOfEach(95, 95, 95), "A"},
		{"B", oneOfEach(85, 85, 85), "B"},
		{"C", oneOfEach(75, 75, 75), "C"},
		{"D", oneOfEach(65, 65, 65), "D"},
		{"F", oneOfEach(55, 55, 55), "F"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := finalGradeFor(tc.scores)

			if tc.expected != actual {
				t.Errorf("Expected GetFinalGrade to return '%s'; got '%s' instead", tc.expected, actual)
			}
		})
	}
}

func TestGetFinalGradeBoundaries(t *testing.T) {
	tests := []struct {
		name     string
		score    int
		expected string
	}{
		{"100 is the top of A", 100, "A"},
		{"90 is the bottom of A", 90, "A"},
		{"89 is the top of B", 89, "B"},
		{"80 is the bottom of B", 80, "B"},
		{"79 is the top of C", 79, "C"},
		{"70 is the bottom of C", 70, "C"},
		{"69 is the top of D", 69, "D"},
		{"60 is the bottom of D", 60, "D"},
		{"59 is the top of F", 59, "F"},
		{"0 is the bottom of F", 0, "F"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := finalGradeFor(oneOfEach(tc.score, tc.score, tc.score))

			if tc.expected != actual {
				t.Errorf("Expected a score of %d to return '%s'; got '%s' instead", tc.score, tc.expected, actual)
			}
		})
	}
}

func TestGetFinalGradeWeights(t *testing.T) {
	tests := []struct {
		name       string
		assignment int
		exam       int
		essay      int
		expected   string
	}{
		// Each pair sits on either side of 90 and only lands there if the
		// category carries exactly its required weight.
		{"assignments are 50 percent: 40+35+15 = 90", 80, 100, 100, "A"},
		{"assignments are 50 percent: 39.5+35+15 = 89.5", 79, 100, 100, "B"},
		{"exams are 35 percent: 50+25.2+15 = 90.2", 100, 72, 100, "A"},
		{"exams are 35 percent: 50+24.85+15 = 89.85", 100, 71, 100, "B"},
		{"essays are 15 percent: 50+35+5.1 = 90.1", 100, 100, 34, "A"},
		{"essays are 15 percent: 50+35+4.95 = 89.95", 100, 100, 33, "B"},

		{"high assignments outweigh low exams and essays: 50+21+9 = 80", 100, 60, 60, "B"},
		{"high exams do not outweigh low assignments: 30+35+9 = 74", 60, 100, 60, "C"},
		{"high essays barely help: 30+21+15 = 66", 60, 60, 100, "D"},
		{"a fractional weighted grade is not rounded up: 89.65", 90, 89, 90, "B"},
		{"an exact 70 is a C: 30.5+33.95+5.55 = 70", 61, 97, 37, "C"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := finalGradeFor(oneOfEach(tc.assignment, tc.exam, tc.essay))

			if tc.expected != actual {
				t.Errorf("Expected assignment %d, exam %d, essay %d to return '%s'; got '%s' instead",
					tc.assignment, tc.exam, tc.essay, tc.expected, actual)
			}
		})
	}
}

func TestGetFinalGradeMultipleGradesPerCategory(t *testing.T) {
	tests := []struct {
		name     string
		scores   []score
		expected string
	}{
		{
			// averages 90, 90, 90
			name: "grades in a category are averaged",
			scores: []score{
				{100, Assignment}, {80, Assignment},
				{95, Exam}, {85, Exam},
				{90, Essay}, {90, Essay}, {90, Essay},
			},
			expected: "A",
		},
		{
			// averages 80, 70, 100: 40+24.5+15 = 79.5
			name: "categories can hold different numbers of grades",
			scores: []score{
				{70, Assignment}, {90, Assignment}, {80, Assignment},
				{60, Exam}, {80, Exam},
				{100, Essay},
			},
			expected: "C",
		},
		{
			// same grades as above, added in a different order
			name: "the order grades are added in does not matter",
			scores: []score{
				{100, Essay}, {80, Exam}, {70, Assignment},
				{90, Assignment}, {60, Exam}, {80, Assignment},
			},
			expected: "C",
		},
		{
			// averages 50, 100, 100: 25+35+15 = 75
			name: "one low grade pulls down only its own category",
			scores: []score{
				{100, Assignment}, {0, Assignment},
				{100, Exam}, {100, Exam},
				{100, Essay}, {100, Essay},
			},
			expected: "C",
		},
		{
			// pooling all six grades would average 360/6 = 60, a D, but the
			// weighted grade is 50+17.5+9 = 76.5
			name: "categories are weighted rather than pooled",
			scores: []score{
				{100, Assignment},
				{50, Exam}, {50, Exam}, {50, Exam}, {50, Exam},
				{60, Essay},
			},
			expected: "C",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := finalGradeFor(tc.scores)

			if tc.expected != actual {
				t.Errorf("Expected GetFinalGrade to return '%s'; got '%s' instead", tc.expected, actual)
			}
		})
	}
}

func TestGetFinalGradeEmptyCategories(t *testing.T) {
	tests := []struct {
		name     string
		scores   []score
		expected string
	}{
		{"no grades at all", nil, "F"},
		{"no essays: 50+35+0 = 85", []score{{100, Assignment}, {100, Exam}}, "B"},
		{"no exams: 50+0+15 = 65", []score{{100, Assignment}, {100, Essay}}, "D"},
		{"no assignments: 0+35+15 = 50", []score{{100, Exam}, {100, Essay}}, "F"},
		{"only assignments: 50", []score{{100, Assignment}}, "F"},
		{"only exams: 35", []score{{100, Exam}}, "F"},
		{"only essays: 15", []score{{100, Essay}}, "F"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := finalGradeFor(tc.scores)

			if tc.expected != actual {
				t.Errorf("Expected GetFinalGrade to return '%s'; got '%s' instead", tc.expected, actual)
			}
		})
	}
}

func TestAddGradeIgnoresUnknownGradeType(t *testing.T) {
	expected_value := "A"

	gradeCalculator := NewGradeCalculator()

	gradeCalculator.AddGrade("open source assignment", 100, Assignment)
	gradeCalculator.AddGrade("exam 1", 100, Exam)
	gradeCalculator.AddGrade("essay on ai ethics", 100, Essay)
	gradeCalculator.AddGrade("not a real category", 0, GradeType(99))

	actual_value := gradeCalculator.GetFinalGrade()

	if expected_value != actual_value {
		t.Errorf("Expected a grade with an unknown type to be ignored and return '%s'; got '%s' instead", expected_value, actual_value)
	}
}

func TestGradeTypeString(t *testing.T) {
	tests := []struct {
		gradeType GradeType
		expected  string
	}{
		{Assignment, "assignment"},
		{Exam, "exam"},
		{Essay, "essay"},
		{GradeType(99), ""},
	}

	for _, tc := range tests {
		actual := tc.gradeType.String()

		if tc.expected != actual {
			t.Errorf("Expected GradeType(%d).String() to return '%s'; got '%s' instead", int(tc.gradeType), tc.expected, actual)
		}
	}
}
