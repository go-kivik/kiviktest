package db

import (
	"context"

	"github.com/blang/semver"
	"github.com/flimzy/diff"

	"github.com/flimzy/kivik"
	"github.com/go-kivik/kiviktest/kt"
)

func init() {
	kt.Register("GetIndexes", getIndexes)
}

func getIndexes(ctx *kt.Context) {
	ctx.RunAdmin(func(ctx *kt.Context) {
		ctx.Parallel()
		roGetIndexesTests(ctx, ctx.Admin)
	})
	ctx.RunNoAuth(func(ctx *kt.Context) {
		ctx.Parallel()
		roGetIndexesTests(ctx, ctx.NoAuth)
	})
	ctx.RunRW(func(ctx *kt.Context) {
		ctx.RunAdmin(func(ctx *kt.Context) {
			ctx.Parallel()
			rwGetIndexesTests(ctx, ctx.Admin)
		})
		ctx.RunNoAuth(func(ctx *kt.Context) {
			ctx.Parallel()
			rwGetIndexesTests(ctx, ctx.NoAuth)
		})
	})
}

func roGetIndexesTests(ctx *kt.Context, client *kivik.Client) {
	databases := ctx.MustStringSlice("databases")
	for _, dbname := range databases {
		func(dbname string) {
			ctx.Run(dbname, func(ctx *kt.Context) {
				ctx.Parallel()
				testGetIndexes(ctx, ctx.Admin, dbname, ctx.Interface("indexes"))
			})
		}(dbname)
	}
}

func rwGetIndexesTests(ctx *kt.Context, client *kivik.Client) {
	dbname := ctx.TestDB()
	defer ctx.Admin.DestroyDB(context.Background(), dbname, ctx.Options("db")) // nolint: errcheck
	dba, err := ctx.Admin.DB(context.Background(), dbname, ctx.Options("db"))
	if err != nil {
		ctx.Fatalf("Failed to open db as admin: %s", err)
	}
	if err = dba.CreateIndex(context.Background(), "foo", "bar", `{"fields":["foo"]}`); err != nil {
		ctx.Fatalf("Failed to create index: %s", err)
	}
	indexDef := map[string]interface{}{
		"fields": []map[string]string{
			{"foo": "asc"},
		},
	}
	if partialFilterSupported(ctx) {
		indexDef["partial_filter_selector"] = map[string]interface{}{}
	}
	testGetIndexes(ctx, client, dbname, []kivik.Index{
		kt.AllDocsIndex,
		{
			DesignDoc:  "_design/foo",
			Name:       "bar",
			Type:       "json",
			Definition: indexDef,
		},
	})
}

func partialFilterSupported(ctx *kt.Context) bool {
	meta, err := ctx.Admin.Version(context.Background())
	if err != nil {
		ctx.Fatalf("Failed to detect db version: %s", err)
	}
	if meta.Vendor == "PouchDB" {
		return true
	}
	if meta.Vendor != "The Apache Software Foundation" {
		return false
	}
	ver, err := semver.New(meta.Version)
	if err != nil {
		ctx.Fatalf("Invalid version '%s' reported by serer: %s", meta.Version, err)
	}
	filteredIndexVersion := semver.MustParse("2.1.1")
	if ver.GE(filteredIndexVersion) {
		return true
	}
	return false
}

func testGetIndexes(ctx *kt.Context, client *kivik.Client, dbname string, expected interface{}) {
	db, err := client.DB(context.Background(), dbname, ctx.Options("db"))
	if err != nil {
		ctx.Fatalf("Failed to open db: %s", err)
	}
	indexes, err := db.GetIndexes(context.Background())
	if !ctx.IsExpectedSuccess(err) {
		return
	}
	if d := diff.AsJSON(expected, indexes); d != nil {
		ctx.Errorf("Indexes differ from expectation:\n%s\n", d)
	}
}
