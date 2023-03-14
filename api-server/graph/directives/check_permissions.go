package directives

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
)

func (d *Directive) CheckPermissions(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {

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
	d.C.Log.Info("scope directive successful")

	// or let it pass through
	return next(ctx)
}
