import { httpService } from '@/api/http-service';

const AUTH_SERVICE_URL = 'http://localhost:8080/v1';

const signIn = () => {
  return httpService.post(`${AUTH_SERVICE_URL}/sign-in`, {
    id: '67309604673069089',
  });
};

export const authService = {
  signIn,
};
