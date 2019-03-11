package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateVPCName(t *testing.T) {
	answers := map[string]string{
		"us-east-1":      "use1",
		"us-east-2":      "use2",
		"us-west-1":      "usw1",
		"us-west-2":      "usw2",
		"ap-south-1":     "aps1",
		"ap-northeast-2": "apne2",
		"ap-southeast-1": "apse1",
		"ap-southeast-2": "apse2",
		"ap-northeast-1": "apne1",
		"ca-central-1":   "cac1",
		"eu-central-1":   "euc1",
		"eu-west-1":      "euw1",
		"eu-west-2":      "euw2",
		"eu-west-3":      "euw3",
		"eu-north-1":     "eun1",
		"sa-east-1":      "sae1",
	}
	for region, vpcPrefix := range answers {
		assert.Equal(t, vpcPrefix, generateVPCNamePrefix(region))
	}
}
