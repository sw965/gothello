package policy

import (
	"fmt"
	"github.com/sw965/gothello"
)

func Legal(board []int) []int {
	for i, disc := range board {
		switch disc {
			case 0:
				left := board[:i]
				fmt.Println(left)
				right := board[i+1:]
				fmt.Println(right)
			case 1:

			case 2:
		}
	}
	return nil
}

type YasuriIndexer map[int]map[gothello.BitBoard]map[gothello.Feature]int

// func NewYasuriIndexer() {
// 	cornerMoveMasks := []gothello.BitBoard{
// 		//左上
// 		gothello.UpSideMask | gothello.LeftSideMask,

// 		//右上
// 		gothello.UpSideMask | gothello.RightSideMask,

// 		//左下
// 		gothello.DownSideMask | gothello.LeftSideMask,

// 		//右下
// 		gothello.DownSideMask | gothello.RightSideMask,
// 	}


// }