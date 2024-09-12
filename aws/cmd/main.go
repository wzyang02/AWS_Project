package main

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"


	"fmt"
	"log"
	"flag"
)

func createEC2(name string) string {

	sess, err := session.NewSession(&aws.Config{
        Region: aws.String("us-west-2")},
    )

    if err != nil{
		fmt.Println("seesion failed", err)
	}

	_, err = sess.Config.Credentials.Get()

	if err != nil{
		fmt.Println("not found", err)
	} 

    // Create EC2 service client
    svc := ec2.New(sess)

    // Specify the details of the instance that you want to create.
    runResult, err := svc.RunInstances(&ec2.RunInstancesInput{
        // An Amazon Linux AMI ID for t2.micro instances in the us-west-2 region
        ImageId:      aws.String("ami-e7527ed7"),
        InstanceType: aws.String("t2.micro"),
        MinCount:     aws.Int64(1),
        MaxCount:     aws.Int64(1),
    })

    if err != nil {
        fmt.Println("Could not create instance", err)
        return ""
    }

    fmt.Println("Created instance", *runResult.Instances[0].InstanceId)

    // Add tags to the created instance
    _, errtag := svc.CreateTags(&ec2.CreateTagsInput{
        Resources: []*string{runResult.Instances[0].InstanceId},
        Tags: []*ec2.Tag{
            {
                Key:   aws.String("Name"),
                Value: aws.String(name),
            },
        },
    })
    if errtag != nil {
        log.Println("Could not create tags for instance", runResult.Instances[0].InstanceId, errtag)
        return ""
    }

	fmt.Println("Successfully tagged instance")
	instance_id := *runResult.Instances[0].InstanceId
	input := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{aws.String(instance_id)},
	}
	err2 := svc.WaitUntilInstanceRunning(input)

	if err2 != nil {
		return "Error"
	}
	return instance_id

	
}


func createKeyPair() {
	sess, err := session.NewSession(&aws.Config{
        Region: aws.String("us-west-2")},
    )

    if err != nil{
		fmt.Println("seesion failed", err)
	}

	_, err = sess.Config.Credentials.Get()

	if err != nil{
		fmt.Println("not found", err)
	} 

    // Create EC2 service client
	svc := ec2.New(sess)
	input := &ec2.CreateKeyPairInput{
		KeyName: aws.String("my-key-pair"),
	}
	
	result, err := svc.CreateKeyPair(input)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	
	fmt.Println(result)
}
func createVolume() string {
	sess, err := session.NewSession(&aws.Config{
        Region: aws.String("us-west-2")},
	)
	if err != nil{
		fmt.Println("seesion failed", err)
	}

	_, err = sess.Config.Credentials.Get()

	if err != nil{
		fmt.Println("not found", err)
	} 

	svc := ec2.New(sess)
	input := &ec2.CreateVolumeInput{
    AvailabilityZone: aws.String("us-west-2a"),
    Size:             aws.Int64(1),
	VolumeType:       aws.String("gp2"),
	}

	result, err := svc.CreateVolume(input)
	if err != nil {
		fmt.Println(err.Error())
		return ""
    }
	fmt.Println(*result.VolumeId)
	vol_id := *result.VolumeId
	input2 := &ec2.DescribeVolumesInput{
		VolumeIds: []*string{aws.String(vol_id)},
	}
	err2 := svc.WaitUntilVolumeAvailable(input2)

	if err2 != nil {
		return "Error"
	}

	
	return vol_id
}

func attachVolume(device string, instance_id string, vol_id string) {
	sess, err := session.NewSession(&aws.Config{
        Region: aws.String("us-west-2")},
	)
	if err != nil{
		fmt.Println("seesion failed", err)
	}

	_, err = sess.Config.Credentials.Get()

	if err != nil{
		fmt.Println("not found", err)
	} 

	svc := ec2.New(sess)
	
	input := &ec2.AttachVolumeInput{
		Device:     aws.String(device),
		InstanceId: aws.String(instance_id),
		VolumeId:   aws.String(vol_id),
	}
	
	result, err := svc.AttachVolume(input)
	if err != nil {
		fmt.Println(err.Error())
	
		return
	}
	fmt.Println(result)
}

func detachVolume(vol_id string) {
	sess, err := session.NewSession(&aws.Config{
        Region: aws.String("us-west-2")},
	)
	if err != nil{
		fmt.Println("seesion failed", err)
	}

	_, err = sess.Config.Credentials.Get()

	if err != nil{
		fmt.Println("not found", err)
	} 

	svc := ec2.New(sess)
	input := &ec2.DetachVolumeInput{
		VolumeId: aws.String(vol_id),
	}
	
	result, err := svc.DetachVolume(input)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	
	fmt.Println(result)
}

func deleteVolume(vol_id string) {
	sess, err := session.NewSession(&aws.Config{
        Region: aws.String("us-west-2")},
	)
	if err != nil{
		fmt.Println("seesion failed", err)
	}

	_, err = sess.Config.Credentials.Get()

	if err != nil{
		fmt.Println("not found", err)
	} 

	svc := ec2.New(sess)
	
	input := &ec2.DeleteVolumeInput{
		VolumeId: aws.String(vol_id),
	}
	
	result, err := svc.DeleteVolume(input)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	
	fmt.Println(result)
	
}


