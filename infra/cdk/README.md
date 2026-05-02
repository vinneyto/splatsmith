## Splatmaker Viewer CDK

Deploys minimal jobs viewer backend on ECS Fargate behind ALB + Cognito auth.

Required parameters:
- `ApiEcrRepositoryName`
- `ApiImageTag`
- `JobsTableName`
- `ResultBucketName`

Run:
```bash
cd infra/cdk
npm i
npm run synth
npx cdk deploy
```
