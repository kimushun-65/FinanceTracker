import React from 'react';

const Loading: React.FC = () => {
  return (
    <div className='flex h-screen items-center justify-center'>
      <div className='size-16 animate-spin rounded-full border-t-4 border-blue-600'></div>
    </div>
  );
};

export default Loading;
