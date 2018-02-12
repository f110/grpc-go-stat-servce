package statservice

import (
	"context"
	"github.com/golang/protobuf/ptypes/any"
	"runtime"
	"sync"
	"time"
)

type Service struct {
	mutex        *sync.Mutex
	internalFunc func() *any.Any
}

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
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	t, _ := time.Now().MarshalText()
	return &RuntimeStat{
		Time:             string(t),
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
	}
}
