#!/bin/bash
set -e

USER_NAME="bot-mgmt-admin"
POLICY_NAME="BotMgmtClusterAdminPolicy"
CLUSTER_NAME="smishing-analysis-cluster"

echo "Updating IAM policy for $USER_NAME..."

# Updated Policy JSON: Added ecs:RunTask to TaskDefinitionGlobal or a new statement.
cat <<EOF > policy.update.json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "ClusterScopedAdmin",
            "Effect": "Allow",
            "Action": [
                "ecs:*"
            ],
            "Resource": [
                "arn:aws:ecs:*:*:cluster/${CLUSTER_NAME}",
                "arn:aws:ecs:*:*:service/${CLUSTER_NAME}/*",
                "arn:aws:ecs:*:*:task/${CLUSTER_NAME}/*",
                "arn:aws:ecs:*:*:container-instance/${CLUSTER_NAME}/*"
            ]
        },
        {
            "Sid": "TaskDefinitionOperations",
            "Effect": "Allow",
            "Action": [
                "ecs:RegisterTaskDefinition",
                "ecs:ListTaskDefinitions",
                "ecs:DescribeTaskDefinition",
                "ecs:DeregisterTaskDefinition",
                "ecs:RunTask"
            ],
            "Resource": "*"
        },
        {
            "Sid": "PassRole",
            "Effect": "Allow",
            "Action": "iam:PassRole",
            "Resource": "*"
        },
        {
             "Sid": "Logging",
             "Effect": "Allow",
             "Action": [
                "logs:CreateLogGroup",
                "logs:CreateLogStream",
                "logs:PutLogEvents",
                "logs:DescribeLogGroups",
                "logs:DescribeLogStreams",
                "logs:GetLogEvents",
                "logs:FilterLogEvents"

             ],
             "Resource": "*"
        },
        {
             "Sid": "S3Access",
             "Effect": "Allow",
             "Action": [
                "s3:ListAllMyBuckets",
                "s3:ListBucket",
                "s3:GetObject"
             ],
             "Resource": "*"
        }
    ]
}

EOF

# Get Account ID & Policy ARN
ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)
POLICY_ARN="arn:aws:iam::${ACCOUNT_ID}:policy/${POLICY_NAME}"

# Update Policy Version
# We need to delete the specific version if we reached the limit (5), but for now let's just try creation.
# If it fails due to limit, we might need a cleanup logic. 
# Simplest fix for "LimitExceeded" in dev scripts is to delete the oldest non-default version, but let's try just set-as-default first.

echo "Creating new policy version for $POLICY_ARN..."
aws iam create-policy-version --policy-arn "$POLICY_ARN" --policy-document file://policy.update.json --set-as-default

# Validation
echo "Policy updated successfully."
rm policy.update.json
