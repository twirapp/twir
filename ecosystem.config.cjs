module.exports = {
  apps: [
    {
      name: 'api',
      script: 'pnpm',
      args: 'dev:api',
    },
    {
      name: 'bots',
      script: 'pnpm',
      args: 'dev:bots',
    },
    {
      name: 'dota',
      script: 'pnpm',
      args: 'dev:dota',
    },
    {
      name: 'eventsub',
      script: 'pnpm',
      args: 'dev:eventsub',
    },
    {
      name: 'scheduler',
      script: 'pnpm',
      args: 'dev:scheduler',
    },
    {
      name: 'streamstatus',
      script: 'pnpm',
      args: 'dev:streamstatus',
    },
    {
      name: 'eval',
      script: 'pnpm',
      args: 'dev:eval',
    },
    {
      name: 'parser',
      script: 'pnpm',
      args: 'dev:parser',
    },
    // {
    //   name: "frontend",
    //   script: "pnpm",
    //   args: "dev:frontend"
    // },
  ],
};
