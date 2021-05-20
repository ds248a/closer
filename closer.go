package closer

import (
	"context"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"sync"
	"time"

	"github.com/google/uuid"
)

// ----------------------
//   Closer
// ----------------------

var closer *Closer

type Action func(context.Context, *sync.WaitGroup, os.Signal)

type Closer struct {
	mu      *sync.RWMutex
	timeout time.Duration
	log     Logger
	actions map[string]Action
}

func init() {
	closer = NewCloser()
}

func NewCloser() *Closer {
	return &Closer{
		mu:      &sync.RWMutex{},
		timeout: 20 * time.Second,
		log:     DefaultLogger(),
		actions: make(map[string]Action),
	}
}

// ----------------------

// Переопределяет время, в течение которого соединения могут быть закрыты
func SetTimeout(t time.Duration) {
	closer.SetTimeout(t)
}

// Переопределяет библиотеку логирования
func SetLogger(l Logger) {
	closer.log = l
}

// Обработка системных прерываний.
// Cписок значений: os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT
func ListenSignal(signals ...os.Signal) {
	closer.ListenSignal(signals...)
}

// Закрытие открытых соединений с ограничением по времени исполнения
func Close(s os.Signal) {
	closer.Close(s)
}

// Регистрация функции закрытия соединения или уничтожения объекта
func Add(f func()) string {
	return closer.Add(f)
}

// Возвращает полный список обрабтчиков
func Actions() map[string]Action {
	return closer.Actions()
}

// Удаляет указанный обработчик из планировщика
func Remove(key string) Action {
	return closer.Remove(key)
}

// Сброс обработчиков
func Reset() {
	closer.Reset()
}

// ----------------------

// Переопределяет время, в течение которого соединения могут быть закрыты
func (c *Closer) SetTimeout(t time.Duration) {
	c.timeout = t
}

// Переопределяет библиотеку логирования
func (c *Closer) SetLogger(l Logger) {
	c.log = l
}

// Обработка системных прерываний
// список значений: os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT
func (c *Closer) ListenSignal(signals ...os.Signal) {
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, signals...)
	s := <-sigChannel
	c.Close(s)
}

// Закрытие открытых соединений с ограничением по времени исполнения
func (c *Closer) Close(s os.Signal) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(len(c.actions))

	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, action := range c.actions {
		go func(a Action) {
			a(ctx, &wg, s)
		}(action)
	}

	wg.Wait()
}

// Регистрация функции закрытия соединения или уничтожения объекта
func (c *Closer) Add(f func()) string {
	key := c.callOnExit(func(ctx context.Context, wg *sync.WaitGroup, s os.Signal) {
		defer wg.Done()

		dst := make(chan bool)
		go func() {
			f()
			dst <- true
		}()

	loop:
		for {
			select {
			case <-ctx.Done():
				c.log.Error("Failed to close: ", funcName(f), " at signal:", s.String())
				break loop
			case <-dst:
				break loop
			}
		}
	})

	return key
}

func (c *Closer) callOnExit(action Action) string {
	key := uuid.New().String()
	c.mu.Lock()
	c.actions[key] = action
	c.mu.Unlock()
	return key
}

// Возвращает полный список обрабтчиков
func (c *Closer) Actions() map[string]Action {
	c.mu.RLock()
	defer c.mu.RUnlock()
	actions := make(map[string]Action)
	for k, a := range c.actions {
		actions[k] = a
	}
	return actions
}

// Удаляет указанный обработчик из планировщика
func (c *Closer) Remove(key string) Action {
	c.mu.Lock()
	defer c.mu.Unlock()
	if action, ok := c.actions[key]; ok {
		delete(c.actions, key)
		return action
	}
	return nil
}

// Сброс обработчиков
func (c *Closer) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for k, _ := range c.actions {
		delete(c.actions, k)
	}
}

// Возвращает наименования структуры и исполняемой функции
func funcName(i interface{}) (name string) {
	t := reflect.TypeOf(i)
	if t.Kind() != reflect.Func {
		return "not defined"
	}
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
