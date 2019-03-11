package logic

import (
	"errors"
	"fmt"

	"github.com/theonlyjohnny/go-logger/logger"
	"github.com/theonlyjohnny/phoenix/common"
	"github.com/theonlyjohnny/phoenix/utils"
)

var (
	log *logger.Logger
)

// Controller is the entrpoint to all the logic -- it coordinates the various controllers/actions
type Controller struct {
	args    common.RunArgs
	vpcName string
	ec2     ec2Controller
}

//NewController creates and returns a new controller
func NewController(args common.RunArgs, passedLog *logger.Logger) Controller {
	log = passedLog
	ec2Controller := newEC2Controller(args.Region)
	var vpcName string
	if args.VPCName != "" {
		vpcName = args.VPCName
	} else {
		vpcName = utils.GenerateVPCName(args.Region)
	}
	return Controller{
		args,
		vpcName,
		ec2Controller,
	}
}

//Act starts the controller and act on the controller.action
func (c *Controller) Act() error {
	if c.args.Action == "create" {
		c.create()
	} else {
		return errors.New("only create is implemented right now")
	}
	return nil
}

func (c *Controller) create() error {
	exists := c.ec2.checkVPCExists(c.vpcName)
	log.Debugf("%s exists? %t", c.vpcName, exists)
	if exists {
		return fmt.Errorf("VPC %s already exists", c.vpcName)
	}
	log.Debugf("VPC %s doesn't exist yet, creating", c.vpcName)
	return nil
}
