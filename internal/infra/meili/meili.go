package meili

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"manga/internal/domain/dtos"
	"manga/internal/domain/params"
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
func (t *Manga) Search(ctx context.Context, args params.SearchParams) (any,any, error) {
	// ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Task.Search")
	// defer span.End()

	// if args.IsZero() {
	// 	return internal.SearchResults{}, nil
	// }

	should := make(map[string]interface{})

	if args.Genre != nil {
		should["genres.name"] = *args.Genre
		
	}

	if args.Status != nil {
		should["status"] = *args.Status
		
	}

	if args.Type != nil {
		should["type"] = *args.Status
	}

	idx := t.client.Index(t.index)
    fmt.Print(should)
	resp, err := idx.SearchWithContext(ctx,*args.Query,&meilisearch.SearchRequest{
		Offset: 0,
		Limit: 20,
		Filter: "genre.name IN [Action,Fantasy] AND type IN [Manga]",
	})
	if err != nil {
		return nil,nil, pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "SearchRequest.Do")
	}
	
	// //nolint: tagliatelle
	// var hits struct {
	// 	Hits struct {
	// 		Total struct {
	// 			Value int64 `json:"value"`
	// 		} `json:"total"`
	// 		Hits []struct {
	// 			Source indexedTask `json:"_source"`
	// 		} `json:"hits"`
	// 	} `json:"hits"`
	// }

	// if err := json.NewDecoder(resp.Body).Decode(&hits); err != nil {
	// 	return internal.SearchResults{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "json.NewDecoder.Decode")
	// }

	// res := make([]internal.Task, len(hits.Hits.Hits))

	// for i, hit := range hits.Hits.Hits {
	// 	res[i].ID = hit.Source.ID
	// 	res[i].Description = hit.Source.Description
	// 	res[i].Priority = hit.Source.Priority
	// 	res[i].Dates.Due = time.Unix(0, hit.Source.DateDue).UTC()
	// 	res[i].Dates.Start = time.Unix(0, hit.Source.DateStart).UTC()
	// }

	return resp.Hits, nil,nil
}
