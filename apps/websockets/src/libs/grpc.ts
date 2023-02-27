import * as Websocket from '@tsuwari/grpc/generated/websockets/websockets';
import { PORTS } from '@tsuwari/grpc/servers/constants';
import { createServer } from 'nice-grpc';

import {
  onAudioDecrease,
  onAudioDisable,
  onAudioEnable,
  onAudioIncrease,
  onSetAudio,
  onSetScene,
  onToggleAudioSource,
  onToggleSource,
  onStartStream,
  onStopStream,
} from '../namespaces/obs.js';
import { onTtsSay, onTtsSkip } from '../namespaces/tts.js';
import { onAddRequest, onRemoveRequest } from '../namespaces/youtube.js';

export const grpcServer = createServer({
  'grpc.keepalive_time_ms': 1 * 60 * 1000,
});

const websocketService: Websocket.WebsocketServiceImplementation = {
  youtubeAddSongToQueue: onAddRequest,
  youtubeRemoveSongToQueue: onRemoveRequest,
  obsSetScene: onSetScene,
  obsToggleSource: onToggleSource,
  obsToggleAudio: onToggleAudioSource,
  obsAudioSetVolume: onSetAudio,
  obsAudioDecreaseVolume: onAudioDecrease,
  obsAudioIncreaseVolume: onAudioIncrease,
  obsAudioDisable: onAudioDisable,
  obsAudioEnable: onAudioEnable,
  obsStartStream: onStartStream,
  obsStopStream: onStopStream,
  textToSpeechSay: onTtsSay,
  textToSpeechSkip: onTtsSkip,
};

grpcServer.add(Websocket.WebsocketDefinition, websocketService);

export const listen = () => grpcServer.listen(`0.0.0.0:${PORTS.WEBSOCKET_SERVER_PORT}`);
