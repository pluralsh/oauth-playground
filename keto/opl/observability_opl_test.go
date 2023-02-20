package opl_test

import (
	"context"

	"github.com/google/go-cmp/cmp"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	px "github.com/ory/x/pointerx"
	"github.com/pluralsh/oauth-playground/keto/client"
	"google.golang.org/protobuf/testing/protocmp"

	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
)

var _ = Describe("Verify expected behaviour of the opl configuration.", func() {
	var _ = Describe("Scenario to cover most constellations.", func() {
		BeforeEach(func() {
			//set up database before each test
			err := kcl.CreateTuples(context.Background(), observability_scenario_admins)
			if err != nil {
				panic("Encountered error: " + err.Error())
			}
			err = kcl.CreateTuples(context.Background(), observability_scenario_users_groups)
			if err != nil {
				panic("Encountered error: " + err.Error())
			}
			err = kcl.CreateTuples(context.Background(), observability_scenario_clients_tenants)
			if err != nil {
				panic("Encountered error: " + err.Error())
			}

		})
		AfterEach(func() {
			//tear down database entries after each test
			query := rts.RelationQuery{
				Namespace: nil,
				Object:    nil,
				Relation:  nil,
				Subject:   nil,
			}
			err := kcl.DeleteAllTuples(context.Background(), &query)
			if err != nil {
				panic("Encountered error: " + err.Error())
			}
		})
		It("should be able to list admins", func() {
			query := rts.RelationQuery{
				Namespace: px.Ptr("Organization"),
				Object:    px.Ptr("main"),
				Relation:  px.Ptr("admins"),
				Subject:   nil,
			}

			respTuples, err := kcl.QueryAllTuples(context.Background(), &query, 100)
			if err != nil {
				panic("Encountered error: " + err.Error())
			}
			// need to sort the tuples to make sure they are in the same order as the expected tuples
			client.SortRelationTuples(respTuples)

			client.PrintTableFromRelationTuples(respTuples, GinkgoWriter)
			Expect(respTuples).To(HaveLen(2))

			equal := cmp.Equal(respTuples, observability_scenario_admins, protocmp.Transform())

			if !equal {
				diff := cmp.Diff(respTuples, observability_scenario_admins, protocmp.Transform())
				GinkgoWriter.Printf("Diff: %s", diff)
			}

			Expect(equal).To(BeTrue())
		})
		It("check if a user is an admin", func() {
			query := rts.RelationTuple{
				Namespace: "Organization",
				Object:    "main",
				Relation:  "admins",
				Subject: rts.NewSubjectSet(
					"Admin",
					"david",
					"",
				),
			}

			ok, err := kcl.Check(context.Background(), &query)
			if err != nil {
				panic("Encountered error: " + err.Error())
			}
			GinkgoWriter.Printf("Admin `david` is an admin: %v\n", ok)
			Expect(ok).To(BeTrue())
		})
		It("admin user should be able to create new group", func() {
			query := rts.RelationTuple{
				Namespace: "Group",
				Object:    "",
				Relation:  "create",
				Subject: rts.NewSubjectSet(
					"Admin",
					"david",
					"",
				),
			}

			ok, err := kcl.Check(context.Background(), &query)
			if err != nil {
				panic("Encountered error: " + err.Error())
			}
			Expect(ok).To(BeTrue())
			GinkgoWriter.Printf("Admin `david` can create new Group: %v\n", ok)
		})
		It("admin user should be able to edit existing group", func() {
			query := rts.RelationTuple{
				Namespace: "Group",
				Object:    "MainCluster",
				Relation:  "edit",
				Subject: rts.NewSubjectSet(
					"Admin",
					"david",
					"",
				),
			}

			ok, err := kcl.Check(context.Background(), &query)
			if err != nil {
				panic("Encountered error: " + err.Error())
			}
			Expect(ok).To(BeTrue())
			GinkgoWriter.Printf("Admin `david` can edit Group `MainCluster`: %v\n", ok)
		})
		It("admin user should be able to delete existing group", func() {
			query := rts.RelationTuple{
				Namespace: "Group",
				Object:    "MainCluster",
				Relation:  "delete",
				Subject: rts.NewSubjectSet(
					"Admin",
					"david",
					"",
				),
			}

			ok, err := kcl.Check(context.Background(), &query)
			if err != nil {
				panic("Encountered error: " + err.Error())
			}
			Expect(ok).To(BeTrue())
			GinkgoWriter.Printf("Admin `david` can delete Group `MainCluster`: %v\n", ok)
		})
		It("admin user should not be able to delete the AllUsers Group", func() {
			query := rts.RelationTuple{
				Namespace: "Group",
				Object:    "AllUsers",
				Relation:  "delete",
				Subject: rts.NewSubjectSet(
					"Admin",
					"david",
					"",
				),
			}

			ok, err := kcl.Check(context.Background(), &query)
			if err != nil {
				panic("Encountered error: " + err.Error())
			}
			Expect(ok).To(BeFalse())
			GinkgoWriter.Printf("Admin `david` can't delete Group `AllUsers`: %v\n", ok)
		})
		It("admin user should be able to create new user", func() {
			query := rts.RelationTuple{
				Namespace: "User",
				Object:    "",
				Relation:  "create",
				Subject: rts.NewSubjectSet(
					"Admin",
					"hans",
					"",
				),
			}

			ok, err := kcl.Check(context.Background(), &query)
			if err != nil {
				panic("Encountered error: " + err.Error())
			}
			Expect(ok).To(BeTrue())
			GinkgoWriter.Printf("Admin `hans` can create new User: %v\n", ok)
		})
		It("admin user should be able to edit existing user", func() {
			query := rts.RelationTuple{
				Namespace: "User",
				Object:    "sam",
				Relation:  "edit",
				Subject: rts.NewSubjectSet(
					"Admin",
					"hans",
					"",
				),
			}

			ok, err := kcl.Check(context.Background(), &query)
			if err != nil {
				panic("Encountered error: " + err.Error())
			}
			Expect(ok).To(BeTrue())
			GinkgoWriter.Printf("Admin `hans` can edit User `sam`: %v\n", ok)
		})
		It("admin user should be able to delete existing user", func() {
			query := rts.RelationTuple{
				Namespace: "User",
				Object:    "nick",
				Relation:  "delete",
				Subject: rts.NewSubjectSet(
					"Admin",
					"hans",
					"",
				),
			}

			ok, err := kcl.Check(context.Background(), &query)
			if err != nil {
				panic("Encountered error: " + err.Error())
			}
			Expect(ok).To(BeTrue())
			GinkgoWriter.Printf("Admin `hans` can delete User `nick`: %v\n", ok)
		})
		It("admin user should be able to create a new OAuth2Client", func() {
			query := rts.RelationTuple{
				Namespace: "OAuth2Client",
				Object:    "",
				Relation:  "create",
				Subject: rts.NewSubjectSet(
					"Admin",
					"david",
					"",
				),
			}

			ok, err := kcl.Check(context.Background(), &query)
			if err != nil {
				panic("Encountered error: " + err.Error())
			}
			Expect(ok).To(BeTrue())
			GinkgoWriter.Printf("Admin `david` can create new OAuth2Client: %v\n", ok)
		})
		It("admin user should be able to edit an existing OAuth2Client", func() {
			query := rts.RelationTuple{
				Namespace: "OAuth2Client",
				Object:    "MainClusterGrafana",
				Relation:  "edit",
				Subject: rts.NewSubjectSet(
					"Admin",
					"david",
					"",
				),
			}

			ok, err := kcl.Check(context.Background(), &query)
			if err != nil {
				panic("Encountered error: " + err.Error())
			}
			Expect(ok).To(BeTrue())
			GinkgoWriter.Printf("Admin `david` can edit OAuth2Client `MainClusterGrafana`: %v\n", ok)
		})
		It("admin user should be able to delete an existing OAuth2Client", func() {
			query := rts.RelationTuple{
				Namespace: "OAuth2Client",
				Object:    "MainClusterGrafana",
				Relation:  "delete",
				Subject: rts.NewSubjectSet(
					"Admin",
					"david",
					"",
				),
			}

			ok, err := kcl.Check(context.Background(), &query)
			if err != nil {
				panic("Encountered error: " + err.Error())
			}
			Expect(ok).To(BeTrue())
			GinkgoWriter.Printf("Admin `david` can delete OAuth2Client `MainClusterGrafana`: %v\n", ok)
		})
		It("admin user should be able to login through an existing OAuth2Client", func() {
			query := rts.RelationTuple{
				Namespace: "OAuth2Client",
				Object:    "MainClusterGrafana",
				Relation:  "login",
				Subject: rts.NewSubjectSet(
					"Admin",
					"david",
					"",
				),
			}

			ok, err := kcl.Check(context.Background(), &query)
			if err != nil {
				panic("Encountered error: " + err.Error())
			}
			Expect(ok).To(BeTrue())
			GinkgoWriter.Printf("Admin `david` can login to OAuth2Client `MainClusterGrafana`: %v\n", ok)
		})
		It("User in Group with login binding to OAuth2Client should be able to login", func() {
			query := rts.RelationTuple{
				Namespace: "OAuth2Client",
				Object:    "MainClusterGrafana",
				Relation:  "login",
				Subject: rts.NewSubjectSet(
					"User",
					"sam",
					"",
				),
			}

			ok, err := kcl.Check(context.Background(), &query)
			if err != nil {
				panic("Encountered error: " + err.Error())
			}
			Expect(ok).To(BeTrue())
			GinkgoWriter.Printf("User `sam` can login to OAuth2Client `MainClusterGrafana`: %v\n", ok)
		})
		It("User with login binding to OAuth2Client should be able to login", func() {
			query := rts.RelationTuple{
				Namespace: "OAuth2Client",
				Object:    "MainClusterGrafana",
				Relation:  "login",
				Subject: rts.NewSubjectSet(
					"User",
					"nick",
					"",
				),
			}

			ok, err := kcl.Check(context.Background(), &query)
			if err != nil {
				panic("Encountered error: " + err.Error())
			}
			Expect(ok).To(BeTrue())
			GinkgoWriter.Printf("User `nick` can login to OAuth2Client `MainClusterGrafana`: %v\n", ok)
		})
		It("User in Group without login binding to OAuth2Client should not be able to login", func() {
			query := rts.RelationTuple{
				Namespace: "OAuth2Client",
				Object:    "MainClusterGrafana",
				Relation:  "login",
				Subject: rts.NewSubjectSet(
					"User",
					"aaron",
					"",
				),
			}

			ok, err := kcl.Check(context.Background(), &query)
			if err != nil {
				panic("Encountered error: " + err.Error())
			}
			Expect(ok).To(BeFalse())
			GinkgoWriter.Printf("User `arron` can't login to OAuth2Client `MainClusterGrafana`: %v\n", ok)
		})
		It("User without login binding to OAuth2Client should not be able to login", func() {
			query := rts.RelationTuple{
				Namespace: "OAuth2Client",
				Object:    "MainClusterGrafana",
				Relation:  "login",
				Subject: rts.NewSubjectSet(
					"User",
					"cris",
					"",
				),
			}

			ok, err := kcl.Check(context.Background(), &query)
			if err != nil {
				panic("Encountered error: " + err.Error())
			}
			Expect(ok).To(BeFalse())
			GinkgoWriter.Printf("User `cris` can't login to OAuth2Client `MainClusterGrafana`: %v\n", ok)
		})
		It("admin user should be able to create a new ObservabilityTenant", func() {
			query := rts.RelationTuple{
				Namespace: "ObservabilityTenant",
				Object:    "",
				Relation:  "create",
				Subject: rts.NewSubjectSet(
					"Admin",
					"david",
					"",
				),
			}

			ok, err := kcl.Check(context.Background(), &query)
			if err != nil {
				panic("Encountered error: " + err.Error())
			}
			Expect(ok).To(BeTrue())
			GinkgoWriter.Printf("Admin `david` can create new ObservabilityTenant: %v\n", ok)
		})
		It("admin user should be able to edit an existing ObservabilityTenant", func() {
			query := rts.RelationTuple{
				Namespace: "ObservabilityTenant",
				Object:    "MainCluster",
				Relation:  "edit",
				Subject: rts.NewSubjectSet(
					"Admin",
					"david",
					"",
				),
			}

			ok, err := kcl.Check(context.Background(), &query)
			if err != nil {
				panic("Encountered error: " + err.Error())
			}
			Expect(ok).To(BeTrue())
			GinkgoWriter.Printf("Admin `david` can edit ObservabilityTenant `MainCluster`: %v\n", ok)
		})
		It("admin user should be able to delete an existing ObservabilityTenant", func() {
			query := rts.RelationTuple{
				Namespace: "ObservabilityTenant",
				Object:    "MainCluster",
				Relation:  "delete",
				Subject: rts.NewSubjectSet(
					"Admin",
					"david",
					"",
				),
			}

			ok, err := kcl.Check(context.Background(), &query)
			if err != nil {
				panic("Encountered error: " + err.Error())
			}
			Expect(ok).To(BeTrue())
			GinkgoWriter.Printf("Admin `david` can delete ObservabilityTenant `MainCluster`: %v\n", ok)
		})
		It("admin user should be able to view an existing ObservabilityTenant", func() {
			query := rts.RelationTuple{
				Namespace: "ObservabilityTenant",
				Object:    "MainCluster",
				Relation:  "view",
				Subject: rts.NewSubjectSet(
					"Admin",
					"david",
					"",
				),
			}

			ok, err := kcl.Check(context.Background(), &query)
			if err != nil {
				panic("Encountered error: " + err.Error())
			}
			Expect(ok).To(BeTrue())
			GinkgoWriter.Printf("Admin `david` can view to ObservabilityTenant `MainCluster`: %v\n", ok)
		})
		It("User with view permissions should be able to view an existing ObservabilityTenant", func() {
			query := rts.RelationTuple{
				Namespace: "ObservabilityTenant",
				Object:    "MainCluster",
				Relation:  "view",
				Subject: rts.NewSubjectSet(
					"User",
					"sam",
					"",
				),
			}

			ok, err := kcl.Check(context.Background(), &query)
			if err != nil {
				panic("Encountered error: " + err.Error())
			}
			Expect(ok).To(BeTrue())
			GinkgoWriter.Printf("User `sam` can view ObservabilityTenant `MainCluster`: %v\n", ok)
		})
		It("User with view permissions should not be able to edit an existing ObservabilityTenant", func() {
			query := rts.RelationTuple{
				Namespace: "ObservabilityTenant",
				Object:    "MainCluster",
				Relation:  "edit",
				Subject: rts.NewSubjectSet(
					"User",
					"sam",
					"",
				),
			}

			ok, err := kcl.Check(context.Background(), &query)
			if err != nil {
				panic("Encountered error: " + err.Error())
			}
			Expect(ok).To(BeFalse())
			GinkgoWriter.Printf("User `sam` can't edit ObservabilityTenant `MainCluster`: %v\n", ok)
		})
		It("User in Group with edit permissions should be able to view an existing ObservabilityTenant", func() {
			query := rts.RelationTuple{
				Namespace: "ObservabilityTenant",
				Object:    "MainCluster",
				Relation:  "view",
				Subject: rts.NewSubjectSet(
					"User",
					"nick",
					"",
				),
			}

			ok, err := kcl.Check(context.Background(), &query)
			if err != nil {
				panic("Encountered error: " + err.Error())
			}
			Expect(ok).To(BeTrue())
			GinkgoWriter.Printf("User `nick` can view ObservabilityTenant `MainCluster`: %v\n", ok)
		})
		It("User in Group with edit permissions should be able to edit an existing ObservabilityTenant", func() {
			query := rts.RelationTuple{
				Namespace: "ObservabilityTenant",
				Object:    "MainCluster",
				Relation:  "edit",
				Subject: rts.NewSubjectSet(
					"User",
					"nick",
					"",
				),
			}

			ok, err := kcl.Check(context.Background(), &query)
			if err != nil {
				panic("Encountered error: " + err.Error())
			}
			Expect(ok).To(BeTrue())
			GinkgoWriter.Printf("User `nick` can edit ObservabilityTenant `MainCluster`: %v\n", ok)
		})
		It("OAuth2Client with view permissions should be able to view an existing ObservabilityTenant", func() {
			query := rts.RelationTuple{
				Namespace: "ObservabilityTenant",
				Object:    "MainCluster",
				Relation:  "view",
				Subject: rts.NewSubjectSet(
					"OAuth2Client",
					"MainClusterAgent",
					"",
				),
			}

			ok, err := kcl.Check(context.Background(), &query)
			if err != nil {
				panic("Encountered error: " + err.Error())
			}
			Expect(ok).To(BeTrue())
			GinkgoWriter.Printf("OAuth2Client `MainClusterAgent` can view ObservabilityTenant `MainCluster`: %v\n", ok)
		})
		It("OAuth2Client with view permissions should not be able to edit an existing ObservabilityTenant", func() {
			query := rts.RelationTuple{
				Namespace: "ObservabilityTenant",
				Object:    "MainCluster",
				Relation:  "edit",
				Subject: rts.NewSubjectSet(
					"OAuth2Client",
					"MainClusterAgent",
					"",
				),
			}

			ok, err := kcl.Check(context.Background(), &query)
			if err != nil {
				panic("Encountered error: " + err.Error())
			}
			Expect(ok).To(BeFalse())
			GinkgoWriter.Printf("OAuth2Client `MainClusterAgent` can't edit ObservabilityTenant `MainCluster`: %v\n", ok)
		})
		It("Should be able to expand all viewers of an ObservabilityTenants", func() {

			query := rts.NewSubjectSet(
				"ObservabilityTenant",
				"MainCluster",
				"viewers",
			)

			respSubTree, err := kcl.Expand(context.Background(), query, 100)
			if err != nil {
				panic("Encountered error: " + err.Error())
			}

			var expected = []*rts.SubjectTree{
				{
					NodeType: rts.NodeType_NODE_TYPE_UNION,
					Subject: rts.NewSubjectSet(
						"Group",
						"MainCluster",
						"members",
					),
					Tuple: &rts.RelationTuple{
						Subject: rts.NewSubjectSet(
							"Group",
							"MainCluster",
							"members",
						),
					},
					Children: []*rts.SubjectTree{
						{
							NodeType: rts.NodeType_NODE_TYPE_LEAF,
							Subject: rts.NewSubjectSet(
								"User",
								"sam",
								"",
							),
							Tuple: &rts.RelationTuple{
								Subject: rts.NewSubjectSet(
									"User",
									"sam",
									"",
								),
							},
						},
					},
				},
				{
					NodeType: rts.NodeType_NODE_TYPE_LEAF,
					Subject: rts.NewSubjectSet(
						"OAuth2Client",
						"MainClusterAgent",
						"",
					),
					Tuple: &rts.RelationTuple{
						Subject: rts.NewSubjectSet(
							"OAuth2Client",
							"MainClusterAgent",
							"",
						),
					},
				},
			}

			client.SortSubjectTrees(respSubTree.Children)

			equal := cmp.Equal(respSubTree.Children, expected, protocmp.Transform())

			if !equal {
				diff := cmp.Diff(respSubTree.Children, expected, protocmp.Transform())
				GinkgoWriter.Printf("Diff: %s", diff)
			}

			Expect(equal).To(BeTrue())

		})
		It("Should be able to list all ObservabilityTenants a Group has permissions for", func() {

			query := rts.RelationQuery{
				Namespace: px.Ptr("ObservabilityTenant"),
				// Object:    nil,
				// Relation: px.Ptr("viewers"),
				Subject: rts.NewSubjectSet("Group", "MainCluster", "members"),
			}

			respTuples, err := kcl.QueryAllTuples(context.Background(), &query, 100)
			if err != nil {
				panic("Encountered error: " + err.Error())
			}
			// need to sort the tuples to make sure they are in the same order as the expected tuples
			client.SortRelationTuples(respTuples)

			var expected = []*rts.RelationTuple{
				{
					Namespace: "ObservabilityTenant",
					Object:    "FirstWorkloadCluster",
					Relation:  "viewers",
					Subject: rts.NewSubjectSet(
						"Group",
						"MainCluster",
						"members",
					),
				},
				{
					Namespace: "ObservabilityTenant",
					Object:    "MainCluster",
					Relation:  "viewers",
					Subject: rts.NewSubjectSet(
						"Group",
						"MainCluster",
						"members",
					),
				},
			}

			equal := cmp.Equal(respTuples, expected, protocmp.Transform())

			if !equal {
				diff := cmp.Diff(respTuples, expected, protocmp.Transform())
				GinkgoWriter.Printf("Diff: %s", diff)
			}

			Expect(equal).To(BeTrue())
		})
		It("Should be able to list all ObservabilityTenants a User has permissions for", func() {

			query := rts.RelationQuery{
				Namespace: px.Ptr("ObservabilityTenant"),
				// Object:    nil,
				// Relation: px.Ptr("viewers"),
				Subject: rts.NewSubjectSet("User", "nick", ""),
			}

			respTuples, err := kcl.QueryAllTuples(context.Background(), &query, 100)
			if err != nil {
				panic("Encountered error: " + err.Error())
			}
			// need to sort the tuples to make sure they are in the same order as the expected tuples
			client.SortRelationTuples(respTuples)

			var expected = []*rts.RelationTuple{
				{
					Namespace: "ObservabilityTenant",
					Object:    "MainCluster",
					Relation:  "editors",
					Subject: rts.NewSubjectSet(
						"User",
						"nick",
						"",
					),
				},
			}

			equal := cmp.Equal(respTuples, expected, protocmp.Transform())

			if !equal {
				diff := cmp.Diff(respTuples, expected, protocmp.Transform())
				GinkgoWriter.Printf("Diff: %s", diff)
			}

			Expect(equal).To(BeTrue())
		})
		It("Should be able to list all Groups a user is a member of", func() {

			query := rts.RelationQuery{
				Namespace: px.Ptr("Group"),
				// Object:    nil,
				Relation: px.Ptr("members"),
				Subject:  rts.NewSubjectSet("User", "aaron", ""),
			}

			respTuples, err := kcl.QueryAllTuples(context.Background(), &query, 100)
			if err != nil {
				panic("Encountered error: " + err.Error())
			}
			// need to sort the tuples to make sure they are in the same order as the expected tuples
			client.SortRelationTuples(respTuples)

			var expected = []*rts.RelationTuple{
				{
					Namespace: "Group",
					Object:    "AllUsers",
					Relation:  "members",
					Subject: rts.NewSubjectSet(
						"User",
						"aaron",
						"",
					),
				},
				{
					Namespace: "Group",
					Object:    "FirstWorkloadCluster",
					Relation:  "members",
					Subject: rts.NewSubjectSet(
						"User",
						"aaron",
						"",
					),
				},
			}

			equal := cmp.Equal(respTuples, expected, protocmp.Transform())

			if !equal {
				diff := cmp.Diff(respTuples, expected, protocmp.Transform())
				GinkgoWriter.Printf("Diff: %s", diff)
			}

			Expect(equal).To(BeTrue())
		})
	})
})
