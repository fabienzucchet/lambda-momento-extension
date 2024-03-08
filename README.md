# lambda-momento-extension

This is a demo Lambda extension to showcase how Lambda extensions can make caching easy for Lambda functions no matter what language you are using. With this extension, developers can easily cache dynamodb GetItem responses in Momento without having to handle the caching logic in their Lambda function.

> This is a demo extension and is intended to be used for educational purposes. The code is not production-ready and should not be used in a production environment.

## How does it work?

The extension exposes a HTTP webserver on the port 4000 that can be used to retrieve items from DynamoDB/Momento. The extension is handling all the caching mechanism and the Lambda function only needs to make a request to the extension to retrieve the item.

You can use the following code to retrieve an item with pk `test` and sk `demo`:

```python
import requests

response = requests.get('http://localhost:4000/get-item?pk=test&sk=demo')
print(response.json())
```

## How to publish the extension?

You can use the makefile to build and publish the extension to AWS. You can use the following commands to build and publish the extension:

```bash
make deploy
```

This will return the ARN of the version of the extension that has been published to AWS.

## How to use the extension?

In order to use the extension in your lambda, you need to:

1. add the ARN of the extension to the `layers` property of the lambda function.
2. You also need to grant permissions to the lambda function to read from DynamoDB
3. You need to set the following environment variables:
   - `MOMENTO_CACHE_NAME`: The name of the cache in Momento
   - `MOMENTO_TOKEN`: The API key to access the cache in Momento
   - `DYNAMODB_TABLE_NAME`: The name of the DynamoDB table

## Demo infrastructure

You can check the `demo` folder to see a very basic CDK code to deploy two lambda functions with the lambda-momento extension. Both functions are running the same logic (retrieve an item from the cache or a DynamoDB table) but one of them is in using Python while the other one is using Typescript.
