import { ChannelCustomvar } from '@tsuwari/typeorm/entities/ChannelCustomvar';
import { ref } from 'vue';

export const variables = ref<ChannelCustomvar[]>([{ 'id':'b5d20375-cf44-4c6a-83d3-1639e197511a', 'name':'qwe', 'description':'', 'type':'SCRIPT', 'evalValue':'return Date.now();', 'response':null, 'channelId':'128644134' }] as any);