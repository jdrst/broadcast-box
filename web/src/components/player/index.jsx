import React, { useContext, useEffect, useMemo, useState } from 'react'
import { parseLinkHeader } from '@web3-storage/parse-link-header'
import { useLocation, useSearchParams } from 'react-router-dom'
import ErrorHeader from '../error-header'

export const CinemaModeContext = React.createContext(null);

export function CinemaModeProvider({ children }) {
  const [searchParams] = useSearchParams();
  const cinemaModeInUrl = searchParams.get("cinemaMode") === "true"
  const [cinemaMode, setCinemaMode] = useState(() => cinemaModeInUrl || localStorage.getItem("cinema-mode") === "true")

  const state = useMemo(() => ({
    cinemaMode,
    setCinemaMode,
    toggleCinemaMode: () => setCinemaMode((prev) => !prev),
  }), [cinemaMode, setCinemaMode]);

  useEffect(() => localStorage.setItem("cinema-mode", cinemaMode), [cinemaMode]);
  return (
    <CinemaModeContext.Provider value={state}>
      {children}
    </CinemaModeContext.Provider>
  );
}

function PlayerPage() {
  const { cinemaMode, toggleCinemaMode } = useContext(CinemaModeContext);
  const [peerConnectionDisconnected, setPeerConnectionDisconnected] = React.useState(false)

  return (
    <>
      {peerConnectionDisconnected && <ErrorHeader> WebRTC has disconnected or failed to connect at all 😭 </ErrorHeader>}
      <div className={`flex flex-col items-center ${!cinemaMode && 'mx-auto px-2 py-2 container'}`}>
        <Player cinemaMode={cinemaMode} peerConnectionDisconnected={peerConnectionDisconnected} setPeerConnectionDisconnected={setPeerConnectionDisconnected} />
        <button className='bg-blue-900 px-4 py-2 rounded-lg mt-6' onClick={toggleCinemaMode}>
          {cinemaMode ? "Disable cinema mode" : "Enable cinema mode"}
        </button>
      </div>
    </>
  )
}

function Player({ cinemaMode, peerConnectionDisconnected, setPeerConnectionDisconnected }) {
  const videoRef = React.createRef()
  const location = useLocation()
  const [videoLayers, setVideoLayers] = React.useState([]);
  const [mediaSrcObject, setMediaSrcObject] = React.useState(null);
  const [layerEndpoint, setLayerEndpoint] = React.useState('');

  const onLayerChange = event => {
    fetch(layerEndpoint, {
      method: 'POST',
      body: JSON.stringify({ mediaId: '1', encodingId: event.target.value }),
      headers: {
        'Content-Type': 'application/json'
      }
    })
  }

  React.useEffect(() => {
    if (videoRef.current) {
      videoRef.current.srcObject = mediaSrcObject
    }
  }, [mediaSrcObject, videoRef])

  React.useEffect(() => {
    const peerConnection = new RTCPeerConnection() // eslint-disable-line

    peerConnection.ontrack = function (event) {
      setMediaSrcObject(event.streams[0])
    }

    peerConnection.oniceconnectionstatechange = () => {
      if (peerConnection.iceConnectionState === 'connected' || peerConnection.iceConnectionState === 'completed') {
        setPeerConnectionDisconnected(false)
      } else if (peerConnection.iceConnectionState === 'disconnected' ||  peerConnection.iceConnectionState === 'failed') {
        setPeerConnectionDisconnected(true)
      }
    }

    peerConnection.addTransceiver('audio', { direction: 'recvonly' })
    peerConnection.addTransceiver('video', { direction: 'recvonly' })

    peerConnection.createOffer().then(offer => {
      offer["sdp"] = offer["sdp"].replace("useinbandfec=1", "useinbandfec=1;stereo=1")
      peerConnection.setLocalDescription(offer)

      const apiPath = import.meta.env.VITE_API_PATH ?? (() => {
      console.warn('[broadcast box] REACT_APP_API_PATH is deprecated, please use VITE_API_PATH instead');
      return import.meta.env.REACT_APP_API_PATH;
      })();

      fetch(`${apiPath}/whep/${location.pathname.split('/').pop()}/`, {
        method: 'POST',
        body: offer.sdp,
        headers: {
          'Content-Type': 'application/sdp'
        }
      }).then(r => {
        console.log(`fetched: ${apiPath}/whep/${location.pathname.split('/').pop()}`)
        const parsedLinkHeader = parseLinkHeader(r.headers.get('Link'))
        setLayerEndpoint(`${window.location.protocol}//${parsedLinkHeader['urn:ietf:params:whep:ext:core:layer'].url}`)

        const evtSource = new EventSource(`${window.location.protocol}//${parsedLinkHeader['urn:ietf:params:whep:ext:core:server-sent-events'].url}`)
        evtSource.onerror = err => evtSource.close();

        evtSource.addEventListener("layers", event => {
          const parsed = JSON.parse(event.data)
          setVideoLayers(parsed['1']['layers'].map(l => l.encodingId))
        })


        return r.text()
      }).then(answer => {
        peerConnection.setRemoteDescription({
          sdp: answer,
          type: 'answer'
        })
      })
    })

    return function cleanup() {
      peerConnection.close()
    }
  }, [location.pathname, setPeerConnectionDisconnected])

  return (
    <>
      <video
        ref={videoRef}
        autoPlay
        muted
        controls
        playsInline
        className={`bg-black w-full ${cinemaMode && "h-full"}`}
        style={cinemaMode ? {
          maxHeight: '100vh',
          maxWidth: '100vw'
        } : {}}
      />

      {videoLayers.length >= 2 &&
        <select defaultValue="disabled" onChange={onLayerChange} className="appearance-none border w-full py-2 px-3 leading-tight focus:outline-hidden focus:shadow-outline bg-gray-700 border-gray-700 text-white rounded-sm shadow-md placeholder-gray-200">
          <option value="disabled" disabled={true}>Choose Quality Level</option>
          {videoLayers.map(layer => {
            return <option key={layer} value={layer}>{layer}</option>
          })}
        </select>
      }
    </>
  )
}

export default PlayerPage
