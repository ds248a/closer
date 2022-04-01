# closer

Обработка функций прерывания основного процесса импортированных пакетов.

Поддержка различного формата функций завершения.
Поддержка ключей слежения, реализующая возможность прекращения процесса отслеживания.
Поддержка обработки ошибок функций завршения.

### Основной функционал
```go
import (
  "fmt"
  "syscall"
  "time"
  "github.com/ds248a/closer"
)

// Далее представлен список формата функций завершения основной оперции
// из импортированного (внешнего) пакета, подлежащие обработке

type Redis struct {}
type Postgre struct {}
type Cache struct {}

// 1. Функция возвращающает ошибку
func (c *Redis) Close() error {
  time.Sleep(2 * time.Second)
  return fmt.Errorf("%s", fmt.Sprintln("Redis err"))
}

// 2. Функция ничего не возвращает
func (c *Postgre) Close() {
  time.Sleep(2 * time.Second)
}

// 3. Функция завершения с произвольным наименованием
func (c *Cache) Clear() {
  time.Sleep(2 * time.Second)
}

func main() {
  cc := &Cache{}
  pg := &Postgre{}
  rd := &Redis{}

  // Ограничение по времени исполнения всем обработчикам
  closer.SetTimeout(10 * time.Second)

  // Регистрация обработчика, с сохранением ключа
  pKey := closer.Add(pg.Close)

  // Простая регистрация обработчика, без сохранения ключа
  closer.Add(cc.Clear)

  // Отмена обработки (прерывание отслеживания) по ключу
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
