import type { Settings } from '@twir/grpc/generated/api/api/overlays_kappagen';
import type { Emote, KappagenAnimations } from 'kappagen';

export type KappagenSettings = Settings & { channelName: string, channelId: string };

export type KappagenCallback = (emotes: Emote[], animation: KappagenAnimations) => void;
export type SpawnCallback = (emotes: Emote[]) => void;
export type SetSettingsCallback = (settings: KappagenSettings) => void;
