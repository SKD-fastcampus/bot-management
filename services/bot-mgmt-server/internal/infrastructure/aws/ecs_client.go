package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/SKD-fastcampus/bot-management/services/bot-mgmt-server/internal/domain"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

type ECSClient struct {
	client        *ecs.Client
	cluster       string
	taskDef       string
	containerName string
	subnets       []string
	secGroupID    string
}

func NewECSClient(cfg aws.Config, cluster, taskDef, containerName string, subnets []string, secGroupID string) *ECSClient {
	return &ECSClient{
		client:        ecs.NewFromConfig(cfg),
		cluster:       cluster,
		taskDef:       taskDef,
		containerName: containerName,
		subnets:       subnets,
		secGroupID:    secGroupID,
	}
}

func (c *ECSClient) RunBot(ctx context.Context, task *domain.AnalysisTask) (string, error) {
	// Prepare environment overrides or command overrides if needed
	// Passing URL and UUID as environment variables
	runTaskInput := &ecs.RunTaskInput{
		Cluster:        aws.String(c.cluster),
		TaskDefinition: aws.String(c.taskDef),
		LaunchType:     types.LaunchTypeFargate,
		NetworkConfiguration: &types.NetworkConfiguration{
			AwsvpcConfiguration: &types.AwsVpcConfiguration{
				Subnets:        c.subnets,
				SecurityGroups: []string{c.secGroupID},
				AssignPublicIp: types.AssignPublicIpEnabled, // Assuming public access needed for external URL
			},
		},

		Overrides: &types.TaskOverride{
			ContainerOverrides: []types.ContainerOverride{
				{
					Name: aws.String(c.containerName),
					Environment: []types.KeyValuePair{
						{Name: aws.String("TARGET_URL"), Value: aws.String(task.URL)},
						{Name: aws.String("REQUEST_UUID"), Value: aws.String(task.RequestUUID)},
					},
				},
			},
		},
	}

	out, err := c.client.RunTask(ctx, runTaskInput)
	if err != nil {
		return "", err
	}

	if len(out.Tasks) == 0 {
		return "", fmt.Errorf("no tasks started")
	}

	return *out.Tasks[0].TaskArn, nil
}

func (c *ECSClient) GetBotStatus(ctx context.Context, externalID string) (domain.TaskStatus, error) {
	out, err := c.client.DescribeTasks(ctx, &ecs.DescribeTasksInput{
		Cluster: aws.String(c.cluster),
		Tasks:   []string{externalID},
	})
	if err != nil {
		return "", err
	}

	if len(out.Tasks) == 0 {
		return "", fmt.Errorf("task not found")
	}

	task := out.Tasks[0]
	// Map AWS status to Domain status
	switch *task.LastStatus {
	case "PROVISIONING", "PENDING", "ACTIVATING":
		return domain.TaskStatusRunning, nil // Treat provisioning as running so we keep polling it
	case "RUNNING":
		return domain.TaskStatusRunning, nil
	case "DEACTIVATING", "STOPPING", "DEPROVISIONING":
		return domain.TaskStatusRunning, nil // Still shutting down
	case "STOPPED":
		// Check container exit code
		for _, container := range task.Containers {
			// If we can't find exit code, it's suspicious, but if it is present and 0, success.
			if container.ExitCode != nil && *container.ExitCode != 0 {
				return domain.TaskStatusFailed, nil
			}
		}
		// If all containers have exit code 0 (or no exit code but stopped?), we consider it completed?
		// Usually stopped with no exit code means infrastructure issue or user stop.
		// Let's rely on ExitCode != 0 check for failure.

		return domain.TaskStatusCompleted, nil

	default:
		return domain.TaskStatusPending, nil
	}
}
