import type { Settings } from '@twir/grpc/generated/api/api/overlays_be_right_back';

export type SetSettings = (s: Settings) => void;
export type OnStart = (minutes: number, incomingText?: string) => void;
export type OnStop = () => void;
