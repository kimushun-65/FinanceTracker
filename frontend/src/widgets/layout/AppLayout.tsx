'use client';

import React from 'react';
import { Sidebar } from './Sidebar';
import { Header } from './Header';

interface AppLayoutProps {
  children: React.ReactNode;
  title?: string;
}

export const AppLayout: React.FC<AppLayoutProps> = ({ children, title }) => {
  return (
    <div className='flex h-screen bg-gray-50'>
      <Sidebar />
      <div className='flex flex-1 flex-col overflow-hidden'>
        <Header title={title} />
        <main className='flex-1 overflow-auto p-6'>{children}</main>
      </div>
    </div>
  );
};
