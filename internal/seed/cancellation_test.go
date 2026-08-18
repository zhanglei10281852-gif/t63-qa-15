package seed

import (
	"context"
	"testing"
	"time"

	"sanitation-operations/internal/clock"
	"sanitation-operations/internal/domain/vehicle"
	"sanitation-operations/internal/pagination"
	"sanitation-operations/internal/repository"
)

type cancellationSeedStore struct {
	repository.Store
	started chan struct{}
}

func (s cancellationSeedStore) ListVehicles(context.Context, repository.VehicleFilter, pagination.Query) (pagination.Result[vehicle.Vehicle], error) {
	return pagination.Result[vehicle.Vehicle]{}, nil
}

func (s cancellationSeedStore) SaveVehicle(ctx context.Context, _ vehicle.Vehicle, _ int) error {
	close(s.started)
	<-ctx.Done()
	return ctx.Err()
}

type fixedSeedIDs struct{}

func (fixedSeedIDs) NewID(string) string { return "seed-id" }

func TestSeedStopsWhenBootstrapContextIsCancelled(t *testing.T) {
	started := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- Ensure(ctx, cancellationSeedStore{started: started}, fixedSeedIDs{}, clock.Fixed{Current: time.Now()})
	}()
	<-started
	cancel()
	select {
	case err := <-done:
		if err == nil || err != context.Canceled {
			t.Fatalf("seed error = %v", err)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("seed did not stop after cancellation")
	}
}
