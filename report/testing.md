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
