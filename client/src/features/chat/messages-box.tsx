import React from 'react';
import { useChatStore } from '@/features/chat/use-chat-store';

const Bubble = () => {};

export const MessagesBox = () => {
  const { messages } = useChatStore();

  return <div>{}</div>;
};
