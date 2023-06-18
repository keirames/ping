'use client';
import React, { useState } from 'react';
import { MessagesBox } from '@/features/chat/messages-box';
import { useChatStore } from '@/features/chat/use-chat-store';

export const Main = () => {
  const [value, setValue] = useState<string>('');
  const { send } = useChatStore();

  return (
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
  );
};
