import { HelixStream } from '@twurple/api';
import { rawDataSymbol } from '@twurple/common';

export type CachedStream = HelixStream[typeof rawDataSymbol] & { parsedMessages?: number }