import { atom, createStore } from 'jotai';
import OBSWebSocket from 'obs-websocket-js';
import { Socket } from 'socket.io-client';

export const obsStore = createStore();

export class MyOBSWebsocket extends OBSWebSocket {
  connected = false;

  async connect(url: string, password: string) {
    const connection = await super.connect(url, password);
    this.connected = true;
    return connection;
  }

  async disconnect() {
    await super.disconnect();
    this.connected = false;
  }
}

export const externalObsWsAtom = atom<MyOBSWebsocket | null>(null);
export const internalObsWsAtom = atom<Socket | null>(null);