package services

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hegdeshashank73/glamr-backend/common"
	"github.com/hegdeshashank73/glamr-backend/entities"
	"github.com/hegdeshashank73/glamr-backend/errors"
	"github.com/hegdeshashank73/glamr-backend/utils"
	"github.com/hegdeshashank73/glamr-backend/vendors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func GeneratePresignedURL(person entities.Person, req entities.AssetUploadHandlerReq) (entities.AssetUploadHandlerRes, errors.GlamrError) {
	st := time.Now()
	defer utils.LogTimeTaken("services.GeneratePresignedURL", st)

	res := entities.AssetUploadHandlerRes{}

	var key = ""
	if req.EntityType == entities.AssetsEntityType_CLOTHING {
		key = fmt.Sprintf("%s/clothing-uploads/%s.%s", person.Id, common.GenerateSnowflake(), req.FileExt)
		res.AccessURL = fmt.Sprintf("%s/%s", viper.GetString("USER_BASE_URL"), key)
		res.Key = key

	} else {
		return res, errors.GlamrErrorGeneralBadRequest("unsupported entity type")
	}
	sreq, _ := vendors.S3Client.PutObjectRequest(&s3.PutObjectInput{
		Bucket:      aws.String(viper.GetString("USER_ASSETS_BUCKET")),
		Key:         aws.String(key),
		ContentType: aws.String(req.ContentType),
	})
	str, err := sreq.Presign(time.Minute * 15)
	if err != nil {
		logrus.Error(err)
		return res, errors.GlamrErrorGeneralBadRequest("unable to upload at this time")
	}
	res.UploadURL = str
	return res, nil
}
