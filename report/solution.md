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


