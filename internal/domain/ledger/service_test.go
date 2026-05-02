package ledger

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestService_CreateScript(t *testing.T) {
	tests := []struct {
		name  string
		setup func(ctrl *gomock.Controller) *Service
		args  func() Config
		want  func(t *testing.T, scr Config, e error)
	}{
		{
			name: "when create new script with success",
			setup: func(ctrl *gomock.Controller) *Service {
				r := NewMockRepository(ctrl)
				r.EXPECT().FindConfigByLevel(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(Config{}, ErrConfigNotFound{}).Times(1)
				r.EXPECT().SaveConfig(gomock.Any(), gomock.Any()).Return(nil).Times(1)
				return New(r)
			},
			args: func() Config {
				return fakeScript(ProgramLevel, "201", "PAGAMENTO A VISTA")
			},
			want: func(t *testing.T, scr Config, e error) {
				assert.Nil(t, e)
			},
		},
		{
			name: "when duplicate entry",
			setup: func(ctrl *gomock.Controller) *Service {
				r := NewMockRepository(ctrl)
				r.EXPECT().FindConfigByLevel(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(fakeScript(ProgramLevel, "201", "PAGAMENTO A VISTA"), nil).Times(1)
				return New(r)
			},
			args: func() Config {
				fs := fakeScript(ProgramLevel, "201", "PAGAMENTO A VISTA")
				fs.Scripts = append(fs.Scripts, Script{
					ScriptID:    401,
					Description: "IOF",
					Expression:  "amount",
				})
				return fs
			},
			want: func(t *testing.T, scr Config, e error) {
				assert.NotNil(t, e)
				assert.Equal(t, "config was created with id: script-1234", e.Error())
			},
		},
		{
			name: "when receive repository error",
			setup: func(ctrl *gomock.Controller) *Service {
				r := NewMockRepository(ctrl)
				r.EXPECT().FindConfigByLevel(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(Config{}, ErrConfigNotFound{}).Times(1)
				r.EXPECT().SaveConfig(gomock.Any(), gomock.Any()).Return(errors.New("any repository error")).Times(1)
				return New(r)
			},
			args: func() Config {
				return fakeScript(ProgramLevel, "201", "PAGAMENTO A VISTA")
			},
			want: func(t *testing.T, scr Config, e error) {
				assert.NotNil(t, e)
				assert.Equal(t, "any repository error", e.Error())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			s := tt.setup(ctrl)
			cid := uuid.NewString()

			scr, err := s.CreateScript(context.Background(), cid, tt.args())
			tt.want(t, scr, err)
		})
	}
}

func TestService_UpdateScript(t *testing.T) {
	tests := []struct {
		name  string
		setup func(ctrl *gomock.Controller) *Service
		args  func() (string, Config)
		want  func(t *testing.T, scr Config, e error)
	}{
		{
			name: "when update saved script with success",
			setup: func(ctrl *gomock.Controller) *Service {
				r := NewMockRepository(ctrl)
				r.EXPECT().FindConfigByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(fakeScript(ProgramLevel, "201", "PAGAMENTO A VISTA"), nil).Times(1)
				r.EXPECT().UpdateConfig(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, s Config) error {
					if s.Description != "Changed Description" {
						return errors.New("script dont changed")
					}
					if s.Version != 2 {
						return errors.New("script dont changed")
					}
					return nil
				})
				return New(r)
			},
			args: func() (string, Config) {
				changed := fakeScript(PlatformLevel, "201", "Changed Description")
				changed.Description = "Changed Description"
				return "uuid-12345", changed
			},
			want: func(t *testing.T, scr Config, e error) {
				assert.Nil(t, e)
			},
		},
		{
			name: "when duplicate entry",
			setup: func(ctrl *gomock.Controller) *Service {
				r := NewMockRepository(ctrl)
				r.EXPECT().FindConfigByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(fakeScript(ProgramLevel, "201", "PAGAMENTO A VISTA"), nil).Times(1)
				return New(r)
			},
			args: func() (string, Config) {
				fs := fakeScript(PlatformLevel, "201", "PAGAMENTO A VISTA")
				fs.Scripts = append(fs.Scripts, Script{
					ScriptID:    401,
					Flow:        Regular,
					Description: "IOF",
					Expression:  "amount",
				})
				return "uuid-12345", fs
			},
			want: func(t *testing.T, scr Config, e error) {
				assert.NotNil(t, e)
				assert.Equal(t, "duplicated script: regular - 401 - IOF", e.Error())
			},
		},
		{
			name: "when receive not found script error",
			setup: func(ctrl *gomock.Controller) *Service {
				r := NewMockRepository(ctrl)
				r.EXPECT().FindConfigByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(Config{}, ErrConfigNotFound{}).Times(1)
				return New(r)
			},
			args: func() (string, Config) {
				return "uuid-12345", fakeScript(PlatformLevel, "201", "Changed Description")
			},
			want: func(t *testing.T, scr Config, e error) {
				assert.NotNil(t, e)
				assert.Equal(t, "ledger config not found", e.Error())
			},
		},
		{
			name: "when receive repository error",
			setup: func(ctrl *gomock.Controller) *Service {
				r := NewMockRepository(ctrl)
				r.EXPECT().FindConfigByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(fakeScript(ProgramLevel, "201", "PAGAMENTO A VISTA"), nil).Times(1)
				r.EXPECT().UpdateConfig(gomock.Any(), gomock.Any()).Return(errors.New("any repository error"))
				return New(r)
			},
			args: func() (string, Config) {
				return "uuid-12345", fakeScript(PlatformLevel, "201", "Changed Description")
			},
			want: func(t *testing.T, scr Config, e error) {
				assert.NotNil(t, e)
				assert.Equal(t, "any repository error", e.Error())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			s := tt.setup(ctrl)
			cid := uuid.NewString()
			id, sc := tt.args()
			scr, err := s.UpdateScript(context.Background(), cid, id, sc)
			tt.want(t, scr, err)
		})
	}
}

func TestService_FindConfigByLevel(t *testing.T) {
	tests := []struct {
		name  string
		setup func(ctrl *gomock.Controller) *Service
		args  func() (string, string, int64)
		want  func(t *testing.T, saved Config, e error)
	}{
		{
			name: "when retrieve a program level config with success",
			setup: func(ctrl *gomock.Controller) *Service {
				r := NewMockRepository(ctrl)
				r.EXPECT().FindConfigByLevel(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, level string, eventTypeID string, orgID string, programID *int64) (Config, error) {
					if level == string(ProgramLevel) && orgID == "TN-77add76c-e395-446b-b306-1a0f9cb99a31" && *programID == int64(1) {
						return fakeScript(ProgramLevel, eventTypeID, "PAGAMENTO A VISTA"), nil
					}
					return Config{}, ErrConfigNotFound{}
				}).Times(1)
				return New(r)
			},
			args: func() (string, string, int64) {
				return "201", "TN-77add76c-e395-446b-b306-1a0f9cb99a31", 1
			},
			want: func(t *testing.T, sc Config, e error) {
				assert.Nil(t, e)
				assert.Equal(t, ProgramLevel, sc.Level)
				assert.Equal(t, "201", sc.ProcessCode)
				assert.Equal(t, "TN-77add76c-e395-446b-b306-1a0f9cb99a31", sc.OrgID)
				assert.Equal(t, int64(1), sc.ProgramID)
				assert.Equal(t, int64(1), sc.Version)
			},
		},
		{
			name: "when retrieve a org level config with success",
			setup: func(ctrl *gomock.Controller) *Service {
				r := NewMockRepository(ctrl)
				r.EXPECT().FindConfigByLevel(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, level string, eventTypeID string, orgID string, programID *int64) (Config, error) {
					if level == string(TenantLevel) && orgID == "TN-77add76c-e395-446b-b306-1a0f9cb99a31" && *programID == int64(1) {
						return fakeScript(TenantLevel, eventTypeID, "PAGAMENTO A VISTA"), nil
					}
					return Config{}, ErrConfigNotFound{}
				}).Times(2)
				return New(r)
			},
			args: func() (string, string, int64) {
				return "201", "TN-77add76c-e395-446b-b306-1a0f9cb99a31", 1
			},
			want: func(t *testing.T, sc Config, e error) {
				assert.Nil(t, e)
				assert.Equal(t, TenantLevel, sc.Level)
				assert.Equal(t, "201", sc.ProcessCode)
				assert.Equal(t, "TN-77add76c-e395-446b-b306-1a0f9cb99a31", sc.OrgID)
				assert.Equal(t, int64(1), sc.Version)
			},
		},
		{
			name: "when receive a ledger config not found error",
			setup: func(ctrl *gomock.Controller) *Service {
				r := NewMockRepository(ctrl)
				r.EXPECT().FindConfigByLevel(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(Config{}, ErrConfigNotFound{}).Times(2)
				return New(r)
			},
			args: func() (string, string, int64) {
				return "201", "TN-77add76c-e395-446b-b306-1a0f9cb99a31", 1
			},
			want: func(t *testing.T, sc Config, e error) {
				assert.NotNil(t, e)
				assert.Equal(t, "ledger config not found", e.Error())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			s := tt.setup(ctrl)
			e, o, p := tt.args()
			saved, err := s.FindScriptByLevel(context.Background(), uuid.NewString(), e, o, p)
			tt.want(t, saved, err)
		})
	}
}

func TestService_FindAllScripts(t *testing.T) {
	tests := []struct {
		name  string
		setup func(ctrl *gomock.Controller) *Service
		args  func() (string, int64)
		want  func(t *testing.T, scripts []Config, e error)
	}{
		{
			name: "when retrieve all scritps with success",
			setup: func(ctrl *gomock.Controller) *Service {
				r := NewMockRepository(ctrl)
				r.EXPECT().FindAllConfigs(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, orgID string, programID *int64) ([]Config, error) {
					if orgID == "TN-77add76c-e395-446b-b306-1a0f9cb99a31" && *programID == int64(1) {
						return fakeSliceScripts(ProgramLevel, 100), nil
					}
					return nil, ErrConfigNotFound{}
				}).Times(1)
				return New(r)
			},
			args: func() (string, int64) {
				return "TN-77add76c-e395-446b-b306-1a0f9cb99a31", 1
			},
			want: func(t *testing.T, scripts []Config, e error) {
				assert.Nil(t, e)
				assert.Equal(t, 100, len(scripts))
			},
		},
		{
			name: "when retrieve not found script error",
			setup: func(ctrl *gomock.Controller) *Service {
				r := NewMockRepository(ctrl)
				r.EXPECT().FindAllConfigs(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, orgID string, programID *int64) ([]Config, error) {
					return nil, ErrConfigNotFound{}
				}).Times(1)
				return New(r)
			},
			args: func() (string, int64) {
				return "TN-77add76c-e395-446b-b306-1a0f9cb99a31", 1
			},
			want: func(t *testing.T, scripts []Config, e error) {
				assert.Nil(t, scripts)
				assert.NotNil(t, e)
				assert.Equal(t, "ledger config not found", e.Error())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			s := tt.setup(ctrl)
			o, p := tt.args()
			saved, err := s.FindAllScripts(context.Background(), uuid.NewString(), o, &p)
			tt.want(t, saved, err)
		})
	}
}

func fakeSliceScripts(level Level, quant int) []Config {
	result := make([]Config, 0, quant)
	for i := 1; i <= quant; i++ {
		id := i + 100
		result = append(result, fakeScript(level, fmt.Sprint(id), fmt.Sprintf("Transaction %v", id)))
	}
	return result
}

func fakeScript(level Level, event string, description string) Config {
	e, _ := strconv.ParseInt(event, 10, 64)
	return Config{
		ConfigID:    "script-1234",
		Level:       level,
		ProcessCode: event,
		OrgID:       "TN-77add76c-e395-446b-b306-1a0f9cb99a31",
		ProgramID:   1,
		Description: description,
		Scripts: []Script{
			{
				ScriptID:    e,
				Flow:        Regular,
				Description: description,
				Expression:  "Amount.amount",
			},
			{
				ScriptID:    401,
				Flow:        Regular,
				Description: "IOF",
				Expression:  "Amount.amount",
			},
		},
		Enable:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Version:   1,
	}
}
