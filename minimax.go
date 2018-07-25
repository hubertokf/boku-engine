// Package minimax implements the minimax algorithm
// Minimax (sometimes MinMax or MM[1]) is a decision rule used in decision theory,
// game theory, statistics and philosophy for minimizing the possible loss for
// a worst case (maximum loss) scenario
// See for more details: https://en.wikipedia.org/wiki/Minimax
package main

import (
	_ "fmt"
	_ "strconv"
	_ "github.com/mohae/deepcopy"
	_ "runtime"
)

// Node represents an element in the decision tree
type Node struct {
	// Score is available when supplied by an evaluation function or when calculated
	Score      *int
	parent     *Node
	children   []*Node
	isOpponent bool
	Movement [2]int

	// Data field can be used to store additional information by the consumer of the
	// algorithm
	// Data [][]int
}

// New returns a new minimax structure
func NewNode() Node {
	n := Node{isOpponent: false}
	return n
}

// GetBestChildNode returns the first child node with the matching score
func (node *Node) GetBestChildNode() *Node {
	for _, cn := range node.children {
		if cn.Score == node.Score {
			return cn
		}
	}

	return nil
}

// Evaluate runs through the tree and caculates the score from the terminal nodes
// all the the way up to the root node
func (node *Node) Evaluate(board [][]int, plies int, player int) {
	eval := false
	if plies == 0{
		eval = true

	}
	node.generateChildMoves(board, player, eval)
	// PrintMemUsage()
	// fmt.Println(node.children)
	for _, cn := range node.children {
		// fmt.Println(cn.Data)
		doMove(board, player, cn.Movement)
		if plies != 0 {
			if player == 1{
				player = 2
			}else{
				player = 1
			}
			cn.Evaluate(board, plies-1, player)
		}

		if cn.parent.Score == nil {
			cn.parent.Score = cn.Score
		} else if cn.isOpponent && *cn.Score > *cn.parent.Score {
			cn.parent.Score = cn.Score
		} else if !cn.isOpponent && *cn.Score < *cn.parent.Score {
			cn.parent.Score = cn.Score
		}

		undoMove(board, cn.Movement)
	}
}

// AddTerminal adds a terminal node (or leave node).  These nodes
// should contain a score and no children
func (node *Node) AddTerminal(score int, movement [2]int) *Node {
	return node.add(&score, movement)
}

// Add a new node to structure, this node should have children and
// an unknown score
func (node *Node) Add(movement [2]int) *Node {
	return node.add(nil, movement)
}

func (node *Node) add(score *int, movement [2]int) *Node {
	// board := copy_board(data)

	childNode := Node{parent: node, Score: score, Movement: movement}

	childNode.isOpponent = !node.isOpponent
	node.children = append(node.children, &childNode)
	// fmt.Println(node.children)
	return &childNode
}

func (node *Node) isTerminal() bool {
	return len(node.children) == 0
}

func neighbors(board [][]int, pos [2]int) [][]int {
	column := pos[0]
	line := pos[1]

	var position []int
	var l [][]int

	if line < len(board[column])-1 { // DOWN
		position = []int{column, line + 1}
		l = append(l, position)
	}

	if column <= len(board)-1 && column != 0 { // DIAGONAL L/D
		if line != len(board[column])-1 && column < 6 {
			position = []int{column - 1, line}
			l = append(l, position)
		} else if column >= 6 {
			position = []int{column - 1, line + 1}
			l = append(l, position)
		}
	}

	if column <= len(board)-1 && column != 0 { // DIAGONAL L/U
		if column < 6 && line != 0 {
			position = []int{column - 1, line - 1}
			l = append(l, position)
		} else if column >= 6 {
			position = []int{column - 1, line}
			l = append(l, position)
		}
	}

	if line != 0 {
		position = []int{column, line - 1} // UP
		l = append(l, position)
	}

	if column < len(board)-1 { // DIAGONAL R/U
		if column < 5 {
			position = []int{column + 1, line}
			l = append(l, position)
		} else if column >= 5 && line != 0 {
			position = []int{column + 1, line - 1}
			l = append(l, position)
		}
	}

	if column < len(board)-1 { // DIAGONAL R/D
		if column < 5 {
			position = []int{column + 1, line + 1}
			l = append(l, position)
		} else if column >= 5 && line != len(board[column+1]) {
			position = []int{column + 1, line}
			l = append(l, position)
		}
	}

	return l
}

func heuristic(board [][]int, position [2]int) int {
	counter := 0
	v := neighbors(board, position)

	for _, vizinho := range v {
		if board[vizinho[0]][vizinho[1]] == 0 {
			counter = counter + 10
		}
		if board[vizinho[0]][vizinho[1]] == 1 { // VERIFICAR o 1
			counter = counter + 20
		}
		if board[vizinho[0]][vizinho[1]] == 2 {
			counter = counter - 50
		}
	}
	// node.Score = counter
	return counter
}

func doMove(board [][]int, player int, movement [2]int){
	board[movement[0]][movement[1]] = player
}

func undoMove(board [][]int, movement [2]int){
	board[movement[0]][movement[1]] = 0
}

func (node *Node) generateChildMoves(board [][]int, player int, eval bool) *Node {
	for k_col, col := range board {
		for k_cell, cell := range col {
			if cell == 0 {
				// node.Data[k_col][k_cell] = player
				// board := copy_board(node.Data)
				// board := deepcopy.Copy(node.Data)
				
				// PrintMemUsage()
				
				if eval == true {
					// evaluate state
					score := heuristic(board, [2]int{k_col,k_cell})
					// add child node to node child list
					node.AddTerminal(score, [2]int{k_col,k_cell})
					}else{
						// add child node to node child list
						node.Add([2]int{k_col,k_cell})
					}
				}

				// node.Data[k_col][k_cell] = 0
		}
	}

	return node
}