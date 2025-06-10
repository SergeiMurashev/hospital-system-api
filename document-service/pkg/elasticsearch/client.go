package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/yourusername/hospital-system-api/document-service/internal/domain"
)

type Client interface {
	IndexDocument(doc *domain.Document) error
	SearchDocuments(query string) ([]*domain.Document, error)
	DeleteDocument(id uint) error
}

type client struct {
	es *elasticsearch.Client
}

func NewClient() (Client, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{os.Getenv("ELASTICSEARCH_URL")},
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create elasticsearch client: %v", err)
	}

	return &client{es: es}, nil
}

func (c *client) IndexDocument(doc *domain.Document) error {
	body, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("failed to marshal document: %v", err)
	}

	res, err := c.es.Index(
		"documents",
		bytes.NewReader(body),
		c.es.Index.WithDocumentID(fmt.Sprintf("%d", doc.ID)),
		c.es.Index.WithContext(context.Background()),
		c.es.Index.WithRefresh("true"),
	)
	if err != nil {
		return fmt.Errorf("failed to index document: %v", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error indexing document: %s", res.String())
	}

	return nil
}

func (c *client) SearchDocuments(query string) ([]*domain.Document, error) {
	var buf bytes.Buffer
	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  query,
				"fields": []string{"data", "room"},
			},
		},
	}

	if err := json.NewEncoder(&buf).Encode(searchQuery); err != nil {
		return nil, fmt.Errorf("failed to encode search query: %v", err)
	}

	res, err := c.es.Search(
		c.es.Search.WithContext(context.Background()),
		c.es.Search.WithIndex("documents"),
		c.es.Search.WithBody(&buf),
		c.es.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to search documents: %v", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error searching documents: %s", res.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode search response: %v", err)
	}

	hits := result["hits"].(map[string]interface{})["hits"].([]interface{})
	documents := make([]*domain.Document, 0, len(hits))

	for _, hit := range hits {
		source := hit.(map[string]interface{})["_source"].(map[string]interface{})
		doc := &domain.Document{
			ID:         uint(source["id"].(float64)),
			Date:       parseTime(source["date"].(string)),
			PatientID:  uint(source["patient_id"].(float64)),
			HospitalID: uint(source["hospital_id"].(float64)),
			DoctorID:   uint(source["doctor_id"].(float64)),
			Room:       source["room"].(string),
			Data:       source["data"].(string),
		}
		documents = append(documents, doc)
	}

	return documents, nil
}

func (c *client) DeleteDocument(id uint) error {
	res, err := c.es.Delete(
		"documents",
		fmt.Sprintf("%d", id),
		c.es.Delete.WithContext(context.Background()),
		c.es.Delete.WithRefresh("true"),
	)
	if err != nil {
		return fmt.Errorf("failed to delete document: %v", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error deleting document: %s", res.String())
	}

	return nil
}

func parseTime(timeStr string) time.Time {
	t, _ := time.Parse(time.RFC3339, timeStr)
	return t
}
