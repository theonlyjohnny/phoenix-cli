package utils

import (
	"strconv"
	"strings"
	"time"
)

// US East (N. Virginia) us-east-1
// US East (Ohio) us-east-2
// US West (N. California) us-west-1
// US West (Oregon) us-west-2
// Asia Pacific (Mumbai) ap-south-1
// Asia Pacific (Seoul) ap-northeast-2
// Asia Pacific (Singapore) ap-southeast-2
// Asia Pacific (Sydney) ap-southeast-1
// Asia Pacific (Tokyo) ap-northeast-1
// Canada (Central) ca-central-1
// EU (Frankfurto) eu-central-1
// EU (Ireland) eu-west-1
// EU (London) eu-west-2
// EU (Paris) eu-west-3
// EU (Stockholm) eu-north-1
// South America (São Paulo) sa-east-1

//GenerateVPCName generates a unique VPC Name. It will be assigned to the VPC tag "Name",
//and will be used as the primary key throughout phoenix.
func GenerateVPCName(region string) string {
	prefix := generateVPCNamePrefix(region)
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	return prefix + timestamp
}

func generateVPCNamePrefix(region string) string {
	var vpcName string
	split := strings.Split(region, "-")
	vpcName += split[0]
	if strings.Contains(region, "north") {
		vpcName += "n"
	}

	if strings.Contains(region, "south") {
		vpcName += "s"
	}

	if strings.Contains(region, "east") {
		vpcName += "e"
	}

	if strings.Contains(region, "west") {
		vpcName += "w"
	}

	if strings.Contains(region, "central") {
		vpcName += "c"
	}

	vpcName += split[2]
	return vpcName
}
