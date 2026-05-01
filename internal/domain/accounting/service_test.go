package accounting

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
		args  func() Script
		want  func(t *testing.T, scr Script, e error)
	}{
		{
			name: "when create new script with success",
			setup: func(ctrl *gomock.Controller) *Service {
				r := NewMockRepository(ctrl)
				r.EXPECT().FindScriptByLevel(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(Script{}, ErrScriptNotFound{}).Times(1)
				r.EXPECT().SaveScript(gomock.Any(), gomock.Any()).Return(nil).Times(1)
				return New(r)
			},
			args: func() Script {
				return fakeScript(ProgramLevel, "201", "PAGAMENTO A VISTA")
			},
			want: func(t *testing.T, scr Script, e error) {
				assert.Nil(t, e)
			},
		},
		{
			name: "when duplicate entry",
			setup: func(ctrl *gomock.Controller) *Service {
				r := NewMockRepository(ctrl)
				r.EXPECT().FindScriptByLevel(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(fakeScript(ProgramLevel, "201", "PAGAMENTO A VISTA"), nil).Times(1)
				return New(r)
			},
			args: func() Script {
				fs := fakeScript(ProgramLevel, "201", "PAGAMENTO A VISTA")
				fs.Entries = append(fs.Entries, Entry{
					EntryTypeID: 401,
					Description: "IOF",
					AmountName:  "amount",
				})
				return fs
			},
			want: func(t *testing.T, scr Script, e error) {
				assert.NotNil(t, e)
				assert.Equal(t, "script was created with id: script-1234", e.Error())
			},
		},
		{
			name: "when receive repository error",
			setup: func(ctrl *gomock.Controller) *Service {
				r := NewMockRepository(ctrl)
				r.EXPECT().FindScriptByLevel(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(Script{}, ErrScriptNotFound{}).Times(1)
				r.EXPECT().SaveScript(gomock.Any(), gomock.Any()).Return(errors.New("any repository error")).Times(1)
				return New(r)
			},
			args: func() Script {
				return fakeScript(ProgramLevel, "201", "PAGAMENTO A VISTA")
			},
			want: func(t *testing.T, scr Script, e error) {
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
		args  func() (string, Script)
		want  func(t *testing.T, scr Script, e error)
	}{
		{
			name: "when update saved script with success",
			setup: func(ctrl *gomock.Controller) *Service {
				r := NewMockRepository(ctrl)
				r.EXPECT().FindScriptByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(fakeScript(ProgramLevel, "201", "PAGAMENTO A VISTA"), nil).Times(1)
				r.EXPECT().UpdateScript(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, s Script) error {
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
			args: func() (string, Script) {
				changed := fakeScript(PlatformLevel, "201", "Changed Description")
				changed.Description = "Changed Description"
				return "uuid-12345", changed
			},
			want: func(t *testing.T, scr Script, e error) {
				assert.Nil(t, e)
			},
		},
		{
			name: "when duplicate entry",
			setup: func(ctrl *gomock.Controller) *Service {
				r := NewMockRepository(ctrl)
				r.EXPECT().FindScriptByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(fakeScript(ProgramLevel, "201", "PAGAMENTO A VISTA"), nil).Times(1)
				return New(r)
			},
			args: func() (string, Script) {
				fs := fakeScript(PlatformLevel, "201", "PAGAMENTO A VISTA")
				fs.Entries = append(fs.Entries, Entry{
					EntryTypeID: 401,
					Description: "IOF",
					AmountName:  "amount",
				})
				return "uuid-12345", fs
			},
			want: func(t *testing.T, scr Script, e error) {
				assert.NotNil(t, e)
				assert.Equal(t, "duplicated entry: 401 - IOF", e.Error())
			},
		},
		{
			name: "when receive not found script error",
			setup: func(ctrl *gomock.Controller) *Service {
				r := NewMockRepository(ctrl)
				r.EXPECT().FindScriptByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(Script{}, ErrScriptNotFound{}).Times(1)
				return New(r)
			},
			args: func() (string, Script) {
				return "uuid-12345", fakeScript(PlatformLevel, "201", "Changed Description")
			},
			want: func(t *testing.T, scr Script, e error) {
				assert.NotNil(t, e)
				assert.Equal(t, "accounting script not found", e.Error())
			},
		},
		{
			name: "when receive repository error",
			setup: func(ctrl *gomock.Controller) *Service {
				r := NewMockRepository(ctrl)
				r.EXPECT().FindScriptByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(fakeScript(ProgramLevel, "201", "PAGAMENTO A VISTA"), nil).Times(1)
				r.EXPECT().UpdateScript(gomock.Any(), gomock.Any()).Return(errors.New("any repository error"))
				return New(r)
			},
			args: func() (string, Script) {
				return "uuid-12345", fakeScript(PlatformLevel, "201", "Changed Description")
			},
			want: func(t *testing.T, scr Script, e error) {
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

func TestService_FindScriptByLevel(t *testing.T) {
	tests := []struct {
		name  string
		setup func(ctrl *gomock.Controller) *Service
		args  func() (string, string, int64)
		want  func(t *testing.T, saved Script, e error)
	}{
		{
			name: "when retrieve a program level scritp with success",
			setup: func(ctrl *gomock.Controller) *Service {
				r := NewMockRepository(ctrl)
				r.EXPECT().FindScriptByLevel(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, level string, eventTypeID string, orgID string, programID *int64) (Script, error) {
					if level == string(ProgramLevel) && orgID == "TN-77add76c-e395-446b-b306-1a0f9cb99a31" && *programID == int64(1) {
						return fakeScript(ProgramLevel, eventTypeID, "PAGAMENTO A VISTA"), nil
					}
					return Script{}, ErrScriptNotFound{}
				}).Times(1)
				return New(r)
			},
			args: func() (string, string, int64) {
				return "201", "TN-77add76c-e395-446b-b306-1a0f9cb99a31", 1
			},
			want: func(t *testing.T, sc Script, e error) {
				assert.Nil(t, e)
				assert.Equal(t, ProgramLevel, sc.Level)
				assert.Equal(t, "201", sc.EventTypeID)
				assert.Equal(t, "TN-77add76c-e395-446b-b306-1a0f9cb99a31", sc.OrgID)
				assert.Equal(t, int64(1), sc.ProgramID)
				assert.Equal(t, int64(1), sc.Version)
			},
		},
		{
			name: "when retrieve a org level scritp with success",
			setup: func(ctrl *gomock.Controller) *Service {
				r := NewMockRepository(ctrl)
				r.EXPECT().FindScriptByLevel(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, level string, eventTypeID string, orgID string, programID *int64) (Script, error) {
					if level == string(TenantLevel) && orgID == "TN-77add76c-e395-446b-b306-1a0f9cb99a31" && *programID == int64(1) {
						return fakeScript(TenantLevel, eventTypeID, "PAGAMENTO A VISTA"), nil
					}
					return Script{}, ErrScriptNotFound{}
				}).Times(2)
				return New(r)
			},
			args: func() (string, string, int64) {
				return "201", "TN-77add76c-e395-446b-b306-1a0f9cb99a31", 1
			},
			want: func(t *testing.T, sc Script, e error) {
				assert.Nil(t, e)
				assert.Equal(t, TenantLevel, sc.Level)
				assert.Equal(t, "201", sc.EventTypeID)
				assert.Equal(t, "TN-77add76c-e395-446b-b306-1a0f9cb99a31", sc.OrgID)
				assert.Equal(t, int64(1), sc.Version)
			},
		},
		{
			name: "when retrieve a platform level scritp with success",
			setup: func(ctrl *gomock.Controller) *Service {
				r := NewMockRepository(ctrl)
				r.EXPECT().FindScriptByLevel(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, level string, eventTypeID string, orgID string, programID *int64) (Script, error) {
					if level == string(PlatformLevel) && orgID == "PISMO" && *programID == int64(1) {
						return fakeScript(PlatformLevel, eventTypeID, "PAGAMENTO A VISTA"), nil
					}
					return Script{}, ErrScriptNotFound{}
				}).Times(3)
				return New(r)
			},
			args: func() (string, string, int64) {
				return "201", "TN-77add76c-e395-446b-b306-1a0f9cb99a31", 1
			},
			want: func(t *testing.T, sc Script, e error) {
				assert.Nil(t, e)
				assert.Equal(t, PlatformLevel, sc.Level)
				assert.Equal(t, "201", sc.EventTypeID)
				assert.Equal(t, int64(1), sc.Version)
			},
		},
		{
			name: "when receive a script not found error",
			setup: func(ctrl *gomock.Controller) *Service {
				r := NewMockRepository(ctrl)
				r.EXPECT().FindScriptByLevel(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, level string, eventTypeID string, orgID string, programID *int64) (Script, error) {
					return Script{}, ErrScriptNotFound{}
				}).Times(3)
				return New(r)
			},
			args: func() (string, string, int64) {
				return "201", "TN-77add76c-e395-446b-b306-1a0f9cb99a31", 1
			},
			want: func(t *testing.T, sc Script, e error) {
				assert.NotNil(t, e)
				assert.Equal(t, "accounting script not found", e.Error())
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
		want  func(t *testing.T, scripts []Script, e error)
	}{
		{
			name: "when retrieve all scritps with success",
			setup: func(ctrl *gomock.Controller) *Service {
				r := NewMockRepository(ctrl)
				r.EXPECT().FindAllScripts(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, orgID string, programID *int64) ([]Script, error) {
					if orgID == "TN-77add76c-e395-446b-b306-1a0f9cb99a31" && *programID == int64(1) {
						return fakeSliceScripts(ProgramLevel, 100), nil
					}
					return nil, ErrScriptNotFound{}
				}).Times(1)
				return New(r)
			},
			args: func() (string, int64) {
				return "TN-77add76c-e395-446b-b306-1a0f9cb99a31", 1
			},
			want: func(t *testing.T, scripts []Script, e error) {
				assert.Nil(t, e)
				assert.Equal(t, 100, len(scripts))
			},
		},
		{
			name: "when retrieve not found script error",
			setup: func(ctrl *gomock.Controller) *Service {
				r := NewMockRepository(ctrl)
				r.EXPECT().FindAllScripts(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, orgID string, programID *int64) ([]Script, error) {
					return nil, ErrScriptNotFound{}
				}).Times(1)
				return New(r)
			},
			args: func() (string, int64) {
				return "TN-77add76c-e395-446b-b306-1a0f9cb99a31", 1
			},
			want: func(t *testing.T, scripts []Script, e error) {
				assert.Nil(t, scripts)
				assert.NotNil(t, e)
				assert.Equal(t, "accounting script not found", e.Error())
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

func fakeSliceScripts(level Level, quant int) []Script {
	result := make([]Script, 0, quant)
	for i := 1; i <= quant; i++ {
		id := i + 100
		result = append(result, fakeScript(level, fmt.Sprint(id), fmt.Sprintf("Transaction %v", id)))
	}
	return result
}

func fakeScript(level Level, event string, description string) Script {
	e, _ := strconv.ParseInt(event, 10, 64)
	return Script{
		ScriptID:    "script-1234",
		Level:       level,
		EventTypeID: event,
		OrgID:       "TN-77add76c-e395-446b-b306-1a0f9cb99a31",
		ProgramID:   1,
		Description: description,
		Company: &Company{
			Code: "1234",
			Type: "BR",
		},
		Entries: []Entry{
			{
				EntryTypeID: e,
				Description: description,
				AmountName:  "amount",
			},
			{
				EntryTypeID: 401,
				Description: "IOF",
				AmountName:  "amount",
			},
		},
		Enable:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Version:   1,
	}
}
