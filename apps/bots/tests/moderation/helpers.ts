import { parseTwitchMessage } from '@twurple/chat';
import type { TwitchPrivateMessage } from '@twurple/chat/lib/commands/TwitchPrivateMessage.js';

type CreateState = {
  mod?: boolean,
  broadcaster?: boolean,
  sub?: boolean,
  vip?: boolean
}

export const createState = (opts?: CreateState) => parseTwitchMessage(`@badge-info=;badges=${opts?.mod && 'moderator/1'},${opts?.sub && 'subscriber/1'},${opts?.broadcaster && 'broadcaster/1'},${opts?.vip && 'vip/1'};client-nonce=31d6f57ff590d3c34eebbcce5812c51c;color=#456073;display-name=Bot_stop;emotes=;first-msg=0;flags=30-34:I.3/P.6;id=977867a2-ae1f-4bc6-9838-674c8e9d98ae;mod=0;reply-parent-display-name=SickestMans;reply-parent-msg-id=bdfb6c55-3553-46dd-bbaa-a942ee512f0b;reply-parent-user-id=786718415;reply-parent-user-login=sickestmans;room-id=128644134;subscriber=${opts?.sub && '1'};tmi-sent-ts=1654697687343;turbo=0;user-id=1;user-type= :bot_stop!bot_stop@bot_stop.tmi.twitch.tv PRIVMSG #sadisnamenya :@SickestMans Стример, или кто пизда 26 градусов haHAA`) as TwitchPrivateMessage;