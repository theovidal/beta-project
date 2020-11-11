package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/patrikeh/go-deep"
	"github.com/patrikeh/go-deep/training"
)

func parseHolds(holds []byte) []float64 {
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

	var result []float64
	for _, xrow := range data {
		for _, y := range xrow {
			if y {
				result = append(result, 1.0)
			} else {
				result = append(result, 1.0)
			}
		}
	}
	return result
}

type Problem struct {
	Holds []byte
	Grade int
}

var (
	scaleX = 2.62
	scaleY = 2.62

	startX = 94 * scaleX
	startY = 936 * scaleY
	deltaX = 50 * scaleX
	deltaY = 50 * scaleY
)

func main() {
	db, err := gorm.Open("sqlite3", "assets/problems.sqlite")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	var trainProblems []Problem
	db.Order("id").Limit(7000).Find(&trainProblems)

	var testProblems []Problem
	db.Raw("SELECT * from problems ORDER BY id LIMIT 7000,-1;").Find(&testProblems)

	var data training.Examples
	for _, problem := range trainProblems {
		grades := make([]float64, 16)
		grades[problem.Grade] = 1.0
		fmt.Println(grades)
		data = append(data, training.Example{
			Input:    parseHolds(problem.Holds),
			Response: grades,
		})
	}

	neural := deep.NewNeural(&deep.Config{
		Inputs: 198,
		Layout: []int{60, 16},
		/* Activation functions: Sigmoid, Tanh, ReLU, Linear */
		Activation: deep.ActivationReLU,
		/* Determines output layer activation & loss function:
		ModeRegression: linear outputs with MSE loss
		ModeMultiClass: softmax output with Cross Entropy loss
		ModeMultiLabel: sigmoid output with Cross Entropy loss
		ModeBinary: sigmoid output with binary CE loss */
		Mode: deep.ModeMultiClass,
		/* Weight initializers: {deep.NewNormal(μ, σ), deep.NewUniform(μ, σ)} */
		Weight: deep.NewNormal(1.0, 0.1),
		Bias:   true,
	})

	trainer := training.NewBatchTrainer(training.NewAdam(0.02, 0.9, 0.999, 1e-8), 1, 200, 8)

	training, validation := data.Split(0.8)
	trainer.Train(neural, training, validation, 10) // training, validation, iterations

	for index, problem := range testProblems {
		predicate := neural.Predict(parseHolds(problem.Holds))

		var bestValue float64
		var bestGrade int
		for grade, value := range predicate {
			if value >= bestValue {
				bestGrade = grade
			}
		}
		fmt.Printf("Test n°%d; Announced: %d, Predicted: %d;\n%v\n", index, problem.Grade, bestGrade, predicate)
	}
}
