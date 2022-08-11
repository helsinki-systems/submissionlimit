package submissionlimit

import (
	"context"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

// TODO: Maybe add support for TorDNSEL: https://2019.www.torproject.org/projects/tordnsel.html.en

const (
	torExitListURL = "https://check.torproject.org/torbulkexitlist"
)

var (
	ErrIPIsTorExitNode = errors.New("IP is from a Tor exit node")

	TorDefaultRefreshInterval = 30 * time.Minute
)

type TorConfig struct {
	ListRefreshInterval         time.Duration
	NoWaitForInitialListRefresh bool
}

type torIP string

func (ti torIP) String() string {
	return string(ti)
}

type torLimiter struct {
	list   []torIP
	listMu sync.RWMutex

	initialDone chan struct{}
}

func (tl *torLimiter) refresh(ctx context.Context) error {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		torExitListURL,
		http.NoBody,
	)
	if err != nil {
		return err
	}

	// TODO: Allow to configure HTTP client
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	const nl = "\n"
	lines := strings.Split(string(body), nl)

	list := make([]torIP, 0, len(lines))
	for _, l := range lines {
		if l == "" {
			continue
		}

		if strings.HasPrefix(l, "#") {
			log.Printf("tor: refresh: ignoring IP: %s\n", l)
			continue
		}

		if pip := net.ParseIP(l); pip == nil {
			log.Printf("tor: refresh: not an IP: %s\n", l)
			continue
		}

		list = append(list, torIP(l))
	}

	var isInitialRefresh bool
	tl.listMu.Lock()
	if tl.list == nil {
		isInitialRefresh = true
	}
	tl.list = list
	tl.listMu.Unlock()

	if isInitialRefresh {
		close(tl.initialDone)
	}

	return nil
}

func (tl *torLimiter) waitForInitialListRefresh() {
	<-tl.initialDone
}

func (tl *torLimiter) IsLimited(ip string) error {
	tl.listMu.RLock()
	defer tl.listMu.RUnlock()

	for _, e := range tl.list {
		if ip == e.String() {
			return ErrIPIsTorExitNode
		}
	}

	return nil
}

func WithTor(tc TorConfig) option {
	tl := &torLimiter{
		initialDone: make(chan struct{}),
	}

	if tc.ListRefreshInterval.Seconds() <= 0 {
		tc.ListRefreshInterval = TorDefaultRefreshInterval
	}

	t := time.NewTicker(tc.ListRefreshInterval)
	go func() {
		for {
			if err := tl.refresh(context.TODO()); err != nil {
				log.Printf("tor: failed to refresh: %v\n", err)
			}

			<-t.C
		}
	}()

	if !tc.NoWaitForInitialListRefresh {
		tl.waitForInitialListRefresh()
	}

	return func(l *Limiter) {
		l.torLimiter = tl
	}
}
