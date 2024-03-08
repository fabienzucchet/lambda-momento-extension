# Demo for the lambda-momento extension

This folder contains a very basic CDK code to deploy two lambda functions with the lambda-momento extension. Both functions are running the same logic (retrieve an item from the cache or a DynamoDB table) but one of them is in using Python while the other one is using Typescript.

## How to deploy the demo?

1. Clone the repository
2. Go to the `demo` folder
3. Run `npm install`
4. Copy the `.env.template` file to `.env`
5. Create a cache in [Momento](https://momento.com) and update the value for `MOMENTO_CACHE_NAME` in the `.env` file
6. Create an API key in [Momento](https://momento.com) and update the value for `MOMENTO_TOKEN` in the `.env` file
7. Run `cdk deploy`

## How to test the demo?

First of all, you need to populate the DynamoDB table with some items of your choice.

Once the database is populated you can test the demo by invoking the lambda functions through API Gateway. For example:

```bash
# Invoke the typescript lambda function using the API Gateway
curl https://<api-gateway-id>.execute-api.eu-west-2.amazonaws.com/prod/typescript\?pk\=test\&sk\=demo
# Invoke the python lambda function using the API Gateway
curl https://<api-gateway-id>.execute-api.eu-west-2.amazonaws.com/prod/python\?pk\=test\&sk\=demo
```

On the first invocation, the lambda function will retrieve the item from the DynamoDB table and store it in the cache. On the second invocation, the lambda function will retrieve the item from the cache.
