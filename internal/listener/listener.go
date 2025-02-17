package listener

import (
	"encoding/base64"
	"encoding/json"
	"example_listener/internal/config"
	"example_listener/internal/server"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	listener "github.com/PicoTools/pico/pkg/proto/listener/v1"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var srv *server.Server

func Start(s *server.Server) error {
	srv = s

	err := srv.RegisterListener()
	if err != nil {
		return err
	}

	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/task", taskHandler)
	http.HandleFunc("/output", outputHandler)

	log.Printf("Starting listener on %s", config.ListenerAddr)
	err = http.ListenAndServe(config.ListenerAddr, nil)
	if err != nil {
		return err
	}

	return nil
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	encodedMsg := r.URL.Query().Get("data")
	if encodedMsg == "" {
		http.Error(w, "", http.StatusForbidden)
		return
	}

	decodedMsg, err := messageDecode(encodedMsg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("New registration request: %s", decodedMsg)

	metadata := &NewBeacon{}
	if err := json.Unmarshal(decodedMsg, metadata); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = srv.NewBeacon(&listener.RegisterAgentRequest{
		Id:       metadata.Id,
		Os:       metadata.Os,
		Arch:     metadata.Arch,
		Sleep:    metadata.Sleep,
		Jitter:   metadata.Jitter,
		Caps:     metadata.Caps,
		Hostname: wrapperspb.String(metadata.Hostname),
		Username: wrapperspb.String(metadata.Username),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	encoded := r.URL.Query().Get("data")
	if encoded == "" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	beaconId, err := strconv.Atoi(string(decoded))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task, err := srv.GetTask(uint32(beaconId))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if task == nil {
		return
	}

	msg, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("New task for %s: %s\n", string(decoded), msg)

	encodedMsg := messageEncode(msg)
	fmt.Fprint(w, encodedMsg)
	return
}

func outputHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	decoded, err := base64.StdEncoding.DecodeString(string(body) + "===")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "/output: %s", string(decoded))
}
