package meili

import (
	"bytes"
	"context"
	"encoding/json"
	"manga/internal/domain/dtos"
	"manga/pkg"

	"github.com/meilisearch/meilisearch-go"
)

type Manga struct {
	client meilisearch.ServiceManager
	index  string
}


// NewTask instantiates the Task repository.
func NewManga(client meilisearch.ServiceManager) *Manga {
	return &Manga{
		client: client,
		index:  "manga",
	}
}

// Index creates or updates a task in an index.
func (t *Manga) Index(ctx context.Context, manga dtos.IndexedManga) error {

	

	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(manga); err != nil {
		return pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "json.NewEncoder.Encode")
	}
	q := t.client.Index(t.index)
	_, err := q.AddDocumentsNdjsonFromReaderWithContext(ctx,&buf)
	if err != nil {
		return pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "IndexRequest.Do")
	}

	return nil
}

// Delete removes a task from the index.
func (t *Manga) Delete(ctx context.Context, id string) error {
	q := t.client.Index(t.index)

    
	_, err := q.DeleteDocumentsByFilterWithContext(ctx,id)
	if err != nil {
		return pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "DeleteRequest.Do")
	}
	
	return nil
}

// Search returns tasks matching a query.
// nolint: funlen, cyclop
// func (t *Task) Search(ctx context.Context, args internal.SearchParams) (internal.SearchResults, error) {
// 	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Task.Search")
// 	defer span.End()

// 	if args.IsZero() {
// 		return internal.SearchResults{}, nil
// 	}

// 	should := make([]interface{}, 0, 3)

// 	if args.Description != nil {
// 		should = append(should, map[string]interface{}{
// 			"match": map[string]interface{}{
// 				"description": *args.Description,
// 			},
// 		})
// 	}

// 	if args.Priority != nil {
// 		should = append(should, map[string]interface{}{
// 			"match": map[string]interface{}{
// 				"priority": *args.Priority,
// 			},
// 		})
// 	}

// 	if args.IsDone != nil {
// 		should = append(should, map[string]interface{}{
// 			"match": map[string]interface{}{
// 				"is_done": *args.IsDone,
// 			},
// 		})
// 	}

// 	var query map[string]interface{}

// 	if len(should) > 1 {
// 		query = map[string]interface{}{
// 			"query": map[string]interface{}{
// 				"bool": map[string]interface{}{
// 					"should": should,
// 				},
// 			},
// 		}
// 	} else {
// 		query = map[string]interface{}{
// 			"query": should[0],
// 		}
// 	}

// 	query["sort"] = []interface{}{
// 		"_score",
// 		map[string]interface{}{"id": "asc"},
// 	}

// 	query["from"] = args.From
// 	query["size"] = args.Size

// 	var buf bytes.Buffer
// meilisearch
// 	if err := json.NewEncoder(&buf).Encode(query); err != nil {
// 		return internal.SearchResults{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "json.NewEncoder.Encode")
// 	}

// 	req := esv7api.SearchRequest{
// 		Index: []string{t.index},
// 		Body:  &buf,
// 	}

// 	resp, err := req.Do(ctx, t.client)
// 	if err != nil {
// 		return internal.SearchResults{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "SearchRequest.Do")
// 	}
// 	defer resp.Body.Close()

// 	if resp.IsError() {
// 		return internal.SearchResults{}, internal.NewErrorf(internal.ErrorCodeUnknown, "SearchRequest.Do %d", resp.StatusCode)
// 	}

// 	//nolint: tagliatelle
// 	var hits struct {
// 		Hits struct {
// 			Total struct {
// 				Value int64 `json:"value"`
// 			} `json:"total"`
// 			Hits []struct {
// 				Source indexedTask `json:"_source"`
// 			} `json:"hits"`
// 		} `json:"hits"`
// 	}

// 	if err := json.NewDecoder(resp.Body).Decode(&hits); err != nil {
// 		return internal.SearchResults{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "json.NewDecoder.Decode")
// 	}

// 	res := make([]internal.Task, len(hits.Hits.Hits))

// 	for i, hit := range hits.Hits.Hits {
// 		res[i].ID = hit.Source.ID
// 		res[i].Description = hit.Source.Description
// 		res[i].Priority = hit.Source.Priority
// 		res[i].Dates.Due = time.Unix(0, hit.Source.DateDue).UTC()
// 		res[i].Dates.Start = time.Unix(0, hit.Source.DateStart).UTC()
// 	}

// 	return internal.SearchResults{
// 		Tasks: res,
// 		Total: hits.Hits.Total.Value,
// 	}, nil
// }
