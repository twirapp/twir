import { staticApi } from '../../bots.js';
import { prisma } from '../../libs/prisma.js';
import { Module } from '../index.js';


const getTop = async (channelId: string, type: 'watched' | 'messages', page: string | number = 1) => {
  const offset = (Number(page) - 1) * 10;
  const limit = 10;

  const stats = await prisma.userStats.findMany({
    where: {
      channelId,
    },
    take: limit,
    skip: offset,
    orderBy: {
      [type]: 'desc',
    },
  });

  const users = await staticApi.users.getUsersByIds(stats.map(s => s.userId));

  return stats.map(stat => {
    const user = users.find(u => u.id === stat.userId);

    if (!user || Number(stat[type]) === 0) return;

    return { displayName: user.displayName, userName: user.name, value: stat[type] };
  });
};

const getPage = (msg?: string) => {
  let page = 1;

  if (!msg) return page;

  if (isNaN(Number(msg))) page = 1;
  else page = Number(msg);
  if (Number(page) <= 0) page = 1;

  return page;
};

export const top: Module[] = [
  {
    key: 'top.messages',
    description: 'Top users by messages',
    handler: async (_, state, params?, message?) => {
      if (!state.channelId) return;

      const page = getPage(message);
      const top = await getTop(state.channelId, 'messages', page);
      console.log(top);
      return top
        .map((u, index) => {
          const name = u?.displayName.toLowerCase() === u?.userName ? u?.displayName : `${u?.userName}`;
          return `${index + 1 + (page - 1) * 10}. ${name} - ${Number(u?.value)}`;
        })
        .join(', ');
    },
  },
];