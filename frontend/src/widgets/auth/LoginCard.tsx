import React from 'react';
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  Button,
} from '../../shared/ui';
import { useLogin } from '../../features/auth';

export const LoginCard: React.FC = () => {
  const { login, isLoading } = useLogin();

  return (
    <Card className='mx-auto w-full max-w-md'>
      <CardHeader>
        <CardTitle>Welcome to FinSight</CardTitle>
      </CardHeader>
      <CardContent className='space-y-4'>
        <p className='text-muted-foreground'>
          Please sign in to access your financial dashboard
        </p>
        <Button
          onClick={login}
          disabled={isLoading}
          className='w-full'
          size='lg'
        >
          {isLoading ? 'Signing in...' : 'Sign In'}
        </Button>
      </CardContent>
    </Card>
  );
};
