package main

import (
    "fmt"
    "sync"
    "time"
)

type Chopstick struct {
    mutex sync.Mutex
}

func (c *Chopstick) PickUp(wg *sync.WaitGroup) {
    defer wg.Done()
    for {
        if c.mutex.TryLock() {
            fmt.Println("Chopstick picked up")
            time.Sleep(500 * time.Millisecond)
            c.mutex.Unlock()
            break
        } else {
            fmt.Println("Chopstick busy, retrying...")
            time.Sleep(100 * time.Millisecond)
        }
    }
}

func main() {
    var wg sync.WaitGroup
    cs := &Chopstick{}

    wg.Add(2)
    go cs.PickUp(&wg)
    go cs.PickUp(&wg)

    wg.Wait()
}

