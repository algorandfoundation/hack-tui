package fortiter

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"github.com/algorand/go-algorand/logging/logspec"
	"github.com/jmoiron/sqlx"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

func LogFile(filename *string) (io.ReadCloser, error) {
	var inputStream io.ReadCloser = os.Stdin
	if *filename == "" {
		return inputStream, errors.New("no input file specified")
	}
	if *filename != "" {
		f, err := os.Open(*filename)
		if err != nil {
			return nil, err
		}
		// Close the handle - we just wanted to verify it was valid
		f.Close()
		cmd := exec.Command("tail", "-n", "-1000", "-F", *filename)
		inputStream, err = cmd.StdoutPipe()
		if err != nil {
			return nil, err
		}
		err = cmd.Start()
		if err != nil {
			return nil, err
		}
	}
	return inputStream, nil
}

func Sync(ctx context.Context, filepath string, db *sqlx.DB) error {
	stats := make(map[string]Stats)

	go Watch(filepath, func(line string, lm map[string]interface{}, err error) {
		if lm["Context"] != nil && lm["Context"] == "Agreement" {
			if lm["Hash"] != nil {
				var event logspec.AgreementEvent
				dec := json.NewDecoder(strings.NewReader(line))
				_ = dec.Decode(&event)

				t, _ := time.Parse(time.RFC3339, lm["time"].(string))
				err = SaveAgreement(AgreementEvent{
					AgreementEvent: event,
					Message:        lm["msg"].(string),
					Time:           t,
				}, db)

				if event.Sender != "" {
					var stat Stats
					stat, ok := stats[event.Sender]
					if !ok {
						stats[event.Sender] = Stats{
							Address:  event.Sender,
							Sent:     0,
							Received: 0,
							Failed:   0,
							Success:  0,
						}
						stat = stats[event.Sender]
					} else {
						stat.Received++
					}
					stat.SaveStats(*db)
				}
			}
		}
	})

	return nil
}

func Watch(filepath string, cb func(line string, lm map[string]interface{}, err error)) {
	inputStream, err := LogFile(&filepath)
	if err != nil {
		cb("", make(map[string]interface{}), err)
	}
	scanner := bufio.NewScanner(inputStream)
	for scanner.Scan() {
		line := scanner.Text()
		var event map[string]interface{}
		dec := json.NewDecoder(strings.NewReader(line))
		err := dec.Decode(&event)
		if err != nil {
			cb("", event, err)
			break
		} else {
			cb(line, event, nil)
		}
	}
}
