import { RemovalPolicy, Stack, StackProps } from "aws-cdk-lib";
import { LambdaIntegration, RestApi } from "aws-cdk-lib/aws-apigateway";
import { AttributeType, BillingMode, Table } from "aws-cdk-lib/aws-dynamodb";
import {
  Architecture,
  LayerVersion,
  Runtime,
  Function,
  AssetCode,
} from "aws-cdk-lib/aws-lambda";
import { NodejsFunction } from "aws-cdk-lib/aws-lambda-nodejs";
import { Construct } from "constructs";

export class DemoStack extends Stack {
  constructor(scope: Construct, id: string, props?: StackProps) {
    super(scope, id, props);

    // RestAPI to trigger the lambda functions
    const api = new RestApi(this, "Api", {
      restApiName: "DemoLambdaCacheExtensionApi",
    });

    // DynamoDB table to store the data
    const table = new Table(this, "Table", {
      partitionKey: { name: "PK", type: AttributeType.STRING },
      sortKey: { name: "SK", type: AttributeType.STRING },
      billingMode: BillingMode.PAY_PER_REQUEST,
      removalPolicy: RemovalPolicy.DESTROY,
    });

    const environment = {
      MOMENTO_CACHE_NAME: process.env.MOMENTO_CACHE_NAME ?? "",
      MOMENTO_TOKEN: process.env.MOMENTO_TOKEN ?? "",
      DYNAMODB_TABLE_NAME: table.tableName,
      CACHING_DISABLED: "false",
    };

    const lambdaMomentoExtensionLayer = new LayerVersion(
      this,
      "LambdaMomentoExtensionLayer",
      {
        removalPolicy: RemovalPolicy.DESTROY,
        code: new AssetCode("lib/layers/lambda-momento-extension.zip"),
      }
    );

    const layers = [lambdaMomentoExtensionLayer];

    // Lambda function with a typescript handler
    const typescriptLambda = new NodejsFunction(this, "TypescriptLambda", {
      entry: "lib/handlers/typescript-handler.ts",
      handler: "handler",
      runtime: Runtime.NODEJS_18_X,
      architecture: Architecture.ARM_64,
      environment,
      layers,
    });

    const pythonLambda = new Function(this, "PythonLambda", {
      runtime: Runtime.PYTHON_3_10,
      architecture: Architecture.ARM_64,
      handler: "index.handler",
      code: new AssetCode("lib/handlers/python-handler.zip"),
      environment,
      layers,
    });

    // Grant access to the DynamoDB table for all the lambdas
    table.grantReadData(typescriptLambda);
    table.grantReadWriteData(pythonLambda);

    // API Gateway integration
    api.root
      .addResource("typescript")
      .addMethod("GET", new LambdaIntegration(typescriptLambda));
    api.root
      .addResource("python")
      .addMethod("GET", new LambdaIntegration(pythonLambda));
  }
}
