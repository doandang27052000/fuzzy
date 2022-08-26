package fuzzy

import "fmt"

const (
	MustBeLowWorkload   = 10
	MustNotBeLowWorload = 20

	MustNotBeMiddleLowWorkload  = 10
	MustBeMiddleLowWorkload     = 20
	MustBeMiddleHighWorkload    = 70
	MustNotBeMiddleHighWorkload = 85

	MustNotBeHighWorkload = 70
	MustBeHighWorkload    = 85
)

const (
	MustBeBadResponeTime    = 20
	MustNotBeBadResponeTime = 30

	MustNotBeLOkayResponeTime = 20
	MustBeLOkayResponeTime    = 30
	MustBeROkayResponeTime    = 65
	MustNotBeROkayResponeTime = 80

	MustNotBeGoodResponeTime = 65
	MustBeGoodResponeTime    = 80
)
const (
	remove2vm = -2
	remove1vm = -1
	add1vm    = 1
	add2vm    = 2
	nochange  = 0
)

// Fuzzy is the main interface of a fuzzy logic algorithm
type Fuzzy interface {
	Fuzzification(number *FuzzyNumber) error
	Defuzzification(number *FuzzyNumber) error
	Inference(number *FuzzyNumber) error
}

// Family is the data needed for this programe.
type State struct {
	Number      string
	Workload    float64
	ResponeTine float64
}

// FuzzyNumber is the struct thatwill continually holds any fuzzy data.
type FuzzyNumber struct {
	State State

	WorkloadMembership    []float64
	ResponeTimeMembership []float64
	FiringLevel           [3][3]float64
	CrispValue            float64
}

type SA struct {
}

// Defuzzification is a function that will transfer fuzzy linguistic to crisp data.
func (b *SA) Defuzzification(number *FuzzyNumber) error {
	Action := [3][3]float64{
		{1, 1, 2},
		{-1, 0, 1},
		{-2, -1, 0},
	}
	number.CrispValue = 0
	var i, j int
	for i = 0; i < 3; i++ {
		for j = 0; j < 3; j++ {
			number.CrispValue += number.FiringLevel[i][j] * Action[i][j]
		}
	}
	return nil
}

func (b *SA) Inference(number *FuzzyNumber) error {
	var i, j int
	for i = 0; i < 3; i++ {
		for j = 0; j < 3; j++ {
			number.FiringLevel[i][j] = number.ResponeTimeMembership[i] * number.WorkloadMembership[j]
		}
	}

	return nil
}

func (b *SA) Fuzzification(number *FuzzyNumber) error {
	number.WorkloadMembership = append(number.WorkloadMembership, b.WorkloadLow(number.State.Workload))
	number.WorkloadMembership = append(number.WorkloadMembership, b.WorkloadMiddle(number.State.Workload))
	number.WorkloadMembership = append(number.WorkloadMembership, b.WorkloadHigh(number.State.Workload))

	number.ResponeTimeMembership = append(number.ResponeTimeMembership, b.ResponeTimeBad(number.State.ResponeTine))
	number.ResponeTimeMembership = append(number.ResponeTimeMembership, b.ResponeTimeOkay(number.State.ResponeTine))
	number.ResponeTimeMembership = append(number.ResponeTimeMembership, b.ResponeTimeGood(number.State.ResponeTine))
	return nil
}

func (b *SA) WorkloadLow(workload float64) float64 {
	if workload <= MustBeLowWorkload {
		return 1
	} else if workload > MustNotBeLowWorload {
		return 0
	}
	return 1 - (float64(workload-MustBeLowWorkload) / float64(MustNotBeLowWorload-MustBeLowWorkload))
}

func (b *SA) WorkloadMiddle(workload float64) float64 {
	if workload > MustNotBeMiddleLowWorkload && workload <= MustBeMiddleHighWorkload {
		return 1
	} else if workload < MustNotBeMiddleLowWorkload || workload > MustNotBeMiddleHighWorkload {
		return 0
	} else if workload < MustBeMiddleLowWorkload && workload >= MustNotBeMiddleLowWorkload {
		return float64(workload-MustNotBeMiddleLowWorkload) / float64(MustBeMiddleLowWorkload-MustNotBeMiddleLowWorkload)
	}

	return 1 - float64(workload-MustBeMiddleHighWorkload)/float64(MustNotBeMiddleHighWorkload-MustBeMiddleHighWorkload)
}

func (b *SA) WorkloadHigh(workload float64) float64 {
	if workload <= MustNotBeHighWorkload {
		return 0
	} else if workload > MustBeHighWorkload {
		return 1
	}
	return float64(workload-MustNotBeHighWorkload) / float64(MustBeHighWorkload-MustNotBeHighWorkload)
}

func (b *SA) ResponeTimeBad(responeTime float64) float64 {
	if responeTime <= MustBeBadResponeTime {
		return 1
	} else if responeTime > MustNotBeBadResponeTime {
		return 0
	}
	return 1 - (float64(responeTime-MustBeBadResponeTime) / float64(MustNotBeBadResponeTime-MustBeBadResponeTime))
}

func (b *SA) ResponeTimeOkay(responeTime float64) float64 {
	if responeTime > MustBeLOkayResponeTime && responeTime <= MustBeROkayResponeTime {
		return 1
	} else if responeTime < MustNotBeLOkayResponeTime || responeTime > MustNotBeROkayResponeTime {
		return 0
	} else if responeTime < MustBeLOkayResponeTime && responeTime >= MustNotBeLOkayResponeTime {
		return float64(responeTime-MustNotBeLOkayResponeTime) / float64(MustBeLOkayResponeTime-MustNotBeLOkayResponeTime)
	}

	return 1 - (float64(responeTime-MustBeROkayResponeTime) / float64(MustNotBeROkayResponeTime-MustBeROkayResponeTime))
}

func (b *SA) ResponeTimeGood(responeTime float64) float64 {
	if responeTime <= MustNotBeGoodResponeTime {
		return 0
	} else if responeTime > MustBeGoodResponeTime {
		return 1
	}
	return float64(responeTime-MustNotBeGoodResponeTime) / float64(MustBeGoodResponeTime-MustNotBeGoodResponeTime)
}

func test(state State, a FuzzyNumber, b SA) {
	state.ResponeTine = 28.12
	state.Workload = 77.34
	b.Fuzzification(&a)
	b.Inference(&a)
	b.Defuzzification(&a)
	fmt.Println(state)
	fmt.Println(a.ResponeTimeMembership, a.WorkloadMembership)
	fmt.Println(a.FiringLevel)
	fmt.Println(a.CrispValue)
}
