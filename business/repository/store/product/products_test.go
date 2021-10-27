package product

import (
	"context"
	"github.com/Fiiii/WT/foundation/logger"
	"github.com/google/go-cmp/cmp"
	"testing"
	"time"
)

// Success and failure markers.
const (
	Success = "\u2713"
	Failed  = "\u2717"
)

func TestProduct(t *testing.T) {
	log, err := logger.New("TEST")
	if err != nil {
		t.Fatalf("creating logger error")

	}
	store := NewStore(log, nil)

	t.Log("Given the need to work with Product records.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen handling single Product.", testID)
		{
			ctx := context.Background()
			now := time.Date(2021, time.October, 1, 0, 0, 0, 0, time.UTC)

			np := NewProduct{
				Name:     "test product",
				Cost:     1,
				Quantity: 2,
			}

			product, err := store.Create(ctx, np, now)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to create product : %s.", Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to create product.", Success, testID)

			saved, err := store.QueryByID(ctx, product.ID)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve product by ID: %s.", Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve product by ID.", Success, testID)

			if diff := cmp.Diff(product, saved); diff != "" {
				t.Fatalf("\t%s\tTest %d:\tShould get back the same product. Diff:\n%s", Failed, testID, diff)
			}
			t.Logf("\t%s\tTest %d:\tShould get back the same product.", Success, testID)
		}
	}
}
