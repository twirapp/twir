'use client';

import { APITypes, PlyrInstance, PlyrProps, usePlyr } from 'plyr-react';
import React, { forwardRef, useEffect } from 'react';

type Props = PlyrProps & {
  onReady?: () => void
  onCanPlay?: () => void
}

export const CustomPlyrInstance = forwardRef<APITypes, Props>(
  (props, ref) => {
    const { source, options = null } = props;
    const raptorRef = usePlyr(ref, { options, source });

    // Do all api access here, it is guaranteed to be called with the latest plyr instance
    useEffect(() => {
      /**
       * Fool react for using forward ref as normal ref
       * NOTE: in a case you don't need the forward mechanism and handle everything via props
       * you can create the ref inside the component by yourself
       */
      const { current } = ref as React.MutableRefObject<APITypes>;
      if (current.plyr.source === null) return;

      const api = current as { plyr: PlyrInstance };
      // api.plyr.on('ready', () => console.log('I\'m ready'));
      // api.plyr.on('canplay', () => {
      //   // NOTE: browser may pause you from doing so:  https://goo.gl/xX8pDD
      //   api.plyr.play();
      //   console.log('duration of audio is', api.plyr.duration);
      // });
      // api.plyr.on('ended', () => console.log('I\'m Ended'));
      api.plyr.on('ready', () => props.onReady ? props.onReady() : null);
      api.plyr.on('canplay', () => props.onCanPlay ? props.onCanPlay() : null);
    });

    return (
      <video
        ref={raptorRef as React.MutableRefObject<HTMLVideoElement>}
        className="plyr-react plyr"
      />
    );
  },
);