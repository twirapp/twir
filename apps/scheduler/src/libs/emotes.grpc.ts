import { config } from '@tsuwari/config';
import { createEmotesCacher } from '@tsuwari/grpc/clients/emotes_cacher';

export const emotesCacherGrpcClient = await createEmotesCacher(config.NODE_ENV);
