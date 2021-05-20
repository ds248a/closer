# closer

### Основной функционал

```go
import (
	"fmt"
	"syscall"
	"time"

	"github.com/ds248a/closer"
)

type Redis struct {
}
type Postgre struct {
}
type Cache struct {
}

// Финкция возвращает ошибку
func (c *Redis) Close() error {
	time.Sleep(2 * time.Second)
	return fmt.Errorf("%s", fmt.Sprintln("Redis err"))
}

// Функция ничего не возвращает
func (c *Postgre) Close() {
	time.Sleep(2 * time.Second)
}

// Функция с произвольным наименованием
func (c *Cache) Clear() {
	time.Sleep(2 * time.Second)
	fmt.Println("cache close")
}

func main() {
	cc := &Cache{}
	pg := &Postgre{}
	rd := &Redis{}

  // Ограничение по времени исполнения всем обработчикам
  closer.SetTimeout(10 * time.Second)

	// Простая регистрация обработчика, без сохранения ключа
	closer.Add(cc.Clear)

	// Регистрация обработчика, с сохранением ключа
	pKey := closer.Add(pg.Close)

	// Отмена обработки по ключу
	closer.Remove(pKey)

	// Обработка ошибок на стороне приложения, с сохранением ключа
	rKey := closer.Add(func() {
		err := rd.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	})

	// Out: c159f74d-8a6c-49fd-a181-83edc1d5d595
	fmt.Println(rKey)

	// Сброс всех заданий на обработку
	closer.Reset()

	// Вызываться в самом конце функции 'main()'
	closer.ListenSignal(syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
}

```

### Объектная и кроспакетная обработка

Пример выше - образец кроспакетной обработки. Следующий вариант - объектный.

```go
import (
	"github.com/ds248a/closer"
)

func main() {
  c := closer.NewCloser()
  c.Add(cc.Clear)
  c.Add(pg.Close)
  c.ListenSignal(syscall.SIGTERM)
}
```
