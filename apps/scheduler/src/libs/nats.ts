import { config } from '@tsuwari/config';
import * as Scheduler from '@tsuwari/nats/scheduler';
import { connect } from 'nats';

import { DefaultCommandsCreatorService } from '../default-commands-creator/default-commands-creator.service.js';
import { app } from '../index.js';

export const nats = await connect({
  servers: [config.NATS_URL],
});

export async function listenForDefaultCommands() {
  for await (const event of nats.subscribe(Scheduler.SUBJECTS.CREATE_DEFAULT_COMMANDS)) {
    const data = Scheduler.CreateDefaultCommandsRequest.fromBinary(event.data);
    const service = app.get(DefaultCommandsCreatorService);
    service.createDefaultCommands([data.userId]);
    event.respond();
  }
}
