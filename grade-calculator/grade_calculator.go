package esepunittests

type GradeCalculator struct {
	assignments []Grade
	exams       []Grade
	essays      []Grade
}

type GradeType int

const (
	Assignment GradeType = iota
	Exam
	Essay
)

var gradeTypeName = map[GradeType]string{
	Assignment: "assignment",
	Exam:       "exam",
	Essay:      "essay",
}

func (gt GradeType) String() string {
	return gradeTypeName[gt]
}

type Grade struct {
	Name  string
	Grade int
	Type  GradeType
}

func NewGradeCalculator() *GradeCalculator {
	return &GradeCalculator{
		assignments: make([]Grade, 0),
		exams:       make([]Grade, 0),
		essays:      make([]Grade, 0),
	}
}

func (gc *GradeCalculator) GetFinalGrade() string {
	numericalGrade := gc.calculateNumericalGrade()

	if numericalGrade >= 90 {
		return "A"
	} else if numericalGrade >= 80 {
		return "B"
	} else if numericalGrade >= 70 {
		return "C"
	} else if numericalGrade >= 60 {
		return "D"
	}

	return "F"
}

func (gc *GradeCalculator) AddGrade(name string, grade int, gradeType GradeType) {
	switch gradeType {
	case Assignment:
		gc.assignments = append(gc.assignments, Grade{
			Name:  name,
			Grade: grade,
			Type:  Assignment,
		})
	case Exam:
		gc.exams = append(gc.exams, Grade{
			Name:  name,
			Grade: grade,
			Type:  Exam,
		})
	case Essay:
		gc.essays = append(gc.essays, Grade{
			Name:  name,
			Grade: grade,
			Type:  Essay,
		})
	}
}

func (gc *GradeCalculator) calculateNumericalGrade() int {
	assignment_average := computeAverage(gc.assignments)
	exam_average := computeAverage(gc.exams)
	essay_average := computeAverage(gc.essays)

	// Weights are whole percentages so the math stays in integers; float
	// weights such as .35 round an exact 70 down to 69.
	weighted_grade := assignment_average*50 + exam_average*35 + essay_average*15

	return weighted_grade / 100
}

func computeAverage(grades []Grade) int {
	if len(grades) == 0 {
		return 0
	}

	sum := 0

	for _, grade := range grades {
		sum += grade.Grade
	}

	return sum / len(grades)
}
