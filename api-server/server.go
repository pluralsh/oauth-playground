package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"time"

	gqlHandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/apollotracing"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/pluralsh/oauth-playground/api-server/clients"
	"github.com/pluralsh/oauth-playground/api-server/graph/directives"
	"github.com/pluralsh/oauth-playground/api-server/graph/generated"
	"github.com/pluralsh/oauth-playground/api-server/graph/resolvers"
	"github.com/pluralsh/oauth-playground/api-server/handlers"
	"github.com/rs/cors"

	_ "github.com/mattn/go-sqlite3"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

const defaultPort = "8082"

var (
	setupLog = ctrl.Log.WithName("setup")
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	opts := zap.Options{}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()
	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	kratosAdminClient, err := clients.NewKratosAdminClient()
	if err != nil {
		setupLog.Error(err, "An admin address for kratos must be configured")
		panic(err)
	}

	kratosPublicClient, err := clients.NewKratosPublicClient()
	if err != nil {
		setupLog.Error(err, "An admin address for kratos must be configured")
		panic(err)
	}

	conndetails := clients.NewKetoConnectionDetailsFromEnv()
	ketoClient, err := clients.NewKetoGrpcClient(context.Background(), conndetails)
	if err != nil {
		setupLog.Error(err, "Failed to setup Keto gRPC client")
		panic(err)
	}

	hydraAdminClient, err := clients.NewHydraAdminClient()
	if err != nil {
		setupLog.Error(err, "An admin address for hydra must be configured")
		panic(err)
	}

	clientWrapper := &clients.ClientWrapper{
		KratosAdminClient:  kratosAdminClient,
		KratosPublicClient: kratosPublicClient,
		KetoClient:         ketoClient,
		HydraClient:        hydraAdminClient,
		Log:                ctrl.Log.WithName("clients"),
	}

	resolver := &resolvers.Resolver{
		C: clientWrapper,
	}

	directives := &directives.Directive{
		C: clientWrapper,
	}

	handlers := &handlers.Handler{
		C:   clientWrapper,
		Log: ctrl.Log.WithName("handlers"),
	}

	if err := serve(ctx, resolver, directives, handlers); err != nil {
		setupLog.Error(err, "failed to serve")
	}
}

func serve(ctx context.Context, resolver *resolvers.Resolver, directives *directives.Directive, handlers *handlers.Handler) (err error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router := chi.NewRouter()

	router.Use(handlers.Middleware())

	// Add CORS middleware around every request
	// See https://github.com/rs/cors for full option listing
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:8082"}, // TODO: add config for actual hostname
		AllowCredentials: true,
		Debug:            false,
	}).Handler)

	gqlConfig := generated.Config{Resolvers: resolver}

	gqlConfig.Directives.IsAuthenticated = directives.IsAuthenticated

	//TODO: change all create and delete mutations so that name and namespace are used directly rather than the wrapped in the input
	gqlConfig.Directives.CheckPermissions = directives.CheckPermissions

	gqlSrv := gqlHandler.NewDefaultServer(generated.NewExecutableSchema(gqlConfig))
	gqlSrv.Use(apollotracing.Tracer{})
	// gqlSrv.AddTransport(&transport.Websocket{
	//     Upgrader: websocket.Upgrader{
	//         CheckOrigin: func(r *http.Request) bool {
	//             // Check against your desired domains here
	// TODO: add domain to Kubricks Config CRD
	//              return r.Host == foundKubricksConfig.Spec.Domain
	//         },
	//         ReadBufferSize:  1024,
	//         WriteBufferSize: 1024,
	//     },
	// })

	router.Handle("/graphiql", playground.Handler("GraphQL playground", "/graphql"))
	router.Handle("/graphql", gqlSrv)
	router.Post("/tenant-hydrator", handlers.HydrateObservabilityTenants)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			setupLog.Error(err, "failed to serve GraphQL API")
		}
	}()

	setupLog.Info("server started")
	setupLog.Info("connect to http://localhost:" + port + "/graphiql for GraphQL playground")

	<-ctx.Done()

	setupLog.Info("server stopped")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err = srv.Shutdown(ctxShutDown); err != nil {
		setupLog.Error(err, "server shutdown failed")
	}

	setupLog.Info("server exited properly")

	if err == http.ErrServerClosed {
		err = nil
	}

	return

}
