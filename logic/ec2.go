package logic

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var ec2Client *ec2.EC2

type ec2Controller struct {
	ec2 *ec2.EC2
}

func newEC2Controller(region string) ec2Controller {
	if ec2Client == nil {
		ec2Client = ec2.New(session.New(), aws.NewConfig().WithRegion(region))
	}
	c := ec2Controller{
		ec2Client,
	}
	return c
}

func (c *ec2Controller) getAWSVPCbyVPCName(vpcName string) *ec2.Vpc {

	tagName := "tag:Name"

	input := &ec2.DescribeVpcsInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: &tagName,
				Values: []*string{
					&vpcName,
				},
			},
		},
	}

	result, err := c.ec2.DescribeVpcs(input)
	if err != nil {
		log.Errorf("Couldn't getAWSVPCbyVPCName: %s", err.Error())
		return nil
	}

	log.Debugf("getAWSVPCbyVPCName %s result: %#v", vpcName, result)

	resultLength := len(result.Vpcs)

	if resultLength == 1 {
		return result.Vpcs[0]
	} else if resultLength > 1 {
		log.Warnf(">1 VPCs w/ same VPC Name? vpcName: %s, result: %v", vpcName, result)
		return nil
	}
	return nil
}

func (c *ec2Controller) checkVPCExists(vpcName string) bool {
	return c.getAWSVPCbyVPCName(vpcName) != nil
}
