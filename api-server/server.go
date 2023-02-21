package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/apollotracing"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/pluralsh/oauth-playground/api-server/auth"
	"github.com/pluralsh/oauth-playground/api-server/clients"
	"github.com/pluralsh/oauth-playground/api-server/graph/generated"
	"github.com/pluralsh/oauth-playground/api-server/graph/resolvers"
	"github.com/rs/cors"

	kratos "github.com/ory/kratos-client-go"
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

	resolver := &resolvers.Resolver{
		C: &clients.ClientWrapper{
			KratosClient: kratosAdminClient,
			KetoClient:   ketoClient,
			HydraClient:  hydraAdminClient,
			Log:          ctrl.Log.WithName("clients"),
		},
	}

	if err := serve(ctx, kratosAdminClient, resolver); err != nil {
		setupLog.Error(err, "failed to serve")
	}
}

func serve(ctx context.Context, kratosClient *kratos.APIClient, resolver *resolvers.Resolver) (err error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router := chi.NewRouter()

	router.Use(auth.Middleware(kratosClient, ctrl.Log.WithName("auth").WithName("middleware")))

	// Add CORS middleware around every request
	// See https://github.com/rs/cors for full option listing
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:8082"}, // TODO: add config for actual hostname
		AllowCredentials: true,
		Debug:            false,
	}).Handler)

	gqlConfig := generated.Config{Resolvers: resolver}

	gqlConfig.Directives.IsAuthenticated = func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
		// userCtx := auth.ForContext(ctx)
		// if userCtx != nil {

		// 	if userCtx.KratosSession != nil {

		// 		session := userCtx.KratosSession

		// 		if !session.GetActive() || session.ExpiresAt.Before(time.Now()) {
		// 			return nil, fmt.Errorf("Access denied: User session not active or has expired.")
		// 		}
		// 	} else {
		// 		// TODO: remove debug log message
		// 		setupLog.Info("auth directive failed")
		// 		return nil, fmt.Errorf("Access denied: User must be logged in.")
		// 	}
		// } else {
		// 	// TODO: remove debug log message
		// 	setupLog.Info("auth directive failed")
		// 	return nil, fmt.Errorf("Access denied: User must be logged in.")
		// }

		// TODO: remove debug log message
		setupLog.Info("auth directive successful")

		// let it pass through if user is authenticated
		return next(ctx)
	}

	//TODO: change all create and delete mutations so that name and namespace are used directly rather than the wrapped in the input
	gqlConfig.Directives.CheckPermissions = func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {

		// setupLog.Info("Scope directive", "object", graphql.GetFieldContext(ctx).Parent)

		// var namespace string

		// namespace = ""

		// // TODO: check that this still works when no variables are used
		// if graphql.GetFieldContext(ctx).Field.Arguments.ForName("namespace") != nil {
		// 	namespaceArg := graphql.GetFieldContext(ctx).Field.Arguments.ForName("namespace")
		// 	namespaceValue, err := namespaceArg.Value.Value(graphql.GetOperationContext(ctx).Variables)
		// 	if err != nil {
		// 		return nil, err
		// 	}
		// 	namespace = namespaceValue.(string)
		// }

		// var operation string

		// if graphql.GetFieldContext(ctx).Object == "Query" || graphql.GetFieldContext(ctx).Object == "Mutation" {

		// 	operation = graphql.GetFieldContext(ctx).Field.Name
		// } else if graphql.GetFieldContext(ctx).Parent != nil {
		// 	// this is for nested objects
		// 	operation = graphql.GetFieldContext(ctx).Object
		// } else {
		// 	return nil, fmt.Errorf("Access denied: The SubjectAccessReview has not yet been implemented for this scenario")
		// }

		// // get the verb and type by splitting the camelCase
		// // for example, `getStorageClass` becomes `get` and `StorageClass`
		// splitOperation := splitCamelCase(operation)

		// if len(splitOperation) != 2 {
		// 	return nil, fmt.Errorf("Access denied: Something went wrong when parsing the operation for Kubernetes RBAC")
		// }

		// var verb, ObjectType string

		// ObjectType = splitOperation[1]

		// // check if first string is empty
		// if splitOperation[0] == "" {
		// 	verb = "get"
		// } else {
		// 	verb = splitOperation[0]
		// }

		// // check if the ObjectType is a plural
		// // if it is, we need to convert it to a singular
		// // for example, `StorageClasses` becomes `StorageClass`
		// pluralize := pluralize.NewClient()
		// if pluralize.IsPlural(ObjectType) {
		// 	ObjectType = pluralize.Singular(ObjectType)
		// }

		// // get the TypeSar from the ObjectType
		// TypeSar, err := sarLookupFunc(ObjectType)
		// if err != nil {
		// 	return nil, fmt.Errorf("Access denied: Failed to check user permissions. %s", err)
		// }

		// if err := auth.UserAuthz(
		// 	ctx,
		// 	kubeClient,
		// 	setupLog,
		// 	TypeSar.Group,
		// 	verb,
		// 	TypeSar.Resource,
		// 	TypeSar.Version,
		// 	namespace); err != nil {

		// 	// TODO: remove debug log message
		// 	setupLog.Info("scope directive failed")
		// 	return nil, fmt.Errorf("Access denied: User is not allowed to '%s' '%s' in namespace '%s'. Error: %s", verb, TypeSar.Resource, namespace, err)
		// }

		// TODO: remove debug log message
		setupLog.Info("scope directive successful")

		// or let it pass through
		return next(ctx)
	}

	gqlSrv := handler.NewDefaultServer(generated.NewExecutableSchema(gqlConfig))
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

	router.Handle("/", playground.Handler("GraphQL playground", "/graphql"))
	router.Handle("/graphql", gqlSrv)

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
	setupLog.Info("connect to http://localhost:" + port + "/ for GraphQL playground")

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
