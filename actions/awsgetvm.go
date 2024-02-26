package actions

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/gobuffalo/buffalo"
)

func GetVmHandler(c buffalo.Context) error {
	MyAccessKeyId := c.Param("AccessKeyId")
	MySecretAccessKey := c.Param("SecretAccessKey")
	MySessionToken := c.Param("SessionToken")

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	}))

	creds := credentials.NewStaticCredentials(
		*aws.String(MyAccessKeyId),
		*aws.String(MySecretAccessKey),
		*aws.String(MySessionToken),
	)

	ec2Svc := ec2.New(sess, &aws.Config{
		Region:      aws.String("ap-northeast-2"),
		Credentials: creds,
	})

	ec2Params := &ec2.DescribeInstancesInput{}

	result, err := ec2Svc.DescribeInstances(ec2Params)
	if err != nil {
		fmt.Println("Error describing instances:", err)
		c.Flash().Add("danger", "Error describing instances:"+err.Error())
		return c.Redirect(302, "/")
	}

	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			fmt.Println("Instance ID:", *instance.InstanceId)
			fmt.Println("Instance Type:", *instance.InstanceType)
			fmt.Println("Public DNS:", *instance.PublicDnsName)
			fmt.Println("State:", *instance.State.Name)
			fmt.Println("-----")
		}
	}

	return c.Render(http.StatusOK, r.HTML("main/index.plush.html"))
}
