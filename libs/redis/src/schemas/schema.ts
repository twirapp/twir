import { Entity, Schema } from 'redis-om';

class Album extends Entity {
  artist: string;
  genres: string[];
}

export const streamS = new Schema(Album, {
  id: { type: 'string' },
  artist: { type: 'string', indexed: true },
  genres: { type: 'string[]' },
}, {
  prefix: 'albumSchemaPrefix',
  indexedDefault: true,
});
