package meilisearch

import (
	"context"
	"fmt"
	"net/http"

	"github.com/coscms/webfront/library/search"
	"github.com/meilisearch/meilisearch-go"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

// https://docs.meilisearch.com/learn/getting_started/quick_start.html#download-and-launch

func New(cfg search.Config) (search.Searcher, error) {
	client := meilisearch.New(
		cfg.Host,
		meilisearch.WithAPIKey(cfg.Password),
		meilisearch.WithCustomClient(&http.Client{Timeout: cfg.Timeout}),
	)
	return &MeiliSearch{
		ServiceManager: client,
	}, nil
}

type MeiliSearch struct {
	meilisearch.ServiceManager
}

func (m *MeiliSearch) getIndex(index string) meilisearch.IndexManager {
	return m.ServiceManager.Index(index)
}

func (m *MeiliSearch) Add(ctx context.Context, index string, primaryKey *string, docs ...interface{}) error {
	indexM := m.getIndex(index)
	opts := meilisearch.DocumentOptions{
		PrimaryKey:         primaryKey,
		SkipCreation:       false,
		TaskCustomMetadata: "",
	}
	t, err := indexM.AddDocuments(ctx, docs, &opts)
	if err != nil {
		return err
	}
	if !isOk(t.Status) {
		err = fmt.Errorf(`%+v`, *t)
	}
	return err
}

func (m *MeiliSearch) Update(ctx context.Context, index string, primaryKey *string, docs ...interface{}) error {
	indexM := m.getIndex(index)
	opts := meilisearch.DocumentOptions{
		PrimaryKey:         primaryKey,
		SkipCreation:       false,
		TaskCustomMetadata: "",
	}
	t, err := indexM.UpdateDocuments(ctx, docs, &opts)
	if err != nil {
		return err
	}
	if !isOk(t.Status) {
		err = fmt.Errorf(`%+v`, *t)
	}
	return err
}

func (m *MeiliSearch) GetTask(ctx context.Context, index string, taskUID int64) error {
	indexM := m.getIndex(index)
	task, err := indexM.GetTask(ctx, taskUID)
	com.Dump(task)
	return err
}

func (m *MeiliSearch) Delete(ctx context.Context, index string, ids ...string) error {
	indexM := m.getIndex(index)
	opts := meilisearch.DocumentOptions{
		//PrimaryKey:         primaryKey,
		SkipCreation:       false,
		TaskCustomMetadata: "",
	}
	_, err := indexM.DeleteDocuments(ctx, ids, &opts)
	return err
}

func (m *MeiliSearch) Flush() error {
	return nil
}

func (m *MeiliSearch) DeleteIndex(ctx context.Context, index string) error {
	_, err := m.ServiceManager.DeleteIndex(ctx, index)
	return err
}

func (m *MeiliSearch) InitIndex(ctx context.Context, cfg *search.IndexConfig) error {
	// m.DeleteIndex(cfg.Index)
	if _, err := m.getIndex(cfg.Index).FetchInfo(ctx); err != nil {
		m.ServiceManager.CreateIndex(ctx, &meilisearch.IndexConfig{
			Uid:        cfg.Index,
			PrimaryKey: cfg.PrimaryKey,
		})
		searchableAttributes := cfg.SearchableAttributes
		sortableAttributes := cfg.SortableAttributes
		filterableAttributes := cfg.FilterableAttributes
		index := m.getIndex(cfg.Index)
		if len(searchableAttributes) > 0 {
			index.UpdateSearchableAttributes(ctx, &searchableAttributes)
		}
		if len(sortableAttributes) > 0 {
			index.UpdateSortableAttributes(ctx, &sortableAttributes)
		}
		if len(filterableAttributes) > 0 {
			index.UpdateFilterableAttributes(ctx, &filterableAttributes)
		}
	}
	return nil
}

func (m *MeiliSearch) Search(ctx context.Context, index string, keywords string, options *search.SearchRequest) (int64, []echo.H, error) {
	searchRes, err := m.getIndex(index).Search(
		ctx,
		keywords,
		&meilisearch.SearchRequest{
			Offset: options.Offset,
			Limit:  options.Limit,
			Filter: options.Filter,
			Sort:   options.SortFields,
		},
	)
	if err != nil {
		return 0, nil, err
	}

	var rows []echo.H
	err = searchRes.Hits.Decode(&rows)
	return searchRes.EstimatedTotalHits, rows, err
}
