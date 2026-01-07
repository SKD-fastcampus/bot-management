package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/SKD-fastcampus/bot-management/services/bot-mgmt-server/internal/domain"
)

type ECSClient struct {
	client     *ecs.Client
	cluster    string
	taskDef    string
	subnetID   string
	secGroupID string
}

func NewECSClient(cfg aws.Config, cluster, taskDef, subnetID, secGroupID string) *ECSClient {
	return &ECSClient{
		client:     ecs.NewFromConfig(cfg),
		cluster:    cluster,
		taskDef:    taskDef,
		subnetID:   subnetID,
		secGroupID: secGroupID,
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
				Subnets:        []string{c.subnetID},
				SecurityGroups: []string{c.secGroupID},
				AssignPublicIp: types.AssignPublicIpEnabled, // Assuming public access needed for external URL
			},
		},
		Overrides: &types.TaskOverride{
			ContainerOverrides: []types.ContainerOverride{
				{
					Name: aws.String("bot-container"), // TODO: Make configurable
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
		return domain.TaskStatusPending, nil
	case "RUNNING":
		return domain.TaskStatusRunning, nil
	case "DEACTIVATING", "STOPPING", "DEPROVISIONING":
		return domain.TaskStatusRunning, nil // Still shutting down
	case "STOPPED":
		// Check exit code
		if task.StopCode != types.TaskStopCodeEssentialContainerExited {
			return domain.TaskStatusFailed, nil // Infrastructure failure?
		}
		// Check container exit code
		for _, container := range task.Containers {
			if container.ExitCode != nil && *container.ExitCode != 0 {
				return domain.TaskStatusFailed, nil
			}
		}
		return domain.TaskStatusCompleted, nil
	default:
		return domain.TaskStatusPending, nil
	}
}
