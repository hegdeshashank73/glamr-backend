package services

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/hegdeshashank73/glamr-backend/entities"
	"github.com/hegdeshashank73/glamr-backend/errors"
	"github.com/hegdeshashank73/glamr-backend/repository"
	"github.com/hegdeshashank73/glamr-backend/utils"
	"github.com/hegdeshashank73/glamr-backend/vendors"
	g "github.com/serpapi/google-search-results-golang"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func GetSerpAPISearchResults(req entities.SearchOptionsReq) (entities.SearchOptionsRes, errors.GlamrError) {
	st := time.Now()
	defer utils.LogTimeTaken("services.GetSerpAPISearchResults", st)

	res := entities.SearchOptionsRes{}
	tx, err := vendors.DBMono.Begin()
	if err != nil {
		logrus.Error(err)
		return res, errors.GlamrErrorDatabaseIssue()
	}
	parameter := map[string]string{
		"engine":  "google_lens",
		"url":     fmt.Sprintf("%s/%s", viper.GetString("USER_BASE_URL"), req.S3Key),
		"country": req.CountryCode,
	}

	search := g.NewGoogleSearch(parameter, viper.GetString("SERP_API"))
	var results map[string]interface{}
	results, err = search.GetJSON()
	if err != nil {
		logrus.Error("Failed to fetch data from SERP API", err)
		return res, errors.GlamrErrorInternalServerError()
	}
	ret, derr := repository.CreatePersonSearch(tx, entities.SearchOptionsArg(req), results)
	if derr != nil {
		tx.Rollback()
		return res, derr
	}

	var serpAPIObject entities.SerpApiObject
	if err := json.Unmarshal(ret.APIResponse, &serpAPIObject); err != nil {
		logrus.Error("failed to unmarshal json,", err)
		return res, errors.GlamrErrorInternalServerError()
	}
	if serpAPIObject.SearchMetadata.Status != "Success" {
		return res, errors.GlamrErrorGeneralBadRequest("Please retry the search")
	}
	var searchOptions []entities.SearchOptions
	for _, visualMatch := range serpAPIObject.VisualMatches {
		if visualMatch.Price.Currency != "$" || visualMatch.Price.ExtractedPrice == 0.0 {
			continue
		}
		option := entities.SearchOptions{
			Title:      visualMatch.Title,
			Link:       visualMatch.Link,
			Source:     visualMatch.Source,
			SourceIcon: visualMatch.SourceIcon,
			InStock:    visualMatch.InStock,
			Image:      visualMatch.Thumbnail,
			Currency:   visualMatch.Price.Currency,
			Price:      int(visualMatch.Price.ExtractedPrice),
		}
		searchOptions = append(searchOptions, option)
	}

	sort.Slice(searchOptions, func(i, j int) bool {
		return searchOptions[i].Price < searchOptions[j].Price
	})
	res.SearchOptions = searchOptions

	derr = repository.CreateSearchOptions(tx, entities.CreateSearchOptionsArg{
		ID:            ret.ID,
		SearchOptions: searchOptions,
	})
	if derr != nil {
		tx.Rollback()
		return res, derr
	}

	err = tx.Commit()
	if err != nil {
		logrus.Error(err)
		return res, errors.GlamrErrorDatabaseIssue()
	}

	return res, nil
}

func GetSearchHistory(person *entities.Person) (entities.GetSearchHistoryRes, errors.GlamrError) {
	st := time.Now()
	defer utils.LogTimeTaken("services.GetSearchHistory", st)

	res := entities.GetSearchHistoryRes{}
	tx, err := vendors.DBMono.Begin()
	if err != nil {
		logrus.Error(err)
		return res, errors.GlamrErrorDatabaseIssue()
	}

	ret, derr := repository.GetSearchHistory(tx, *person)
	if derr != nil {
		logrus.Error(derr)
		tx.Rollback()
		return res, derr
	}

	err = tx.Commit()
	if err != nil {
		logrus.Error(err)
		return res, errors.GlamrErrorDatabaseIssue()
	}
	res.SearchHistory = ret.SearchHistory
	return res, nil
}

func GetSearchHistoryOptions(person *entities.Person, req entities.GetSearchHistoryOptionsReq) (entities.GetSearchHistoryOptionsRes, errors.GlamrError) {
	st := time.Now()
	defer utils.LogTimeTaken("services.GetSearchHistoryOptions", st)

	res := entities.GetSearchHistoryOptionsRes{}
	tx, err := vendors.DBMono.Begin()
	if err != nil {
		logrus.Error(err)
		return res, errors.GlamrErrorDatabaseIssue()
	}

	ret, derr := repository.GetSearchHistoryOptions(tx, entities.GetSearchHistoryOptionsArg(req))
	if derr != nil {
		logrus.Error(derr)
		tx.Rollback()
		return res, derr
	}

	err = tx.Commit()
	if err != nil {
		logrus.Error(err)
		return res, errors.GlamrErrorDatabaseIssue()
	}
	res.SearchOptions = ret.SearchOptions
	return res, nil
}
