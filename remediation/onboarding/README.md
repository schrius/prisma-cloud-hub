# Prisma Onboard Lambda
Lambda function used to onboard new accounts to Primsa
Current version only support AWS

### 
Apply the following command to create the role in the new account:
```bash
aws cloudformation deploy --template-file prisma-cloud-read-only-role.yaml  --stack-name Prisma-Onboard  --capabilities CAPABILITY_NAMED_IAM
```
The command create a Prisma-Onboard Role that allow Prisma to have assume role access to the new account.
### AWS CLI Invoke
Example
```bash
aws lambda invoke --function-name "PrismaOnboard" --payload '{"name": "doitt", "cloudType": "aws", "accountId": "123456789", "groupNames": ["AWS Agency Subscriptions"]}'
```
You should replace the paylay with your new account information, name of the account, accoundId, and groupNames

Done.