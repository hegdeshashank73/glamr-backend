package entities

import (
	"encoding/json"
	"fmt"

	"github.com/hegdeshashank73/glamr-backend/common"
	"github.com/hegdeshashank73/glamr-backend/errors"
)

type AssetUploadHandlerReq struct {
	EntityID    common.Snowflake `json:"entity_id"`
	EntityType  AssetsEntityType `json:"entity_type"`
	ContentType string           `json:"content_type"`
	FileExt     string           `json:"ext"`
}

func (r *AssetUploadHandlerReq) Validate() errors.GlamrError {

	if r.EntityType == 0 {
		return errors.GlamrErrorGeneralBadRequest("missing entity type")
	}

	switch r.ContentType {
	case "image/jpeg":
		r.FileExt = "jpg"
	case "image/png":
		r.FileExt = "png"
	case "application/pdf":
		r.FileExt = "pdf"
	case "video/mp4":
		r.FileExt = "mp4"
	default:
		return errors.GlamrErrorGeneralBadRequest("unsupported content type")
	}

	return nil
}

type AssetUploadHandlerRes struct {
	UploadURL string `json:"upload_url"`
	AccessURL string `json:"access_url"`
	Key       string `json:"key"`
}

type AssetsEntityType int8

const AssetsEntityType_CLOTHING AssetsEntityType = 1

var AssetsEntityTypesStringMap = map[AssetsEntityType]string{
	AssetsEntityType_CLOTHING: "clothing",
}
var StringAssetsEntityTypesMap = map[string]AssetsEntityType{
	"clothing": AssetsEntityType_CLOTHING,
}

func (ci AssetsEntityType) MarshalJSON() ([]byte, error) {
	return json.Marshal(AssetsEntityTypesStringMap[ci])
}

func (ci *AssetsEntityType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	if v, ok := StringAssetsEntityTypesMap[s]; !ok {
		return fmt.Errorf("invalid asset entity type")
	} else {
		*ci = v
	}
	return nil
}
