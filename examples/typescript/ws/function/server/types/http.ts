export interface HttpResult {
  statusCode: number;
  headers?: Record<string, string>;
  multiValueHeaders?: Record<string, string[]>;
  body?: string;
  isBase64Encoded?: boolean;
}

export interface RequestContext {
  identity: {
    sourceIp: string;
    userAgent: string;
  };
  httpMethod: string;
  requestId: string;
  requestTime: string;
  requestTimeEpoch: number;
}
