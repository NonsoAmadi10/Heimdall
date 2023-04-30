import React from 'react'
import { Chart as ChartJS, ArcElement, Tooltip, Legend } from 'chart.js';
import { Doughnut } from 'react-chartjs-2';

ChartJS.register(ArcElement, Tooltip, Legend);



function PeerChart({ data }) {
    console.log(data)
    const dataset = {
        labels: ['Lightning',  'Bitcoin'],
        datasets: [
          {
            label: 'Number of Peers',
            data: data,
            backgroundColor: [
                'rgba(	59, 130, 246, 0.2)',
              'rgba(255, 206, 86, 0.2)',
            
            ],
            borderColor: [
                'rgba(54, 162, 235, 1)',
    
              'rgba(255, 206, 86, 1)',
            
            ],
            borderWidth: 1,
          },
        ],
      };
  return (
    <div className="w-3/6  p-6">
        <p className="text-gray-500 py-8 ">Number of Peer Connections:</p>
        <Doughnut data={dataset} />
    </div>
  )
}

export default PeerChart