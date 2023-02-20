package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	px "github.com/ory/x/pointerx"
	ketocl "github.com/pluralsh/oauth-playground/keto/client"

	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
)

func main() {
	//conn, err := grpc.Dial("127.0.0.1:4467", grpc.WithTransportCredentials(insecure.NewCredentials()))
	err := godotenv.Load("./.env") // 👈 load .env file
	if err != nil {
		log.Fatal(err)
	}
	conndetails := ketocl.NewConnectionDetailsFromEnv()
	kcl, err := ketocl.NewGrpcClient(context.Background(), conndetails)
	if err != nil {
		panic("Encountered error: " + err.Error())
	}

	//directories:/photos#owner@maureen
	//files:/photos/beach.jpg#owner@maureen
	//files:/photos/mountains.jpg#owner@laura
	//directories:/photos#access@laura
	//directories:/photos#access@(directories:/photos#owner)
	//files:/photos/beach.jpg#access@(files:/photos/beach.jpg#owner)
	//files:/photos/beach.jpg#access@(directories:/photos#access)
	//files:/photos/mountains.jpg#access@(files:/photos/mountains.jpg#owner)
	//files:/photos/mountains.jpg#access@(directories:/photos#access)

	tuples := []*rts.RelationTuple{
		{
			Namespace: "Organization",
			Object:    "main",
			Relation:  "admin",
			Subject: rts.NewSubjectSet(
				"Admin",
				"david",
				"",
			),
		},
	}

	kcl.CreateTuples(context.Background(), tuples)

	// kcl.DeleteTuple(context.Background(), tuples[0])
	// kcl.DeleteTuple(context.Background(), tuples[0])

	relationQuery := rts.RelationQuery{
		Namespace: px.Ptr("Organization"),
		Object:    px.Ptr("main"),
		Relation:  px.Ptr("admin"),
	}

	existingTuples, err := kcl.QueryAllTuples(context.Background(), &relationQuery, 100)

	if err != nil {
		panic("Encountered error: " + err.Error())
	}

	if existingTuples != nil {
		log.Println("getTuples is not nil")
		ketocl.PrintTableFromRelationTuples(existingTuples, os.Stdout)
	} else {
		log.Println("getTuples is nil")
	}

	// should be subject set rewrite
	// owners have access
	//for _, o := range []struct{ n, o string }{
	//	{"files", "/photos/beach.jpg"},
	//	{"files", "/photos/mountains.jpg"},
	//	{"directories", "/photos"},
	//} {
	//	tuples = append(tuples, &rts.RelationTuple{
	//		Namespace: o.n,
	//		Object:    o.o,
	//		Relation:  "access",
	//		Subject:   rts.NewSubjectSet(o.n, o.o, "owner"),
	//	})
	//}
	// should be subject set rewrite
	// access on parent means access on child
	//for _, obj := range []string{"/photos/beach.jpg", "/photos/mountains.jpg"} {
	//	tuples = append(tuples, &rts.RelationTuple{
	//		Namespace: "files",
	//		Object:    obj,
	//		Relation:  "access",
	//		Subject:   rts.NewSubjectSet("directories", "/photos", "access"),
	//	})
	//}

}
