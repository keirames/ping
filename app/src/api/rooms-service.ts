import { httpService } from '@/api/http-service';

const ROOMS_SERVICE_URL = 'http://localhost:8080/v1';

export interface Room {
  id: string;
  name: string;
}

export interface Paginated<T> {
  page: number;
  limit: number;
  data: T;
}

const rooms = async (page: number) => {
  const res = await httpService.get<Paginated<Room[]>>(
    `${ROOMS_SERVICE_URL}/rooms`,
    {
      params: { page },
    },
  );

  return res.data;
};

export const roomsService = {
  rooms,
};
