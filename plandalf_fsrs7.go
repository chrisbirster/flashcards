package main

import (
	"encoding/binary"
	"errors"
	"math"
)

const plandalfMillisecondsPerDay int64 = 86_400_000

var plandalfFSRS7DefaultWeights = [35]float64{
	0.041, 2.4175, 4.1283, 11.9709, 5.6385, 0.4468, 3.262, 2.3054,
	0.1688, 1.3325, 0.3524, 0.0049, 0.7503, 0.0896, 0.6625, 1.3,
	0.882, 0.3072, 3.5875, 0.303, 0.0107, 0.2279, 2.6413, 0.5594,
	1.3, 2.5, 1.0, 0.0723, 0.1634, 0.5, 0.9555, 0.2245, 0.6232,
	0.1362, 0.3862,
}

type plandalfFSRS7Parameters struct {
	Weights             [35]float64
	DesiredRetention    float64
	MinimumIntervalDays float64
	MaximumIntervalDays float64
}

func defaultPlandalfFSRS7Parameters() plandalfFSRS7Parameters {
	return plandalfFSRS7Parameters{
		Weights:             plandalfFSRS7DefaultWeights,
		DesiredRetention:    0.9,
		MinimumIntervalDays: 1.0 / 86_400.0,
		MaximumIntervalDays: 36_500.0,
	}
}

type plandalfRating uint8

const (
	plandalfAgain plandalfRating = 1
	plandalfHard  plandalfRating = 2
	plandalfGood  plandalfRating = 3
	plandalfEasy  plandalfRating = 4
)

func parsePlandalfRating(value int) (plandalfRating, error) {
	if value < int(plandalfAgain) || value > int(plandalfEasy) {
		return 0, errors.New("rating must be 1 (Again), 2 (Hard), 3 (Good), or 4 (Easy)")
	}
	return plandalfRating(value), nil
}

type plandalfHistoryEntry struct {
	Rating       plandalfRating
	ReviewedAtMs int64
}

type plandalfMemoryState struct {
	StabilityDays float64
	Difficulty    float64
}

type plandalfReplayState struct {
	Memory           plandalfMemoryState
	LastReviewedAtMs int64
}

type plandalfCandidate struct {
	Rating       plandalfRating `json:"rating"`
	DueAtMs      int64          `json:"due_at_ms"`
	IntervalDays float64        `json:"interval_days"`
}

type plandalfSchedule struct {
	Again plandalfCandidate `json:"again"`
	Hard  plandalfCandidate `json:"hard"`
	Good  plandalfCandidate `json:"good"`
	Easy  plandalfCandidate `json:"easy"`
}

func clampFloat(value, minValue, maxValue float64) float64 {
	return math.Max(minValue, math.Min(value, maxValue))
}

func plandalfInitialDifficulty(rating plandalfRating, p plandalfFSRS7Parameters) float64 {
	value := p.Weights[4] - math.Exp(p.Weights[5]*(float64(rating)-1.0)) + 1.0
	return clampFloat(value, 1.0, 10.0)
}

func plandalfInitialMemoryState(rating plandalfRating, p plandalfFSRS7Parameters) plandalfMemoryState {
	return plandalfMemoryState{
		StabilityDays: p.Weights[int(rating)-1],
		Difficulty:    plandalfInitialDifficulty(rating, p),
	}
}

type plandalfCurve struct {
	Retention  float64
	Derivative float64
}

func plandalfForgettingCurve(elapsedDays, stabilityDays float64, p plandalfFSRS7Parameters) plandalfCurve {
	w := p.Weights
	decay1 := -w[27]
	decay2 := -w[28]
	base1 := w[29]
	base2 := w[30]
	c1 := math.Pow(base1, 1.0/decay1) - 1.0
	c2 := math.Pow(base2, 1.0/decay2) - 1.0
	tOverS := elapsedDays / stabilityDays
	inner1 := math.Max(1.0+c1*tOverS, 1e-12)
	inner2 := math.Max(1.0+c2*tOverS, 1e-12)
	r1 := math.Pow(inner1, decay1)
	r2 := math.Pow(inner2, decay2)
	weight1 := w[31] * math.Pow(stabilityDays, -w[33])
	weight2 := w[32] * math.Pow(stabilityDays, w[34])
	weightSum := math.Max(weight1+weight2, 1e-12)
	retention := clampFloat((weight1*r1+weight2*r2)/weightSum, 0.0, 1.0)
	dr1dt := decay1 * math.Pow(inner1, decay1-1.0) * (c1 / stabilityDays)
	dr2dt := decay2 * math.Pow(inner2, decay2-1.0) * (c2 / stabilityDays)
	derivative := math.Min((weight1*dr1dt+weight2*dr2dt)/weightSum, 0.0)
	return plandalfCurve{Retention: retention, Derivative: derivative}
}

