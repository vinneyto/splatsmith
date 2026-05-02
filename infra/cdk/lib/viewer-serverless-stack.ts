import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as cloudfront from 'aws-cdk-lib/aws-cloudfront';
import * as origins from 'aws-cdk-lib/aws-cloudfront-origins';
import * as dynamodb from 'aws-cdk-lib/aws-dynamodb';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as lambdaNodejs from 'aws-cdk-lib/aws-lambda-nodejs';
import * as s3 from 'aws-cdk-lib/aws-s3';

export class ViewerServerlessStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const jobsTableName = new cdk.CfnParameter(this, 'JobsTableName', { type: 'String' });
    const resultBucketName = new cdk.CfnParameter(this, 'ResultBucketName', { type: 'String' });
    const frontendBucketName = new cdk.CfnParameter(this, 'FrontendBucketName', { type: 'String' });

    const jobsTable = dynamodb.Table.fromTableName(this, 'JobsTable', jobsTableName.valueAsString);
    const resultBucket = s3.Bucket.fromBucketName(this, 'ResultBucket', resultBucketName.valueAsString);
    const frontendBucket = s3.Bucket.fromBucketName(this, 'FrontendBucket', frontendBucketName.valueAsString);

    const jobsApi = new lambdaNodejs.NodejsFunction(this, 'ViewerJobsApiFn', {
      runtime: lambda.Runtime.NODEJS_20_X,
      entry: 'lambda/jobs-api/index.ts',
      handler: 'handler',
      timeout: cdk.Duration.seconds(10),
      memorySize: 256,
      bundling: {
        target: 'node20',
        format: lambdaNodejs.OutputFormat.ESM,
      },
      environment: {
        JOBS_TABLE_NAME: jobsTableName.valueAsString,
      },
    });

    jobsTable.grantReadData(jobsApi);

    const jobsApiUrl = jobsApi.addFunctionUrl({ authType: lambda.FunctionUrlAuthType.NONE });
    const jobsApiDomain = cdk.Fn.select(2, cdk.Fn.split('/', jobsApiUrl.url));

    const distribution = new cloudfront.Distribution(this, 'ViewerDistribution', {
      defaultRootObject: 'index.html',
      defaultBehavior: {
        origin: origins.S3BucketOrigin.withOriginAccessControl(frontendBucket),
        viewerProtocolPolicy: cloudfront.ViewerProtocolPolicy.REDIRECT_TO_HTTPS,
        cachePolicy: cloudfront.CachePolicy.CACHING_OPTIMIZED,
      },
      additionalBehaviors: {
        'v1/*': {
          origin: new origins.HttpOrigin(jobsApiDomain),
          allowedMethods: cloudfront.AllowedMethods.ALLOW_ALL,
          cachePolicy: cloudfront.CachePolicy.CACHING_DISABLED,
          viewerProtocolPolicy: cloudfront.ViewerProtocolPolicy.REDIRECT_TO_HTTPS,
          originRequestPolicy: cloudfront.OriginRequestPolicy.ALL_VIEWER_EXCEPT_HOST_HEADER,
        },
        healthz: {
          origin: new origins.HttpOrigin(jobsApiDomain),
          allowedMethods: cloudfront.AllowedMethods.ALLOW_ALL,
          cachePolicy: cloudfront.CachePolicy.CACHING_DISABLED,
          viewerProtocolPolicy: cloudfront.ViewerProtocolPolicy.REDIRECT_TO_HTTPS,
          originRequestPolicy: cloudfront.OriginRequestPolicy.ALL_VIEWER_EXCEPT_HOST_HEADER,
        },
        'results/*': {
          origin: origins.S3BucketOrigin.withOriginAccessControl(resultBucket),
          viewerProtocolPolicy: cloudfront.ViewerProtocolPolicy.REDIRECT_TO_HTTPS,
          cachePolicy: cloudfront.CachePolicy.CACHING_OPTIMIZED,
        },
      },
      errorResponses: [
        {
          httpStatus: 403,
          responseHttpStatus: 200,
          responsePagePath: '/index.html',
          ttl: cdk.Duration.minutes(1),
        },
        {
          httpStatus: 404,
          responseHttpStatus: 200,
          responsePagePath: '/index.html',
          ttl: cdk.Duration.minutes(1),
        },
      ],
    });

    jobsApi.addEnvironment('RESULTS_BASE_URL', `https://${distribution.distributionDomainName}/results`);

    new cdk.CfnOutput(this, 'ViewerServerlessURL', {
      value: `https://${distribution.distributionDomainName}`,
    });
    new cdk.CfnOutput(this, 'ViewerJobsApiFunctionUrl', {
      value: jobsApiUrl.url,
    });
  }
}
