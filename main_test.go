package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadCSV(t *testing.T) {
	a := assert.New(t)

	students, err := readCSV("grades.csv")
	a.NoError(err, "CSV parsing should not return an error")
	a.Equal(30, len(students), "Student list size should be 30")

	expectedFirst := Student{"Kaylen", "Johnson", "Duke University", []int{52, 47, 35, 38}, 0, ""}
	expectedLast := Student{"Solomon", "Hunter", "Boston University", []int{45, 62, 32, 58}, 0, ""}

	a.Equal(expectedFirst, students[0], "First student should be Kaylen")
	a.Equal(expectedLast, students[29], "Last student should be Solomon")
}


func TestCalculateFinalScore(t *testing.T) {
	a := assert.New(t)

	students, err := readCSV("grades.csv")
	a.NoError(err)

	expectedScores := []float64{43.00, 59.25, 53.00, 58.25, 52.25, 50.75, 54.75, 49.25, 64.75, 43.25, 68.50, 57.75, 68.25, 66.75, 45.50, 45.75, 45.50, 58.00, 56.00, 60.25, 61.00, 62.50, 80.50, 53.00, 30.75, 57.50, 70.75, 48.50, 60.25, 49.25}

	for i, student := range students {
		finalScore := calculateFinalScore(student.TestScores)
		a.InDelta(expectedScores[i], finalScore, 0.01, "Final score mismatch for student %v", student.FirstName)
	}
}

func TestDetermineGrade(t *testing.T) {
	a := assert.New(t)

	testCases := []struct {
		score float64
		grade string
	}{
		{30.0, "F"},
		{40.0, "C"},
		{55.0, "B"},
		{75.0, "A"},
		{89.0, "A"},
	}

	for _, tc := range testCases {
		a.Equal(tc.grade, determineGrade(tc.score), "Expected grade %v for score %.2f, but got %v", tc.grade, tc.score, determineGrade(tc.score))
	}
}

func TestFindOverallTopper(t *testing.T) {
	a := assert.New(t)

	students, err := readCSV("grades.csv")
	a.NoError(err)

	for i := range students {
		students[i].FinalScore = calculateFinalScore(students[i].TestScores)
		students[i].Grade = determineGrade(students[i].FinalScore)
	}

	overallTopper := findOverallTopper(students)

	expectedTopper := Student{"Bernard", "Wilson", "Boston University", []int{90, 85, 76, 71}, 80.5, "A"}
	a.Equal(expectedTopper, overallTopper, "Overall topper mismatch")
}

func TestFindTopperPerUniversity(t *testing.T) {
	a := assert.New(t)

	students, err := readCSV("grades.csv")
	a.NoError(err)

	for i := range students {
		students[i].FinalScore = calculateFinalScore(students[i].TestScores)
		students[i].Grade = determineGrade(students[i].FinalScore)
	}

	universityToppers := findTopperPerUniversity(students)

	expectedToppers := map[string]Student{
		"Boston University":      {"Bernard", "Wilson", "Boston University", []int{90, 85, 76, 71}, 80.5, "A"},
		"Duke University":        {"Tamara", "Webb", "Duke University", []int{73, 62, 90, 58}, 70.75, "A"},
		"Union College":          {"Izayah", "Hunt", "Union College", []int{29, 78, 41, 85}, 58.25, "B"},
		"University of California": {"Karina", "Shaw", "University of California", []int{69, 78, 56, 70}, 68.25, "B"},
		"University of Florida":  {"Nathan", "Gordon", "University of Florida", []int{53, 79, 84, 51}, 66.75, "B"},
	}

	for uni, expectedTopper := range expectedToppers {
		a.Equal(expectedTopper, universityToppers[uni], "University topper mismatch for %v", uni)
	}
}
