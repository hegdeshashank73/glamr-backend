package entities

import (
	"encoding/json"
	"fmt"

	"github.com/hegdeshashank73/glamr-backend/common"
	"github.com/spf13/viper"
)

type CountryCurrency string

const CountryCurrency_IN CountryCurrency = "IN"

var CountryCurrencyStringMap = map[CountryCurrency]string{
	CountryCurrency_IN: "₹",
}
var StringAssetsCountryCurrencyMap = map[string]CountryCurrency{
	"₹": CountryCurrency_IN,
}

type SearchOptionsReq struct {
	S3Key       string `json:"s3_key"`
	CountryCode string `json:"country_code"`
}

type SearchOptionsArg SearchOptionsReq

type SearchOptionsRes struct {
	ID            common.Snowflake `json:"id"`
	SearchOptions []SearchOptions  `json:"search_options"`
}

type SearchOptions struct {
	Title      string `json:"title"`
	Link       string `json:"link"`
	Source     string `json:"source"`
	SourceIcon string `json:"source_icon"`
	InStock    bool   `json:"in_stock"`
	Price      int    `json:"price"`
	Image      string `json:"image"`
	Currency   string `json:"currency"`
}

type Price struct {
	Currency       string `json:"currency"`
	ExtractedPrice string `json:"extracted_price"`
	Price          int    `json:"total_price"`
}

type SerpApiObject struct {
	SearchMetadata SearchMetadata      `json:"search_metadata"`
	VisualMatches  []SearchVisualMatch `json:"visual_matches"`
}

type SearchMetadata struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type SearchVisualMatch struct {
	Title      string           `json:"title"`
	Link       string           `json:"link"`
	Source     string           `json:"source"`
	SourceIcon string           `json:"source_icon"`
	InStock    bool             `json:"in_stock"`
	Thumbnail  string           `json:"thumbnail"`
	Price      VisualMatchPrice `json:"price"`
}

type VisualMatchPrice struct {
	Currency       string  `json:"currency"`
	ExtractedPrice float64 `json:"extracted_value"`
	Value          string  `json:"value"`
}

type CreatePersonSearchRet struct {
	ID          common.Snowflake `json:"id"`
	APIResponse []byte           `json:"api_response"`
}

type CreateSearchOptionsArg SearchOptionsRes

type GetSearchHistoryRet struct {
	SearchHistory []SearchHistory `json:"search_history"`
}

type SearchHistory struct {
	ID        common.Snowflake `json:"id"`
	CreatedAt int64            `json:"created_at"`
	S3Key     string           `json:"-"`
}

func (d SearchHistory) MarshalJSON() ([]byte, error) {
	type SearchHistoryAlias SearchHistory
	return json.Marshal(&struct {
		SearchHistoryAlias
		Image string `json:"image"`
	}{
		SearchHistoryAlias: (SearchHistoryAlias)(d),
		Image:              fmt.Sprintf("%s/%s", viper.GetString("USER_BASE_URL"), d.S3Key),
	})
}

type GetSearchHistoryRes GetSearchHistoryRet

type GetSearchHistoryOptionsReq struct {
	SearchID common.Snowflake `json:"search_id"`
}

type GetSearchHistoryOptionsArg GetSearchHistoryOptionsReq

type GetSearchHistoryOptionsRes struct {
	SearchOptions []SearchOptions `json:"search_options"`
}
