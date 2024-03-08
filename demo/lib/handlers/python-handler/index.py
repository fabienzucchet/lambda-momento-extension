import json
import requests

def handler(event, context):
    try:
        params = event.get('queryStringParameters', {})
        pk = params.get('pk')
        sk = params.get('sk')

        if not pk or not sk:
            return {
                'statusCode': 400,
                'body': json.dumps({
                    'message': 'Both pk and sk query parameters are required.'
                })
            }

        # Make a request to localhost/cache with pk and sk
        response = requests.get(f'http://localhost:4000/cache?pk={pk}&sk={sk}')
        data = response.json()

        return {
            'statusCode': 200,
            'body': json.dumps(data)
        }

    except Exception as e:
        print(e)
        return {
            'statusCode': 500,
            'body': json.dumps({'message': 'Internal server error.'})
        }
