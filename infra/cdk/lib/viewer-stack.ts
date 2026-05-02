import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as ec2 from 'aws-cdk-lib/aws-ec2';
import * as ecs from 'aws-cdk-lib/aws-ecs';
import * as ecsPatterns from 'aws-cdk-lib/aws-ecs-patterns';
import * as ecr from 'aws-cdk-lib/aws-ecr';
import * as iam from 'aws-cdk-lib/aws-iam';
import * as cognito from 'aws-cdk-lib/aws-cognito';
import * as elbv2 from 'aws-cdk-lib/aws-elasticloadbalancingv2';
import * as elbv2Actions from 'aws-cdk-lib/aws-elasticloadbalancingv2-actions';

export class ViewerStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const vpc = new ec2.Vpc(this, 'ViewerVpc', { maxAzs: 2 });
    const cluster = new ecs.Cluster(this, 'ViewerCluster', { vpc });

    const repoName = new cdk.CfnParameter(this, 'ApiEcrRepositoryName', { type: 'String' });
    const imageTag = new cdk.CfnParameter(this, 'ApiImageTag', { type: 'String', default: 'latest' });
    const jobsTable = new cdk.CfnParameter(this, 'JobsTableName', { type: 'String' });
    const resultBucket = new cdk.CfnParameter(this, 'ResultBucketName', { type: 'String' });

    const userPool = new cognito.UserPool(this, 'ViewerUserPool', {
      selfSignUpEnabled: false,
      signInAliases: { email: true },
    });
    const userPoolClient = new cognito.UserPoolClient(this, 'ViewerUserPoolClient', {
      userPool,
      generateSecret: true,
      oAuth: { flows: { authorizationCodeGrant: true }, scopes: [cognito.OAuthScope.OPENID, cognito.OAuthScope.EMAIL] },
    });
    const domainPrefix = 'splatmaker-viewer-auth';
    const userPoolDomain = new cognito.UserPoolDomain(this, 'ViewerUserPoolDomain', {
      userPool,
      cognitoDomain: { domainPrefix },
    });

    const taskRole = new iam.Role(this, 'ViewerTaskRole', {
      assumedBy: new iam.ServicePrincipal('ecs-tasks.amazonaws.com'),
    });
    taskRole.addToPolicy(new iam.PolicyStatement({
      actions: ['dynamodb:Scan', 'dynamodb:GetItem', 'dynamodb:Query'],
      resources: [`arn:aws:dynamodb:${this.region}:${this.account}:table/${jobsTable.valueAsString}`],
    }));
    taskRole.addToPolicy(new iam.PolicyStatement({
      actions: ['s3:GetObject'],
      resources: [`arn:aws:s3:::${resultBucket.valueAsString}/*`],
    }));

    const service = new ecsPatterns.ApplicationLoadBalancedFargateService(this, 'ViewerService', {
      cluster,
      publicLoadBalancer: true,
      cpu: 512,
      memoryLimitMiB: 1024,
      desiredCount: 1,
      taskImageOptions: {
        image: ecs.ContainerImage.fromEcrRepository(
          ecr.Repository.fromRepositoryName(this, 'ApiRepo', repoName.valueAsString),
          imageTag.valueAsString,
        ),
        containerPort: 8080,
        taskRole,
        environment: {
          SPLATMAKER_MODE: 'aws',
          SPLATMAKER_API_ADDR: ':8080',
          SPLATMAKER_AWS_REGION: this.region,
          SPLATMAKER_AWS_JOBS_TABLE: jobsTable.valueAsString,
          SPLATMAKER_AWS_RESULT_BUCKET: resultBucket.valueAsString,
        },
      },
    });

    service.targetGroup.configureHealthCheck({ path: '/healthz' });
    service.listener.addAction('CognitoAuth', {
      priority: 10,
      conditions: [elbv2.ListenerCondition.pathPatterns(['/*'])],
      action: new elbv2Actions.AuthenticateCognitoAction({
        userPool,
        userPoolClient,
        userPoolDomain,
        next: elbv2.ListenerAction.forward([service.targetGroup]),
      }),
    });

    new cdk.CfnOutput(this, 'ViewerURL', { value: `https://${service.loadBalancer.loadBalancerDnsName}` });
  }
}
