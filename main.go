package main

import (
	"autentification_service/startup"
	cfg "autentification_service/startup/config"
	"context"
	"log"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func main() {

	log.SetOutput(os.Stderr)
	// OpenTelemetry
	var err error
	tp, err = initTracer()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	config := cfg.NewConfig()
	log.Println("Starting server Autentification Service...")
	server := startup.NewServer(config)
	server.Start()
}
