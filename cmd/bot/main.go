package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/darmayasa221/polymarket-go/internal/infrastructures/clob"
	"github.com/darmayasa221/polymarket-go/internal/infrastructures/config"
	"github.com/darmayasa221/polymarket-go/internal/infrastructures/container/botcontainer"
	"github.com/darmayasa221/polymarket-go/internal/interfaces/signing"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
	log.Println("main: shutdown complete")
}

// run is the real entry point — separating it from main allows defer to work correctly.
func run() error {
	cfg, err := config.Load(".env")
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	bc, err := botcontainer.Build(botcontainer.FromAppConfig(cfg))
	if err != nil {
		return fmt.Errorf("botcontainer: %w", err)
	}
	defer func() { _ = bc.Close() }()

	if err := bc.RunMigration(context.Background()); err != nil {
		return fmt.Errorf("migration: %w", err)
	}

	// PRIVATE_KEY is read from env, never stored in config struct.
	signer, err := signing.NewSigner(os.Getenv("PRIVATE_KEY"))
	if err != nil {
		return fmt.Errorf("signer: %w", err)
	}

	clobConfig := clob.Config{
		BaseURL:       cfg.CLOB.BaseURL,
		Address:       cfg.CLOB.Address,
		APIKey:        cfg.CLOB.APIKey,
		APISecret:     cfg.CLOB.APISecret,
		APIPassphrase: cfg.CLOB.APIPassphrase,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	r := newRunner(bc, signer, clobConfig)
	return r.run(ctx)
}
