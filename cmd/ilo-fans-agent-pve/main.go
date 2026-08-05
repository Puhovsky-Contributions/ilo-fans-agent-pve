package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/playtika/ilo-fans-agent-pve/internal/auth"
	"github.com/playtika/ilo-fans-agent-pve/internal/config"
	"github.com/playtika/ilo-fans-agent-pve/internal/httpapi"
)

func main() {
	configPath := flag.String("config", "/etc/ilo-fans-agent-pve/config.yaml", "path to config.yaml")
	flag.Parse()

	args := flag.Args()
	if len(args) > 0 && args[0] == "token" {
		runTokenCLI(*configPath, args[1:])
		return
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	store := auth.NewStore(cfg.TokenFile)
	if err := store.Load(); err != nil {
		log.Fatalf("token: %v (run: ilo-fans-agent-pve token init)", err)
	}

	srv := httpapi.New(cfg, store)
	httpServer := &http.Server{
		Addr:    cfg.Listen,
		Handler: srv.Handler(),
	}

	reload := make(chan os.Signal, 1)
	signal.Notify(reload, syscall.SIGHUP)
	go func() {
		for range reload {
			if err := store.Load(); err != nil {
				log.Printf("SIGHUP token reload failed: %v", err)
			} else {
				log.Printf("token reloaded (%s)", store.ShowMasked())
			}
		}
	}()

	log.Printf("listening on %s", cfg.Listen)
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func runTokenCLI(configPath string, sub []string) {
	if len(sub) == 0 {
		fmt.Fprintln(os.Stderr, "usage: token init|rotate|show")
		os.Exit(2)
	}
	cfg, err := config.Load(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "config: %v\n", err)
		os.Exit(1)
	}
	store := auth.NewStore(cfg.TokenFile)
	_ = store.Load()

	switch sub[0] {
	case "init":
		tok, err := store.Init()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		fmt.Println(tok)
	case "rotate":
		tok, err := store.Rotate()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		fmt.Println(tok)
	case "show":
		if err := store.Load(); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		fmt.Println(store.ShowMasked())
	default:
		fmt.Fprintln(os.Stderr, "unknown token subcommand")
		os.Exit(2)
	}
}
