'use client';
import React, { useEffect, useState } from 'react';
import { MessagesBox } from '@/features/chat/messages-box';
import { useChatStore } from '@/features/chat/use-chat-store';
import {
  QueryClient,
  QueryClientProvider,
  useMutation,
  useQuery,
} from 'react-query';
import axios from 'axios';

const queryClient = new QueryClient();

const Auth = (props: { set: any }) => {
  const { data, mutate } = useMutation('auth', () =>
    axios.post('http://localhost:3000/v1/sign-in', { id: '1231232' }),
  );

  useEffect(() => {
    if (data) {
      props.set();
    }
  }, [data, props]);

  return (
    <div>
      <button
        onClick={() => {
          mutate();
        }}
      >
        sign in
      </button>
    </div>
  );
};

export const Main = () => {
  const [value, setValue] = useState<string>('');
  const [isAuth, setIsAuth] = useState<boolean>(false);
  const { send } = useChatStore();

  return (
    <QueryClientProvider client={queryClient}>
      {!isAuth ? (
        <Auth set={() => setIsAuth(true)} />
      ) : (
        <div className="w-full h-full flex flex-col">
          <div className="h-full overflow-x-hidden overflow-y-scroll">
            <MessagesBox />
          </div>
          <div className="w-full">
            <input
              className="w-full dark:bg-black bg-black border-blue-100"
              value={value}
              onChange={(e) => {
                setValue(e.currentTarget.value);
              }}
              onKeyDown={(e) => {
                if (e.key === 'Enter' && value !== '') {
                  setValue('');
                  send(value);
                }
              }}
            />
          </div>
        </div>
      )}
    </QueryClientProvider>
  );
};
