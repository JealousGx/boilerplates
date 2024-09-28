import * as cdk from "aws-cdk-lib";
import { LambdaIntegration, RestApi } from "aws-cdk-lib/aws-apigateway";
import * as lambda from "aws-cdk-lib/aws-lambda";
import { Construct } from "constructs";

export class ServerlessAwsCdkStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const gateway = new RestApi(this, "myGateway", {
      defaultCorsPreflightOptions: {
        allowOrigins: ["*"],
        allowMethods: ["GET", "POST", "PUT", "DELETE", "OPTIONS"],
      },
    });

    const baseResource = gateway.root.addResource("api/v1");

    const myFunc = new lambda.Function(this, "MyCDKLambda", {
      code: lambda.Code.fromAsset("out/lambdas/hello"),
      handler: "main",
      runtime: lambda.Runtime.PROVIDED_AL2023,
    });

    const integration = new LambdaIntegration(myFunc);
    const proxyResource = baseResource.addResource("test/{proxy+}");
    proxyResource.addMethod("ANY", integration);
  }
}
