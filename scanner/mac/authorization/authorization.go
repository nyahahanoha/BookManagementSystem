package authorization

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
)

func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", url}
	case "darwin": // macOS
		cmd = "open"
		args = []string{url}
	default: // Linux, BSDなど
		cmd = "xdg-open"
		args = []string{url}
	}

	return exec.Command(cmd, args...).Start()
}

type PomeriumSession struct {
	JWT string `json:"jwt"`
}

type SessionMessage []byte

func Authorization(api url.URL, callbackPort uint16) (string, error) {
	if api == (url.URL{}) || callbackPort == 0 {
		return "", fmt.Errorf("api or callback port is empty")
	}

	redirectURI := fmt.Sprintf("http://localhost:%d/callback", callbackPort)
	loginURI := url.URL{
		Scheme: api.Scheme,
		Host:   api.Host,
		Path:   "/.pomerium/api/v1/login",
		RawQuery: url.Values{
			"pomerium_redirect_uri": {redirectURI},
		}.Encode(),
	}

	resp, err := http.Get(loginURI.String())
	if err != nil {
		return "", fmt.Errorf("failed to request login URL: %w", err)
	}
	defer resp.Body.Close()
	bodyBytes, _ := io.ReadAll(resp.Body)
	signInURL := string(bodyBytes)
	if err := openBrowser(signInURL); err != nil {
		return "", fmt.Errorf("failed to open browser: %w", err)
	}

	sessionCh := make(chan SessionMessage)
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		jwt := r.URL.Query().Get("pomerium_jwt")
		if jwt == "" {
			http.Error(w, "missing token", http.StatusBadRequest)
			return
		}

		sessionCh <- []byte(jwt)
		sessionJSON, err := json.Marshal([]byte(jwt))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(sessionJSON)
	})

	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%d", callbackPort), nil); err != nil {
			log.Fatal("callback server error: %w", err)
		}
	}()

	token := string(<-sessionCh)

	return token, nil
}
