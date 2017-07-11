package main

import (
	"fmt"
)

func min(x, y int16) int16 {
	if x < y {
		return x
	}
	return y
}

func max(x, y int16) int16 {
	if x > y {
		return x
	}
	return y
}

func pow(x, y int) int {
	res := 1
	for i := 0; i < y; i++ {
		res = res * x
	}
	return res
}

type movementGrid []int16

type grid interface {
	TakeDamage(dmg int16, col int) grid
	HasMovement(mov movementGrid) bool
	RemainingHealth() int16
}

type simpleGrid int16

func (g simpleGrid) TakeDamage(dmg int16, col int) grid {
	d := max(0, int16(g)-dmg)
	return simpleGrid(d)
}

func (g simpleGrid) HasMovement(mov movementGrid) bool {
	return true
}

func (g simpleGrid) RemainingHealth() int16 {
	return int16(g)
}

type complexGrid []int16

func (g complexGrid) TakeDamage(dmg int16, col int) grid {
	remDmg := dmg
	res := make([]int16, len(g))
	for i := col; i < len(g); i++ {
		possibleDmg := min((g)[i], remDmg)
		res[i] = (g)[i] - possibleDmg
		remDmg = remDmg - possibleDmg
	}

	for i := 0; i < col; i++ {
		possibleDmg := min((g)[i], remDmg)
		res[i] = (g)[i] - possibleDmg
		remDmg = remDmg - possibleDmg
	}

	return complexGrid(res)
}

func (g complexGrid) HasMovement(mov movementGrid) bool {
	// has no movement if f.a. either no mov in col
	// or
	hasNoMovement := true
	for i := 0; i < len(g); i++ {
		if mov[i] > 0 {
			if g[i] > 0 {
				hasNoMovement = false
				break
			}
		}
	}

	return !hasNoMovement
}

func (g complexGrid) RemainingHealth() int16 {
	var res int16
	for i := 0; i < len(g); i++ {
		res = res + g[i]
	}
	return res
}

type randomVariable []float32

func (pv randomVariable) GreaterEqual(n int) float32 {
	res := float32(0.0)
	for i := n - 1; i < len(pv); i++ {
		res = res + (pv)[i]
	}
	return res
}

func (pv randomVariable) FirstNonZeroIndex() int {
	for i := 0; i < len(pv); i++ {
		if pv[i] > 0.0 {
			return i
		}
	}
	return -1
}

func sum(ls []int) int {
	res := 0
	for i := 0; i < len(ls); i++ {
		res = res + ls[i]
	}
	return res
}

func createRandomVariableNDiceSum(n int) randomVariable {
	combs := generateNDiceCombinations(n)

	l := len(combs)

	res := make([]float32, n*6)

	for i := 0; i < l; i++ {
		s := sum(combs[i]) - 1

		res[s] = res[s] + 1
	}

	for i := 0; i < len(res); i++ {
		res[i] = res[i] / float32(l)
	}

	return res
}

func generateNDiceCombinations(n int) [][]int {
	length := pow(6, n)
	res := make([][]int, length)

	for i := 0; i < length; i++ {
		subres := make([]int, n)
		res[i] = subres
	}

	s := length
	for i := 0; i < n; i++ {
		s = s / 6
		h := 1
		for j := 0; j < length; {
			if h > 6 {
				h = 1
			}
			for k := 0; k < s; k++ {
				res[j][i] = h
				j = j + 1
			}
			h = h + 1
		}
	}

	return res
}

type attacker struct {
	Mat int16
	PS  int16
}

type defender struct {
	Def int16
	Arm int16
}

type state struct {
	grid             grid
	remainingAttacks int
	probability      float32
}

func doAttack(s state, a attacker, d defender,
	hit randomVariable, dmg randomVariable, col randomVariable) []state {
	if s.remainingAttacks == 0 {
		res := []state{s}
		return res
	}
	fmt.Printf("Attack No %v\n", s.remainingAttacks)
	// Hit
	toHit := max(d.Def-a.Mat, 3)
	fmt.Printf("To Hit %v\n", toHit)

	pToHit := hit.GreaterEqual(int(toHit))

	hitSuccesState := state{s.grid, s.remainingAttacks - 1, s.probability * pToHit}
	hitFailiureState := state{s.grid, s.remainingAttacks - 1, s.probability * (1.0 - pToHit)}

	// Damage + Column
	armdif := a.PS - d.Arm
	fmt.Printf("Arm Malus %v\n", armdif)

	mindmg := int16(dmg.FirstNonZeroIndex()) + 1 + armdif
	fmt.Printf("Minimum Damage %v\n", mindmg)

	maxdmg := int16(len(dmg)) + armdif
	fmt.Printf("Maximum Damage %v\n", maxdmg)

	maxdmgpossible := min(maxdmg, s.grid.RemainingHealth())
	fmt.Printf("Maximum possible Damage %v\n", maxdmgpossible)

	stateno := max(1, max(0, maxdmgpossible-mindmg+1)*6) + 1
	fmt.Printf("States No %v\n", stateno)

	res := make([]state, stateno)

	for i := 0; i < int(stateno-1); i++ {
		res[i] = hitSuccesState
		//dmg := mindmg + int16(i)
		for j := 0; j < 6; j++ {

		}
	}
	res[stateno-1] = hitFailiureState

	return res
}

func main() {
	fmt.Printf("Hallo Welt.")
	fmt.Println()

	sg := simpleGrid(20)
	fmt.Printf("%v", sg.TakeDamage(15, 0))

	fmt.Println()

	cg := complexGrid([]int16{3, 4, 6, 6, 5, 3})
	fmt.Printf("%v", cg.TakeDamage(15, 3))

	fmt.Println()

	rv := createRandomVariableNDiceSum(2)
	sh := rv.GreaterEqual(7)
	fmt.Printf("%v\n", sh)

	dmgrv := createRandomVariableNDiceSum(4)
	colrv := createRandomVariableNDiceSum(1)

	chicken := defender{14, 14}
	bane := attacker{6, 13}

	chickenGrid := simpleGrid(20)
	initialState := state{chickenGrid, 1, 1.0}

	doAttack(initialState, bane, chicken, rv, dmgrv, colrv)
}
