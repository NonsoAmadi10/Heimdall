import React, { useState, useEffect } from 'react'
import PeerChart from './PeerChart'
import Bandwidth from './Bandwidth';
import moment from 'moment';


function Conn() {
    const [metrics, setMetrics] = useState([]);
    const [peers, setPeers] = useState([]);
    const [bandWidthIn, setBandIn] = useState([])
    const [bandWidthOut, setBandOut] = useState([])
    const [labels, setLabel] = useState([])

    useEffect(() => {
        const fetchData = async () => {
          try {
            const response = await fetch('http://localhost:1700/conn-metrics');
            const json = await response.json();
            setMetrics(json.data);
            const num_of_btc_peers = [...new Set(json.data.map(obj => obj.num_of_btc_peers))];
            const num_of_lnd_peers = [...new Set(json.data.map(obj => obj.num_lnd_peers))];
            const bandIn = json.data.map(obj => obj.btc_bandwidth_in).slice(-12)
            const bandOut = json.data.map(obj => obj.btc_bandwidth_out).slice(-12)
            const label = json.data.map( time =>moment(time.timestamp).fromNow()).slice(-12)

            setPeers([num_of_lnd_peers[num_of_lnd_peers.length - 1], num_of_btc_peers[num_of_btc_peers.length - 1]])
            setBandIn(bandIn)
            setBandOut(bandOut)
            setLabel(label)
          } catch (error) {
            console.error('Error fetching data:', error);
          }
        };
      
        fetchData();
      }, []);

     
     
  return (
    <div className="mx-auto max-w-screen-xl  flex px-4 py-12 sm:px-6 md:py-16 lg:px-8">
        <PeerChart data={peers} />
        <Bandwidth dataset={bandWidthIn} datasetOut={bandWidthOut} time={labels}/>
    </div>
  )
}

export default Conn