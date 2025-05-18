package api

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/sync/errgroup"
)

type oauth2params struct {
	code  string
	state string
	err   string
}

func authorize(config *oauth2.Config) (string, error) {
	l, err := net.Listen("tcp", ":0")

	if err != nil {
		return "", err
	}

	mux := http.NewServeMux()
	server := &http.Server{Handler: mux}
	ch := make(chan *oauth2params)

	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		params := &oauth2params{
			code:  query.Get("code"),
			state: query.Get("state"),
			err:   query.Get("error"),
		}

		if params.err != "" {
			fmt.Fprint(w, params.err)
		} else {
			fmt.Fprint(w, "OK")
		}

		ch <- params
	})

	eg, _ := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		err := server.Serve(l)
		close(ch)

		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}

		return nil
	})

	port := l.Addr().(*net.TCPAddr).Port
	config.RedirectURL = fmt.Sprintf("http://localhost:%d/callback", port)
	state := uuid.NewString()
	authURL := config.AuthCodeURL(state, oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser: \n%s\n", authURL)

	params := <-ch
	server.Shutdown(context.Background())
	err = eg.Wait()

	if err != nil {
		return "", err
	}

	if params.state != state {
		panic("invalid oauth2 state")
	}

	if params.err != "" {
		return "", errors.New(params.err)
	}

	return params.code, nil
}
