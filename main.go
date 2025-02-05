package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Student struct {
	FirstName  string
	LastName   string
	University string
	TestScores []int
	FinalScore float64
	Grade      string
}

func readCSV(filename string) ([]Student, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var students []Student
	for i, row := range rows {
		if i == 0 { // Skip header row
			continue
		}

		student, err := createStudentFromRow(row)
		if err != nil {
			return nil, err
		}
		students = append(students, student)
	}

	return students, nil
}

func createStudentFromRow(row []string) (Student, error) {
	if len(row) < 7 {
		return Student{}, fmt.Errorf("invalid row: %v", row)
	}

	testScores := make([]int, 4)
	for i := 0; i < 4; i++ {
		score, err := strconv.Atoi(strings.TrimSpace(row[3+i]))
		if err != nil {
			return Student{}, fmt.Errorf("invalid score in row: %v", row)
		}
		testScores[i] = score
	}

	finalScore := calculateFinalScore(testScores)
	grade := determineGrade(finalScore)

	return Student{
		FirstName:  row[0],
		LastName:   row[1],
		University: row[2],
		TestScores: testScores,
		FinalScore: finalScore,
		Grade:      grade,
	}, nil
}

func calculateFinalScore(scores []int) float64 {
	total := 0
	for _, score := range scores {
		total += score
	}
	return float64(total) / float64(len(scores))
}

func determineGrade(score float64) string {
	switch {
	case score < 35:
		return "F"
	case score >= 35 && score < 50:
		return "C"
	case score >= 50 && score < 70:
		return "B"
	default:
		return "A"
	}
}

func findOverallTopper(students []Student) Student {
	if len(students) == 0 {
		return Student{}
	}

	topStudent := students[0]
	for _, student := range students {
		if student.FinalScore > topStudent.FinalScore {
			topStudent = student
		}
	}
	return topStudent
}

func findTopperPerUniversity(students []Student) map[string]Student {
	toppers := make(map[string]Student)

	for _, student := range students {
		topStudent, exists := toppers[student.University]
		if !exists || student.FinalScore > topStudent.FinalScore {
			toppers[student.University] = student
		}
	}

	return toppers
}

func main() {
	filename := "grades.csv"

	students, err := readCSV(filename)
	if err != nil {
		fmt.Println("Error reading CSV:", err)
		return
	}

	overallTopper := findOverallTopper(students)
	fmt.Printf("Overall Topper: %s %s (%.2f - Grade %s)\n",
		overallTopper.FirstName, overallTopper.LastName, overallTopper.FinalScore, overallTopper.Grade)
	
	universityToppers := findTopperPerUniversity(students)
	fmt.Println("\nUniversity-wise Toppers:")
	for uni, topper := range universityToppers {
		fmt.Printf("%s: %s %s (%.2f - Grade %s)\n",
			uni, topper.FirstName, topper.LastName, topper.FinalScore, topper.Grade)
	}
}