func plandalfNextDifficulty(state plandalfMemoryState, rating plandalfRating, p plandalfFSRS7Parameters) float64 {
	delta := -p.Weights[6] * (float64(rating) - 3.0)
	damped := delta * (10.0 - state.Difficulty) / 9.0
	current := state.Difficulty + damped
	return clampFloat(0.01*plandalfInitialDifficulty(plandalfEasy, p)+0.99*current, 1.0, 10.0)
}

func plandalfComponentStability(state plandalfMemoryState, retention float64, rating plandalfRating, p plandalfFSRS7Parameters, base int) float64 {
	w := p.Weights
	stability := state.StabilityDays
	difficulty := state.Difficulty
	failed := w[base+3] * math.Pow(difficulty, -w[base+4]) *
		(math.Pow(stability+1.0, w[base+5])-1.0) * math.Exp((1.0-retention)*w[base+6])
	postLapse := math.Min(stability, failed)
	if rating == plandalfAgain {
		return postLapse
	}

	hardPenalty := 1.0
	if rating == plandalfHard {
		hardPenalty = w[base+7]
	}
	easyBonus := 1.0
	if rating == plandalfEasy {
		easyBonus = w[base+8]
	}
	increase := 1.0 + math.Exp(w[base]-1.5)*(11.0-difficulty)*
		math.Pow(stability, -w[base+1])*(math.Exp(math.Min((1.0-retention)*w[base+2], 30.0))-1.0)*
		hardPenalty*easyBonus
	return math.Max(postLapse, stability*increase)
}

func plandalfNextMemoryState(state plandalfMemoryState, elapsedDays float64, rating plandalfRating, p plandalfFSRS7Parameters) (plandalfMemoryState, error) {
	if math.IsNaN(elapsedDays) || math.IsInf(elapsedDays, 0) || elapsedDays < 0 {
		return plandalfMemoryState{}, errors.New("invalid elapsed time")
	}
	if math.IsNaN(state.StabilityDays) || math.IsInf(state.StabilityDays, 0) || state.StabilityDays <= 0 {
		return plandalfMemoryState{}, errors.New("invalid stability")
	}
	retention := plandalfForgettingCurve(elapsedDays, state.StabilityDays, p).Retention
	longTerm := plandalfComponentStability(state, retention, rating, p, 7)
	shortTerm := plandalfComponentStability(state, retention, rating, p, 16)
	coefficient := clampFloat(1.0-p.Weights[26]*math.Exp(-p.Weights[25]*elapsedDays), 0.0, 1.0)
	stability := clampFloat(coefficient*longTerm+(1.0-coefficient)*shortTerm, 0.0001, 36_500.0)
	return plandalfMemoryState{StabilityDays: stability, Difficulty: plandalfNextDifficulty(state, rating, p)}, nil
}

func plandalfIntervalForRetention(stabilityDays, target float64, p plandalfFSRS7Parameters) (float64, error) {
	if stabilityDays <= 0 || target <= 0 || target >= 1 {
		return 0, errors.New("invalid FSRS interval input")
	}
	minInterval := p.MinimumIntervalDays
	maxInterval := p.MaximumIntervalDays
	u := math.Log(math.Max(stabilityDays, 1e-10))
	minU := math.Log(minInterval)
	maxU := math.Log(maxInterval)
	for i := 0; i < 12; i++ {
		u = clampFloat(u, minU, maxU)
		interval := clampFloat(math.Exp(u), minInterval, maxInterval)
		curve := plandalfForgettingCurve(interval, stabilityDays, p)
		dfdu := math.Min(curve.Derivative*interval, -1e-12)
		u -= (curve.Retention - target) / dfdu
	}
	return clampFloat(math.Exp(clampFloat(u, minU, maxU)), minInterval, maxInterval), nil
}

