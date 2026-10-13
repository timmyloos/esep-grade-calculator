# ESEP Grade Calculator

A small Go library that calculates a student's final letter grade from
their assignment, exam, and essay scores. It is used in ESEP as an exercise
in unit testing and code coverage.

## Project layout

| File | Purpose |
| --- | --- |
| `README.MD` | This document. |
| `grade-calculator/go.mod` | Go module definition (`esep/grade-calculator`, Go 1.24.4). |
| `grade-calculator/grade_calculator.go` | The library: `GradeCalculator`, `Grade`, and `GradeType`. |
| `grade-calculator/grade_calculator_test.go` | Unit tests for letter grades, weights, and edge cases. |
| `grade-calculator/grade_calculator_passfail_test.go` | Unit tests for Pass/Fail grading. |

## Requirements

- Go 1.24.4 or newer (the version declared in `go.mod`)
- Git, to clone the repository

There are no third-party dependencies. The tests use Go's built-in
`testing` package.

## Installation

1. Install Go from <https://go.dev/dl/> and confirm it is on your `PATH`:

   ```
   go version
   ```

2. Clone the repository and move into the module directory:

   ```
   git clone https://github.com/timmyloos/esep-grade-calculator.git
   cd esep-grade-calculator/grade-calculator
   ```

All `go` commands below are run from the `grade-calculator/` directory,
because that is where `go.mod` lives.

## Running the program

This project is a library, not an executable. There is no `main` package,
so there is nothing to start with `go run`. The code is exercised through
its unit tests, or by calling its public API from other Go code (see
[Example usage](#example-usage)).

To check that the package compiles:

```
go build ./...
```

## Running the tests

```
cd grade-calculator/
go test ./...
```

Add `-v` to see each test by name:

```
go test -v ./...
```

Run a single test by name:

```
go test -v -run TestGetGradeA ./...
```

## Running the tests with code coverage

Print the overall statement coverage:

```
go test -cover ./...
```

Write a coverage profile and show coverage per function:

```
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

Open an HTML report that highlights covered and uncovered lines:

```
go tool cover -html=coverage.out
```

Coverage profiles are generated output and should not be committed.

## How the grade calculator works

1. Create a calculator with `NewGradeCalculator()`.
2. Record each score with `AddGrade(name, grade, gradeType)`. The grade
   type decides which category the score belongs to.
3. Call `GetFinalGrade()` to get the final grade.

`GetFinalGrade()` averages the scores in each category, combines the three
category averages using the weights below, and converts the weighted
numerical grade to a letter using the grading scale.

### Grade weights

| Category | `GradeType` | Weight |
| --- | --- | --- |
| Assignments | `Assignment` | 50% |
| Exams | `Exam` | 35% |
| Essays | `Essay` | 15% |

```
numerical grade = assignment average * 0.50
                + exam average       * 0.35
                + essay average      * 0.15
```

### Grading scale

| Letter | Numerical grade |
| --- | --- |
| A | 90 - 100 |
| B | 80 - 89 |
| C | 70 - 79 |
| D | 60 - 69 |
| F | Below 60 |

A grade of C or better (70 or above) is passing.

A category with no grades averages to 0, so it contributes nothing to the
numerical grade. The numerical grade is truncated, not rounded, so 89.65
is a B.

### Grading modes

| `GradingMode` | `GetFinalGrade()` returns |
| --- | --- |
| `LetterGrade` (default) | `"A"`, `"B"`, `"C"`, `"D"`, or `"F"` |
| `PassFail` | `"Pass"` for 70 or above, `"Fail"` below 70 |

`NewGradeCalculator()` always uses `LetterGrade`. To choose a mode, create
the calculator with `NewGradeCalculatorWithMode(mode)`.

## Public API

| Name | Description |
| --- | --- |
| `NewGradeCalculator() *GradeCalculator` | Creates an empty calculator in `LetterGrade` mode. |
| `NewGradeCalculatorWithMode(mode GradingMode) *GradeCalculator` | Creates an empty calculator in the given grading mode. |
| `(*GradeCalculator).AddGrade(name string, grade int, gradeType GradeType)` | Records one score in the given category. |
| `(*GradeCalculator).GetFinalGrade() string` | Returns the final grade: a letter in `LetterGrade` mode, or `"Pass"` / `"Fail"` in `PassFail` mode. |
| `GradingMode` | How the final grade is reported: `LetterGrade` or `PassFail`. |
| `GradeType` | Category of a grade: `Assignment`, `Exam`, or `Essay`. |
| `(GradeType).String() string` | Returns `"assignment"`, `"exam"`, or `"essay"`. |
| `Grade` | A single recorded score with `Name`, `Grade`, and `Type` fields. |

## Example usage

The package is named `esepunittests`. The example below is written as code
inside that package (for example, in a test file):

```go
gradeCalculator := NewGradeCalculator()

gradeCalculator.AddGrade("open source assignment", 80, Assignment)
gradeCalculator.AddGrade("exam 1", 81, Exam)
gradeCalculator.AddGrade("essay on ai ethics", 85, Essay)

// 80*0.50 + 81*0.35 + 85*0.15 = 81.1, which is a B
finalGrade := gradeCalculator.GetFinalGrade() // "B"

fmt.Println(Essay.String()) // "essay"
```

Pass/Fail grading with the same scores:

```go
gradeCalculator := NewGradeCalculatorWithMode(PassFail)

gradeCalculator.AddGrade("open source assignment", 80, Assignment)
gradeCalculator.AddGrade("exam 1", 81, Exam)
gradeCalculator.AddGrade("essay on ai ethics", 85, Essay)

finalGrade := gradeCalculator.GetFinalGrade() // "Pass"
```
