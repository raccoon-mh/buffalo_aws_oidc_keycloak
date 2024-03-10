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

type VmData struct {
	InstanceID   string
	InstanceType string
	PublicDNS    string
	State        string
}

func GetVmListPageHandler(c buffalo.Context) error {
	c.Set("pretitle", "EC2 instacne")
	c.Set("title", "List")

	return c.Render(http.StatusOK, tr.HTML("ec2instance/list.html"))
}

func VmlifecyclePageHandler(c buffalo.Context) error {
	c.Set("pretitle", "EC2 instacne")
	c.Set("title", "LifeCycle")

	return c.Render(http.StatusOK, tr.HTML("ec2instance/list.html"))
}

func VMCreatePageHandler(c buffalo.Context) error {
	c.Set("pretitle", "EC2 instacne")
	c.Set("title", "Create")

	return c.Render(http.StatusOK, tr.HTML("ec2instance/list.html"))
}

func GetVmListHandler(c buffalo.Context) error {
	accessKeyId := c.Param("accessKeyId")
	secretAccessKey := c.Param("secretAccessKey")
	sessionToken := c.Param("sessionToken")

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	}))

	creds := credentials.NewStaticCredentials(
		*aws.String(accessKeyId),
		*aws.String(secretAccessKey),
		*aws.String(sessionToken),
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

	var vmList []VmData

	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			vm := VmData{
				InstanceID:   *instance.InstanceId,
				InstanceType: *instance.InstanceType,
				PublicDNS:    *instance.PublicDnsName,
				State:        *instance.State.Name,
			}
			vmList = append(vmList, vm)
		}
	}

	// c.Set("vmList", vmList)

	return c.Render(http.StatusOK, r.JSON(vmList))
}
