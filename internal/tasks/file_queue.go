package tasks

import (
	"encoding/json"
	"os"
	"sync"
)

// FileQueue provides a simple, file-based queue for logging failed tasks.
type FileQueue struct {
	filePath string
	mu       sync.Mutex
}

// NewFileQueue creates a new FileQueue.
func NewFileQueue(filePath string) *FileQueue {
	return &FileQueue{
		filePath: filePath,
	}
}

// LogFailedTask serializes the given data to JSON and appends it to the log file.
func (q *FileQueue) LogFailedTask(taskData interface{}) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	// Marshal the task data to JSON
	data, err := json.Marshal(taskData)
	if err != nil {
		return err
	}

	// Open the file in append mode, create it if it doesn't exist
	file, err := os.OpenFile(q.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the JSON data followed by a newline
	if _, err := file.Write(append(data, '\n')); err != nil {
		return err
	}

	return nil
}
