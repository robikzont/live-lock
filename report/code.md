Решение проблемы в коде
Чтобы решить проблему Live Lock, можно применить следующие подходы:

1.	Добавить случайную задержку перед повторной попыткой
func (c *Chopstick) PickUpFixed(wg *sync.WaitGroup) {
    defer wg.Done()
    for {
        if c.mutex.TryLock() {
            fmt.Println("Chopstick picked up")
            time.Sleep(500 * time.Millisecond)
            c.mutex.Unlock()
            break
        } else {
            delay := time.Duration(rand.Intn(300)+100) * time.Millisecond
            fmt.Printf("Chopstick busy, backing off for %v\n", delay)
            time.Sleep(delay)
        }
    }
}

2. Использовать ограничения на количество попыток
Можно добавить счётчик попыток и принудительно завершать работу, если ресурс недоступен слишком долго.
3. Использовать контекст с таймером
Для долгих операций полезно использовать context.WithTimeout.

6. Пример тестирования с параллелизмом

func TestLivelock(t *testing.T) {
    var wg sync.WaitGroup
    cs := &Chopstick{mutex: sync.Mutex{}}

    runTest := func(id int, t *testing.T) {
        for i := 0; i < 100; i++ {
            if cs.mutex.TryLock() {
                time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
                cs.mutex.Unlock()
                break
            } else {
                time.Sleep(time.Duration(rand.Intn(300)+100) * time.Millisecond)
            }
        }
        t.Logf("Goroutine %d finished\n", id)
    }

    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            runTest(id, t)
        }(i)
    }

    wg.Wait()
}

Такой тест эмулирует случайное поведение и помогает проверить стабильность решения.


