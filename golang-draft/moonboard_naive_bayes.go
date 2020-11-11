package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/navossoc/bayesian"
)

const (
	C_Unknown bayesian.Class = "Unknown"
	C_6Ap     bayesian.Class = "C_6Ap"
	C_6B      bayesian.Class = "C_6B"
	C_6Bp     bayesian.Class = "C_6Bp"
	C_6C      bayesian.Class = "C_6C"
	C_6Cp     bayesian.Class = "C_6Cp"
	C_7A      bayesian.Class = "C_7A"
	C_7Ap     bayesian.Class = "C_7Ap"
	C_7B      bayesian.Class = "C_7B"
	C_7Bp     bayesian.Class = "C_7Bp"
	C_7C      bayesian.Class = "C_7C"
	C_7Cp     bayesian.Class = "C_7Cp"
	C_8A      bayesian.Class = "C_8A"
	C_8Ap     bayesian.Class = "C_8Ap"
	C_8B      bayesian.Class = "C_8B"
	C_8Bp     bayesian.Class = "C_8Bp"

	Yes bayesian.Class = "Yes"
	No  bayesian.Class = "No"

	scaleX = 2.62
	scaleY = 2.62

	startX = 94 * scaleX
	startY = 936 * scaleY
	deltaX = 50 * scaleX
	deltaY = 50 * scaleY
)

var classes = []bayesian.Class{
	C_Unknown,
	C_6Ap, C_6B, C_6Bp, C_6C, C_6Cp,
	C_7A, C_7Ap, C_7B, C_7Bp, C_7C, C_7Cp,
	C_8A, C_8Ap, C_8B, C_8Bp,
}

type Problem struct {
	Name  string
	Holds []byte
	Grade int
}

func parseHolds(holds []byte) string {
	var data [11][18]bool
	for i := 0; i < len(holds)/2; i++ {
		holdIndex := holds[i*2] & 0xFF
		_ = holds[i*2+1] & 0xFF

		x := holdIndex % 11
		y := holdIndex / 11

		_ = startX + float64(x)*deltaX
		_ = startY - float64(y)*deltaY

		data[x][y] = true
	}

	var result string
	for _, xrow := range data {
		for _, y := range xrow {
			if y {
				result += "1"
			} else {
				result += "0"
			}
		}
	}
	return result
}

func main() {
	db, err := gorm.Open("sqlite3", "assets/problems.sqlite")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	var trainProblems []Problem
	db.Order("id").Limit(3740).Find(&trainProblems)

	var testProblems []Problem
	db.Raw("SELECT * from problems ORDER BY id LIMIT 3740,-1;").Find(&testProblems)

	trainings := make(map[bayesian.Class][]string)

	classifier := bayesian.NewClassifier(Yes, No)
	for index, problem := range trainProblems {
		println("Train n°", index)
		if problem.Grade > 7 {
			trainings[Yes] = append(trainings[Yes], parseHolds(problem.Holds))
		} else {
			trainings[No] = append(trainings[No], parseHolds(problem.Holds))
		}
	}

	for class, training := range trainings {
		classifier.Learn(training, class)
	}

	for index, problem := range testProblems {
		probs, likely, _ := classifier.ProbScores([]string{parseHolds(problem.Holds)})
		fmt.Printf("Test n°%d; Real: %d; Announced: %d; %v;\n", index, problem.Grade, likely, probs)
	}

}
