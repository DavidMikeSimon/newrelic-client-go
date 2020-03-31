package nerdstorage

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/mitchellh/mapstructure"
	"github.com/newrelic/newrelic-client-go/pkg/config"
)

type AccountScopedDoc struct {
	MyField string
}

func Example_accountScope() {
	// Initialize the client configuration.  A Personal API key is required to
	// communicate with the backend API.
	cfg := config.New()
	cfg.PersonalAPIKey = os.Getenv("NEW_RELIC_API_KEY")

	// Initialize the client.
	client := New(cfg)

	accountID := 12345678
	packageID := "ecaeb28c-7b3f-4932-9e33-7385980efa1c"

	// Write a NerdStorage document with account scope.
	writeDocumentInput := WriteDocumentInput{
		PackageID:  packageID,
		Collection: "myCol",
		DocumentID: "myDoc",
		Document: AccountScopedDoc{
			MyField: "myValue",
		},
	}

	_, err := client.WriteDocumentWithAccountScope(accountID, writeDocumentInput)
	if err != nil {
		log.Fatal("error writing NerdStorage document:", err)
	}

	// Write a second NerdStorage document to the same collection with account scope.
	writeAlternateDocumentInput := writeDocumentInput
	writeAlternateDocumentInput.DocumentID = "myOtherDoc"

	_, err = client.WriteDocumentWithAccountScope(accountID, writeAlternateDocumentInput)
	if err != nil {
		log.Fatal("error writing NerdStorage document:", err)
	}

	// Get a NerdStorage collection with account scope.
	getCollectionInput := GetCollectionInput{
		PackageID:  packageID,
		Collection: "myCol",
	}

	collection, err := client.GetCollectionWithAccountScope(accountID, getCollectionInput)
	if err != nil {
		log.Fatal("error retrieving NerdStorage collection:", err)
	}

	fmt.Printf("Collection length: %v\n", len(collection))

	// Get a NerdStorage document with account scope.
	getDocumentInput := GetDocumentInput{
		PackageID:  packageID,
		Collection: "myCol",
		DocumentID: "myDoc",
	}

	rawDoc, err := client.GetDocumentWithAccountScope(accountID, getDocumentInput)
	if err != nil {
		log.Fatal("error retrieving NerdStorage document:", err)
	}

	// Convert the document to a struct.
	var myDoc AccountScopedDoc

	// Method 1:
	marshalled, err := json.Marshal(rawDoc)
	if err != nil {
		log.Fatal("error marshalling NerdStorage document to json:", err)
	}

	err = json.Unmarshal(marshalled, &myDoc)
	if err != nil {
		log.Fatal("error unmarshalling NerdStorage document to struct:", err)
	}

	fmt.Printf("Document: %v\n", myDoc)

	// Method 2:
	err = mapstructure.Decode(rawDoc, &myDoc)
	if err != nil {
		log.Fatal("error converting NerdStorage document to struct:", err)
	}

	fmt.Printf("Document: %v\n", myDoc)

	// Delete a NerdStorage document with account scope.
	deleteDocumentInput := DeleteDocumentInput{
		PackageID:  packageID,
		Collection: "myCol",
		DocumentID: "myDoc",
	}

	ok, err := client.DeleteDocumentWithAccountScope(accountID, deleteDocumentInput)

	if !ok || err != nil {
		log.Fatal("error deleting NerdStorage document:", err)
	}

	// Delete a NerdStorage collection with account scope.
	deleteCollectionInput := DeleteCollectionInput{
		PackageID:  packageID,
		Collection: "myCol",
	}

	ok, err = client.DeleteCollectionWithAccountScope(accountID, deleteCollectionInput)
	if err != nil {
		log.Fatal("error deleting NerdStorage collection:", err)
	}
}
