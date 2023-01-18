package provision

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func SpinNewVm() (string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-west-2"),
		Credentials: credentials.NewSharedCredentials("", ""),
	})

	// Create EC2 service client
	svc := ec2.New(sess)
	launchTemplateId := ""      // Replace from config file
	launchTmeplateVersion := "" // Replace from config file

	launchTemplate := &ec2.LaunchTemplateSpecification{
		LaunchTemplateId: &launchTemplateId,
		Version:          &launchTmeplateVersion,
	}

	// Specify the details of the instance that you want to create.
	runResult, err := svc.RunInstances(&ec2.RunInstancesInput{
		// An Amazon Linux AMI ID for t2.micro instances in the us-west-2 region
		LaunchTemplate: launchTemplate,
		MinCount:       aws.Int64(1),
		MaxCount:       aws.Int64(1),
	})

	log.Info.Println("Creating new instance *************")

	if err != nil {
		log.Info.Println("Could not create instance", err)
		return "", err
	}

	log.Info.Println("Created instance, Instacne ID: ", *runResult.Instances[0].InstanceId)
	private_ip := *runResult.Instances[0].PrivateIpAddress
	log.Info.Println("Created instance, Private IP: ", *runResult.Instances[0].PrivateIpAddress)

	// Add tags to the created instance
	_, errtag := svc.CreateTags(&ec2.CreateTagsInput{
		Resources: []*string{runResult.Instances[0].InstanceId},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("User"),
				Value: aws.String("meghana.r"),
			},
			{
				Key:   aws.String("Production"),
				Value: aws.String("Yes"),
			},
			{
				Key:   aws.String("Project"),
				Value: aws.String("Dev"),
			},
		},
	})
	if errtag != nil {
		log.Error.Println("Could not create tags for instance", runResult.Instances[0].InstanceId, errtag)
		return "", errtag
	}
	log.Info.Println("Tagged the instance")

	allInstances := true

	log.Info.Println("Waiting until instanceStatus to be Ok.......")
	errRunning := svc.WaitUntilInstanceStatusOk(&ec2.DescribeInstanceStatusInput{
		InstanceIds:         []*string{runResult.Instances[0].InstanceId},
		IncludeAllInstances: &allInstances,
	})
	if errRunning != nil {
		log.Error.Println("Instance state is nt okay even after maximum wait window")
		return "", errRunning
	}
	return private_ip, nil

}

