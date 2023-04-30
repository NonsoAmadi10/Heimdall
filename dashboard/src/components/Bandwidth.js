import React from 'react'

import {
    Chart as ChartJS,
    CategoryScale,
    LinearScale,
    BarElement,
    Title,
    Tooltip,
    Legend,
  } from 'chart.js';
  import { Bar } from 'react-chartjs-2';
 
  
  ChartJS.register(
    CategoryScale,
    LinearScale,
    BarElement,
    Title,
    Tooltip,
    Legend
  );

function Bandwidth({ dataset, time, datasetOut }) {

    const options = {
        indexAxis: 'y',
        elements: {
          bar: {
            borderWidth: 4,
          },
        },
        responsive: true,
        plugins: {
          legend: {
            position: 'right',
          },
          title: {
            display: true,
            text: 'Network Bandwidth',
          },
        },
      };

      const labels = [...time];

      const data = {
        labels,
        datasets: [
          {
            fill: true,
            label: 'Bandwidth In (Bytes)',
            data: dataset,
            borderColor: 'rgb(255, 99, 132)',
            backgroundColor: 'rgba(255, 99, 132, 0.5)',
            tension: 0.1
          },
          {
            fill: true,
            label: 'Bandwidth Out (Bytes)',
            data: datasetOut,
            borderColor: 'rgb(53, 162, 235)',
            backgroundColor: 'rgba(53, 162, 235, 0.5)',
            tension: 0.1
          },
        ],
      };
      
  return (
    <div className="w-3/6 p-8 ">
    <p className="text-gray-500 py-8">Network Bandwith:</p>
    
    <div className='mx-auto flex items-center justify-center'>
    <Bar options={options} data={data} />
    </div>
</div>
  )
}

export default Bandwidth