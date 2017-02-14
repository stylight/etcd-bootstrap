package aws

import (
	"errors"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/autoscaling/autoscalingiface"
)

type AutoScalingGroup autoscaling.Group
type autoScalingGroupsService autoscalingiface.AutoScalingAPI

type AutoScalingGroupHelper struct {
	service autoScalingGroupsService
}

func NewAutoScallingService(region string) (*AutoScalingGroupHelper, error) {
	svc := autoscaling.New(
		session.New(&aws.Config{Region: aws.String(region)}),
	)

	return &AutoScalingGroupHelper{
		service: svc,
	}, nil
}

func (as *AutoScalingGroupHelper) GetAutoScallingGroupOfInstance(region string, instanceIDs []*string) (*AutoScalingGroup, error) {
	out, err := as.service.DescribeAutoScalingInstances(&autoscaling.DescribeAutoScalingInstancesInput{
		InstanceIds: instanceIDs,
	})

	if err != nil {
		log.Println("Could not get auto scaling group: ", err)

		return nil, err
	}

	if len(out.AutoScalingInstances) < 1 {
		return nil, errors.New("This instance is not part of Autoscaling :P")
	}

	a := out.AutoScalingInstances[0]

	asgs, err := as.service.DescribeAutoScalingGroups(&autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []*string{a.AutoScalingGroupName},
	})
	if err != nil || len(asgs.AutoScalingGroups) < 1 {
		log.Println("Could not get ASG information: ", err)
		return nil, errors.New("Could not get ASG information")
	}

	res := AutoScalingGroup(*asgs.AutoScalingGroups[0])
	return &res, nil
}
