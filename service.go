package statservice

import (
	"context"
	"runtime"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes/any"
)

type Service struct {
	mutex        *sync.Mutex
	internalFunc func() *any.Any
}

var (
	lastSampleTime time.Time
	lastPauseNs    uint64  = 0
	lastNumGc      uint32  = 0
	nsInMs         float64 = float64(time.Millisecond)
)

func New(internalFunc func() *any.Any) *Service {
	return &Service{mutex: &sync.Mutex{}, internalFunc: internalFunc}
}

func (s *Service) Get(ctx context.Context, req *GetRequest) (*GetResponse, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return &GetResponse{
		Time:         time.Now().Format(time.RFC3339Nano),
		GoVersion:    runtime.Version(),
		GoOs:         runtime.GOOS,
		GoArch:       runtime.GOARCH,
		CpuNum:       int32(runtime.NumCPU()),
		GoMaxProcs:   int32(runtime.GOMAXPROCS(0)),
		RuntimeStat:  getRuntimeStat(),
		InternalStat: s.internalFunc(),
	}, nil
}

func getRuntimeStat() *RuntimeStat {
	now := time.Now()
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	var gcPausePerSecond float64

	if lastPauseNs > 0 {
		pauseSinceLastSample := mem.PauseTotalNs - lastPauseNs
		gcPausePerSecond = float64(pauseSinceLastSample) / nsInMs
	}

	lastPauseNs = mem.PauseTotalNs
	countGc := int(mem.NumGC - lastNumGc)

	var gcPerSecond float64
	if lastNumGc > 0 {
		diff := float64(countGc)
		diffTime := now.Sub(lastSampleTime).Seconds()
		gcPerSecond = diff / diffTime
	}

	if countGc > 256 {
		countGc = 256
	}

	gcPause := make([]float64, countGc)
	for i := 0; i < countGc; i++ {
		idx := int((mem.NumGC-uint32(i))+255) % 256
		pause := float64(mem.PauseNs[idx])
		gcPause[i] = pause / nsInMs
	}

	lastNumGc = mem.NumGC
	lastSampleTime = time.Now()

	return &RuntimeStat{
		Time:             now.Format(time.RFC3339Nano),
		GoroutineNum:     int32(runtime.NumGoroutine()),
		CgoCallNum:       int32(runtime.NumCgoCall()),
		MemoryAlloc:      mem.Alloc,
		MemoryTotalAlloc: mem.TotalAlloc,
		MemorySys:        mem.Sys,
		MemoryLookups:    mem.Lookups,
		MemoryMallocs:    mem.Mallocs,
		MemoryFrees:      mem.Frees,
		StackInUse:       mem.StackInuse,
		HeapAlloc:        mem.HeapAlloc,
		HeapSys:          mem.HeapSys,
		HeapIdle:         mem.HeapIdle,
		HeapInUse:        mem.HeapInuse,
		HeapReleased:     mem.HeapReleased,
		HeapObjects:      mem.HeapObjects,
		GcNext:           mem.NextGC,
		GcLast:           mem.LastGC,
		GcNum:            mem.NumGC,
		GcPause:          gcPause,
		GcPerSecond:      gcPerSecond,
		GcPausePerSecond: gcPausePerSecond,
	}
}
