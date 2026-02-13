package game_test

import (
	"github.com/sw965/gothello"
	"github.com/sw965/gothello/game"
	"github.com/sw965/crow/pucb"
	mcts "github.com/sw965/crow/mcts/puct"
	"testing"
	"fmt"
	"github.com/sw965/omw/mathx/randx"
	"github.com/sw965/crow/game/sequential"
	"math/rand/v2"
)

func Test(t *testing.T) {
	rng := randx.NewPCGFromGlobalSeed()
	engine := game.NewEngine()
	mctsEngine := mcts.Engine[gothello.State, gothello.BitBoard, gothello.Turn]{
		Game:engine,
		PUCBFunc:pucb.NewAlphaGoFunc(5.0),
		NextNodesCap:64,
		VirtualValue:0.0,
	}

	randActor := sequential.NewRandomActor[gothello.State, gothello.BitBoard, gothello.Turn]("randActor")
	mctsEngine.SetUniformPolicyFunc()
	mctsEngine.SetPlayout(randActor, rng)

	rngs := make([]*rand.Rand, 6)
	for i := range rngs {
		rngs[i] = randx.NewPCGFromGlobalSeed()
	}

	policyFunc := mctsEngine.NewPolicy(2560, rngs)
	mctsActor := sequential.Actor[gothello.State, gothello.BitBoard, gothello.Turn]{
		Name:"mctsActor",
		PolicyFunc:policyFunc,
		SelectFunc:sequential.MaxSelectFunc[gothello.BitBoard, gothello.Turn],
	}

	inits := make([]gothello.State, 100)
	for i := range inits {
		inits[i] = gothello.NewInitState()
	}
	results, err := engine.CrossPlayouts(inits, []sequential.Actor[gothello.State, gothello.BitBoard, gothello.Turn]{
		randActor, mctsActor,
	}, []*rand.Rand{rng})

	if err != nil {
		panic(err)
	}

// 各アクターの合計スコアとプレイ回数を記録するマップ
	totalScores := make(map[string]float32)
	playCounts := make(map[string]int)

	for _, result := range results {
		for _, final := range result.Finals {
			// 最終盤面から各エージェント（Turn）のスコアを評価
			scores, err := engine.EvaluateResultScoreByAgent(final)
			if err != nil {
				panic(err)
			}

			// エージェントごとのスコアを、対応するアクター（randActor/mctsActor）に加算
			for agent, score := range scores {
				actor := result.ActorByAgent[agent]
				totalScores[actor.Name] += score
				playCounts[actor.Name]++ // 今回のCrossPlayoutsの仕様上、各アクターはエージェント数分プレイする
			}
		}
	}

	// 結果の出力
	fmt.Println("=== 対戦結果 ===")
	for name, count := range playCounts {
		winRate := float32(0)
		if count > 0 {
			// 勝率（平均スコア）を算出
			winRate = totalScores[name] / float32(count)
		}
		fmt.Printf("Actor: %-10s | 勝率(平均スコア): %5.1f%% | 試合数: %d\n", name, winRate*100, count)
	}
}