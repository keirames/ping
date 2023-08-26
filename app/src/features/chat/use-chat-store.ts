import { create } from 'zustand';
import { devtools } from 'zustand/middleware';

interface ChatState {
  messages: string[];
  send: (m: string) => void;
}

export const useChatStore = create<ChatState>()(
  devtools((set) => ({
    messages: [
      'qwdqwdqwd',
      'qwdqwdqwd',
      'qwdqwdqwd',
      'qwdqwdqwd',
      'qwdqwdqwd',
      'qwdqwdqwd',
      'qwdqwdqwd',
      'qwdqwdqwd',
      'qwdqwdqwd',
      'qwdqwdqwd',
      'qwdqwdqwd',
      'qwdqwdqwd',
      'qwdqwdqwd',
      'qwdqwdqwd',
      'qdqwdqw',
      'qwdqwd',
      'qwdqwd',
      'qwdqwdqwdqwdqwd',
      'qwdqwdqwdqwdqwd',
      'qwdqwdqwdqwdqwd',
      'qwdqwdqwdqwdqwd',
      'qwdqwdqwdqwdqwd',
      'qwdqwdqwdqwdqwd',
      'qwdqwdqwdqwdqwd',
      'qwdqwdqwdqwdqwd',
      'qwdqwdqwdqwdqwd',
      'qwdqwdqwdqwdqwd',
      'qwdqwdqwdqwdqwd',
      'qwdqwdqwdqwdqwd',
      'qwdqwdqwdqwdqwd',
      'qwdqwdqwdqwdqwd',
      'qwdqwdqwdqwdqwd',
      'qwdqwdqwdqwdqwd',
      'qwdqwdqwdqwdqwd',
      'qwdqwdqwdqwdqwd',
      'qwdqwdqwdqwdqwd',
      'qwdqwdqwdqwdqwd',
      'qwdqwdqwdqwdqwd',
      'qwdqwdqwdqwdqwd',
      'qwdqwdqwdqwdqwd',
      'qwdqwdqwdqwdqwd',
      'qwdqwdqwdqwdqwd',
      'qwdqwdqwdqwdqwd',
      'qwdqwdqwdqwdqwd',
      'qwdqwdqwdqwdqwd',
      'qwdqwdqwdqwdqwd',
      'qwdqwdqwdqwdqwd',
      'qwdqwdqwdqwdqwd',
      'qwdqwdqwdqwdqwd',
      'qwdqwdqwdqwdqwd',
      'qwdqwdqwdqwdqwd',
    ],
    send: (m) => {
      set((state) => ({ messages: [...state.messages, m] }));
    },
  })),
);
