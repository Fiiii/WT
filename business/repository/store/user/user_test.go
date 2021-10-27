package user

import (
	"context"
	"testing"
	"time"

	"github.com/Fiiii/WT/foundation/logger"
	"github.com/google/go-cmp/cmp"
)

// Success and failure markers.
const (
	Success = "\u2713"
	Failed  = "\u2717"
)

func TestUser(t *testing.T) {
	log, err := logger.New("TEST")
	if err != nil {
		t.Fatalf("creating logger error")

	}
	store := NewStore(log, nil)

	t.Log("Given the need to work with User records.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen handling single User.", testID)
		{
			ctx := context.Background()
			now := time.Date(2021, time.October, 1, 0, 0, 0, 0, time.UTC)

			nu := NewUser{
				Name:  "Fii",
				Email: "fii@fii.com",
				Roles: []string{"Admin"},
			}

			user, err := store.Create(ctx, nu, now)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to create user : %s.", Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to create user.", Success, testID)

			saved, err := store.Query(ctx, user.ID)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve user by ID: %s.", Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve user by ID.", Success, testID)

			if diff := cmp.Diff(user, saved); diff != "" {
				t.Fatalf("\t%s\tTest %d:\tShould get back the same user. Diff:\n%s", Failed, testID, diff)
			}
			t.Logf("\t%s\tTest %d:\tShould get back the same user.", Success, testID)

			upd := UpdateUser{
				Name:  StringPointer("Updated Fii"),
				Email: StringPointer("updated@fii.com"),
			}

			if err := store.Update(ctx, user.ID, upd, now); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to update user : %s.", Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to update user.", Success, testID)

			if saved.Name != *upd.Name {
				t.Errorf("\t%s\tTest %d:\tShould be able to see updates to Name.", Failed, testID)
				t.Logf("\t\tTest %d:\tGot: %v", testID, saved.Name)
				t.Logf("\t\tTest %d:\tExp: %v", testID, *upd.Name)
			} else {
				t.Logf("\t%s\tTest %d:\tShould be able to see updates to Name.", Success, testID)
			}

			if saved.Email != *upd.Email {
				t.Errorf("\t%s\tTest %d:\tShould be able to see updates to Email.", Failed, testID)
				t.Logf("\t\tTest %d:\tGot: %v", testID, saved.Email)
				t.Logf("\t\tTest %d:\tExp: %v", testID, *upd.Email)
			} else {
				t.Logf("\t%s\tTest %d:\tShould be able to see updates to Email.", Success, testID)
			}

			if err := store.Delete(ctx, user.ID); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to delete user : %s.", Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to delete user.", Success, testID)
		}
	}
}

func StringPointer(s string) *string {
	return &s
}
