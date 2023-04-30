import React from 'react';
import VideoJS from '../components/videojs'
import videojs from 'video.js';
import {Container} from "@mui/material";

function Video(props) {
    const playerRef = React.useRef(null);

    const videoJsOptions = {
        autoplay: true,
        controls: true,
        responsive: true,
        fluid: true,
        sources: [{
            src: '/api/video/downloadM3u8?id=1',
            type: 'application/x-mpegURL'
        }]
    };

    const handlePlayerReady = (player) => {
        playerRef.current = player;

        // You can handle player events here, for example:
        player.on('waiting', () => {
            videojs.log('player is waiting');
        });

        player.on('dispose', () => {
            videojs.log('player will dispose');
        });
    };

    return (
        <Container>
            <VideoJS options={videoJsOptions} onReady={handlePlayerReady} />
        </Container>
    );
}

export default Video;