func stopEC2(instanceID string) {
  // Create a new AWS session
  sess, err := session.NewSession(&aws.Config{
 
    Region: aws.String("us-west-2")},
 
  )
 
 
  if err != nil {
 
    fmt.Println("Error creating session:", err)
 
    return
 
  }
 
 
  // Create a new EC2 service client
  svc := ec2.New(sess)
 
 
  // Specify the instance ID to stop
 
 
  // Stop the instance
  input := &ec2.StopInstancesInput{
 
    InstanceIds: []*string{
 
      aws.String(instanceID),
 
    },
 
  }
 
 
  result, err := svc.StopInstances(input)
 
  if err != nil {
 
    fmt.Println("Error stopping instance:", err)
 
    return
 
  }
 
 
  fmt.Println("Instance stopped:", result.StoppingInstances)
 
}

func startEC2(instanceID string) {
	sess, err := session.NewSession(&aws.Config{
 
		Region: aws.String("us-west-2")},
	 
	  )
	 
	 
	  if err != nil {
	 
		fmt.Println("Error creating session:", err)
	 
		return
	 
	  }
	 
	 
	  // Create a new EC2 service client
	  svc := ec2.New(sess)
	 
	 
	  // Specify the instance ID to stop

	  input := &ec2.StartInstancesInput{
 
		InstanceIds: []*string{
	 
		  aws.String(instanceID),
	 
		},
	 
	  }
	  result, err := svc.StartInstances(input)
 
	  if err != nil {
	 
		fmt.Println("Error starting instance:", err)
	 
		return
	 
	  }
	 
	 
	  fmt.Println("Instance started:", result.StartingInstances)
}

func terminateEC2(instanceID string) {
	sess, err := session.NewSession(&aws.Config{
 
		Region: aws.String("us-west-2")},
	 
	  )	 
	  if err != nil {
		fmt.Println("Error starting session:", err)
		return	 
	  }
	 
	  // Create a new EC2 service client
	  svc := ec2.New(sess)
	 
	 
	  // Specify the instance ID to stop

	  input := &ec2.TerminateInstancesInput{
 
		InstanceIds: []*string{
	 
		  aws.String(instanceID),
	 
		},
	 
	  }
	  result, err := svc.TerminateInstances(input)
 
	  if err != nil {
	 
		fmt.Println("Error terminating instance:", err)
	 
		return
	 
	  }
	 
	 
	  fmt.Println("Instance terminated:", result.TerminatingInstances)
}

func statusEC2(instanceID string) {
	sess, err := session.NewSession(&aws.Config{
 
		Region: aws.String("us-west-2")},
	 
	  )
	 
	 
	  if err != nil {
	 
		fmt.Println("Error starting session:", err)
	 
		return
	 
	  }
	 
	 
	  // Create a new EC2 service client
	  svc := ec2.New(sess)
	 
	 
	  // Specify the instance ID to stop

	  input := &ec2.DescribeInstanceStatusInput{
 
		InstanceIds: []*string{
	 
		  aws.String(instanceID),
	 
		},
	 
	  }
	  result, err := svc.DescribeInstanceStatus(input)
 
	  if err != nil {
	 
		fmt.Println("Error terminating instance:", err)
	 
		return
	 
	  }
	 
	 
	  fmt.Println(result)
}

func describeInstance() (string, []string) {
	var vol_id []string
	sess, err := session.NewSession(&aws.Config{
 
		Region: aws.String("us-west-2")},
	 
	  )
	 
	 
	  if err != nil {
	 
		fmt.Println("Error starting session:", err)
	 
		return "" , vol_id
	 
	  }
	 
	 
	// Create a new EC2 service client
	svc := ec2.New(sess)
	 
	result, err := svc.DescribeInstances(nil)
    if err != nil {
        fmt.Println("Error", err)
    }
	
	instance := result.Reservations[0].Instances[0]


	instance_id := *instance.InstanceId

	blockdevicemappings := result.Reservations[0].Instances[0].BlockDeviceMappings
	for _ , mappings := range(blockdevicemappings) {
		if *mappings.DeviceName != "/dev/xvda" && *mappings.Ebs.DeleteOnTermination == false {
			vol_id = append(vol_id, *mappings.Ebs.VolumeId)
		}
	}
	
	return instance_id, vol_id
}



func main() {
	var option string
	var EC2_name string
	var EC2id string
	flag.StringVar(&option, "o", "", "help message for flagname")
	flag.StringVar(&EC2_name, "n", "", "help message for flagname")
	flag.Parse()
	fmt.Println(option)
	if option == "Stop" {
		fmt.Println("Stop EC2")
		stopEC2(EC2id)
	} else if option == "Create" {
		instance_id := createEC2(EC2_name)
		vol_id := createVolume()
		attachVolume("/dev/xvdb", instance_id, vol_id)
	} else if option == "Start" {
		startEC2(EC2id)
	} else if option == "Terminate" {
		EC2id, vol_id := describeInstance()
		for _ , id := range(vol_id) {
			detachVolume(id)
		}
		terminateEC2(EC2id)
		for _ , id := range(vol_id) {
			deleteVolume(id)
		}
	} else if option == "Status" {
		statusEC2(EC2id)
	} else {
		return
	}

}
