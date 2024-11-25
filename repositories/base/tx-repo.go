package base_repository

import (
	"context"

	"gorm.io/gorm"
)

// BaseTxRepository интерфейс представляет базовый набор методов для работы с сущностями.
//
// Пример/Шаблон применения:
//
//	tx, err := h.repo.Example.BeginTx(ctx.Context())
//	if err != nil {
//		// Если транзакцию не удалось начать, возвращаем ошибку.
//		return apu.MapErrorToHTTPStatus(ctx, err)
//	}
//
//	defer func() {
//		if r := recover(); r != nil {
//			// Если в коде внутри транзакции произошла паника (например, необработанный runtime exception),
//			// транзакция будет отменена (Rollback), чтобы не оставить базу данных в неконсистентном состоянии.
//			h.repo.Example.RollbackTx(tx)
//			panic(r) // Паника перенаправляется выше, чтобы код за пределами транзакции тоже мог обработать её.
//		} else if err != nil {
//			// Если в процессе выполнения транзакции была возвращена ошибка,
//			// транзакция также будет отменена (Rollback).
//			h.repo.Example.RollbackTx(tx)
//		}
//	}()
//
// // Код выполнения операций внутри транзакции
//
//	err = h.repo.Example.DoSomething(tx, someData)
//	if err != nil {
//	    return err
//	}
//
//	// Если всё прошло успешно, фиксируем транзакцию (commit).
//	if err := h.repo.Example.CommitTx(tx); err != nil {
//		return err
//	}
//
//	return nil
//
// Пояснение работы defer и обработки panic:
//
// 1. Если в транзакции произошла ошибка, функция `RollbackTx` отменяет все изменения, выполненные в транзакции, и возвращает базу данных в исходное состояние.
//
// 2. Если в коде внутри транзакции возникает паника (например, необработанная ошибка или критический сбой), она перехватывается с помощью `recover()`.
//   - В этом случае транзакция будет отменена (Rollback), чтобы база данных осталась в целостном состоянии.
//   - Затем паника будет повторно вызвана с помощью `panic(r)` для обработки на более высоком уровне.
//
// 3. Если выполнение транзакции завершилось успешно (без ошибок и паники), вызывается `CommitTx`, фиксирующий все изменения.
type BaseTxRepository interface {
	// Начать транзакцию
	BeginTx(ctx context.Context) (*gorm.DB, error)

	// Завершить транзакцию (commit)
	CommitTx(tx *gorm.DB) error

	// Отменить транзакцию (rollback)
	RollbackTx(tx *gorm.DB) error
}

// BaseTxRepository представляет базовую структуру для репозиториев
type baseTxRepository struct {
	db *gorm.DB
}

// NewBaseTxRepository создает новый экземпляр базового репозитория
func NewBaseTxRepository(db *gorm.DB) BaseTxRepository {
	return &baseTxRepository{
		db: db,
	}
}

// Начать транзакцию
func (repo *baseTxRepository) BeginTx(ctx context.Context) (*gorm.DB, error) {
	tx := repo.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}

// Завершить транзакцию (commit)
func (repo *baseTxRepository) CommitTx(tx *gorm.DB) error {
	return tx.Commit().Error
}

// Отменить транзакцию (rollback)
func (b *baseTxRepository) RollbackTx(tx *gorm.DB) error {
	return tx.Rollback().Error
}
