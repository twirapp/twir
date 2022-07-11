export interface RichPresence {
  status: string;
  version: string;
  time?: string;
  'game:state': string;
  steam_display: string;
  connect: string;
  watching_server?: string,
  param0?: string, // lobby type
  param1?: string,
  param2?: string, // hero
  WatchableGameID?: string,
}