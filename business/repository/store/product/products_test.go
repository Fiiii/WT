package product

import (
	"github.com/Fiiii/WT/foundation/logger"
	"testing"
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

		}
	}
}
