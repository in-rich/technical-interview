package dao_test

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"strings"
)

// CleanFirestore deletes all test data for the local firestore instance.
func CleanFirestore(client *firestore.Client) error {
	collections, err := client.Collections(context.Background()).GetAll()
	if err != nil {
		return err
	}

	for _, collection := range collections {
		// Only clear test collections.
		if !strings.HasPrefix(collection.ID, "test-") {
			continue
		}

		documents, err := collection.Documents(context.Background()).GetAll()
		if err != nil {
			return fmt.Errorf("failed to get documents from collection %s: %w", collection.ID, err)
		}

		for _, document := range documents {
			if _, err := document.Ref.Delete(context.Background()); err != nil {
				return fmt.Errorf("failed to delete document %s: %w", document.Ref.ID, err)
			}
		}
	}

	return nil
}
