package communicationService

import (
	"github.com/BargheNo/Backend/bootstrap"
	"github.com/kavenegar/kavenegar-go"
)

type SMSService struct {
	providerConfig *bootstrap.SMSGateway
	smsTemplates   *bootstrap.SMSTemplates
}

func NewSMSService(
	providerConfig *bootstrap.SMSGateway,
	smsTemplates *bootstrap.SMSTemplates,
) *SMSService {
	return &SMSService{
		providerConfig: providerConfig,
		smsTemplates:   smsTemplates,
	}
}

func (smsService *SMSService) SendOTP(receptor, token string) {
	api := kavenegar.New(smsService.providerConfig.APIKey)
	template := smsService.smsTemplates.OTP
	params := &kavenegar.VerifyLookupParam{}
	if _, err := api.Verify.Lookup(receptor, template, token, params); err != nil {
		panic(err)
	}
}
