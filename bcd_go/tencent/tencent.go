package tencent

import (
	"github.com/pkg/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/regions"
	ocrTable "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ocr/v20181119"
)

var ClientOcrTable *ocrTable.Client

func Init() error {
	credential := common.NewCredential("AKIDIPFo79uObYDbiY76oFQKwHMAhwuc3C4g", "NTIr9AuNZT6SYNuAAtY5NQEm47uRikX4")
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqTimeout = 30
	var err error
	ClientOcrTable, err = ocrTable.NewClient(credential, regions.Beijing, cpf)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
