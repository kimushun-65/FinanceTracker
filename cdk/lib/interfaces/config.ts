export interface EnvironmentConfig {
  environment: string;
  region: string;
  customDomain?: string;
  auth0Domain: string;
  auth0Audience: string;
  auth0ClientId: string;
  githubOwner: string;
  repositoryName: string;
  databaseConfig: DatabaseConfig;
  lambdaConfig: LambdaConfig;
  sesConfig: SesConfig;
}

export interface DatabaseConfig {
  instanceType: string;
  multiAz: boolean;
  deletionProtection: boolean;
}

export interface LambdaConfig {
  memorySize: number;
  timeout: number;
}

export interface SesConfig {
  fromEmail: string;
  sendingQuota: number;
  sendingRate: number;
}