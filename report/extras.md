Следует реализовать систему обработки задач, где воркеры постоянно переходят к следующей задаче, если текущая занята, без ожидания ("перепрыгивают"). Это может привести к Live Lock.

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Task struct {
	ID        int
	mu        sync.Mutex
	processed bool
}

func (t *Task) Process(workerID int) bool {
	locked := t.mu.TryLock()
	if !locked {
		return false
	}
	defer t.mu.Unlock()

	if t.processed {
		return false
	}

	fmt.Printf("Worker %d processing task %d\n", workerID, t.ID)
	time.Sleep(time.Duration(rand.Intn(100)+50) * time.Millisecond)
	t.processed = true

	return true
}

func worker(workerID int, tasks <-chan *Task, wg *sync.WaitGroup) {
	defer wg.Done()
	
	for task := range tasks {
		if task.Process(workerID) {
			continue
		}
		go func(t *Task) { tasks <- t }(task)
	}
}

func main() {
	const numWorkers = 5
	const numTasks = 20

	rand.Seed(time.Now().UnixNano())

	taskList := make([]*Task, numTasks)
	for i := 0; i < numTasks; i++ {
		taskList[i] = &Task{ID: i + 1}
	}

	tasks := make(chan *Task, numTasks)
	for _, task := range taskList {
		tasks <- task
	}

	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for i := 0; i < numWorkers; i++ {
		go worker(i+1, tasks, &wg)
	}

	for {
		allProcessed := true
		for _, task := range taskList {
			task.mu.Lock()
			if !task.processed {
				allProcessed = false
				task.mu.Unlock()
				break
			}
			task.mu.Unlock()
		}

		if allProcessed {
			close(tasks)
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	wg.Wait()
	fmt.Println("All tasks processed")
}

Ключевые моменты решения:
1.	Использование TryLock: Воркеры пытаются захватить задачу без блокировки. Если не получается, они переходят к следующей.
2.	Отметка о выполнении: Каждая задача имеет флаг processed, чтобы избежать повторной обработки.
3.	Возврат задач в очередь: Если задача не была обработана (занята или уже выполнена), она возвращается в канал для повторной обработки.
4.	Контроль завершения: Главная горутина проверяет, все ли задачи выполнены, перед закрытием канала.
5.	Буферизированный канал: Позволяет избежать блокировки при возврате задач в очередь.
Это решение предотвращает live lock, гарантируя, что:
•	Каждая задача будет обработана ровно один раз
•	Воркеры не будут бесконечно "перепрыгивать" между задачами
•	Система корректно завершится после обработки всех задач


________________________________________
 Задача 2:
У нас есть очередь задач и несколько воркеров, которые берут задачи из очереди. Если очередь быстро обновляется, возможен случай, когда воркеры никогда не находят задачи, хотя они есть.

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	const numWorkers = 5
	const queueSize = 100
	const totalTasks = 1000

	var wg sync.WaitGroup
	taskQueue := make(chan int, queueSize)
	rand.Seed(time.Now().UnixNano())

	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func(id int) {
			defer wg.Done()
			for {
				select {
				case task, ok := <-taskQueue:
					if !ok {
						return
					}
					fmt.Printf("Worker %d got task %d\n", id, task)
					time.Sleep(time.Duration(rand.Intn(10))
				default:
					time.Sleep(time.Duration(rand.Intn(10)))
				}
			}
		}(i)
	}

	go func() {
		for i := 0; i < totalTasks; i++ {
			taskQueue <- i
			time.Sleep(time.Duration(rand.Intn(5)))
		}
		close(taskQueue)
	}()

	wg.Wait()
	fmt.Println("All workers done")
}

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	const numWorkers = 5
	const queueSize = 100
	const totalTasks = 1000

	var wg sync.WaitGroup
	taskQueue := make(chan int, queueSize)
	rand.Seed(time.Now().UnixNano())

	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func(id int) {
			defer wg.Done()
			backoff := time.Duration(0)
			for {
				select {
				case task, ok := <-taskQueue:
					if !ok {
						return
					}
					backoff = 0
					fmt.Printf("Worker %d got task %d\n", id, task)
					time.Sleep(time.Duration(rand.Intn(10)))
				default:
					if backoff < time.Millisecond*100 {
						backoff += time.Millisecond * 10
					}
					time.Sleep(backoff)
				}
			}
		}(i)
	}

	go func() {
		for i := 0; i < totalTasks; i++ {
			taskQueue <- i
			time.Sleep(time.Duration(rand.Intn(5)))
		}
		close(taskQueue)
	}()

	wg.Wait()
	fmt.Println("All workers done")
}

Ключевые моменты решения проблемы live lock:

1.	Экспоненциальная задержка (backoff)
o	Воркеры увеличивают время ожидания при пустой очереди
o	Начинают с малой задержки (10ms), постепенно увеличивают до максимума (100ms)
o	Сбрасывают задержку при успешном получении задачи
2.	Приоритет обработки над проверкой очереди
o	Используется select с приоритетным чтением из канала (case task)
o	Только при отсутствии задач переходит к default с задержкой
3.	Контроль скорости потребления
o	Воркеры искусственно замедляют обработку (rand.Intn(10))
o	Позволяет производителю добавлять новые задачи
4.	Гибкое управление потоком
o	Разная скорость производства (5ms) и потребления (10ms) задач
o	Буферизированный канал (size=100) сглаживает пики нагрузки
5.	Гарантия завершения
o	Явное закрытие канала после всех задач
o	sync.WaitGroup для корректного завершения воркеров
Стратегия эффективна потому что:
•	Уменьшает конкуренцию воркеров при малом количестве задач
•	Сохраняет быструю реакцию при появлении новых задач
•	Автоматически адаптируется к нагрузке без ручной настройки

