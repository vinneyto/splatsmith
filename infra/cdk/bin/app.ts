#!/usr/bin/env node
import * as cdk from 'aws-cdk-lib';
import { ViewerStack } from '../lib/viewer-stack';
import { ViewerServerlessStack } from '../lib/viewer-serverless-stack';

const app = new cdk.App();
new ViewerStack(app, 'SplatmakerViewerStack', {});
new ViewerServerlessStack(app, 'SplatmakerViewerServerlessStack', {});
