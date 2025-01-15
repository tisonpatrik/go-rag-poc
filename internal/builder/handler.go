package builder

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"rag-poc/internal/repository"
	"sync"
	"time"
)

// Task represents a unit of work to be processed
type Task struct {
	EmbeddingUUID string `json:"embedding_uuid"`
	ChunkSeq      int    `json:"chunk_seq"`
	Chunk         string `json:"chunk"`
	Embedding     []float32
}

// Handler encapsulates the dependencies for handling builder routes
type Handler struct {
	Queries   *repository.Queries
	TaskQueue chan Task
}

// NewHandler creates a new instance of the Handler
func NewHandler(queries *repository.Queries, queueSize int) *Handler {
	return &Handler{
		Queries:   queries,
		TaskQueue: make(chan Task, queueSize), // Asynchronous queue
	}
}

// ListenHandler starts the LISTEN mechanism
func (h *Handler) ListenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	go h.startListening(r.Context())

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Listen started successfully"))
}

// startListening listens for new notifications and pushes them to the queue
func (h *Handler) startListening(ctx context.Context) {
	log.Println("Listening for new_record notifications...")

	// Call the existing ListenForNewSpeechRecords method
	if err := h.Queries.ListenForNewSpeechRecords(ctx); err != nil {
		log.Printf("Error starting LISTEN: %v\n", err)
		return
	}

	// Loop to handle notifications (simulated here, adapt for actual notification handling)
	for {
		// Example: Simulate receiving a notification (replace with actual pgx.WaitForNotification handling if needed)
		notificationPayload := `{"embedding_uuid": "uuid", "chunk_seq": 1, "chunk": "example"}`
		var task Task
		if err := json.Unmarshal([]byte(notificationPayload), &task); err != nil {
			log.Printf("Error parsing notification payload: %v\n", err)
			continue
		}

		// Push task to the queue
		select {
		case h.TaskQueue <- task:
			log.Printf("Task enqueued: %v\n", task)
		default:
			log.Printf("Task queue is full, dropping task: %v\n", task)
		}
	}
}

// ProcessQueue processes tasks from the queue
func (h *Handler) ProcessQueue(workerCount int) {
	var wg sync.WaitGroup

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			log.Printf("Worker %d started\n", workerID)

			for task := range h.TaskQueue {
				log.Printf("Worker %d processing task: %v\n", workerID, task)
				time.Sleep(2 * time.Second) // Simulate slow processing
				log.Printf("Worker %d completed task: %v\n", workerID, task)
			}
		}(i)
	}

	wg.Wait()
}
