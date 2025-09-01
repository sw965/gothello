package main

import (
	"fmt"
	"github.com/sw965/gothello"
	"github.com/sw965/gothello/bunny"
	"testing"
	"slices"
	"github.com/sw965/crow/model/linear"
	"github.com/sw965/gothello/game"
	//cmath "github.com/sw965/crow/math"
	cgame "github.com/sw965/crow/game/sequential"
	orand "github.com/sw965/omw/math/rand"
	oslices "github.com/sw965/omw/slices"
	"github.com/sw965/omw/funcs"
	"math/rand"
	"runtime"
)

func newActor(model linear.Model) cgame.Actor[gothello.State, gothello.BitBoard, gothello.Disc] {
	pp := func(state gothello.State, legalActions []gothello.BitBoard) cgame.Policy[gothello.BitBoard] {
		input := bunny.NewDefaultInput(state)
		y := model.Predict(input)
		policy := cgame.Policy[gothello.BitBoard]{}
		for i, a := range gothello.Singles {
			policy[a] = y[i]
		}
		return policy
	}

	selector := cgame.WeightedRandomSelector[gothello.BitBoard, gothello.Disc]

	return cgame.Actor[gothello.State, gothello.BitBoard, gothello.Disc]{
		PolicyProvider:pp,
		Selector:selector,
	}
}

func winRateBySelfPlay(gl cgame.Logic[gothello.State, gothello.BitBoard, gothello.Disc], selfActor, oppActor cgame.Actor[gothello.State, gothello.BitBoard, gothello.Disc], n int, rngs []*rand.Rand) (float32) {
	inits := make([]gothello.State, n)
	for i := range inits {
		inits[i] = gothello.NewInitState()
	}

	finalsByAgPerm, err := gl.CrossPlayouts(inits, []cgame.Actor[gothello.State, gothello.BitBoard, gothello.Disc]{selfActor, oppActor}, rngs)
	if err != nil {
		panic(err)
	}

	agentPerms := oslices.Permutations[[]gothello.Disc](gl.Agents, len(gl.Agents))
	sumScore := float32(0.0)
	for i, finals := range finalsByAgPerm {
		selfAgent := agentPerms[i][0]
		for _, final := range finals {
			scores, err := gl.EvaluateResultScoreByAgent(final)
			if err != nil {
				panic(err)
			}
			sumScore += scores[selfAgent]
		}
	}
	return sumScore / float32(n * len(agentPerms))
}

func hard(x float32) float32 {
	return 1 - x
}

func varF(x float32) float32 {
	return x * (1 - x)
}

func Test(t *testing.T) {
	p := runtime.NumCPU() / 2
	playoutRngs := make([]*rand.Rand, p)
	for i := range playoutRngs {
		playoutRngs[i] = orand.NewMt19937()
	}

	gameLogic := game.NewLogic()

	bIdxs := make([]int, gothello.BoardSize)
	for i := range bIdxs {
		for idx, t := range gothello.GroupIndexTable {
			if slices.Contains(t, i) {
				bIdxs[i] = idx
			}
		}
	}

	latestModel := linear.Model{
		Parameter:linear.Parameter{
			Weight:make([]float32, bunny.DefaultFeatureWeightNum),
			Bias:make([]float32, len(gothello.GroupIndexTable)),
		},
		OutputLayer:linear.NewSoftmaxLayer(0.001, 0.99),
		BiasIndices:slices.Clone(bIdxs),
	}

	poolSize := 512
	policyPool := make([]linear.Model, 0, poolSize)
	policyPool = append(policyPool, latestModel.Clone())
	currentOppModel := policyPool[0].Clone()
	
	lossFunc := func(newModel linear.Model, workerIdx int) (float32, error) {
		selfActor := newActor(newModel)
		oppActor := newActor(currentOppModel)
		winRate := winRateBySelfPlay(gameLogic, selfActor, oppActor, 128, playoutRngs)
		loseRate := 1.0 - winRate
		sqSum := float32(0.0)
		for _, wi := range newModel.Parameter.Weight {
			sqSum += wi * wi
		}
		wn := len(newModel.Parameter.Weight)
		return loseRate + (0.001 * sqSum / float32(wn)), nil
	}

	latestWinRates := make([]float32, 0, poolSize)
	latestWinRates = append(latestWinRates, 0.0)

	selectPolicyRng := orand.NewMt19937()

	selectCount := 0
	selectPolicyModel := func() linear.Model {
		var f func(float32) float32
		if selectCount%2 == 0 {
			f = hard
		} else {
			f = varF
		}

		w := funcs.Map(latestWinRates, f)
		idx, err := orand.IntByWeight(w, selectPolicyRng)
		if err != nil {
			panic(err)
		}
		selectCount += 1
		return policyPool[idx]
	}

	runLeague := func() bool {
		for i, model := range policyPool {
			selfActor := newActor(latestModel)
			oppActor := newActor(model)
			winRate := winRateBySelfPlay(gameLogic, selfActor, oppActor, 640, playoutRngs)
			latestWinRates[i] = winRate
			if winRate < 0.55 {
				fmt.Println(i, "番目のモデルへの勝率 =", winRate)
				return false
			}
		}
		return true
	}

	rndModel := latestModel.Clone()
	rndActor := newActor(rndModel)
	checkPoint := 128
	spsaRngs := make([]*rand.Rand, 1)
	spsaRngs[0] = orand.NewMt19937()

	for i := 0; i < 12800000; i++ {
		currentOppModel = selectPolicyModel()
		grad, err := latestModel.EstimateGradBySPSA(0.1, lossFunc, spsaRngs)
		if err != nil {
			panic(err)
		}
		latestModel.Parameter.AxpyGrad(-0.001, grad)

		if i != 0 && i%checkPoint == 0 {
			fmt.Println("i = ", i)
			ok := runLeague()
			fmt.Println("latestWinRates = ", latestWinRates)
			fmt.Println("")

			if ok {
				if poolSize == len(policyPool) {
					policyPool = make([]linear.Model, 0, poolSize)
					latestWinRates = make([]float32, 0, poolSize)
					fmt.Println("ポリシープールが満タンになったので、削除しました。")
				}
				policyPool = append(policyPool, latestModel.Clone())
				latestWinRates = append(latestWinRates, 0.5)
				fmt.Println("ポリシープールに新たなモデルを追加しました。", len(policyPool))
			} else {
				fmt.Println("条件を満たさなかった為、ポリシープールに追加出来ませんでした。", len((policyPool)))
				continue
			}
			fmt.Println("")

			latestModelActor := newActor(latestModel)

			testInits := make([]gothello.State, 1280)
			for i := range testInits {
				testInits[i] = gothello.NewInitState()
			}

			finalsByAgPerm, err := gameLogic.CrossPlayouts(testInits, []cgame.Actor[gothello.State, gothello.BitBoard, gothello.Disc]{latestModelActor, rndActor}, playoutRngs)
			if err != nil {
				panic(err)
			}

			agentPerms := oslices.Permutations[[]gothello.Disc](gameLogic.Agents, len(gameLogic.Agents))

			for i, finals := range finalsByAgPerm {
				selfAgent := agentPerms[i][0]
				sumScore := float32(0.0)
				for _, final := range finals {
					scores, err := gameLogic.EvaluateResultScoreByAgent(final)
					if err != nil {
						panic(err)
					}
					sumScore += scores[selfAgent]
				}
				fmt.Println("agent =", selfAgent, "テスト勝率 =", sumScore / float32(len(finals)))
			}
			fmt.Println("")
		}
	}
}
