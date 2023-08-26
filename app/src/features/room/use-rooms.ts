import { useQuery } from 'react-query';
import { roomsService } from '@/api/rooms-service';

export const useRooms = (page: number) => {
  const res = useQuery({
    queryKey: 'rooms',
    queryFn: () => roomsService.rooms(page),
  });

  return { rooms: res.data?.data || [], ...res };
};
