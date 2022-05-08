# closer

Реализация корректного закрытия процессов импортированных пакетов по завершению работы приложения.

Реализована поддержка:
- обработки системных событий. По умолчанию обрабатываются: os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT
- конкурентный запуск процессов завершения с ограничением времени исполнения;
- возможность пользовательской обработки картежей функций заврешения процесса;
- доступ к обработчикам процесса завершения по ключу.

### Пример использования

```go
import (
  "fmt"
  "syscall"
  "time"
  "github.com/ds248a/closer"
)

// Для примеры представлены структуры с разными форматами реализации функций завершения их процессов.
type Redis struct {}
type Postgre struct {}
type Cache struct {}

// 1. Функция возвращает результат своей работы.
func (c *Redis) Close() error {
  time.Sleep(time.Second)
  return fmt.Errorf("%s", fmt.Sprintln("Redis err"))
}

// 2. Функция ничего не возвращает.
func (c *Postgre) Close() {
  time.Sleep(time.Second)
}

// 3. Имеет произвольное наименование.
func (c *Cache) Clear() {
  time.Sleep(time.Second)
}

func main() {
  cc := &Cache{}
  pg := &Postgre{}
  rd := &Redis{}

  // Ограничение времени исполнения всех зависимых процессов.
  closer.SetTimeout(10 * time.Second)

  // Регистрация функций завершения процессов в планировщике.
  // Возвращает ключ, с поиощью которого возможно отменить обработчик.
  pKey := closer.Add(pg.Close)

  // Простая регистрация обработчика.
  closer.Add(cc.Clear)

  // Отмена обработчика.
  closer.Remove(pKey)

  // Пользовательская обработка ошибки исполнения.
  rKey := closer.Add(func() {
    err := rd.Close()
    if err != nil {
      fmt.Println(err.Error())
    }
  })

  // Вывод: 1
  fmt.Println(rKey)

  // Удаление всех задач с планировщика.
  closer.Reset()

  // Запуск обработки системых вызовов.
  closer.ListenSignal(syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
}

```

### Объектная и крос-пакетная обрабка

```go
func main() {
  c := closer.NewCloser()
  c.Add(cc.Clear)
  c.Add(pg.Close)
  c.ListenSignal(syscall.SIGTERM)
}
```
