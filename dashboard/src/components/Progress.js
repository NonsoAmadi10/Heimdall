import React from 'react';

const CircularProgress = ({ progress, symbol,color, value }) => {
  

  return (
    <div className="w-full bg-neutral-200 dark:bg-neutral-600">
    <div
      className={`${color} p-2 text-center text-sm font-medium leading-none text-primary-100`}
      style={{ width: `${progress}%` }}
    >
      { value }{" "}{symbol}
    </div>
  </div>
  
  );
};



export default CircularProgress;
