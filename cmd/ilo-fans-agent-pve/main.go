package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Puhovsky-Contributions/ilo-fans-agent-pve/internal/auth"
	"github.com/Puhovsky-Contributions/ilo-fans-agent-pve/internal/config"
	"github.com/Puhovsky-Contributions/ilo-fans-agent-pve/internal/httpapi"
	"github.com/Puhovsky-Contributions/ilo-fans-agent-pve/internal/version"
)

func main() {
	if handleVersion(os.Args[1:], os.Stdout) {
		return
	}

	configPath := flag.String("config", "/etc/ilo-fans-agent-pve/config.yaml", "path to config.yaml")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintln(flag.CommandLine.Output(), "\nCommands:")
		fmt.Fprintln(flag.CommandLine.Output(), "  version [--short|-s]    show agent version and build information")
		fmt.Fprintln(flag.CommandLine.Output(), "  token init|rotate|show  manage bearer auth token")
	}
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

func handleVersion(args []string, out io.Writer) bool {
	if len(args) == 0 {
		return false
	}
	var isVersion bool
	var remaining []string
	for _, arg := range args {
		switch arg {
		case "version", "-v", "-version", "--version":
			isVersion = true
		default:
			remaining = append(remaining, arg)
		}
	}
	if !isVersion {
		return false
	}
	runVersionCLI(remaining, out)
	return true
}

func runVersionCLI(args []string, out io.Writer) {
	var short bool
	for _, a := range args {
		if a == "--short" || a == "-short" || a == "-s" || a == "--s" {
			short = true
			break
		}
	}
	info := version.Get()
	if short {
		fmt.Fprintln(out, info.Short())
	} else {
		fmt.Fprintln(out, info.String())
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
