package bunny_test

import (
	"testing"
	"fmt"
	"github.com/sw965/gothello/bunny"
	"github.com/sw965/gothello"
	"slices"
	"github.com/sw965/crow/model/linear"
	"github.com/sw965/gothello/game"
	cgame "github.com/sw965/crow/game/sequential"
	orand "github.com/sw965/omw/math/rand"
	"math/rand"
	"runtime"
	"github.com/sw965/omw/fn"
	cmath "github.com/sw965/crow/math"
)

func Test(t *testing.T) {
	// bunny.MakePatterns()
	// return
	gameLogic := game.NewLogic()

	param := linear.Parameter{
		Weight:make([]float32, len(gothello.GroupIndexTable)),
		Bias:make([]float32, len(gothello.GroupIndexTable)),
	}

	for i := range param.Weight {
		param.Weight[i] = 1.0
	}

	bIdxs := make([]int, gothello.BoardSize)
	for i := range bIdxs {
		for idx, t := range gothello.GroupIndexTable {
			if slices.Contains(t, i) {
				bIdxs[i] = idx
			}
		}
	}

	latestModel := linear.Model{
		Parameter:param,
		OutputLayer:linear.NewSoftmaxLayer(0.001, 0.99),
		BiasIndices:bIdxs,
	}

	newPolicyProvider := func(blackModel, whiteModel linear.Model) cgame.PolicyProvider[gothello.State, gothello.BitBoard] {
		return func(state gothello.State, _ []gothello.BitBoard) cgame.Policy[gothello.BitBoard] {
			model := map[gothello.Color]linear.Model{
				gothello.Black:blackModel,
				gothello.White:whiteModel,
			}[state.Hand]
			y := model.Predict(bunny.NewInput(state))
			m := map[gothello.BitBoard]float32{}
			u := gothello.BitBoard(0b11111111_11111111_11111111_11111111_11111111_11111111_11111111_11111111).ToSingles()
			for i, a := range u {
				m[a] = y[i]
			}
			return m
		}
	}

	p := runtime.NumCPU()
	rngs := make([]*rand.Rand, p)
	for i := range rngs {
		rngs[i] = orand.NewMt19937()
	}

	trainNum := 256000
	oldModel := latestModel.Clone()

	allActions := gothello.BitBoard(0b11111111_11111111_11111111_11111111_11111111_11111111_11111111_11111111).ToSingles()

	// initStates := make(gothello.States, 128)
	// for i := range initStates {
	// 	initStates[i] = gothello.NewInitState()
	// }

	// lossFunc := func(model linear.Model, workerIdx int) (float32, error) {
	// 	pp := newPolicyProvider(model, latestModel)
	// 	blackFinals, err := gameLogic.Playouts(initStates, pp, rngs)
	// 	if err != nil {
	// 		return 0.0, err
	// 	}

	// 	pp = newPolicyProvider(latestModel, model)
	// 	whiteFinals, err := gameLogic.Playouts(initStates, pp, rngs)
	// 	if err != nil {
	// 		return 0.0, err
	// 	}

	// 	counter, err := gothello.States(blackFinals).CountResult()
	// 	if err != nil {
	// 		return 0.0, err
	// 	}

	// 	blackWinRate := counter.BlackWinRate(true)

	// 	counter, err = gothello.States(whiteFinals).CountResult()
	// 	if err != nil {
	// 		return 0.0, err
	// 	}

	// 	whiteWinRate := counter.WhiteWinRate(true)

	// 	winRate := float32(blackWinRate + whiteWinRate) / 2.0
	// 	return 1.0 - winRate, nil
	// }

	// spsaRng := orand.NewMt19937()

	// for i := 0; i < trainNum; i++ {
	// 	grad, err := latestModel.EstimateGradBySPSA(0.2, lossFunc, []*rand.Rand{spsaRng})
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	latestModel.Parameter.AxpyGrad(-0.1, grad)

	// 	if i%1280 == 0 {
	// 		s := ""
	// 		for i, bIdx := range latestModel.BiasIndices {
	// 			s += fmt.Sprintf("%.2f ", latestModel.Parameter.Bias[bIdx])
	// 			if i != 0 && (i+1)%8 == 0 {
	// 				s += "\n"
	// 			}
	// 		}
	// 		fmt.Println(s)

	// 		policyProvider := newPolicyProvider(latestModel, oldModel)
	// 		initStates := make([]gothello.State, 1920)
	// 		for i := range initStates {
	// 			initStates[i] = gothello.NewInitState()
	// 		}

	// 		testBlackFinals, err := gameLogic.Playouts(initStates, policyProvider, rngs)
	// 		if err != nil {
	// 			panic(err)
	// 		}

	// 		counter, err := gothello.States(testBlackFinals).CountResult()
	// 		if err != nil {
	// 			panic(err)
	// 		}

	// 		fmt.Println("黒テスト", counter.BlackWinRate(true))

	// 		policyProvider = newPolicyProvider(oldModel, latestModel)
	// 		testWhiteFinals, err := gameLogic.Playouts(initStates, policyProvider, rngs)
	// 		if err != nil {
	// 			panic(err)
	// 		}

	// 		counter, err = gothello.States(testWhiteFinals).CountResult()
	// 		if err != nil {
	// 			panic(err)
	// 		}

	// 		fmt.Println("白テスト", counter.WhiteWinRate(true))	
	// 	}
	// }

	for i := 0; i < trainNum; i++ {
		initStates := make(gothello.States, 128)
		for i := range initStates {
			initStates[i] = gothello.NewInitState()
		}

		history, err := gameLogic.PlayoutsWithHistory(initStates, newPolicyProvider(latestModel, latestModel), rngs)
		if err != nil {
			panic(err)
		}

		experiences, err := history.ToExperiences(gameLogic)
		if err != nil {
			panic(err)
		}

		stateHistory, _, actionHistory, scores, _ := experiences.Split()

		inputHistory := fn.Map[linear.Inputs](stateHistory, func(state gothello.State) linear.Input {
			return bunny.NewInput(state)
		})

		actionIndexHistory := fn.Map[[]int](actionHistory, func(action gothello.BitBoard) int {
			return slices.Index(allActions, action)
		})

		rewards := fn.Map[[]float32](scores, func(score float32) float32 {
			return cmath.ConvertScale(score, 0.0, 1.0, -1.0, 1.0)
		})

		grad, err := latestModel.ComputeGradByReinforce(inputHistory, actionIndexHistory, rewards, p)
		if err != nil {
			panic(err)
		}

		latestModel.Parameter.AxpyGrad(-0.01, grad)

		if i%2560 == 0 {
			s := ""
			for i, bIdx := range latestModel.BiasIndices {
				s += fmt.Sprintf("%.2f ", latestModel.Parameter.Bias[bIdx])
				if i != 0 && (i+1)%8 == 0 {
					s += "\n"
				}
			}
			fmt.Println(s)

			policyProvider := newPolicyProvider(latestModel, oldModel)
			initStates := make([]gothello.State, 1920)
			for i := range initStates {
				initStates[i] = gothello.NewInitState()
			}

			testBlackFinals, err := gameLogic.Playouts(initStates, policyProvider, rngs)
			if err != nil {
				panic(err)
			}

			counter, err := gothello.States(testBlackFinals).CountResult()
			if err != nil {
				panic(err)
			}

			fmt.Println("黒テスト", counter.BlackWinRate(true))

			policyProvider = newPolicyProvider(oldModel, latestModel)
			testWhiteFinals, err := gameLogic.Playouts(initStates, policyProvider, rngs)
			if err != nil {
				panic(err)
			}

			counter, err = gothello.States(testWhiteFinals).CountResult()
			if err != nil {
				panic(err)
			}

			fmt.Println("白テスト", counter.WhiteWinRate(true))
			//oldModel = latestModel.Clone()
		}
	}
}