import React from 'react'
import { useLocation } from 'react-router-dom'

function Player (props) {
  const videoRef = React.createRef()
  const location = useLocation()

  React.useEffect(() => {
    const peerConnection = new RTCPeerConnection() // eslint-disable-line

    peerConnection.ontrack = function (event) {
      videoRef.current.srcObject = event.streams[0]
    }

    peerConnection.addTransceiver('audio')
    peerConnection.addTransceiver('video')

    peerConnection.createOffer().then(offer => {
      peerConnection.setLocalDescription(offer)

      fetch('http://localhost:8080/api/whep', {
        method: 'POST',
        body: offer.sdp,
        headers: {
          Authorization: `Bearer ${location.pathname.substring(1)}`,
          'Content-Type': 'application/sdp'
        }
      }).then(r => {
        return r.text()
      }).then(answer => {
        peerConnection.setRemoteDescription({
          sdp: answer,
          type: 'answer'
        })
      })
    })

    return function cleanup () {
      peerConnection.close()
    }
  }, [videoRef, location.pathname])

  return (
    <video
      ref={videoRef}
      autoPlay
      controls
    />
  )
}

export default Player
