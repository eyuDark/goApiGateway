package log

import (
	"bytes"
	"fmt"
	"os"
	"sync"
	"testing"
)


func TestConcurrentWrites(t *testing.T) {
    logger, _ := NewRotatingLogger("test.log", 10)
    defer os.Remove("test.log")
    defer logger.Close()

    var wg sync.WaitGroup
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            logger.Write([]byte(fmt.Sprintf("Request %d\n", id)))
        }(i)
    }
    wg.Wait()

    // Verify exactly 100 lines
    data, _ := os.ReadFile("test.log")
    if len(bytes.Split(data, []byte("\n"))) != 101 { // +1 for empty last line
        t.Fatal("Missing log entries in concurrent test")
    }
}