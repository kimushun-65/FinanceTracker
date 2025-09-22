import React from 'react';
import { UserProfile, LogoutButton } from '../auth';

interface HeaderProps {
  title?: string;
}

export const Header: React.FC<HeaderProps> = ({ title = 'Dashboard' }) => {
  return (
    <header className='border-b border-gray-200 bg-white px-6 py-4'>
      <div className='flex items-center justify-between'>
        <div className='flex items-center space-x-4'>
          <button className='rounded-lg p-2 hover:bg-gray-100'>
            <svg
              className='h-5 w-5 text-gray-600'
              fill='none'
              stroke='currentColor'
              viewBox='0 0 24 24'
            >
              <path
                strokeLinecap='round'
                strokeLinejoin='round'
                strokeWidth={2}
                d='M15 19l-7-7 7-7'
              />
            </svg>
          </button>
          <h1 className='text-2xl font-semibold text-gray-900'>{title}</h1>
        </div>

        <div className='flex items-center space-x-4'>
          <UserProfile />
          <LogoutButton />
        </div>
      </div>
    </header>
  );
};
