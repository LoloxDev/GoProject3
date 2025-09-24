package analyzer

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	"loganalyzer/internal/config"
)

const (
	StatusOK     = "OK"
	StatusFailed = "FAILED"
)

var ErrLogFileAccess = errors.New("log file not accessible")

var ErrLogParsing = errors.New("log parsing failed")

type LogFileError struct {
	LogID string
	Path  string
	Err   error
}

func (e *LogFileError) Error() string {
	return fmt.Sprintf("%s (%s): %v", e.LogID, e.Path, e.Err)
}

func (e *LogFileError) Unwrap() error { return e.Err }

func (e *LogFileError) Is(target error) bool {
	return target == ErrLogFileAccess
}

type ParseError struct {
	LogID string
	Err   error
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("%s: %v", e.LogID, e.Err)
}

func (e *ParseError) Unwrap() error { return e.Err }

func (e *ParseError) Is(target error) bool {
	return target == ErrLogParsing
}

type Result struct {
	LogID        string `json:"log_id"`
	FilePath     string `json:"file_path"`
	Status       string `json:"status"`
	Message      string `json:"message"`
	ErrorDetails string `json:"error_details"`
	Err          error  `json:"-"`
}

func Analyze(entries []config.LogConfig) []Result {
	results := make([]Result, len(entries))
	if len(entries) == 0 {
		return results
	}

	type payload struct {
		index  int
		result Result
	}

	ch := make(chan payload, len(entries))
	var wg sync.WaitGroup

	for idx, entry := range entries {
		wg.Add(1)
		go func(i int, item config.LogConfig) {
			defer wg.Done()

			res := Result{
				LogID:    item.ID,
				FilePath: item.Path,
			}

			randSrc := rand.New(rand.NewSource(time.Now().UnixNano() + int64(i)))

			file, err := os.Open(item.Path)
			if err != nil {
				res.Status = StatusFailed
				res.Message = "Fichier introuvable ou illisible."
				res.Err = &LogFileError{LogID: item.ID, Path: item.Path, Err: err}
				if res.Err != nil {
					res.ErrorDetails = res.Err.Error()
				}
			} else {
				_ = file.Close()
				delay := time.Duration(randSrc.Intn(151)+50) * time.Millisecond
				time.Sleep(delay)

				if randSrc.Intn(100) < 10 {
					res.Status = StatusFailed
					res.Message = "Erreur lors de l'analyse du contenu."
					res.Err = &ParseError{
						LogID: item.ID,
						Err:   fmt.Errorf("analyse du log de type %q impossible", item.Type),
					}
					res.ErrorDetails = res.Err.Error()
				} else {
					res.Status = StatusOK
					res.Message = "Analyse terminée avec succès."
				}
			}

			ch <- payload{index: i, result: res}
		}(idx, entry)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for item := range ch {
		results[item.index] = item.result
	}

	return results
}
