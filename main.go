package main

import (
	"fmt"
	"log"

	"etcd-bootstrap/aws"
	"etcd-bootstrap/etcd"

	"golang.org/x/net/context"
)

func main() {
	metadataService := aws.NewEC2MetadataService()
	m, err := metadataService.GetMetadata()

	if err != nil {
		log.Fatal("Are you kidding me? This should be executed inside an EC2 instance")
	}

	fmt.Println(m)
	asgservice, _ := aws.NewAutoScallingService(m.Region)
	a, err := asgservice.GetAutoScallingGroupOfInstance(m.Region, []*string{&m.InstanceID})
	if err != nil {
		log.Fatal(err)
	}

	var ids []*string
	for _, i := range a.Instances {
		ids = append(ids, i.InstanceId)
	}

	fmt.Printf("Found these %d instances in AGS: %s\n", len(ids), (*a.AutoScalingGroupName))

	if len(ids) == 1 {
		fmt.Println("It seems that we are the only memeber of the cluster. So try to create a new cluster!")
	}

	ec2service, _ := aws.NewEC2Service(m.Region)
	insts, err := ec2service.GetEC2Instance(ids...)
	for _, i := range insts {
		fmt.Printf("Checking ETCD instance at %s", *i.PrivateIpAddress)

		e, err := etcd.New(fmt.Sprintf("http://%s:%d", *i.PrivateIpAddress, 2379))

		if err != nil {
			log.Printf("EtcD instance is not responding at: %s", *i.PrivateIpAddress)
		} else {
			mAPI := e.NewMembersAPI()

			log.Print("Trying to find leader...")
			resp, err := mAPI.Leader(context.Background())
			if err != nil {
				log.Println(err)
			} else {
				// print common key info
				log.Printf("Get is done. Metadata is %q\n", resp)
				/*
					for _, m := range resp {
						// print value
						log.Printf("Members: %s (%s)", m.ID, m.Name)
					}
				*/
			}
		}
	}
}
