import { Controller, Logger } from '@nestjs/common';
import { MessagePattern, Payload } from '@nestjs/microservices';
import { ClientProxyCommands, ClientProxyResult } from '@tsuwari/shared';
import { parseTwitchMessage } from '@twurple/chat';
import { TwitchPrivateMessage } from '@twurple/chat/lib/commands/TwitchPrivateMessage.js';
import { of } from 'rxjs';

import { AppService } from './app.service.js';
import { VariablesParser } from './variables/index.js';

@Controller()
export class AppController {
  private logger = new Logger(AppController.name);

  constructor(
    private readonly service: AppService,
    private readonly variablesParser: VariablesParser,
  ) { }

  @MessagePattern('bots.getVariables')
  async getVariables(): Promise<ClientProxyResult<'bots.getVariables'>> {
    const vars = await import('./variables/modules/index.js');
    const variables = Object.values(vars).map(v => {
      const modules = Array.isArray(v) ? v : [v];

      return modules
        .flat()
        .filter(m => typeof m.visible !== 'undefined' ? m.visible : true)
        .map(m => ({ name: m.key, example: m.example, description: m.description }));
    }).flat();

    return of(variables);
  }

  @MessagePattern('bots.getDefaultCommands')
  async getDefaultCommands(): Promise<ClientProxyResult<'bots.getDefaultCommands'>> {
    const commands = await import('./defaultCommands/index.js');

    return of(Object.values(commands)
      .flat()
      .map(c => ({ name: c.name, permission: c.permission, visible: c.visible ?? true, description: c.description, module: c.module })),
    );
  }

  @MessagePattern('parseResponse')
  async parseResponse(@Payload() data: ClientProxyCommands['parseResponse']['input']) {
    const parsedResponses = await this.service.parseResponses(data, {
      responses: [data.text],
      params: '',
    });

    return parsedResponses;
  }


  @MessagePattern('parseChatMessage')
  async parseChatMessage(@Payload() data: ClientProxyCommands['parseChatMessage']['input']) {
    const state = parseTwitchMessage(data) as TwitchPrivateMessage;
    let message = state.content.value;

    const replyTo = state.tags.get('reply-parent-display-name');

    if (replyTo) {
      message = message.replace(`@${replyTo}`, '').trim();
    }

    this.logger.log(`${state.channelId} | ${state.userInfo.userName}: ${message}`);
    if (!message.startsWith('!') || !state.channelId) return;

    const commandData = await this.service.getResponses(message, state);
    if (!commandData) return;

    const parsedResponses = await this.service.parseResponses({
      channelId: state.channelId,
      text: '',
      userId: state.userInfo.userId,
      userName: state.userInfo.userName,
    }, commandData);
    return parsedResponses;
  }
}
