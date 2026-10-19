package app

import (
	"atomisu.com/ocg-statics/cardApiCall/internal/config"
	"atomisu.com/ocg-statics/cardApiCall/internal/usecase/master"
	"atomisu.com/ocg-statics/cardApiCall/internal/usecase/neon"
	"atomisu.com/ocg-statics/cardApiCall/internal/usecase/neuron"
	"atomisu.com/ocg-statics/cardApiCall/internal/usecase/tcgapi"
	_ "github.com/lib/pq" // PostgreSQL ドライバをインポート
	"github.com/samber/do"
)

// SetupDIContainer はDIコンテナをセットアップします
// テストなどでコンテナの中身をオーバーライドしたい場合には、
// samba/doのOverride関数を使って上書きしてください
func SetupDIContainer() *do.Injector {
	injector := do.New()
	// DBコネクション
	do.Provide(injector, config.NewDbConnection)

	// ユースケース
	do.Provide(injector, neuron.NewNeuronUseCase)
	do.Provide(injector, tcgapi.NewTcgUseCase)
	do.Provide(injector, neon.NewNeonUseCase)
	do.Provide(injector, master.NewMasterUseCaseForDI)

	// ハンドラー
	//do.Provide(injector, handler.NewAccountHandler)
	//do.Provide(injector, handler.NewAuthHandler)

	return injector
}
