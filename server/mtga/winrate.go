package mtga

import (
	"math/rand"
	"runtime"
	"sync"

	"github.com/Manbeardo/mtga-helper/server/mtga/formats"
)

type EventSimulationConfig struct {
	Format                  formats.Format
	GameWinRate             float64
	PerEventWinRateVariance float64
}

type EventSimulationResult struct {
	EventsPlayed int
	GamesWon     int
	GamesLost    int
	MatchesWon   int
	MatchesLost  int
	GemsSpent    int
	GemsWon      int
	PacksWon     int
	PacksOpened  int
}

func (r EventSimulationResult) EconomyStats() EventEconomyStats {
	mythicInsertionRate := 1.0 / 7.0
	wildcardInsertionRate := 1.0 / 15.0
	avgPacksOpened := (float64(r.PacksOpened) / float64(r.EventsPlayed))
	avgRaresFromPacksOpened := avgPacksOpened * (1.0 - mythicInsertionRate)
	avgMythicsFromPacksOpened := avgPacksOpened * mythicInsertionRate

	avgPacksWon := (float64(r.PacksWon) / float64(r.EventsPlayed))
	avgRaresFromPacksWon := avgPacksWon * (1.0 - wildcardInsertionRate) * (1.0 - mythicInsertionRate)
	avgMythicsFromPacksWon := avgPacksWon * (1.0 - wildcardInsertionRate) * mythicInsertionRate

	avgGemsFromExcessRares := ((avgRaresFromPacksOpened + avgRaresFromPacksWon) * 20.0) +
		((avgMythicsFromPacksOpened + avgMythicsFromPacksWon) * 40.0)

	return EventEconomyStats{
		AvgGemsWon:             float64(r.GemsWon) / float64(r.EventsPlayed),
		AvgPacksWon:            float64(r.PacksWon) / float64(r.EventsPlayed),
		AvgPacksOpened:         avgPacksOpened,
		AvgGemsFromExcessRares: avgGemsFromExcessRares,
	}
}

type EventEconomyStats struct {
	AvgGemsWon             float64
	AvgPacksWon            float64
	AvgPacksOpened         float64
	AvgGemsFromExcessRares float64
}

func MergeSimulationResults(results ...EventSimulationResult) EventSimulationResult {
	if len(results) == 0 {
		return EventSimulationResult{}
	} else if len(results) == 1 {
		return results[0]
	}
	a, b := results[0], results[1]
	combined := EventSimulationResult{
		EventsPlayed: a.EventsPlayed + b.EventsPlayed,
		GamesWon:     a.GamesWon + b.GamesWon,
		GamesLost:    a.GamesLost + b.GamesLost,
		MatchesWon:   a.MatchesWon + b.MatchesWon,
		MatchesLost:  a.MatchesLost + b.MatchesLost,
		GemsSpent:    a.GemsSpent + b.GemsSpent,
		GemsWon:      a.GemsWon + b.GemsWon,
		PacksWon:     a.PacksWon + b.PacksWon,
		PacksOpened:  a.PacksOpened + b.PacksOpened,
	}
	if len(results) == 2 {
		return combined
	}
	return MergeSimulationResults(append(results[2:], combined)...)
}

func SimulateEvents(cfg EventSimulationConfig, count int) EventSimulationResult {
	feederChan := make(chan struct{})
	accumulationChan := make(chan EventSimulationResult)

	go func() {
		for i := 0; i < count; i++ {
			feederChan <- struct{}{}
		}
		close(feederChan)
	}()

	workerCount := runtime.NumCPU()
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(workerCount)

	for i := 0; i < workerCount; i++ {
		go func() {
			for range feederChan {
				accumulationChan <- SimulateEvent(cfg)
			}
			waitGroup.Done()
		}()
	}

	go func() {
		waitGroup.Wait()
		close(accumulationChan)
	}()

	mergedResult := EventSimulationResult{}
	for result := range accumulationChan {
		mergedResult = MergeSimulationResults(mergedResult, result)
	}

	return mergedResult
}

func SimulateEvent(cfg EventSimulationConfig) EventSimulationResult {
	winRate := cfg.GameWinRate + (cfg.PerEventWinRateVariance / 2.0) - (rand.Float64() * cfg.PerEventWinRateVariance)

	result := EventSimulationResult{
		EventsPlayed: 1,
		GemsSpent:    cfg.Format.EntryFeeGems(),
		PacksOpened:  cfg.Format.PacksOpened(),
	}

	for result.MatchesLost < cfg.Format.MaxLosses() && result.MatchesWon < len(cfg.Format.Prizes())-1 {
		gameWins, gameLosses := 0, 0
		var maxGameWins, maxGameLosses int
		if cfg.Format.MatchKind() == formats.BestOfOne {
			maxGameWins, maxGameLosses = 1, 1
		} else if cfg.Format.MatchKind() == formats.BestOfThree {
			maxGameWins, maxGameLosses = 2, 2
		} else {
			panic("unknown match kind: " + cfg.Format.MatchKind())
		}
		for gameWins < maxGameWins && gameLosses < maxGameLosses {
			if rand.Float64() < winRate {
				gameWins++
			} else {
				gameLosses++
			}
		}
		result.GamesWon += gameWins
		result.GamesLost += gameLosses
		if gameWins == maxGameWins {
			result.MatchesWon++
		} else {
			result.MatchesLost++
		}
	}

	prize := cfg.Format.Prizes()[result.MatchesWon]
	result.GemsWon = prize.Gems()
	result.PacksWon = prize.Packs()

	return result
}
