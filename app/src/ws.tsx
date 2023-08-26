'use client';
import React from 'react';
import { useEffect } from 'react';

export const Ws = () => {
  useEffect(() => {
    const ws = new WebSocket('ws://localhost:3000/v1/ws');

    ws.addEventListener('open', (e) => {
      console.log('ws opened');
    });

    ws.addEventListener('message', (e) => {
      console.log('got msg');
      console.log(e);
    });

    setTimeout(() => {
      ws.send('ha');
    }, 5000);
  }, []);

  return <div>WS</div>;
};