func plandalfReplay(history []plandalfHistoryEntry, p plandalfFSRS7Parameters) (*plandalfReplayState, error) {
	if len(history) == 0 {
		return nil, nil
	}
	state := plandalfInitialMemoryState(history[0].Rating, p)
	previous := history[0].ReviewedAtMs
	for _, entry := range history[1:] {
		if entry.ReviewedAtMs < previous {
			return nil, errors.New("non-monotonic review history")
		}
		elapsedDays := float64(entry.ReviewedAtMs-previous) / float64(plandalfMillisecondsPerDay)
		next, err := plandalfNextMemoryState(state, elapsedDays, entry.Rating, p)
		if err != nil {
			return nil, err
		}
		state = next
		previous = entry.ReviewedAtMs
	}
	return &plandalfReplayState{Memory: state, LastReviewedAtMs: previous}, nil
}

func plandalfScheduleFor(history []plandalfHistoryEntry, nowMs int64, p plandalfFSRS7Parameters) (plandalfSchedule, error) {
	replayed, err := plandalfReplay(history, p)
	if err != nil {
		return plandalfSchedule{}, err
	}
	candidate := func(rating plandalfRating) (plandalfCandidate, error) {
		var state plandalfMemoryState
		if replayed == nil {
			state = plandalfInitialMemoryState(rating, p)
		} else {
			if nowMs < replayed.LastReviewedAtMs {
				return plandalfCandidate{}, errors.New("review time is before previous review")
			}
			elapsedDays := float64(nowMs-replayed.LastReviewedAtMs) / float64(plandalfMillisecondsPerDay)
			state, err = plandalfNextMemoryState(replayed.Memory, elapsedDays, rating, p)
			if err != nil {
				return plandalfCandidate{}, err
			}
		}
		interval, err := plandalfIntervalForRetention(state.StabilityDays, p.DesiredRetention, p)
		if err != nil {
			return plandalfCandidate{}, err
		}
		delta := int64(math.Round(interval * float64(plandalfMillisecondsPerDay)))
		return plandalfCandidate{Rating: rating, DueAtMs: nowMs + delta, IntervalDays: interval}, nil
	}

	again, err := candidate(plandalfAgain)
	if err != nil { return plandalfSchedule{}, err }
	hard, err := candidate(plandalfHard)
	if err != nil { return plandalfSchedule{}, err }
	good, err := candidate(plandalfGood)
	if err != nil { return plandalfSchedule{}, err }
	easy, err := candidate(plandalfEasy)
	if err != nil { return plandalfSchedule{}, err }
	return plandalfSchedule{Again: again, Hard: hard, Good: good, Easy: easy}, nil
}

func plandalfParameterSetID(p plandalfFSRS7Parameters) [32]byte {
	lanes := [4]uint64{0xcbf29ce484222325, 0x84222325cbf29ce4, 0x6a09e667f3bcc909, 0xbb67ae8584caa73b}
	mix := func(lane *uint64, value uint64) {
		*lane ^= value + 0x9e3779b97f4a7c15
		*lane *= 0x100000001b3
		*lane ^= *lane >> 29
		*lane *= 0xbf58476d1ce4e5b9
		*lane ^= *lane >> 31
	}
	for index, weight := range p.Weights {
		bits := math.Float64bits(weight)
		for laneIndex := range lanes {
			position := uint64(index + 1)
			laneNumber := uint64(laneIndex + 1)
			mix(&lanes[laneIndex], bits^(position*(laneNumber*0x517cc1b727220a95)))
		}
	}
	config := [4]uint64{
		math.Float64bits(p.DesiredRetention),
		math.Float64bits(p.MinimumIntervalDays),
		math.Float64bits(p.MaximumIntervalDays),
		7,
	}
	for i := range config {
		mix(&lanes[i], config[i])
	}
	var id [32]byte
	for i, lane := range lanes {
		binary.LittleEndian.PutUint64(id[i*8:(i+1)*8], lane)
	}
	return id
}
