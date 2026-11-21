package processor

import (
	"context"
	"sync"
)

func Wrap(ctx context.Context, wg *sync.WaitGroup, cb func(context.Context)) {
	// TODO: Wrap запускает callback в отдельной горутине с учётом WaitGroup и контекста.
	//
	// Ожидаемое поведение:
	// - Если wg передан (не nil), горутина должна быть учтена в WaitGroup
	//   (увеличить счётчик перед запуском, уменьшить при завершении).
	// - Перед выполнением cb проверяется контекст: если он уже отменён,
	//   callback не выполняется — горутина сразу завершается.
	// - Если контекст активен — вызывается cb(ctx).
	if wg != nil {
		wg.Add(1)
	}
	go func() {
		if wg != nil {
			defer wg.Done()
		}
		select {
		case <-ctx.Done():
			return
		default:
			cb(ctx)
		}
	}()
}
