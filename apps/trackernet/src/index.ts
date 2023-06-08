import { config } from '@tsuwari/config';
import * as Trackernet from '@tsuwari/grpc/generated/trackernet/trackernet';
import { PORTS } from '@tsuwari/grpc/servers/constants';
import Redis from 'ioredis';
import { ServerError, Status, createServer } from 'nice-grpc';
import { Builder, By, Capabilities, until } from 'selenium-webdriver';

export const redis = new Redis(config.REDIS_URL);


const trackernetServer: Trackernet.TrackernetServiceImplementation = {
  async getRanks(request: Trackernet.GetRanksRequest) {
    if (!request.platform || !request.username) {
      throw new ServerError(Status.INVALID_ARGUMENT, 'Empty platform or username');
    }

    try {
      const cached = await redis.get(`${request.platform}/${request.username}`);
      if (cached) {
        return JSON.parse(cached);
      }

      const caps = Capabilities.chrome();
      const driver = new Builder().usingServer(config.SELENIUM_ADDR).withCapabilities(caps).build();
      await driver.get(
        `https://rocketleague.tracker.network/rocket-league/profile/${request.platform}/${request.username}/overview`,
      );

      await driver.wait(until.elementLocated(By.className('trn-table')), 10000);

      const table = await driver.findElement(By.className('trn-table'));
		
      const rows = await table.findElements(By.tagName('tr'));

      if (!rows) {
        throw new ServerError(Status.INTERNAL, 'Cannot parse tracker.network');
      }

      const rankings: Trackernet.Ranking[] = [];

      for (let i = 1; i < rows.length; i++) {
        const columns = await rows[i]?.findElements(By.tagName('td'));
        if (!columns || columns.length < 6) {
          throw new ServerError(Status.INTERNAL, 'Internal server error');
        }
      
        const rankCol = columns[1];
        const playlist = await rankCol?.findElement(By.className('playlist')).getText();
        const rank = await rankCol?.findElement(By.className('rank')).getText();
        const rating = (await columns[2]?.getText())?.split('\n')[0];
        const ratingStr = rating?.replaceAll(',', '');
        const matches = await columns[5]?.getText();
        const [totalMatches, streak] = matches?.split('\n') || [];

        rankings.push({
          playlist: playlist ?? '',
          rating: ratingStr ? +ratingStr : 0,
          rank: rank ?? '',
          matches: {
            total: totalMatches ? +totalMatches : 0,
            streak: streak ?? '',
          },
        });
      }

      await driver.quit();
      await redis.set(`${request.platform}/${request.username}`, JSON.stringify({
        displayName: request.username,
        rankings: rankings,
      }));
      await redis.expire(`${request.platform}/${request.username}`, 10 * 60);
      console.log({
				displayName: request.username,
				rankings: rankings,
			});
      return {
        displayName: request.username,
        rankings: rankings,
      };
    } catch (e) {
      console.log(e);
      await redis.set(`${request.platform}/${request.username}`, JSON.stringify({displayName: request.username, rankings: [] }));
			await redis.expire(`${request.platform}/${request.username}`, 10 * 60);
      throw new ServerError(Status.INTERNAL, 'Internal server error');
    }
  },
};

const server = createServer({
  'grpc.keepalive_time_ms': 1 * 60 * 1000,
});

server.add(Trackernet.TrackernetDefinition, trackernetServer);

await server.listen(`0.0.0.0:${PORTS.TRACKERNET_SERVER_PORT}`);

process.on('uncaughtException', console.error);
process.on('unhandledRejection', console.error);