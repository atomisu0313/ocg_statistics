package transaction

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"atomisu.com/ocg-statics/cardApiCall/internal/sqlc_gen"
)

type Tx struct {
	db *sql.DB
}

func NewTx(db *sql.DB) *Tx {
	return &Tx{db: db}
}

// ─────────────────────────────────────────────
// ① コア: *sql.Tx を受け取り、任意の型 T を返す（sqlc に依存しない）
// ─────────────────────────────────────────────

// Run はトランザクションを開始し、fn を実行する。
// fn がエラーを返した場合は Rollback、成功した場合は Commit する。
// Commit / Rollback のエラーも呼び出し元へ返す（名前付き戻り値 err を defer 内で更新する）。
// 失敗時の result は常に T のゼロ値になる。
// opts が nil の場合は DB のデフォルト（PostgreSQL では read committed）で開始する。
func Run[T any](
	ctx context.Context,
	db *sql.DB,
	opts *sql.TxOptions,
	fn func(*sql.Tx) (T, error),
) (result T, err error) {
	txx, err := db.BeginTx(ctx, opts)
	if err != nil {
		return result, fmt.Errorf("begin tx: %w", err)
	}

	defer func() {
		// fn 内で panic した場合もトランザクションを閉じてから再 panic する
		if p := recover(); p != nil {
			_ = txx.Rollback()
			panic(p)
		}

		if err != nil {
			// fn のエラーを保持したまま、Rollback のエラーがあれば合成する
			if rbErr := txx.Rollback(); rbErr != nil {
				err = errors.Join(err, fmt.Errorf("rollback tx: %w", rbErr))
			}
			var zero T
			result = zero // 失敗時はゼロ値を返す（中途半端な ID を返さない）
			return
		}

		// Commit 失敗時点で Tx は終了しているため Rollback は呼ばない
		if cmErr := txx.Commit(); cmErr != nil {
			err = fmt.Errorf("commit tx: %w", cmErr)
			var zero T
			result = zero
		}
	}()

	// ここで外から渡した関数を実行する。
	return fn(txx)
}

// ─────────────────────────────────────────────
// ② アダプタ: *sql.Tx から任意の Queries 型 Q を作る（sqlc パッケージを問わない）
// ─────────────────────────────────────────────

// RunWith は bind で *sql.Tx から Q を生成し、Q を受け取る fn をトランザクション内で実行する。
func RunWith[Q any, T any](
	ctx context.Context,
	db *sql.DB,
	opts *sql.TxOptions,
	bind func(*sql.Tx) Q,
	fn func(Q) (T, error),
) (T, error) {
	return Run(ctx, db, opts, func(t *sql.Tx) (T, error) {
		return fn(bind(t))
	})
}

// ─────────────────────────────────────────────
// ③ このパッケージ向けのショートカット
// ─────────────────────────────────────────────

// bindQueries は、トランザクション（t）に紐づいた sqlc 操作オブジェクトを生成します。
// これにより、コールバック関数内で q.InsertCard(...) などのメソッドを同一トランザクション内で呼び出せるようになります。
func bindQueries(t *sql.Tx) *sqlc_gen.Queries { return sqlc_gen.New(t) }

// ExecTx はトランザクション内で fn を実行し、結果とエラーを返します。
// （結果が不要な場合は戻り値を _ で無視できます）
func ExecTx[T any](ctx context.Context, tx *Tx, fn func(*sqlc_gen.Queries) (T, error)) (T, error) {
	return RunWith(ctx, tx.db, nil, bindQueries, fn)
}

// ExecNonTx はトランザクションを開始せずに fn を実行し、結果とエラーを返します。
// （結果が不要な場合は戻り値を _ で無視できます）
func ExecNonTx[T any](ctx context.Context, tx *Tx, fn func(*sqlc_gen.Queries) (T, error)) (T, error) {
	q := sqlc_gen.New(tx.db)
	return fn(q)
}
