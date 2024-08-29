package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/patrikeh/go-deep"
	"github.com/patrikeh/go-deep/training"
)

func parseHolds(holds []byte) (result []float64) {
	var data [11][18]float64
	for i := 0; i < len(holds)/2; i++ {
		holdIndex := holds[i*2] & 0xFF
		// hold type
		_ = holds[i*2+1] & 0xFF

		x := holdIndex % 11
		y := holdIndex / 11

		data[x][y] = 1.0
	}

	for _, xrow := range data {
		for _, y := range xrow {
			result = append(result, y)
		}
	}
	return result
}

type problem struct {
	Holds []byte
	Grade int
}

func main() {
	rand.Seed(0)

	networkFlag := flag.String("network", "data/go_deep_neural.json", "specify another path for the neural network")
	databaseFlag := flag.String("database", "assets/problems.sqlite", "specify another path for the examples database")
	flag.Parse()

	if _, err := os.Stat(*networkFlag); os.IsNotExist(err) {
		_, err := os.Create(*networkFlag)
		if err != nil {
			log.Fatalln("Failed to create the neural network file:", err)
		}
	}

	neuralFile, err := os.OpenFile(*networkFlag, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatalln("Failed to open the neural network file:", err)
	}

	var neural *deep.Neural

	neuralData, err := ioutil.ReadFile(*networkFlag)
	if len(neuralData) == 0 {
		neural = deep.NewNeural(&deep.Config{
			Inputs: 198,
			Layout: []int{60, 16},
			/* Activation functions: Sigmoid, Tanh, ReLU, Linear */
			Activation: deep.ActivationSigmoid,
			/* Determines output layer activation & loss function:
			ModeRegression: linear outputs with MSE loss
			ModeMultiClass: softmax output with Cross Entropy loss
			ModeMultiLabel: sigmoid output with Cross Entropy loss
			ModeBinary: sigmoid output with binary CE loss */
			Mode: deep.ModeBinary,
			/* Weight initializers: {deep.NewNormal(μ, σ), deep.NewUniform(μ, σ)} */
			Weight: deep.NewNormal(1.0, 0.01),
			Bias:   true,
		})
	} else {
		var dump deep.Dump
		err = json.Unmarshal(neuralData, &dump)
		if err != nil {
			log.Fatalln("Failed to read the neural network from file:", err)
		}
		neural = deep.FromDump(&dump)
	}

	trainer := training.NewBatchTrainer(training.NewAdam(0.02, 0.9, 0.999, 1e-8), 1, 200, 8)

	db, err := gorm.Open("sqlite3", *databaseFlag)
	if err != nil {
		log.Panicln("Failed to load database:", err)
	}
	defer db.Close()

	var trainProblems []problem
	db.Order("id").Limit(8200).Find(&trainProblems)

	var testProblems []problem
	db.Raw("SELECT * from problems ORDER BY id LIMIT 8200,-1;").Find(&testProblems)

	var data training.Examples
	for _, problem := range trainProblems {
		grades := make([]float64, 16)
		grades[problem.Grade] = 1.0
		data = append(data, training.Example{
			Input:    parseHolds(problem.Holds),
			Response: grades,
		})
	}

	training, validation := data.Split(0.9)
	trainer.Train(neural, training, validation, 500)

	rawNeural, _ := json.Marshal(neural.Dump())

	neuralFile.Truncate(0)
	neuralFile.Write(rawNeural)

	accurateTests := 0
	for index, problem := range testProblems {
		predicate := neural.Predict(parseHolds(problem.Holds))

		var bestValue float64
		var predictedGrade int
		for grade, value := range predicate {
			if value >= bestValue {
				predictedGrade = grade
				bestValue = value
			}
		}
		accurate := predictedGrade == problem.Grade
		if accurate {
			accurateTests++
		}
		fmt.Printf("Test n°%d --- Expected: %d, Got: %d (Accurate: %v);\n", index, problem.Grade, predictedGrade, accurate)
	}
	fmt.Println("--------------------")
	fmt.Printf("Got %d accurate tests out of %d (proportion: %d%%)\n", accurateTests, len(testProblems), accurateTests*100/len(testProblems))
}
