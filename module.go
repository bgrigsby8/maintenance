package maintenance

import (
	"context"
	"errors"
	"fmt"
	"time"

	sensor "go.viam.com/rdk/components/sensor"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
)

var (
	TimeAllowed      = resource.NewModel("brad-grigsby", "maintenance", "time-allowed")
	errUnimplemented = errors.New("unimplemented")
)

func init() {
	resource.RegisterComponent(sensor.API, TimeAllowed,
		resource.Registration[sensor.Sensor, *Config]{
			Constructor: newMaintenanceTimeAllowed,
		},
	)
}

type Config struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

// Validate ensures all parts of the config are valid and important fields exist.
// Returns three values:
//  1. Required dependencies: other resources that must exist for this resource to work.
//  2. Optional dependencies: other resources that may exist but are not required.
//  3. An error if any Config fields are missing or invalid.
//
// The `path` parameter indicates
// where this resource appears in the machine's JSON configuration
// (for example, "components.0"). You can use it in error messages
// to indicate which resource has a problem.
func (cfg *Config) Validate(path string) ([]string, []string, error) {
	// Add config validation code here
	if cfg.StartTime == "" {
		return nil, nil, resource.NewConfigValidationFieldRequiredError(path, "start_time")
	}
	if cfg.EndTime == "" {
		return nil, nil, resource.NewConfigValidationFieldRequiredError(path, "end_time")
	}
	if _, err := time.Parse("15:04", cfg.StartTime); err != nil {
		return nil, nil, fmt.Errorf("%s.start_time: must be in HH:MM (24-hour) format", path)
	}
	if _, err := time.Parse("15:04", cfg.EndTime); err != nil {
		return nil, nil, fmt.Errorf("%s.end_time: must be in HH:MM (24-hour) format", path)
	}

	return nil, nil, nil
}

type maintenanceTimeAllowed struct {
	resource.AlwaysRebuild

	name resource.Name

	logger logging.Logger
	cfg    *Config

	cancelCtx  context.Context
	cancelFunc func()
}

func newMaintenanceTimeAllowed(ctx context.Context, deps resource.Dependencies, rawConf resource.Config, logger logging.Logger) (sensor.Sensor, error) {
	conf, err := resource.NativeConfig[*Config](rawConf)
	if err != nil {
		return nil, err
	}

	return NewTimeAllowed(ctx, deps, rawConf.ResourceName(), conf, logger)

}

func NewTimeAllowed(ctx context.Context, deps resource.Dependencies, name resource.Name, conf *Config, logger logging.Logger) (sensor.Sensor, error) {

	cancelCtx, cancelFunc := context.WithCancel(context.Background())

	s := &maintenanceTimeAllowed{
		name:       name,
		logger:     logger,
		cfg:        conf,
		cancelCtx:  cancelCtx,
		cancelFunc: cancelFunc,
	}
	return s, nil
}

func (s *maintenanceTimeAllowed) Name() resource.Name {
	return s.name
}

func (s *maintenanceTimeAllowed) Readings(ctx context.Context, extra map[string]interface{}) (map[string]interface{}, error) {
	now := time.Now()

	start, err := time.Parse("15:04", s.cfg.StartTime)
	if err != nil {
		return nil, fmt.Errorf("invalid start_time %q: %w", s.cfg.StartTime, err)
	}
	end, err := time.Parse("15:04", s.cfg.EndTime)
	if err != nil {
		return nil, fmt.Errorf("invalid end_time %q: %w", s.cfg.EndTime, err)
	}

	loc := now.Location()
	y, m, d := now.Date()
	startToday := time.Date(y, m, d, start.Hour(), start.Minute(), 0, 0, loc)
	endToday := time.Date(y, m, d, end.Hour(), end.Minute(), 0, 0, loc)

	var allowed bool
	if !startToday.After(endToday) {
		// Same-day range e.g. 08:00–17:00
		allowed = !now.Before(startToday) && now.Before(endToday)
	} else {
		// Overnight range e.g. 22:00–06:00
		allowed = !now.Before(startToday) || now.Before(endToday)
	}

	return map[string]interface{}{"is_allowed": allowed}, nil
}

func (s *maintenanceTimeAllowed) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *maintenanceTimeAllowed) Status(ctx context.Context) (map[string]interface{}, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *maintenanceTimeAllowed) Close(context.Context) error {
	// Put close code here
	s.cancelFunc()
	return nil
}
