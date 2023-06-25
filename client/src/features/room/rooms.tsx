import { useRooms } from '@/features/room/use-rooms';
import clsx from 'clsx';
import React, { useState } from 'react';

export const Rooms = () => {
  const [page, setPage] = useState<number>(1);
  const [choose, setChoose] = useState<string | null>(null);
  const { rooms } = useRooms(page);

  return (
    <div>
      {rooms.map((r) => {
        return (
          <div
            key={r.id}
            className={clsx({
              'py-2 cursor-pointer': true,
              'bg-slate-300': choose === r.id,
            })}
            onClick={() => setChoose(r.id)}
          >
            <span className="font-bold p-2">#</span>
            <span className="p-2">{r.name}</span>
          </div>
        );
      })}
    </div>
  );
};
