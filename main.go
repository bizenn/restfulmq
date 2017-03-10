package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"

	"github.com/Sirupsen/logrus"
)

// Entry ...
type Entry struct {
	mediaType string
	data      []byte
}

// Queue ...
type Queue interface {
	Enqueue(*Entry)
	Dequeue() *Entry
}

type channelQueue struct {
	Queue
	queue chan *Entry
}

// NewQueue ...
func NewQueue(size int) Queue {
	return &channelQueue{
		queue: make(chan *Entry, size),
	}
}

// Enqueue ...
func (q *channelQueue) Enqueue(e *Entry) {
	q.queue <- e
}

// Dequeue ...
func (q *channelQueue) Dequeue() *Entry {
	return <-q.queue
}

func makeQueueHandler(path string, queue Queue) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			e := queue.Dequeue()
			if e != nil {
				w.Header().Set("Content-Type", e.mediaType)
				w.Write(e.data)
				logrus.Infof("Dequeue %s to %s", path, r.RemoteAddr)
			}
		case http.MethodPost:
			if b, err := ioutil.ReadAll(r.Body); err != nil {
				http.Error(w, "Bad Request", http.StatusBadRequest)
			} else {
				queue.Enqueue(&Entry{r.Header.Get("Content-Type"), b})
				logrus.Infof("Enqueue %s from %s", path, r.RemoteAddr)
				w.WriteHeader(http.StatusAccepted)
			}
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}
}

// QueueConfig ...
type QueueConfig struct {
	Path     string `json:"path"`
	Capacity int    `json:"capacity"`
}

// ServerConfig ...
type ServerConfig struct {
	Host    string        `json:"host"`
	Port    int           `json:"port"`
	LogPath string        `json:"logpath"`
	Queues  []QueueConfig `json:"queues"`
}

var config = ServerConfig{
	Host: "",
	Port: 8888,
	Queues: []QueueConfig{
		QueueConfig{"/", 0},
	},
}

func init() {
	if len(os.Args) > 1 {
		in, err := os.Open(os.Args[1])
		if err != nil {
			panic(err)
		}
		defer in.Close()
		decoder := json.NewDecoder(in)
		err = decoder.Decode(&config)
		if err != nil {
			panic(err)
		}
	}
	if config.LogPath == "" || config.LogPath == "-" {
		logrus.SetOutput(os.Stdout)
	} else {
		if out, err := os.OpenFile(config.LogPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777); err != nil {
			panic(err)
		} else {
			logrus.SetOutput(out)
		}
	}
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	host := fmt.Sprintf("%s:%d", config.Host, config.Port)
	logrus.WithField("listen", host).Info("Start")
	for _, qc := range config.Queues {
		http.HandleFunc(qc.Path, makeQueueHandler(qc.Path, NewQueue(qc.Capacity)))
		logrus.WithFields(logrus.Fields{
			"path":     qc.Path,
			"capacity": qc.Capacity,
		}).Infof("Queue created")
	}
	http.ListenAndServe(host, nil)
	logrus.Info("Finish")
}
