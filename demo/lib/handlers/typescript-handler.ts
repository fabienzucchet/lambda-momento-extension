import { APIGatewayProxyHandler } from "aws-lambda";
import axios from "axios";

export const handler: APIGatewayProxyHandler = async (event) => {
  try {
    const pk = event.queryStringParameters?.pk;
    const sk = event.queryStringParameters?.sk;

    if (!pk || !sk) {
      return {
        statusCode: 400,
        body: JSON.stringify({
          message: "Both pk and sk query parameters are required.",
        }),
      };
    }

    // Make a request to localhost/cache with pk and sk
    const response = await axios.get(
      `http://localhost:4000/cache?pk=${pk}&sk=${sk}`
    );
    const data = await response.data;

    return {
      statusCode: 200,
      body: JSON.stringify(data),
    };
  } catch (error) {
    console.error(error);
    return {
      statusCode: 500,
      body: JSON.stringify({ message: "Internal server error." }),
    };
  }
};
