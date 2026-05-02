#!/usr/bin/env node
import * as cdk from 'aws-cdk-lib';
import { ViewerStack } from '../lib/viewer-stack';

const app = new cdk.App();
new ViewerStack(app, 'SplatmakerViewerStack', {});
