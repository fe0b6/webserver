// webserver
package webserver

import (
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/fe0b6/webobj"
)

var (
	wg        sync.WaitGroup
	exited    bool
	routeFunc func(*webobj.RqObj)
)

func Run(o InitObj) chan bool {
	routeFunc = o.Route

	exitChan := make(chan bool)

	go waitExit(exitChan)

	// Начинаем слушать порт
	go listen(o.Port)

	return exitChan
}

// Ждем сигнал о выходе
func waitExit(exitChan chan bool) {

	_ = <-exitChan

	exited = true

	log.Println("[info]", "Завершаем работу web сервера")

	// Ждем пока все запросы завершатся
	wg.Wait()

	log.Println("[info]", "Работа web сервера завершена корректно")
	exitChan <- true

}

// Начитаем слушать порт
func listen(port int) {

	http.HandleFunc("/", parseRequest)

	log.Fatalln("[fatal]", http.ListenAndServe(":"+strconv.Itoa(port), nil))
}

// Разбираем запрос
func parseRequest(w http.ResponseWriter, r *http.Request) {

	// Если сервер завершает работу
	if exited {
		w.WriteHeader(503)
		w.Write([]byte(http.StatusText(503)))
		return
	}

	// Отмечаем что начался новый запрос
	wg.Add(1)
	// По завершению запроса отмечаем что он закончился
	defer wg.Done()

	ro := &webobj.RqObj{R: r, W: w, TimeStart: time.Now(), FontChan: make(chan string, 1)}
	go ro.GetFonts(r)
	routeFunc(ro)
}
