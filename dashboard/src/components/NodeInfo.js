import React, { useEffect, useState } from 'react'
import CircularProgress from './Progress';

function NodeInfo() {
  const [info, setInfo] = useState({});
  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await fetch('http://localhost:1700/node-info');
        const json = await response.json();
        setInfo(json);
      } catch (error) {
        console.error('Error fetching data:', error);
      }
    };
  
    fetchData();
  }, []);
 
  return (
    <section className="w-full">
  <div className="mx-auto max-w-screen-xl px-4 py-12 sm:px-6 md:py-16 lg:px-8">
    <div className="mx-auto max-w-3xl text-center">
      <h2 className="text-3xl font-bold text-gray-700 sm:text-4xl">
        Your Node Info:
      </h2>
    </div>
    { Object.keys(info).length !== 0 ?
    <div>
    <div className="mt-8 sm:mt-12">
      
    <div
          className="flex flex-col rounded-lg border border-gray-100 px-4 py-8 text-center"
        >
            <dt className="order-last text-lg py-2font-medium text-gray-500">
            Lightning Public Key
          </dt>

          <dd className="text-2xl font-extrabold text-blue-600 md:text-xl">
          { info["lightning"]["pub_key"]}
          </dd>
            </div>
    </div>

    <div className="mt-8 sm:mt-12">
      <dl className="grid grid-cols-1 gap-4 sm:grid-cols-3">
        <div
          className="flex flex-col rounded-lg border border-gray-100 px-4 py-8 text-center"
        >
          <dt className="order-last text-lg font-medium text-gray-500">
            Newtork Capacity
          </dt>

          <dd className="text-2xl font-extrabold text-blue-600 md:text-3xl">
          { info["lightning"]["network_capacity"]}
          </dd>
        </div>

        <div
          className="flex flex-col rounded-lg border border-gray-100 px-4 py-8 text-center"
        >
          <dt className="order-last text-lg font-medium text-gray-500">
            Network Difficulty
          </dt>

          <dd className="text-4xl font-extrabold text-blue-600 md:text-3xl">1</dd>
        </div>

        <div
          className="flex flex-col rounded-lg border border-gray-100 px-4 py-8 text-center"
        >
          <dt className="order-last text-lg font-medium text-gray-500">
            Chain
          </dt>

          <dd className="text-4xl font-extrabold text-blue-600 md:text-3xl">{ info["bitcoin"]["chain"]}net</dd>
        </div>
        <div
          className="flex flex-col rounded-lg border border-gray-100 px-4 py-8 text-center"
        >
          <dt className="order-last text-lg font-medium text-gray-500">
            Bitcoin Client Version
          </dt>

          <dd className="text-4xl font-extrabold text-blue-600 md:text-3xl">{ (info["bitcoin"]["version"] / 10000)}</dd>
        </div>
        <div
          className="flex flex-col rounded-lg border border-gray-100 px-4 py-8 text-center"
        >
          <dt className="order-last text-lg font-medium text-gray-500">
            Number of Blocks
          </dt>

          <dd className="text-4xl font-extrabold text-blue-600 md:text-3xl">{ info["bitcoin"]["no_of_blocks"]}</dd>
        </div>
        <div
          className="flex flex-col rounded-lg border border-gray-100 px-4 py-8 text-center"
        >
          <dt className="order-last text-lg font-medium text-gray-500">
            User Agent
          </dt>

          <dd className="text-4xl font-extrabold text-blue-600 md:text-3xl">{ (info["bitcoin"]["user_agent"]).replace(/\//g, "")}</dd>
        </div>
        
      </dl>
    </div>

    <div className="mx-auto max-w-screen-xl grid gap-x-4 grid-cols-2 py-8">
      
      <div className="w-full py-8 px-8 rounded-lg flex flex-col h-auto border border-gray-100"> 
      <p className="text-gray-500 py-8"> Node Hashrate</p>

      <CircularProgress progress={60} symbol={"TH/s"} color={"bg-blue-600"} value={Math.floor(info["bitcoin"]["hash_rate"] * (10**12))}/>
      
      </div>
      <div className="w-full py-4 px-8 rounded-lg flex flex-col h-auto border border-gray-100">
         <p className="text-gray-500 py-8"> Block Propagation Time</p>
         <CircularProgress progress={80} symbol={"mins"} color={"bg-red-600"} value={Math.floor(info["bitcoin"]["block_propagation"])}/>
      </div>
    </div>
    </div>
    : ""
}
  </div>
</section>
  )
}

export default NodeInfo