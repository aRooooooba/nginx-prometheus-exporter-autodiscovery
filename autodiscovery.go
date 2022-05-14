package main

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

func getTargetIPs(region string, clusterName string, serviceName string) ([]string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return nil, err
	}
	client := ecs.NewFromConfig(cfg)
	listTasksOutput, err := client.ListTasks(context.TODO(), &ecs.ListTasksInput{
		Cluster:     &clusterName,
		ServiceName: &serviceName,
	})
	if err != nil {
		return nil, err
	}
	describeTasksOutput, err := client.DescribeTasks(context.TODO(), &ecs.DescribeTasksInput{
		Tasks:   listTasksOutput.TaskArns,
		Cluster: &clusterName,
	})
	if err != nil {
		return nil, err
	}
	res := []string{}
	for _, task := range describeTasksOutput.Tasks {
		for _, attachment := range task.Attachments {
			if *attachment.Type == "ElasticNetworkInterface" {
				for _, kv := range attachment.Details {
					if *kv.Name == "privateIPv4Address" {
						res = append(res, *kv.Value)
					}
				}
			}
		}
	}
	return res, nil
}
