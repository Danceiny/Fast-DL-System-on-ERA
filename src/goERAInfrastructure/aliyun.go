package goERAInfrastructure

// 使用阿里云的SDK
//https://github.com/aliyun/alibaba-cloud-sdk-go
import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
)

func main() {
	// 创建ecsClient实例
	ecsClient, err := ecs.NewClientWithAccessKey(
		"<your-region-id>",         // 您的可用区ID
		"<your-access-key-id>",     // 您的Access Key ID
		"<your-access-key-secret>") // 您的Access Key Secret
	if err != nil {
		// 异常处理
		panic(err)
	}
	// 创建API请求并设置参数
	request := ecs.CreateDescribeInstancesRequest()
	request.PageSize = "10"
	// 发起请求并处理异常
	response, err := ecsClient.DescribeInstances(request)
	if err != nil {
		// 异常处理
		panic(err)
	}
	fmt.Println(response)
}
