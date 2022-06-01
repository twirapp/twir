export interface ServerToClientEvents {
  noArg: () => void;
  basicEmit: (a: number, b: string, c: Buffer) => void;
  withAck: (d: string, callback: (e: number) => void) => void;
}

export interface ClientToServerEvents {
  example: (d: string, cb?: (e: number) => void) => void;
  isBotMod: (opts: { channelId: string, channelName: string, userId: string }, cb?: (value: { channelId: string, value: boolean }) => void) => void;
  botJoin: (opts: { channelName: string, channelId: string, }, cb?: () => void) => void
  botPart: (opts: { channelName: string, channelId: string, }, cb?: () => void) => void
}

export interface SocketData {
  name: string;
  age: number;
}


export interface EventsMap {
  [event: string]: any;
}
type EventNames<Map extends EventsMap> = keyof Map & (string | symbol);
export type EventParams<
  Map extends EventsMap,
  Ev extends EventNames<Map>
  > = Parameters<Map[Ev]>;