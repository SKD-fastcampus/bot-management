#!/bin/bash
set -e

# Configuration
USER_NAME="bot-mgmt-admin"
PROFILE_NAME="bot-mgmt"
CLUSTER_NAME="smishing-analysis-cluster"
POLICY_NAME="BotMgmtClusterAdminPolicy"

echo "Starting secure credential setup for user: $USER_NAME..."

# 1. Define Policy JSON
# Allows full control over the specific cluster, global task definition access (required for RunTask), and PassRole (required to launch tasks).
cat <<EOF > policy.temp.json
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
            "Sid": "TaskDefinitionGlobal",
            "Effect": "Allow",
            "Action": [
                "ecs:RegisterTaskDefinition",
                "ecs:ListTaskDefinitions",
                "ecs:DescribeTaskDefinitions",
                "ecs:DeregisterTaskDefinition"
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
                "logs:DescribeLogGroups"
             ],
             "Resource": "*"
        }
    ]
}
EOF

# 2. Create User (Idempotent-ish)
if aws iam get-user --user-name "$USER_NAME" >/dev/null 2>&1; then
    echo "User $USER_NAME already exists."
else
    echo "Creating user $USER_NAME..."
    aws iam create-user --user-name "$USER_NAME" >/dev/null
fi

# 3. Create or Get Policy
ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)
POLICY_ARN="arn:aws:iam::${ACCOUNT_ID}:policy/${POLICY_NAME}"

if ! aws iam get-policy --policy-arn "$POLICY_ARN" >/dev/null 2>&1; then
    echo "Creating policy $POLICY_NAME..."
    POLICY_ARN=$(aws iam create-policy --policy-name "$POLICY_NAME" --policy-document file://policy.temp.json --query 'Policy.Arn' --output text)
else
    echo "Policy $POLICY_NAME already exists. Updating version..."
    # Prune old versions if limit reached is not handled here, simplest is to assume it's fine or just use existing.
    # For robust script, we'll CreatePolicyVersion.
    aws iam create-policy-version --policy-arn "$POLICY_ARN" --policy-document file://policy.temp.json --set-as-default >/dev/null || echo "Warning: Could not update policy version, using existing."
fi

# 4. Attach Policy
echo "Attaching policy to user..."
aws iam attach-user-policy --user-name "$USER_NAME" --policy-arn "$POLICY_ARN"

# 5. Create Access Key & Configure Local Profile (Securely)
echo "Generating Access Keys..."
# Using --output json and writing to a file prevents keys from appearing in stdout/logs
aws iam create-access-key --user-name "$USER_NAME" --output json > keys.temp.json

# Parse Keys using grep/cut to avoid jq dependency
ACCESS_KEY=$(grep '"AccessKeyId":' keys.temp.json | cut -d '"' -f 4)
SECRET_KEY=$(grep '"SecretAccessKey":' keys.temp.json | cut -d '"' -f 4)

if [ -z "$ACCESS_KEY" ] || [ -z "$SECRET_KEY" ]; then
    echo "Error: Failed to retrieve access keys."
    rm keys.temp.json policy.temp.json
    exit 1
fi

echo "Configuring local profile '$PROFILE_NAME'..."
aws configure set aws_access_key_id "$ACCESS_KEY" --profile "$PROFILE_NAME"
aws configure set aws_secret_access_key "$SECRET_KEY" --profile "$PROFILE_NAME"
aws configure set region "ap-northeast-2" --profile "$PROFILE_NAME"
aws configure set output "json" --profile "$PROFILE_NAME"

# 6. Cleanup
rm keys.temp.json policy.temp.json

echo "=========================================="
echo "SUCCESS!"
echo "Created IAM User: $USER_NAME"
echo "Created AWS Profile: $PROFILE_NAME"
echo ""
echo "Credential file updated at ~/.aws/credentials"
echo "I have NOT displayed the secret key."
echo "=========================================="
