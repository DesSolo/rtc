package file

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"

	"github.com/fsnotify/fsnotify"

	"rtc/pkg/rtc"
)

// Provider ...
type Provider struct {
	watcher  *fsnotify.Watcher
	reader   Reader
	filePath string

	items    map[rtc.Key]rtc.Value
	itemsMux sync.RWMutex

	callbacks    map[rtc.Key]rtc.ValueChangeCallback
	callbacksMux sync.RWMutex
}

// NewProvider ...
func NewProvider(filePath string, reader Reader) (*Provider, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("fsnotify.NewWatcher: %v", err)
	}

	if err = w.Add(filePath); err != nil {
		return nil, fmt.Errorf("fsnotify.Add: %v", err)
	}

	p := Provider{
		watcher:   w,
		reader:    reader,
		filePath:  filePath,
		items:     make(map[rtc.Key]rtc.Value),
		callbacks: make(map[rtc.Key]rtc.ValueChangeCallback),
	}

	if err := p.read(); err != nil {
		return nil, fmt.Errorf("read: %w", err)
	}

	go p.run()

	return &p, nil
}

// Value ...
func (p *Provider) Value(_ context.Context, key rtc.Key) (rtc.Value, error) {
	p.itemsMux.RLock()
	defer p.itemsMux.RUnlock()

	item, ok := p.items[key]
	if !ok {
		return nil, rtc.ErrNotPresent
	}

	return item, nil
}

// WatchValue ...
func (p *Provider) WatchValue(_ context.Context, key rtc.Key, handler rtc.ValueChangeCallback) error {
	p.callbacksMux.Lock()
	defer p.callbacksMux.Unlock()

	p.callbacks[key] = handler

	return nil
}

// Close ...
func (p *Provider) Close() error {
	return p.watcher.Close()
}

func (p *Provider) read() error {
	data, err := os.ReadFile(p.filePath)
	if err != nil {
		return fmt.Errorf("os.ReadFile: %v", err)
	}

	items, err := p.reader.Read(data)
	if err != nil {
		return fmt.Errorf("reader.Read: %v", err)
	}

	if len(p.callbacks) != 0 {
		p.callbacksMux.RLock()
		defer p.callbacksMux.RUnlock()

		p.itemsMux.Lock()
		defer p.itemsMux.Unlock()

		for cn, cf := range p.callbacks {
			if newVal, ok := items[cn]; ok {
				cf(p.items[cn], newVal)
			}
		}

		p.items = items

		return nil
	}

	p.itemsMux.Lock()
	defer p.itemsMux.Unlock()

	p.items = items

	return nil
}

func (p *Provider) run() {
	for event := range p.watcher.Events {
		if !event.Op.Has(fsnotify.Write) {
			continue
		}

		if err := p.read(); err != nil {
			slog.Error("w.read", "err", err)
		}
	}
}
