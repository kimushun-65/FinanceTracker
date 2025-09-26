import { APIGatewayProxyEvent, APIGatewayProxyResult, Context } from 'aws-lambda';

const API_BASE_URL = process.env.API_BASE_URL || 'http://host.docker.internal:8080';

export const handler = async (
  event: APIGatewayProxyEvent,
  context: Context
): Promise<APIGatewayProxyResult> => {
  try {
    // Health check endpoint
    if (event.path === '/v1/health') {
      return {
        statusCode: 200,
        headers: {
          'Content-Type': 'application/json',
          'Access-Control-Allow-Origin': '*',
        },
        body: JSON.stringify({
          status: 'healthy',
          timestamp: new Date().toISOString(),
          service: 'FinanceTracker API',
          environment: process.env.ENVIRONMENT || 'unknown',
        }),
      };
    }
    // APIエンドポイントを構築（/v1/auth/* → /api/v1/auth/*）
    const apiPath = `/api${event.path}`;
    const url = new URL(apiPath, API_BASE_URL);

    // クエリパラメータを追加
    if (event.queryStringParameters) {
      Object.entries(event.queryStringParameters).forEach(([key, value]) => {
        if (value) url.searchParams.append(key, value);
      });
    }

    // ヘッダーを準備（認証エンドポイントは特別な処理が必要かもしれない）
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      ...event.headers
    };

    // Auth0 IDを追加（認証エンドポイントでは必要ないかもしれない）
    const auth0Id = event.requestContext.authorizer?.userId;
    if (auth0Id && event.path !== '/v1/auth/callback' && event.path !== '/v1/auth/refresh') {
      headers['X-Auth0-ID'] = auth0Id;
    }

    // HTTP APIを呼び出し
    const response = await fetch(url.toString(), {
      method: event.httpMethod,
      headers,
      body: event.body || undefined,
    });

    // レスポンスボディを取得
    const responseBody = await response.text();

    // レスポンスヘッダーを変換
    const responseHeaders: Record<string, string> = {};
    response.headers.forEach((value, key) => {
      responseHeaders[key] = value;
    });

    return {
      statusCode: response.status,
      headers: responseHeaders,
      body: responseBody,
    };
  } catch (error) {
    console.error('Proxy error:', error);
    
    return {
      statusCode: 503,
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        success: false,
        error: 'Service unavailable',
      }),
    };
  }
};