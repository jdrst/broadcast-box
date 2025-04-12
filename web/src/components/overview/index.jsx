import { useEffect, useState } from "react"

const OverviewPage = () => {
  const refreshStreamKeys = () => {
    fetch('api/status')
      .then(resp => {
        resp.json().then(streams => {
          if (streams.length > 0){
            setStreamKeys(streams.filter(s => s.videoStreams.length > 0).map(s => s.streamKey))
          }
        })   
      })
  }
  
  const [streamKeys, setStreamKeys] = useState([])
  useEffect(() => {
    refreshStreamKeys()
    const interval = setInterval(refreshStreamKeys, 1000*60)
    return () => {clearInterval(interval)}
  }, []) 
  
  return (<>
        {
          streamKeys.map(key => 
            <a href={key} key={key}>{key}</a> 
          )
        }
   </>)
}


export default OverviewPage
