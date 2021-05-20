package closer_test

import (
	"fmt"
	"os"
	"syscall"
	"testing"
	"time"

	c "github.com/ds248a/closer"
)

// ----------------------
//   Struct example
// ----------------------

type Postgre struct {
}

func (t *Postgre) Close() {
	time.Sleep(2 * time.Second)
}

type Redis struct {
}

func (t *Redis) Close() error {
	time.Sleep(2 * time.Second)
	return nil
}

type Cache struct {
}

func (t *Cache) Clear() {
	time.Sleep(2 * time.Second)
}

// ----------------------
//   Test
// ----------------------

func Test(t *testing.T) {
	defer c.Reset()
	c.SetTimeout(3 * time.Second)

	//mu := sync.Mutex{}
	testData := make(map[int]int)

	pg := &Postgre{}
	rd := &Redis{}
	cc := &Cache{}

	// Регистрация обработчиков

	pgKey := c.Add(pg.Close)

	c.Add(cc.Clear)

	// Обработка ошибки на стороне прилоения
	c.Add(func() {
		if err := rd.Close(); err != nil {
			fmt.Println(err.Error())
		}
		testData[1]++
	})

	/*
		k4 := sig.Add(func(s os.Signal) {
			mu.Lock()
			defer mu.Unlock()
			if s == syscall.SIGTERM {
				testData[4]++
			}
		})
	*/

	// Удаление обработчика по ключу
	pgAction := c.Remove(pgKey)
	if pgAction == nil {
		t.Errorf("Action key not exist.")
	}

	// Удвление не существующего обработчика
	action := c.Remove("empty")
	if action != nil {
		t.Errorf("Action should not be existing but it does.")
	}

	// Проверка текущего списка обработчиков
	actions := c.Actions()
	if len(actions) != 2 {
		t.Errorf("Invalid number of actions expected %d found %d.", 2, len(actions))
	}

	if len(testData) != 0 {
		t.Errorf("Test data is not empty %+v.", testData)
	}

	// Регистрация дополнительного обработчика
	c.Add(func() {
		pg.Close()
		testData[2]++
	})

	go func() {
		time.Sleep(100 * time.Millisecond)
		p, err := os.FindProcess(os.Getpid())
		if err != nil {
			panic(err.Error())
		}
		p.Signal(syscall.SIGTERM)
	}()

	c.ListenSignal()

	if len(testData) != 2 {
		t.Errorf("Test data invalid expected length %d found %d.", 2, len(testData))
	}

	if testData[1] != 1 || testData[2] != 1 {
		t.Errorf("Invalid test data %+v.", testData)
	}

	/*
		a4(syscall.SIGTERM)
		if testData[4] != 1 {
			t.Errorf("Invalid removed action.")
		}
	*/
}
