import { useRooms } from '@/features/room/use-rooms';
import clsx from 'clsx';
import React, { useEffect, useState } from 'react';

const Ws = () => {
  const [socket, setSocket] = useState<WebSocket | null>(null);

  useEffect(() => {
    const ws = new WebSocket('ws://localhost:8080/v1/ws');
    setSocket(ws);
  }, []);

  useEffect(() => {
    if (socket) {
      socket.addEventListener('open', () => console.log('open connection'));

      socket.addEventListener('message', (ev) => {
        if (ev.data === 'ping') {
          socket.send('pong');
        }
      });
    }
  }, [socket]);

  return (
    <button
      onClick={() => {
        socket?.send(
          JSON.stringify({
            type: 'chat-room/send-message',
            payload: { roomId: '123123123', text: 'some text' },
          }),
        );
      }}
    >
      send message comment
    </button>
  );
};

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
      <Ws />
    </div>
  );
};
