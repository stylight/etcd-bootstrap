package aws

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
)

//AWSService interface for ServiceHelper
type AWSService interface {
	GetAutoScalingSelfInstances() ([]*EC2Instance, error)

	NewEC2MetadataService() *EC2MetadataHelper
	NewAutoScallingService() *AutoScalingGroupHelper
	NewEC2Service() *EC2Helper
}

//AWSServiceHelper provides access to various AWS resources
type AWSServiceHelper struct {
	Session *session.Session
}

//New initializes new AWSServiceHelper object
func New(region string) AWSService {
	return &AWSServiceHelper{
		Session: session.New(&aws.Config{Region: aws.String(region)}),
	}
}

//NewEC2MetadataService initializes new EC2Metadata helper
func (h *AWSServiceHelper) NewEC2MetadataService() *EC2MetadataHelper {
	return &EC2MetadataHelper{
		service: ec2metadata.New(h.Session),
	}
}

//NewAutoScallingService initializes new Autoscaling Helper
func (h *AWSServiceHelper) NewAutoScallingService() *AutoScalingGroupHelper {
	svc := autoscaling.New(h.Session)

	return &AutoScalingGroupHelper{
		service: svc,
	}
}

//NewEC2Service initializes new EC2 Helper
func (h *AWSServiceHelper) NewEC2Service() *EC2Helper {
	svc := ec2.New(h.Session)

	return &EC2Helper{
		service: svc,
	}
}

/*GetAutoScalingSelfInstances returns autoscaling instances that this instance belongs to
 *
 * This is a helper method that finds the current instance's ASG and then fetches all
 * the instances that are active within the ASG including the current instnance
 */
func (h *AWSServiceHelper) GetAutoScalingSelfInstances() ([]*EC2Instance, error) {
	m, err := h.NewEC2MetadataService().GetMetadata()
	if err != nil {
		panic("Are you kidding me? This should be executed inside an EC2 instance")
	}

	a, err := h.NewAutoScallingService().GetAutoScallingGroupOfInstance([]*string{&m.InstanceID})
	if err != nil {
		log.Fatal(err)
	}

	ids := make([]*string, len(a.Instances))
	for x, i := range a.Instances {
		ids[x] = i.InstanceId
	}

	insts, err := h.NewEC2Service().GetRunningEC2Instance(ids...)

	return insts, err
}
