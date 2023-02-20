package opl_test

import (
	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
)

var scenario_1 = []*rts.RelationTuple{
	//-------- create Groups and Users ---------
	// Group: AllUsers
	// User objects are created implicitly through member relation tuples to a group `AllUsers` that contains all users
	// In a live Ory Kratos/Keto setup this group should reflect the users that are registered with Kratos
	{
		Namespace: "Group",
		Object:    "AllUsers",
		Relation:  "usermember",
		Subject: rts.NewSubjectSet(
			"User",
			"Hans",
			"",
		),
	},
	{
		Namespace: "Group",
		Object:    "AllUsers",
		Relation:  "usermember",
		Subject: rts.NewSubjectSet(
			"User",
			"David",
			"",
		),
	},
	{
		Namespace: "Group",
		Object:    "AllUsers",
		Relation:  "usermember",
		Subject: rts.NewSubjectSet(
			"User",
			"Nico",
			"",
		),
	},
	{
		Namespace: "Group",
		Object:    "AllUsers",
		Relation:  "usermember",
		Subject: rts.NewSubjectSet(
			"User",
			"Lianet",
			"",
		),
	},
	{
		Namespace: "Group",
		Object:    "AllUsers",
		Relation:  "usermember",
		Subject: rts.NewSubjectSet(
			"User",
			"Sophie",
			"",
		),
	},
	// Group: Ops
	// Some users are members of the group `Ops` that should contain users with admin access intended for administrative tasks
	{
		Namespace: "Group",
		Object:    "Ops",
		Relation:  "usermember",
		Subject: rts.NewSubjectSet(
			"User",
			"Hans",
			"",
		),
	},
	{
		Namespace: "Group",
		Object:    "Ops",
		Relation:  "usermember",
		Subject: rts.NewSubjectSet(
			"User",
			"David",
			"",
		),
	},
	//-------- create a project ---------
	// access does not entail any permissions, it is just a relation that's there to confirm that a user is registered with the project
	// if a user is not registered with a project they cannot perform any actions on resources inside the project even if there are still roles and policies in place that allow certain actions
	// that way we can for example revoke all access by severing the access to a project without deleting all policies or roles
	// Project: Manhattan

	// all members of group `Ops` have access to the project `Manhattan`
	{
		Namespace: "Project",
		Object:    "Manhattan",
		Relation:  "groupaccess",
		Subject: rts.NewSubjectSet(
			"Group",
			"Ops",
			"",
		),
	},
	// additionally Nico and Lianet have explicit access to the project `Manhattan`
	{
		Namespace: "Project",
		Object:    "Manhattan",
		Relation:  "useraccess",
		Subject: rts.NewSubjectSet(
			"User",
			"Nico",
			"",
		),
	},
	{
		Namespace: "Project",
		Object:    "Manhattan",
		Relation:  "useraccess",
		Subject: rts.NewSubjectSet(
			"User",
			"Lianet",
			"",
		),
	},

	//-------- create roles ---------
	// every action from a project member is performed on behalf of a principal
	// at this time only roles can act as principals

	// Role: Admin
	// we create a role `Admin` as a principal of project `Manhattan` that will get various wide ranging permissions through appropriate policies
	{
		Namespace: "Role",
		Object:    "Admin",
		Relation:  "principal",
		Subject: rts.NewSubjectSet(
			"Project",
			"Manhattan",
			"",
		),
	},
	// Role: Dev
	// we create a role `Dev` as a principal of project `Manhattan` that will get more narrow permissions
	{
		Namespace: "Role",
		Object:    "Dev",
		Relation:  "principal",
		Subject: rts.NewSubjectSet(
			"Project",
			"Manhattan",
			"",
		),
	},

	//-------- Create policies ---------
	// policies bundle permissions that are granted to a role

	// Policy: AdminPolicy
	{
		Namespace: "Policy",
		Object:    "AdminPolicy",
		Relation:  "trust",
		Subject: rts.NewSubjectSet(
			"Role",
			"Admin",
			"",
		),
	},
	// Policy: DevPolicy
	{
		Namespace: "Policy",
		Object:    "DevPolicy",
		Relation:  "trust",
		Subject: rts.NewSubjectSet(
			"Role",
			"Dev",
			"",
		),
	},

	//-------- create resource types and implicitly create permissions ---------
	// The resource types namespaces comprise all resources that are available in the system
	// They act as classes of resources for which you can create permissions by attaching them to policies
	// Any principal with a bound policy that has a permission for a certain resource type can perform the granted actions on all instances of resources of that type
	// It's also the only place where it makes sense to define permission relations for the creation of resources

	// KubernetesResourceType: Service
	// All Kubernetes primitive types (and maybe by extension CRDs, not sure yet) should be defined in the KubernetesResourceType namespace
	// For each action (rule verb) that is defined in the Kubernetes API for a resource type we create an according permission relation tuple

	// create permissions for the `Secret` Kubernetes resource type and bundle them in the `AdminPolicy`
	// this could be reconciled from existing Kubernetes roles
	{
		Namespace: "KubernetesResourceType",
		Object:    "Secret",
		Relation:  "create",
		Subject: rts.NewSubjectSet(
			"Policy",
			"AdminPolicy",
			"",
		),
	},
	{
		Namespace: "KubernetesResourceType",
		Object:    "Secret",
		Relation:  "delete",
		Subject: rts.NewSubjectSet(
			"Policy",
			"AdminPolicy",
			"",
		),
	},
	{
		Namespace: "KubernetesResourceType",
		Object:    "Secret",
		Relation:  "get",
		Subject: rts.NewSubjectSet(
			"Policy",
			"AdminPolicy",
			"",
		),
	},
	{
		Namespace: "KubernetesResourceType",
		Object:    "Secret",
		Relation:  "list",
		Subject: rts.NewSubjectSet(
			"Policy",
			"AdminPolicy",
			"",
		),
	},

	// create certain narrow permissions for the `Secret` Kubernetes resource type and bundle them in the `DevPolicy`
	{
		Namespace: "KubernetesResourceType",
		Object:    "Secret",
		Relation:  "list",
		Subject: rts.NewSubjectSet(
			"Policy",
			"DevPolicy",
			"",
		),
	},

	// KubricksResourceType: MLFlow

	// All Kubricks apps/plugins should be defined in the KubricksResourceType namespace
	// TODO: We shall see if that makes sense or if we need to create a separate namespace for each app/plugin
	// e.g. because we want to define permissions actions or individual http endpoints that don't exist on every app/plugin
	// the `accessapi` shall stand in as a placeholder for all data plane access to the Kubricks app/plugin

	// an MLFlow instance needs a secret to be able to access the object store (for example)
	// if any user can create/delete Secrets in a project, they can also get/set the particular secret of an MLFlow instance
	{
		Namespace: "KubricksResourceType",
		Object:    "MLFlow",
		Relation:  "hassecret",
		Subject: rts.NewSubjectSet(
			"KubernetesResourceType",
			"Secret",
			"",
		),
	},

	// admins can create, manipulate the instance, as well as access the MLFlow API
	{
		Namespace: "KubricksResourceType",
		Object:    "MLFlow",
		Relation:  "create",
		Subject: rts.NewSubjectSet(
			"Policy",
			"AdminPolicy",
			"",
		),
	},
	{
		Namespace: "KubricksResourceType",
		Object:    "MLFlow",
		Relation:  "delete",
		Subject: rts.NewSubjectSet(
			"Policy",
			"AdminPolicy",
			"",
		),
	},
	{
		Namespace: "KubricksResourceType",
		Object:    "MLFlow",
		Relation:  "get",
		Subject: rts.NewSubjectSet(
			"Policy",
			"AdminPolicy",
			"",
		),
	},
	{
		Namespace: "KubricksResourceType",
		Object:    "MLFlow",
		Relation:  "list",
		Subject: rts.NewSubjectSet(
			"Policy",
			"AdminPolicy",
			"",
		),
	},
	{
		Namespace: "KubricksResourceType",
		Object:    "MLFlow",
		Relation:  "update",
		Subject: rts.NewSubjectSet(
			"Policy",
			"AdminPolicy",
			"",
		),
	},
	{
		Namespace: "KubricksResourceType",
		Object:    "MLFlow",
		Relation:  "accessapi",
		Subject: rts.NewSubjectSet(
			"Policy",
			"AdminPolicy",
			"",
		),
	},

	// devs get info on any instance, as well as access to the MLFlow API, they also can get the secret of an instance, but not set it
	{
		Namespace: "KubricksResourceType",
		Object:    "MLFlow",
		Relation:  "accessapi",
		Subject: rts.NewSubjectSet(
			"Policy",
			"DevPolicy",
			"",
		),
	},
	{
		Namespace: "KubricksResourceType",
		Object:    "MLFlow",
		Relation:  "list",
		Subject: rts.NewSubjectSet(
			"Policy",
			"DevPolicy",
			"",
		),
	},
	{
		Namespace: "KubricksResourceType",
		Object:    "MLFlow",
		Relation:  "get",
		Subject: rts.NewSubjectSet(
			"Policy",
			"DevPolicy",
			"",
		),
	},
	{
		Namespace: "KubricksResourceType",
		Object:    "MLFlow",
		Relation:  "getsecret",
		Subject: rts.NewSubjectSet(
			"KubernetesResourceType",
			"Secret",
			"",
		),
	},

	// ServiceResourceType: MLFlowInstance

	// Finally we create an actual instance of the KubricksResourceType: MLFlow
	{
		Namespace: "KubricksResource",
		Object:    "ManhattanMLFlowInstance",
		Relation:  "owner",
		Subject: rts.NewSubjectSet(
			"User",
			"Hans",
			"",
		),
	},
	{
		Namespace: "KubricksResource",
		Object:    "ManhattanMLFlowInstance",
		Relation:  "kbrx_instance",
		Subject: rts.NewSubjectSet(
			"KubricksResourceType",
			"MLFlow",
			"",
		),
	},
	// in addition to any principal of project Manhattan being able to access the MLFlow API, we also want to allow an external user Sophie to access the API
	{
		Namespace: "KubricksResource",
		Object:    "ManhattanMLFlowInstance",
		Relation:  "accessapi",
		Subject: rts.NewSubjectSet(
			"ResourcePolicy",
			"ShareManhattanMLFlowAccess",
			"",
		),
	},
	{
		Namespace: "ResourcePolicy",
		Object:    "ManhattanMLFlowInstance",
		Relation:  "trust",
		Subject: rts.NewSubjectSet(
			"User",
			"Sophie",
			"",
		),
	},
}
