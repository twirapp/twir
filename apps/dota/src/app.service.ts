import { dirname, resolve } from 'path';
import { fileURLToPath } from 'url';

import { Injectable, OnModuleInit } from '@nestjs/common';
import { config } from '@tsuwari/config';
import { RedisService } from '@tsuwari/shared';
import protobufjs from 'protobufjs';
import SteamUser from 'steam-user';

@Injectable()
export class AppService extends SteamUser implements OnModuleInit {
  #watchRoot: protobufjs.Root;
  #clientHelloRoot: protobufjs.Root;

  constructor(private readonly redis: RedisService) {
    super({
      autoRelogin: true,
    });
  }

  async onModuleInit() {
    this.#watchRoot = await new protobufjs.Root().load(resolve(dirname(fileURLToPath(import.meta.url)), '..', 'protos', 'dota2', 'dota_gcmessages_client_watch.proto'), {
      keepCase: true,
    });

    this.#clientHelloRoot = await new protobufjs.Root().load(resolve(dirname(fileURLToPath(import.meta.url)), '..', 'protos', 'dota2', 'gcsdk_gcmessages.proto'), {
      keepCase: true,
    });

    this.logOn({
      accountName: config.STEAM_USERNAME,
      password: config.STEAM_PASSWORD,
    });

    this.on('loggedOn', (e) => {
      console.log(`Logged into Steam as ${e.client_supplied_steamid} ${e.vanity_url}`);
      this.setPersona(1);
      this.gamesPlayed([570], true);
    });
  }

}