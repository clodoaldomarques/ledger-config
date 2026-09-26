package ledger

import (
	"context"
	"testing"

	"github.com/clodoaldomarques/ledger-config/internal/domain/ledger"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

// =============================================================================
// Helpers de setup
// =============================================================================

// newBenchCreate monta o Service para o caminho feliz do CreateConfig.
// Reaproveita fakeScript do service_test.go.
func newBenchCreate(b *testing.B) *Service {
	ctrl := gomock.NewController(b)

	r := NewMockRepository(ctrl)
	r.EXPECT().
		FindConfigByLevel(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(ledger.Config{}, ErrConfigNotFound{}).
		AnyTimes()
	r.EXPECT().
		SaveConfig(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()

	tp := NewMockTopic(ctrl)
	tp.EXPECT().
		Emit(gomock.Any(), gomock.Any(), gomock.All()).
		Return(nil).
		AnyTimes()

	return New(r, tp)
}

// newBenchUpdate monta o Service para o caminho feliz do UpdateConfig.
func newBenchUpdate(b *testing.B) *Service {
	ctrl := gomock.NewController(b)

	r := NewMockRepository(ctrl)
	r.EXPECT().
		FindConfigByID(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(fakeScript(ledger.ProgramLevel, "201", "PAGAMENTO A VISTA"), nil).
		AnyTimes()
	r.EXPECT().
		UpdateConfig(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()

	tp := NewMockTopic(ctrl)
	tp.EXPECT().
		Emit(gomock.Any(), gomock.Any(), gomock.All()).
		Return(nil).
		AnyTimes()

	return New(r, tp)
}

// newBenchFindByLevel monta o Service para o caminho feliz do FindConfigByLevel.
// Retorna program level direto na primeira chamada (não cai no tenant).
func newBenchFindByLevel(b *testing.B) *Service {
	ctrl := gomock.NewController(b)

	r := NewMockRepository(ctrl)
	r.EXPECT().
		FindConfigByLevel(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, cid, level, eventTypeID, orgID string, programID *int64) (ledger.Config, error) {
			if level == string(ledger.ProgramLevel) {
				return fakeScript(ledger.ProgramLevel, eventTypeID, "PAGAMENTO A VISTA"), nil
			}
			return ledger.Config{}, ErrConfigNotFound{}
		}).
		AnyTimes()

	tp := NewMockTopic(ctrl)
	return New(r, tp)
}

// newBenchFindAll monta o Service para o caminho feliz do FindAllConfigs.
func newBenchFindAll(b *testing.B, quant int) *Service {
	ctrl := gomock.NewController(b)

	r := NewMockRepository(ctrl)
	r.EXPECT().
		FindAllConfigs(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(fakeSliceScripts(ledger.ProgramLevel, quant), nil).
		AnyTimes()

	tp := NewMockTopic(ctrl)
	return New(r, tp)
}

// =============================================================================
// Benchmarks — CreateConfig
// =============================================================================

// caminho feliz: busca (não acha) -> salva -> emite
func BenchmarkCreateConfig_Success(b *testing.B) {
	svc := newBenchCreate(b)

	ctx := context.Background()
	cid := "cid-fixo"
	cfg := fakeScript(ledger.ProgramLevel, "201", "PAGAMENTO A VISTA")

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := svc.CreateConfig(ctx, cid, cfg)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// mesma coisa, mas gerando UUID por iteração (cenário realista)
func BenchmarkCreateConfig_Success_WithUUID(b *testing.B) {
	svc := newBenchCreate(b)

	ctx := context.Background()
	cfg := fakeScript(ledger.ProgramLevel, "201", "PAGAMENTO A VISTA")

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := svc.CreateConfig(ctx, uuid.NewString(), cfg)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// =============================================================================
// Benchmarks — UpdateConfig
// =============================================================================

// caminho feliz: busca por ID -> atualiza -> emite
func BenchmarkUpdateConfig_Success(b *testing.B) {
	svc := newBenchUpdate(b)

	ctx := context.Background()
	cid := "cid-fixo"
	id := "uuid-12345"
	changed := fakeScript(ledger.PlatformLevel, "201", "Changed Description")

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := svc.UpdateConfig(ctx, cid, id, changed)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// =============================================================================
// Benchmarks — FindConfigByLevel
// =============================================================================

// caminho feliz: acha program level de primeira (não vai pro tenant)
func BenchmarkFindConfigByLevel_ProgramHit(b *testing.B) {
	svc := newBenchFindByLevel(b)

	ctx := context.Background()
	cid := "cid-fixo"
	event := "201"
	org := "TN-77add76c-e395-446b-b306-1a0f9cb99a31"
	var programID int64 = 1

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := svc.FindConfigByLevel(ctx, cid, event, org, programID)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// =============================================================================
// Benchmarks — FindAllConfigs
// =============================================================================

// varia o volume de configs retornadas para medir o custo de escala
func BenchmarkFindAllConfigs(b *testing.B) {
	for _, quant := range []int{10, 100, 1000} {
		b.Run(nameForQuant(quant), func(b *testing.B) {
			svc := newBenchFindAll(b, quant)

			ctx := context.Background()
			cid := "cid-fixo"
			org := "TN-77add76c-e395-446b-b306-1a0f9cb99a31"
			var programID int64 = 1

			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_, err := svc.FindAllConfigs(ctx, cid, org, &programID)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func nameForQuant(q int) string {
	switch q {
	case 10:
		return "10_configs"
	case 100:
		return "100_configs"
	case 1000:
		return "1000_configs"
	default:
		return "n_configs"
	}
}

// =============================================================================
// Benchmark paralelo — checa escalabilidade e contenção
// =============================================================================

func BenchmarkCreateConfig_Parallel(b *testing.B) {
	svc := newBenchCreate(b)

	ctx := context.Background()
	cid := "cid-fixo"

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		cfg := fakeScript(ledger.ProgramLevel, "201", "PAGAMENTO A VISTA")
		for pb.Next() {
			_, err := svc.CreateConfig(ctx, cid, cfg)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
