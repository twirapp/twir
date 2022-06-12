import { staticApi } from '../../bots.js';
import { prisma } from '../../libs/prisma.js';
import { DefaultCommand } from '../types.js';


export const permit: DefaultCommand = {
  name: 'permit',
  permission: 'MODERATOR',
  handler: async (state, params) => {
    if (!params || !state.channelId) return;
    const paramsArray = params.split(' ');
    const userName = paramsArray[0];
    const count = paramsArray[1];

    if (!userName) return 'you have type user name to permit.';

    const user = await staticApi.users.getUserByName(userName!);
    if (!user) return `user with name ${userName} not found.`;

    const parsedCount = count ? isNaN(parseInt(count, 10)) ? 1 : parseInt(count, 10) : 1;
    if (parsedCount > 100) return 'cannot create more then 100 permits.';

    console.log(user.id, state.channelId);

    await prisma.$transaction([...Array(parsedCount)].map(() => prisma.permit.create({
      data: {
        userId: user.id,
        channelId: state.channelId!,
      },
    })));

    return `you gave out ${parsedCount} ${parsedCount > 1 ? 'permits' : 'permit'} to user ${user.displayName}`;
  },
};