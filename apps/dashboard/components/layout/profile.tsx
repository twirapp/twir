import { Avatar, Image, Menu, ScrollArea, TextInput } from '@mantine/core';
import { IconLogout } from '@tabler/icons';
import { useState } from 'react';

export function Profile() {
  const dashboards = [
    {
      id: '72218eae-5584-4288-b9e6-bc35ac606942',
      channelId: '139336353',
      userId: '128644134',
      twitchUser: {
        id: '139336353',
        login: '7ssk7',
        display_name: '7ssk7',
        type: '',
        broadcaster_type: 'partner',
        description: 'Каво ',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/66cb7060-1a8a-4fca-9ccd-f760b70af636-profile_image-300x300.png',
        offline_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/452001ef-b406-49b2-bf44-c7c7fcae96be-channel_offline_image-1920x1080.jpeg',
        view_count: 10228739,
        email: '',
        created_at: '2016-11-12T15:51:02Z',
      },
    },
    {
      id: 'a2d436d0-b453-407a-9c3a-61151e3a4788',
      channelId: '138780411',
      userId: '128644134',
      twitchUser: {
        id: '138780411',
        login: 'vozdyhatel_',
        display_name: 'vozdyhatel_',
        type: '',
        broadcaster_type: 'affiliate',
        description:
          '4000 игр на Инвокере в прошлом, сейчас не играю на мейне, бущу и учусь на Виверне 3 pos.',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/86f1a112-f4fd-4591-98a1-1e06cff0f410-profile_image-300x300.png',
        offline_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/d54b1e99-2898-491a-b4c6-eb78d0a45e74-channel_offline_image-1920x1080.png',
        view_count: 292180,
        email: '',
        created_at: '2016-11-05T11:15:27Z',
      },
    },
    {
      id: '8cd3186d-13e4-4ec1-bc03-7c47fd2b3e81',
      channelId: '106729827',
      userId: '128644134',
      twitchUser: {
        id: '106729827',
        login: 'rusty',
        display_name: 'RUSTY',
        type: '',
        broadcaster_type: 'partner',
        description:
          'Привет, дорогой зритель, добро пожаловать, чувствуй себя как дома, но не забывай что ты в гостях 🌿 Ознакомься с информацией под стримом и правилами чата ⚠  И приятного просмотра! 🦄 ',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/f848db72-5115-4e60-9f60-4f4112603ba4-profile_image-300x300.png',
        offline_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/82f48484-efb0-42bd-8efc-e85b369ba012-channel_offline_image-1920x1080.jpeg',
        view_count: 1370115,
        email: '',
        created_at: '2015-11-10T13:12:19Z',
      },
    },
    {
      id: 'dff3003e-cfb0-4829-a172-aa93fb3b3dac',
      channelId: '48385787',
      userId: '128644134',
      twitchUser: {
        id: '48385787',
        login: 'promotive',
        display_name: 'PROMOTIVE',
        type: '',
        broadcaster_type: 'partner',
        description: '\\m/',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/c3d1b93c-7734-4898-b798-344c753aeaf2-profile_image-300x300.png',
        offline_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/promotive-channel_offline_image-dc39be92d722f578-1920x1080.png',
        view_count: 2097598,
        email: '',
        created_at: '2013-08-31T17:13:12Z',
      },
    },
    {
      id: '3bd4411f-499c-4a56-9489-44dead2bbc5d',
      channelId: '68948764',
      userId: '128644134',
      twitchUser: {
        id: '68948764',
        login: 'rupero',
        display_name: 'rupero',
        type: '',
        broadcaster_type: '',
        description: '27.',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/bfe30b69-99c5-4e28-bf8c-046c2e177872-profile_image-300x300.png',
        offline_image_url: '',
        view_count: 1067,
        email: '',
        created_at: '2014-08-14T10:56:30Z',
      },
    },
    {
      id: '61ef6237-c007-442a-b54b-6bc89a70ec74',
      channelId: '161293680',
      userId: '128644134',
      twitchUser: {
        id: '161293680',
        login: 'annet_broadway',
        display_name: 'Annet_Broadway',
        type: '',
        broadcaster_type: 'affiliate',
        description: 'Кибертуристка, киберспортсменка, токсик, чсв. ',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/a84e4601-bccd-40dd-94f8-f7db95c6108b-profile_image-300x300.png',
        offline_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/12b2d265-6aa3-46e2-a2ca-55a2ba42dab2-channel_offline_image-1920x1080.jpeg',
        view_count: 19668,
        email: '',
        created_at: '2017-06-23T15:31:30Z',
      },
    },
    {
      id: 'cf8c68cc-a774-4367-ad02-1c9434317b6d',
      channelId: '425431025',
      userId: '128644134',
      twitchUser: {
        id: '425431025',
        login: 'snussed',
        display_name: 'Snussed',
        type: '',
        broadcaster_type: '',
        description: ' 󠀀',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/5a1b492c-486c-49c7-9149-c59944edebbd-profile_image-300x300.png',
        offline_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/9e21e1a8-00e2-4749-a0ed-0e793e531af3-channel_offline_image-1920x1080.jpeg',
        view_count: 1553,
        email: '',
        created_at: '2019-03-23T18:46:19Z',
      },
    },
    {
      id: '1af86e17-718f-47da-a6ce-89c6d6b04d88',
      channelId: '100909459',
      userId: '128644134',
      twitchUser: {
        id: '100909459',
        login: 'bogush',
        display_name: 'Bogush',
        type: '',
        broadcaster_type: 'partner',
        description: '',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/7a342f9d-282e-41aa-ab31-d0165b74f746-profile_image-300x300.png',
        offline_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/1d2564d7-a4c9-46b6-94be-0e90d79b6910-channel_offline_image-1920x1080.jpeg',
        view_count: 2687877,
        email: '',
        created_at: '2015-08-29T11:56:58Z',
      },
    },
    {
      id: '487667d2-cef9-4e1f-96f5-3402de343531',
      channelId: '175657580',
      userId: '128644134',
      twitchUser: {
        id: '175657580',
        login: 'lawsn_',
        display_name: 'Lawsn_',
        type: '',
        broadcaster_type: '',
        description: '',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/01641125-37d8-4c91-9dea-36062062b969-profile_image-300x300.png',
        offline_image_url: '',
        view_count: 306,
        email: '',
        created_at: '2017-09-27T09:07:23Z',
      },
    },
    {
      id: 'e045da9e-e083-4b10-bdeb-a74749761c0e',
      channelId: '189703483',
      userId: '128644134',
      twitchUser: {
        id: '189703483',
        login: 'daetojekara',
        display_name: 'daetojekaRa',
        type: '',
        broadcaster_type: 'affiliate',
        description:
          'Добро пожаловать на самый темный и карающий канал в стрим комьюнити 🔥Campfire. На канале вы сможете увидеть трансляций по ⚡️DOTA 2 и других играх, так же по разработке каких либо 💻 небольших проектов ✏️ образовательных материалов.',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/b73f81e7-3fe1-415b-a543-4fe164d16e56-profile_image-300x300.png',
        offline_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/df874b1d-5c8b-4c6b-b088-f94ef10e5f18-channel_offline_image-1920x1080.jpeg',
        view_count: 10272,
        email: '',
        created_at: '2018-01-09T11:32:26Z',
      },
    },
    {
      id: '1dad453b-a12f-4b5a-93c4-f90157a88d3f',
      channelId: '145291600',
      userId: '128644134',
      twitchUser: {
        id: '145291600',
        login: 'heyitzyztem',
        display_name: 'HeyItZyztem',
        type: '',
        broadcaster_type: '',
        description: 'Your average streamer who streams on weekends! ',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/4b95d66f-b928-4aeb-b13a-d835ca49d6ff-profile_image-300x300.png',
        offline_image_url: '',
        view_count: 219,
        email: '',
        created_at: '2017-01-19T21:50:46Z',
      },
    },
    {
      id: 'aa0994f3-f695-45f9-8a2b-c62a165937a4',
      channelId: '207545504',
      userId: '128644134',
      twitchUser: {
        id: '207545504',
        login: 'poyale',
        display_name: 'poyale',
        type: '',
        broadcaster_type: '',
        description: '',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/dde20ef7-a14e-42f0-8b7c-5b5114076787-profile_image-300x300.png',
        offline_image_url: '',
        view_count: 560,
        email: '',
        created_at: '2018-03-22T17:00:40Z',
      },
    },
    {
      id: 'eccd4ff5-4337-4c3a-9c38-5ce4c6ca92d7',
      channelId: '272645477',
      userId: '128644134',
      twitchUser: {
        id: '272645477',
        login: 'fudzzuzzim',
        display_name: 'fudzzuzzim',
        type: '',
        broadcaster_type: '',
        description: '',
        profile_image_url:
          'https://static-cdn.jtvnw.net/user-default-pictures-uv/13e5fa74-defa-11e9-809c-784f43822e80-profile_image-300x300.png',
        offline_image_url: '',
        view_count: 15,
        email: '',
        created_at: '2018-11-05T18:22:38Z',
      },
    },
    {
      id: '96e19f04-9ccb-4570-ae7d-33bf772bf4a4',
      channelId: '208860955',
      userId: '128644134',
      twitchUser: {
        id: '208860955',
        login: 'smart_iq',
        display_name: 'smart_iq',
        type: '',
        broadcaster_type: '',
        description: '',
        profile_image_url:
          'https://static-cdn.jtvnw.net/user-default-pictures-uv/41780b5a-def8-11e9-94d9-784f43822e80-profile_image-300x300.png',
        offline_image_url: '',
        view_count: 9,
        email: '',
        created_at: '2018-03-26T17:49:31Z',
      },
    },
    {
      id: '5c86257d-da83-43d6-b7ae-458386d68db7',
      channelId: '478006053',
      userId: '128644134',
      twitchUser: {
        id: '478006053',
        login: 'blackwolf_osu',
        display_name: 'blackwolf_osu',
        type: '',
        broadcaster_type: '',
        description: '',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/cd069bcf-247f-4385-a324-1a3d01993909-profile_image-300x300.png',
        offline_image_url: '',
        view_count: 95,
        email: '',
        created_at: '2019-12-15T17:37:18Z',
      },
    },
    {
      id: '7af0afaa-f018-4c2d-8fed-e1f2bbf8a5e4',
      channelId: '704163392',
      userId: '128644134',
      twitchUser: {
        id: '704163392',
        login: 'bogdan_liusik',
        display_name: 'bogdan_liusik',
        type: '',
        broadcaster_type: '',
        description: '',
        profile_image_url:
          'https://static-cdn.jtvnw.net/user-default-pictures-uv/ead5c8b2-a4c9-4724-b1dd-9f00b46cbd3d-profile_image-300x300.png',
        offline_image_url: '',
        view_count: 0,
        email: '',
        created_at: '2021-07-04T23:26:12Z',
      },
    },
    {
      id: '5177ca69-54b0-4208-bcda-14b4f1c5e0a8',
      channelId: '130606008',
      userId: '128644134',
      twitchUser: {
        id: '130606008',
        login: 'desgozhik',
        display_name: 'desgozhik',
        type: '',
        broadcaster_type: 'affiliate',
        description: 'cute dead inside :3',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/8e9650e3-36af-4b2a-8481-3051dfdcf42c-profile_image-300x300.png',
        offline_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/75ca789f-7ad2-4d68-97c2-76428ecde3bc-channel_offline_image-1920x1080.jpeg',
        view_count: 40659,
        email: '',
        created_at: '2016-07-27T11:31:54Z',
      },
    },
    {
      id: '21c3b1cc-ac18-4ae0-ba53-9a090e1eea6b',
      channelId: '81815340',
      userId: '128644134',
      twitchUser: {
        id: '81815340',
        login: 'isnicable',
        display_name: 'isnicable',
        type: '',
        broadcaster_type: 'affiliate',
        description:
          'Next Stream coming soon, hit the follow button to get a notification 🔔  | !next ',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/cb10e029-2051-465c-88be-2697a2f3fb77-profile_image-300x300.png',
        offline_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/5b382e13-ecd3-4b29-bb28-7d29439aa90f-channel_offline_image-1920x1080.jpeg',
        view_count: 13095,
        email: '',
        created_at: '2015-02-04T20:56:09Z',
      },
    },
    {
      id: 'a7dfcd29-3fea-4fd5-ae06-db76b1ef58d0',
      channelId: '840442714',
      userId: '128644134',
      twitchUser: {
        id: '840442714',
        login: 'daxahe2544',
        display_name: 'daxahe2544',
        type: '',
        broadcaster_type: '',
        description: '',
        profile_image_url:
          'https://static-cdn.jtvnw.net/user-default-pictures-uv/ce57700a-def9-11e9-842d-784f43822e80-profile_image-300x300.png',
        offline_image_url: '',
        view_count: 0,
        email: '',
        created_at: '2022-10-25T15:18:43Z',
      },
    },
    {
      id: 'f1f91d2d-b4b1-4c6a-9432-0f0026a499a4',
      channelId: '722533144',
      userId: '128644134',
      twitchUser: {
        id: '722533144',
        login: 'nedojunior',
        display_name: 'NedoJunioR',
        type: '',
        broadcaster_type: 'affiliate',
        description: 'Всем привет, веселюся с Джабаскрипт и катаю в фифулечку фифу. ',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/26ddb2bf-5e4b-4a28-ad79-01dd2743332d-profile_image-300x300.png',
        offline_image_url: '',
        view_count: 9859,
        email: '',
        created_at: '2021-08-30T11:23:14Z',
      },
    },
    {
      id: '83593e80-2550-4a03-b917-dc656106ab7e',
      channelId: '136818541',
      userId: '128644134',
      twitchUser: {
        id: '136818541',
        login: 'try_to_ca7ch',
        display_name: 'try_to_ca7ch',
        type: '',
        broadcaster_type: '',
        description:
          'No flame stream. Only peace and love. Playing anything i want. Welcome traveller)',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/fe87a767-cc29-42db-8632-156eb1c9a30f-profile_image-300x300.png',
        offline_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/3f822b8e-3b6b-487a-8ec3-22425db91ac5-channel_offline_image-1920x1080.jpeg',
        view_count: 41,
        email: '',
        created_at: '2016-10-11T06:56:32Z',
      },
    },
    {
      id: '470d5740-fe06-4e9f-a30f-3da053b41b56',
      channelId: '155934011',
      userId: '128644134',
      twitchUser: {
        id: '155934011',
        login: 'awaayf',
        display_name: 'awaayf',
        type: '',
        broadcaster_type: '',
        description: 'liubliu devochek',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/9eb0c1e4-f115-4b2c-a22a-5b604c60b7cd-profile_image-300x300.png',
        offline_image_url: '',
        view_count: 1268,
        email: '',
        created_at: '2017-05-06T07:23:44Z',
      },
    },
    {
      id: '4c8ddf3a-4aa3-4665-a6e7-d8cd195a7188',
      channelId: '58424295',
      userId: '128644134',
      twitchUser: {
        id: '58424295',
        login: 'morfixx',
        display_name: 'MORFIXX',
        type: '',
        broadcaster_type: 'affiliate',
        description: 'idfk Alex - 25 - ez 😎',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/f6fcfb1ebf395661-profile_image-300x300.jpeg',
        offline_image_url: '',
        view_count: 3148,
        email: '',
        created_at: '2014-03-08T16:11:22Z',
      },
    },
    {
      id: '4021c650-2337-468d-bde5-5f6202a4d7d4',
      channelId: '63980920',
      userId: '128644134',
      twitchUser: {
        id: '63980920',
        login: 'phill67',
        display_name: 'phill67',
        type: '',
        broadcaster_type: '',
        description: '',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/phill67-profile_image-9067c72b79aaf6c1-300x300.jpeg',
        offline_image_url: '',
        view_count: 325,
        email: '',
        created_at: '2014-06-09T10:44:32Z',
      },
    },
    {
      id: 'eb470811-b1b9-4172-ba67-0899179aa487',
      channelId: '218626726',
      userId: '128644134',
      twitchUser: {
        id: '218626726',
        login: 'sp4rco',
        display_name: 'sp4rco',
        type: '',
        broadcaster_type: 'affiliate',
        description: 'SP4RCO - НоуНейм, бог, можно просто amor hower',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/c222f590-2769-450e-b17d-07d7efafda14-profile_image-300x300.png',
        offline_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/874edf98-fe06-40cf-bd30-d89beb4d077f-channel_offline_image-1920x1080.jpeg',
        view_count: 21060,
        email: '',
        created_at: '2018-05-05T13:27:09Z',
      },
    },
    {
      id: 'e5f71bf3-79e9-450e-850f-94994a780751',
      channelId: '786718415',
      userId: '128644134',
      twitchUser: {
        id: '786718415',
        login: 'sickestmans',
        display_name: 'SickestMans',
        type: '',
        broadcaster_type: '',
        description: '',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/942f1690-3864-4ba7-a0a1-190795f0f87a-profile_image-300x300.png',
        offline_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/0ce697d4-4c63-4b60-bc5d-1f46f350d4dd-channel_offline_image-1920x1080.png',
        view_count: 2,
        email: '',
        created_at: '2022-04-09T20:53:31Z',
      },
    },
    {
      id: '27b72f31-b68a-43ba-82d4-ea44fb98b027',
      channelId: '180851298',
      userId: '128644134',
      twitchUser: {
        id: '180851298',
        login: 'dankybat',
        display_name: 'dankYbat',
        type: '',
        broadcaster_type: '',
        description: 'dY~ waste time in myself...☘️ ~ DS: dankY#1157',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/67889c89-dd4b-46bc-bfed-fd6d9ea92d58-profile_image-300x300.jpeg',
        offline_image_url: '',
        view_count: 445,
        email: '',
        created_at: '2017-11-07T20:27:30Z',
      },
    },
    {
      id: 'bd2c6060-881d-4992-ac20-a832c87c8c84',
      channelId: '191990911',
      userId: '128644134',
      twitchUser: {
        id: '191990911',
        login: 'intessssa',
        display_name: 'intessssa',
        type: '',
        broadcaster_type: 'affiliate',
        description: 'Люблю сидеть на стуле и смотреть в экран',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/2ff7412f-1bcb-425a-8563-48ae44dd10f1-profile_image-300x300.png',
        offline_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/a8a1fefa-9560-459b-8dbb-8dd7c9276bf1-channel_offline_image-1920x1080.jpeg',
        view_count: 45092,
        email: '',
        created_at: '2018-01-21T07:15:50Z',
      },
    },
    {
      id: 'd09e0c3d-29db-48f0-9387-d9c25befa19e',
      channelId: '104435562',
      userId: '128644134',
      twitchUser: {
        id: '104435562',
        login: 'qrushcsgo',
        display_name: 'QRUSHcsgo',
        type: '',
        broadcaster_type: 'partner',
        description: 'Легендарный Контерструкер Уфы',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/a477bccc-9b23-44d7-a379-fe64f67898c3-profile_image-300x300.png',
        offline_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/f1f477b9-5d40-47e2-8b03-3a07f4f18519-channel_offline_image-1920x1080.png',
        view_count: 44475094,
        email: '',
        created_at: '2015-10-15T17:01:39Z',
      },
    },
    {
      id: '9e895704-3cb1-4d06-b58a-5b2851047779',
      channelId: '492303803',
      userId: '128644134',
      twitchUser: {
        id: '492303803',
        login: 'rusyahigh',
        display_name: 'Rusyahigh',
        type: '',
        broadcaster_type: 'affiliate',
        description: 'Привет,мое имя - Руся.Играю в дотку.Об остальном спрашивай сам (:',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/b8583c1b-3ae5-4f97-9bc6-69e506f4ae4e-profile_image-300x300.png',
        offline_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/f22cef1b-d12f-4cbd-abe8-b5628523745e-channel_offline_image-1920x1080.png',
        view_count: 241851,
        email: '',
        created_at: '2020-02-13T21:24:10Z',
      },
    },
    {
      id: 'f9a83bad-8974-454c-adec-d1092b920021',
      channelId: '101474138',
      userId: '128644134',
      twitchUser: {
        id: '101474138',
        login: 'vyrts',
        display_name: 'Vyrts',
        type: '',
        broadcaster_type: 'affiliate',
        description: 'Business inquiries - vyrts.business@gmail.com',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/3d4c07b9-b252-48b6-a16f-7f2b3305c3e3-profile_image-300x300.png',
        offline_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/ff2c6206-a808-433e-b795-4cc8a64f1a26-channel_offline_image-1920x1080.png',
        view_count: 68,
        email: '',
        created_at: '2015-09-05T11:05:24Z',
      },
    },
    {
      id: 'cc10ac2d-458e-479f-a599-cd05e462e1e8',
      channelId: '117189279',
      userId: '128644134',
      twitchUser: {
        id: '117189279',
        login: 'gosugdr',
        display_name: 'gosugdr',
        type: '',
        broadcaster_type: '',
        description: '',
        profile_image_url:
          'https://static-cdn.jtvnw.net/user-default-pictures-uv/ebe4cd89-b4f4-4cd9-adac-2f30151b4209-profile_image-300x300.png',
        offline_image_url: '',
        view_count: 91,
        email: '',
        created_at: '2016-02-27T21:42:11Z',
      },
    },
    {
      id: 'e170f942-0eb0-4f48-a89f-9e855ad3ced1',
      channelId: '155664706',
      userId: '128644134',
      twitchUser: {
        id: '155664706',
        login: 'appell__',
        display_name: 'Appell__',
        type: '',
        broadcaster_type: 'affiliate',
        description: '',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/499f4c8b-8a20-47e5-97ad-9c4dccc9fedc-profile_image-300x300.png',
        offline_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/c1dffdf140debea1-channel_offline_image-1920x1080.jpeg',
        view_count: 3446,
        email: '',
        created_at: '2017-05-03T20:34:57Z',
      },
    },
    {
      id: '46daf87b-3490-459f-a860-99260ed66173',
      channelId: '148191620',
      userId: '128644134',
      twitchUser: {
        id: '148191620',
        login: 'rahimoov',
        display_name: 'rahimoov',
        type: '',
        broadcaster_type: '',
        description: '',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/4cf4b45d-be81-46ff-8573-23db68c33119-profile_image-300x300.jpg',
        offline_image_url: '',
        view_count: 186,
        email: '',
        created_at: '2017-02-19T15:28:49Z',
      },
    },
    {
      id: 'e8df4134-cdd1-4b89-8607-f6d945df3beb',
      channelId: '87574843',
      userId: '128644134',
      twitchUser: {
        id: '87574843',
        login: 'sora1moonpma',
        display_name: 'sora1moonPMA',
        type: '',
        broadcaster_type: 'affiliate',
        description: 'Существую',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/39fdbdba-1936-435e-b572-513f70b1ecfb-profile_image-300x300.png',
        offline_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/38c4375f-974e-47fa-ad37-87ba12d9eea2-channel_offline_image-1920x1080.png',
        view_count: 164,
        email: '',
        created_at: '2015-04-05T10:09:47Z',
      },
    },
    {
      id: '5a473294-eb22-4f53-992a-befa070705ad',
      channelId: '136896435',
      userId: '128644134',
      twitchUser: {
        id: '136896435',
        login: 'andy_i_',
        display_name: 'ANDY_I_',
        type: '',
        broadcaster_type: 'affiliate',
        description: '            ｗｉｌｌ　ｂｅ　ｆｉｎｅ ',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/976beec4-d65a-4cb9-a18a-404829131376-profile_image-300x300.png',
        offline_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/5d582cb5-7818-42aa-8725-312b6a4921e1-channel_offline_image-1920x1080.jpeg',
        view_count: 54304,
        email: '',
        created_at: '2016-10-12T12:21:04Z',
      },
    },
    {
      id: 'f1e403bf-3fe7-4739-ba66-1f4feb3c25a0',
      channelId: '40733666',
      userId: '128644134',
      twitchUser: {
        id: '40733666',
        login: 'ninox14',
        display_name: 'NiNoX14',
        type: '',
        broadcaster_type: '',
        description: '',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/8d176add-6105-4ebb-bdb4-d9d53287ccf1-profile_image-300x300.jpg',
        offline_image_url: '',
        view_count: 1640,
        email: '',
        created_at: '2013-02-25T14:46:50Z',
      },
    },
    {
      id: '5c0b6a96-0dbb-4576-b9d1-45f46b1fd30d',
      channelId: '100233760',
      userId: '128644134',
      twitchUser: {
        id: '100233760',
        login: 'fonchik',
        display_name: 'FONCHIK',
        type: '',
        broadcaster_type: '',
        description: '',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/49904d05-6c2c-42cc-a2d9-c5f55220ae66-profile_image-300x300.jpeg',
        offline_image_url: '',
        view_count: 102,
        email: '',
        created_at: '2015-08-22T14:21:05Z',
      },
    },
    {
      id: '07f5250f-a315-40c2-9e5f-dc60bc916489',
      channelId: '614425362',
      userId: '128644134',
      twitchUser: {
        id: '614425362',
        login: 'heytellmesomething',
        display_name: 'heytellmesomething',
        type: '',
        broadcaster_type: '',
        description: '',
        profile_image_url:
          'https://static-cdn.jtvnw.net/user-default-pictures-uv/75305d54-c7cc-40d1-bb9c-91fbe85943c7-profile_image-300x300.png',
        offline_image_url: '',
        view_count: 787,
        email: '',
        created_at: '2020-11-29T08:14:31Z',
      },
    },
    {
      id: 'f4deb802-cb8e-4c6d-a0cb-ea01eaf67d23',
      channelId: '554108347',
      userId: '128644134',
      twitchUser: {
        id: '554108347',
        login: 'kvizyx',
        display_name: 'kvizyx',
        type: '',
        broadcaster_type: '',
        description: '-_- ',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/5359ad5d-b80c-427b-b493-9f01fc968bae-profile_image-300x300.png',
        offline_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/7ede7084-f3a8-4b37-9292-ffdf7f8049be-channel_offline_image-1920x1080.png',
        view_count: 894,
        email: '',
        created_at: '2020-07-13T07:57:37Z',
      },
    },
    {
      id: '231c5961-d03b-463e-b7f1-9f3cdf9fc2eb',
      channelId: '244902384',
      userId: '128644134',
      twitchUser: {
        id: '244902384',
        login: 'mrpandir',
        display_name: 'MrPandir',
        type: '',
        broadcaster_type: '',
        description: 'Зовут Влад, вроде как программист :D',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/404e84dd-a10a-441e-9b1e-86d900abf78e-profile_image-300x300.png',
        offline_image_url: '',
        view_count: 3,
        email: '',
        created_at: '2018-08-02T14:25:12Z',
      },
    },
    {
      id: '7cc5c952-9e3e-4139-b3be-a7c754d5b4a6',
      channelId: '222673956',
      userId: '128644134',
      twitchUser: {
        id: '222673956',
        login: 'mellkam',
        display_name: 'mellkam',
        type: '',
        broadcaster_type: '',
        description: 'Artem, 17 years old web developer and ui/ux designer from Ukraine.',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/9d51611d-a8a6-4d50-8821-83c41666c225-profile_image-300x300.png',
        offline_image_url: '',
        view_count: 24,
        email: '',
        created_at: '2018-05-13T18:46:41Z',
      },
    },
    {
      id: '73d3c2da-da52-4602-ba5e-f56f9a66516a',
      channelId: '277306419',
      userId: '128644134',
      twitchUser: {
        id: '277306419',
        login: 'youtheman12221',
        display_name: 'youtheman12221',
        type: '',
        broadcaster_type: '',
        description: '',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/9519cc68-592c-43fe-aa22-89af8a4f3df0-profile_image-300x300.png',
        offline_image_url: '',
        view_count: 568,
        email: '',
        created_at: '2018-11-22T23:55:01Z',
      },
    },
    {
      id: '4e9ea9d1-46e1-4f88-9471-45f1ac0c290b',
      channelId: '36979094',
      userId: '128644134',
      twitchUser: {
        id: '36979094',
        login: 'kujijiepuk',
        display_name: 'KuJIJIePuK',
        type: '',
        broadcaster_type: '',
        description: '',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/kujijiepuk-profile_image-e075d0ad96e2b768-300x300.jpeg',
        offline_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/6473bc37-8b6c-4d58-87f6-61e562386cb2-channel_offline_image-1920x1080.jpeg',
        view_count: 1402,
        email: '',
        created_at: '2012-10-17T09:25:45Z',
      },
    },
    {
      id: '1a347518-0998-425a-96a4-4e326f2ccf32',
      channelId: '146712489',
      userId: '128644134',
      twitchUser: {
        id: '146712489',
        login: 'qharka',
        display_name: 'qHarka',
        type: '',
        broadcaster_type: 'affiliate',
        description: 'qHarka стримит afk контент.',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/b141020d-a51a-4833-986f-dba095f06a11-profile_image-300x300.png',
        offline_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/939f17b4-39da-46e7-bf6f-9ef9ed57ca95-channel_offline_image-1920x1080.png',
        view_count: 313857,
        email: '',
        created_at: '2017-02-02T13:21:59Z',
      },
    },
    {
      id: 'd69b61db-cd27-4b27-a038-04ea812cdcab',
      channelId: '155644238',
      userId: '128644134',
      twitchUser: {
        id: '155644238',
        login: 'le_xot',
        display_name: 'le_xot',
        type: '',
        broadcaster_type: 'affiliate',
        description: 'Волосатая чубака',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/423e40e6-9534-46ac-9ed8-5714657dd03b-profile_image-300x300.png',
        offline_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/36df4401-b17f-46b8-85c9-eed2799bc42c-channel_offline_image-1920x1080.png',
        view_count: 21863,
        email: '',
        created_at: '2017-05-03T17:31:24Z',
      },
    },
    {
      id: '58e8c6c3-a385-4b76-b4e1-55afc8160081',
      channelId: '813674987',
      userId: '128644134',
      twitchUser: {
        id: '813674987',
        login: 'gowstbot4',
        display_name: 'Gowstbot4',
        type: '',
        broadcaster_type: '',
        description: '',
        profile_image_url:
          'https://static-cdn.jtvnw.net/user-default-pictures-uv/dbdc9198-def8-11e9-8681-784f43822e80-profile_image-300x300.png',
        offline_image_url: '',
        view_count: 0,
        email: '',
        created_at: '2022-07-30T09:48:36Z',
      },
    },
    {
      id: '3626614d-194a-4de5-9fcf-998145ceba69',
      channelId: '39580248',
      userId: '128644134',
      twitchUser: {
        id: '39580248',
        login: 'shhdotcom',
        display_name: 'shhDOTcom',
        type: '',
        broadcaster_type: 'affiliate',
        description: `shhDOTcom streams Counter-Strike: Global Offensive... As long as we're having fun!`,
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/58ef39ae-adf6-4a07-a8c7-1da493837b94-profile_image-300x300.png',
        offline_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/dd433b2b6a36475a-channel_offline_image-1920x1080.jpeg',
        view_count: 49329,
        email: '',
        created_at: '2013-01-22T11:03:10Z',
      },
    },
    {
      id: 'e420aea6-8db9-4677-b1bc-b7218695cc2b',
      channelId: '79788319',
      userId: '128644134',
      twitchUser: {
        id: '79788319',
        login: 'clbagrat',
        display_name: 'clbagrat',
        type: '',
        broadcaster_type: '',
        description: `Hi, my name is Bagrat. I'm streaming frontend development. `,
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/5fc9f62c-15a1-4135-9096-e7184fe0e0f7-profile_image-300x300.jpeg',
        offline_image_url: '',
        view_count: 555,
        email: '',
        created_at: '2015-01-14T10:29:54Z',
      },
    },
    {
      id: 'dc393700-c092-48fa-a526-9c2be920bcb9',
      channelId: '164759125',
      userId: '128644134',
      twitchUser: {
        id: '164759125',
        login: 'dr1sha',
        display_name: 'DR1SHA',
        type: '',
        broadcaster_type: 'affiliate',
        description: 'сигам контра болталка с кентами',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/bc671e8a-4cd6-494a-a1d0-290d78057cf1-profile_image-300x300.jpeg',
        offline_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/a2c302a1-5a79-4368-b834-83726c5fd45d-channel_offline_image-1920x1080.png',
        view_count: 2109,
        email: '',
        created_at: '2017-07-15T07:35:08Z',
      },
    },
    {
      id: '0b434a2f-7b54-4ccf-bcf8-71bd2763e719',
      channelId: '118912176',
      userId: '128644134',
      twitchUser: {
        id: '118912176',
        login: 'nofearj',
        display_name: 'NOFEARj',
        type: '',
        broadcaster_type: 'affiliate',
        description: '',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/0e12f391-f6a3-4cd6-8e8c-0cd7f8224750-profile_image-300x300.png',
        offline_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/7726eb4e-e733-4c88-99a6-ab368f864fe6-channel_offline_image-1920x1080.png',
        view_count: 4331,
        email: '',
        created_at: '2016-03-16T21:50:43Z',
      },
    },
    {
      id: 'b37aefe5-889d-4eea-9d4b-e9fd9a82b5d5',
      channelId: '155976884',
      userId: '128644134',
      twitchUser: {
        id: '155976884',
        login: 'evgenous',
        display_name: 'evgenous',
        type: '',
        broadcaster_type: '',
        description: 'Веб-программирование',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/d304344e-521e-46c1-ad10-d604f6794de6-profile_image-300x300.png',
        offline_image_url: '',
        view_count: 109,
        email: '',
        created_at: '2017-05-06T16:07:58Z',
      },
    },
    {
      id: 'cfad034d-949d-47ea-a75d-5f451baa4987',
      channelId: '67564984',
      userId: '128644134',
      twitchUser: {
        id: '67564984',
        login: 'baldrcoal',
        display_name: 'BaldrCoal',
        type: '',
        broadcaster_type: 'affiliate',
        description: 'я че лох?',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/94b7763f-4f20-475a-a3dd-cf32b5ba410a-profile_image-300x300.png',
        offline_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/898b5140-3346-4b13-ae29-76e7d799fd1b-channel_offline_image-1920x1080.jpeg',
        view_count: 7135,
        email: '',
        created_at: '2014-07-28T21:16:38Z',
      },
    },
    {
      id: '3683aec5-9e30-4d0f-b6c0-581266bb335e',
      channelId: '27618908',
      userId: '128644134',
      twitchUser: {
        id: '27618908',
        login: 'monk324',
        display_name: 'MonK324',
        type: '',
        broadcaster_type: '',
        description: '',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/462a77f9-5a29-4f72-8a0a-62d608d18fc5-profile_image-300x300.png',
        offline_image_url: '',
        view_count: 1564,
        email: '',
        created_at: '2012-01-21T09:01:43Z',
      },
    },
    {
      id: 'a587c825-c167-401c-97c4-0c6b7db5343d',
      channelId: '84446583',
      userId: '128644134',
      twitchUser: {
        id: '84446583',
        login: 'rdocompendium',
        display_name: 'RDOCompendium',
        type: '',
        broadcaster_type: '',
        description: '',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/e94e66f3-3396-479c-a27b-ddaabf905d4a-profile_image-300x300.png',
        offline_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/ddc3861c-280b-47da-8a3d-178addb6f1cb-channel_offline_image-1920x1080.jpeg',
        view_count: 459,
        email: '',
        created_at: '2015-03-05T20:07:05Z',
      },
    },
    {
      id: '3ccc82e1-d483-4d32-86ca-9e4b4d0f5beb',
      channelId: '90470158',
      userId: '128644134',
      twitchUser: {
        id: '90470158',
        login: 'tapakeht',
        display_name: 'TaPaKeHT',
        type: '',
        broadcaster_type: 'affiliate',
        description:
          'Меня зовут Александр, мне 30 лет. Играю в разнообразные игры, но большую часть времени уделяю WoW и Apex Legends. Буду рад с вами поболтать на тему игр, программирования, или еще чего-нибудь...',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/71c7f6b5-b159-4689-8140-f202ff25ba58-profile_image-300x300.png',
        offline_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/2b817dd7-ef54-4a22-9c34-634b2a5aa811-channel_offline_image-1920x1080.jpeg',
        view_count: 21,
        email: '',
        created_at: '2015-05-07T13:38:31Z',
      },
    },
    {
      id: '8884d7d0-bd98-41bd-80d4-c94003006f6c',
      channelId: '827781630',
      userId: '128644134',
      twitchUser: {
        id: '827781630',
        login: 'melkambot',
        display_name: 'melkambot',
        type: '',
        broadcaster_type: '',
        description: '',
        profile_image_url:
          'https://static-cdn.jtvnw.net/user-default-pictures-uv/13e5fa74-defa-11e9-809c-784f43822e80-profile_image-300x300.png',
        offline_image_url: '',
        view_count: 0,
        email: '',
        created_at: '2022-09-22T13:45:57Z',
      },
    },
    {
      id: 'e7e09b75-4970-494c-9875-39d869bcc794',
      channelId: '430720419',
      userId: '128644134',
      twitchUser: {
        id: '430720419',
        login: 'sandakoff',
        display_name: 'sandakoff',
        type: '',
        broadcaster_type: '',
        description: 'Тут просто текст, который нечего не значит',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/8721141e-92cb-478c-871e-08aee7ba8743-profile_image-300x300.png',
        offline_image_url: '',
        view_count: 67,
        email: '',
        created_at: '2019-04-18T05:27:29Z',
      },
    },
    {
      id: 'bffab8d0-ebd1-4d4d-85a5-2cefd54db98c',
      channelId: '155641150',
      userId: '128644134',
      twitchUser: {
        id: '155641150',
        login: 'flashhhhh',
        display_name: 'Flashhhhh',
        type: '',
        broadcaster_type: '',
        description: '',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/570ae213-c5a3-49ab-bc46-eb0042f87637-profile_image-300x300.png',
        offline_image_url: '',
        view_count: 3591,
        email: '',
        created_at: '2017-05-03T17:04:22Z',
      },
    },
    {
      id: '1cc512d4-445d-40a9-a64a-3556af701cb3',
      channelId: '127678844',
      userId: '128644134',
      twitchUser: {
        id: '127678844',
        login: 'stre1ok_',
        display_name: 'Stre1ok_',
        type: '',
        broadcaster_type: 'affiliate',
        description:
          'Приветствую на канале Стрелка! ☕  Вступай в наш Дискорд Разработчиков https://discord.gg/FmrBzNcRHx',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/d5c790b5-eb8b-4dbe-9111-cdabbb5b391a-profile_image-300x300.png',
        offline_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/f564cf92-ff3f-4f2f-a469-ea73677eea8f-channel_offline_image-1920x1080.png',
        view_count: 128966,
        email: '',
        created_at: '2016-06-25T07:53:46Z',
      },
    },
    {
      id: '1cccc7eb-51b5-4dda-934a-369221dabad2',
      channelId: '257316595',
      userId: '128644134',
      twitchUser: {
        id: '257316595',
        login: 'iwlj4s',
        display_name: 'iwlj4s',
        type: '',
        broadcaster_type: '',
        description: '',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/46f29fe4-860d-45af-b03d-ee8103495a17-profile_image-300x300.png',
        offline_image_url: '',
        view_count: 7,
        email: '',
        created_at: '2018-09-11T06:42:37Z',
      },
    },
    {
      id: '90d9f47a-6f08-4203-bc9f-ccc98341fae7',
      channelId: '234946190',
      userId: '128644134',
      twitchUser: {
        id: '234946190',
        login: 'b4dtreeper',
        display_name: 'B4DTREEPER',
        type: '',
        broadcaster_type: 'affiliate',
        description: 'поигрываю в валорантик',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/140fd64f-80bb-477d-b5a2-df9cb20f2f28-profile_image-300x300.png',
        offline_image_url: '',
        view_count: 11718,
        email: '',
        created_at: '2018-06-29T14:17:11Z',
      },
    },
    {
      id: '0273617b-ae74-44a7-9277-1cd04f1e3f1e',
      channelId: '40488774',
      userId: '128644134',
      twitchUser: {
        id: '40488774',
        login: 'stray228',
        display_name: 'Stray228',
        type: '',
        broadcaster_type: 'partner',
        description:
          'Крепкий и скилловый мужчина, который обучает фишкам в Dota2 и не только ;)Адекватная и доброжелательная атмосфера на каждой трансляции, всегда рад тебя поприветствовать на потоке!По вопросам рекламы stray@streamers-alliance.com',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/stray228-profile_image-ceb0393a88eb8286-300x300.jpeg',
        offline_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/9f10a988-a079-4ed6-afe8-a90e9d103f29-channel_offline_image-1920x1080.png',
        view_count: 119956993,
        email: '',
        created_at: '2013-02-19T09:36:53Z',
      },
    },
    {
      id: '04387c75-d5ab-4b06-a630-16cd8ebca014',
      channelId: '87451222',
      userId: '128644134',
      twitchUser: {
        id: '87451222',
        login: 'tiberiankrash',
        display_name: 'Tiberiankrash',
        type: '',
        broadcaster_type: '',
        description: '',
        profile_image_url:
          'https://static-cdn.jtvnw.net/user-default-pictures-uv/215b7342-def9-11e9-9a66-784f43822e80-profile_image-300x300.png',
        offline_image_url: '',
        view_count: 21,
        email: '',
        created_at: '2015-04-04T06:44:59Z',
      },
    },
    {
      id: 'aa3d3ced-906e-4e66-99b4-e2ab4370cb71',
      channelId: '42412771',
      userId: '128644134',
      twitchUser: {
        id: '42412771',
        login: 'ysx7',
        display_name: 'YSX7',
        type: '',
        broadcaster_type: 'affiliate',
        description: '',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/ysx7-profile_image-8842076fc3461ad1-300x300.png',
        offline_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/ysx7-channel_offline_image-a3afe26f7ee38647-1920x1080.jpeg',
        view_count: 7151,
        email: '',
        created_at: '2013-04-13T14:15:57Z',
      },
    },
    {
      id: 'a41d1dc8-8c1c-4c8a-ad48-5250e5223ac1',
      channelId: '789938151',
      userId: '128644134',
      twitchUser: {
        id: '789938151',
        login: 'drinkk',
        display_name: 'drinkk',
        type: '',
        broadcaster_type: '',
        description: '',
        profile_image_url:
          'https://static-cdn.jtvnw.net/user-default-pictures-uv/294c98b5-e34d-42cd-a8f0-140b72fba9b0-profile_image-300x300.png',
        offline_image_url: '',
        view_count: 0,
        email: '',
        created_at: '2022-04-21T21:03:20Z',
      },
    },
    {
      id: 'ff395d6b-f119-4655-bdca-e7d69ac33446',
      channelId: '238943226',
      userId: '128644134',
      twitchUser: {
        id: '238943226',
        login: 'omurilodev',
        display_name: 'oMuriloDev',
        type: '',
        broadcaster_type: '',
        description: 'Meu nome é Murilo, sou programador e não sei fazer live.',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/9b3d1af1-4c58-4cd8-81d1-12a4e7aa7348-profile_image-300x300.png',
        offline_image_url: '',
        view_count: 248,
        email: '',
        created_at: '2018-07-16T00:14:50Z',
      },
    },
    {
      id: '1941f0bd-086b-42f7-a6c6-a69d79d504ca',
      channelId: '98870001',
      userId: '128644134',
      twitchUser: {
        id: '98870001',
        login: 'mayonezn',
        display_name: 'mayonezn',
        type: '',
        broadcaster_type: '',
        description: '',
        profile_image_url:
          'https://static-cdn.jtvnw.net/user-default-pictures-uv/de130ab0-def7-11e9-b668-784f43822e80-profile_image-300x300.png',
        offline_image_url: '',
        view_count: 16,
        email: '',
        created_at: '2015-08-10T23:28:34Z',
      },
    },
    {
      id: '47716538-eb95-4983-a1b0-02e29cd72820',
      channelId: '494748645',
      userId: '128644134',
      twitchUser: {
        id: '494748645',
        login: 'evergowst4',
        display_name: 'Evergowst4',
        type: '',
        broadcaster_type: '',
        description: '',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/5ce256cd-93f2-4505-a560-fc24debab44f-profile_image-300x300.png',
        offline_image_url: '',
        view_count: 319,
        email: '',
        created_at: '2020-02-23T21:39:24Z',
      },
    },
    {
      id: '21d5b0b1-c4d7-4112-8453-65136223bdff',
      channelId: '193030999',
      userId: '128644134',
      twitchUser: {
        id: '193030999',
        login: 'yvnglual',
        display_name: 'yvnglual',
        type: '',
        broadcaster_type: 'affiliate',
        description: 'luvl ?Backend dev ',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/d71053c2-c533-453e-846e-c4c510285a2c-profile_image-300x300.jpeg',
        offline_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/20837f9d-66c7-482d-be81-e904d18b6281-channel_offline_image-1920x1080.png',
        view_count: 1706,
        email: '',
        created_at: '2018-01-26T14:55:40Z',
      },
    },
    {
      id: '2cf0010a-967c-49b8-973f-fa57c87615ff',
      channelId: '469218892',
      userId: '128644134',
      twitchUser: {
        id: '469218892',
        login: 'zekiukas',
        display_name: 'zekiukas',
        type: '',
        broadcaster_type: '',
        description: 'i like to play pew pew games.',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/2d388751-cf57-4520-acdc-38c3d7e69da9-profile_image-300x300.png',
        offline_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/8faa51fc-0c8a-43fc-bf67-9f209f19d660-channel_offline_image-1920x1080.png',
        view_count: 208,
        email: '',
        created_at: '2019-10-26T13:14:34Z',
      },
    },
    {
      id: 'c4e97df7-7652-4d6a-b3fd-11b4a31120d5',
      channelId: '99969771',
      userId: '128644134',
      twitchUser: {
        id: '99969771',
        login: 'pachimariclique',
        display_name: 'PachimariClique',
        type: '',
        broadcaster_type: 'affiliate',
        description: 'Не страшно, когда ты один, страшно, когда ты два',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/dbdb6c50-1353-4c18-ad8c-bc8519ebba45-profile_image-300x300.png',
        offline_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/27da033a-ff23-42f7-9567-5c66b2002324-channel_offline_image-1920x1080.png',
        view_count: 4614,
        email: '',
        created_at: '2015-08-20T18:01:30Z',
      },
    },
    {
      id: '1d96f5cb-1e6a-4f63-b46d-da0ee944d269',
      channelId: '52703474',
      userId: '128644134',
      twitchUser: {
        id: '52703474',
        login: 'sygeman',
        display_name: 'Sygeman',
        type: '',
        broadcaster_type: 'affiliate',
        description: '',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/249adb35-0006-46b9-9c5e-02fc1d975bca-profile_image-300x300.jpeg',
        offline_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/33603aea-0529-4d69-b918-59898c4a7a91-channel_offline_image-1920x1080.jpeg',
        view_count: 286628,
        email: '',
        created_at: '2013-12-05T19:43:08Z',
      },
    },
    {
      id: 'ffa13633-0dad-4a0c-a450-02642ddc8730',
      channelId: '665562197',
      userId: '128644134',
      twitchUser: {
        id: '665562197',
        login: 'lwgerry',
        display_name: 'LWGerry',
        type: '',
        broadcaster_type: '',
        description: 'Backend Node.js developer from Kyiv.',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/a7624030-2edb-4034-a361-9edb2c91f24c-profile_image-300x300.png',
        offline_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/6262a955-e277-412b-86a4-57fbd1064f38-channel_offline_image-1920x1080.png',
        view_count: 1055,
        email: '',
        created_at: '2021-03-22T15:24:30Z',
      },
    },
    {
      id: 'a8d2c7c5-971f-46ea-8a8f-ef73ad612947',
      channelId: '591764846',
      userId: '128644134',
      twitchUser: {
        id: '591764846',
        login: 'kawaguchie',
        display_name: 'kawaguchie',
        type: '',
        broadcaster_type: 'affiliate',
        description: '',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/d7ee61b2-7fa3-4de6-8a11-51c8f72f3899-profile_image-300x300.png',
        offline_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/7b2ae171-e813-46f5-8547-2f7d021f36e6-channel_offline_image-1920x1080.png',
        view_count: 9360,
        email: '',
        created_at: '2020-10-04T10:55:17Z',
      },
    },
    {
      id: 'e550438a-325c-40ce-96c6-1d7c36b50711',
      channelId: '70930005',
      userId: '128644134',
      twitchUser: {
        id: '70930005',
        login: 'rprtr258',
        display_name: 'rprtr258',
        type: '',
        broadcaster_type: 'affiliate',
        description:
          'Однозначно, стрим не для тех кому меньше 9-и. Тут надо додумывать мысль стримера и просто думать, а не плыть по течению трансляции. Сразу можно сказать, стрим тяжелый, скорее драма. Сижу и думаю, какое впечатление оставил стример.... но точно, зацепил. А это уже что-то значит. Эта трансляция оставля',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/6a70368f-8c23-4c17-947a-75e324d9198b-profile_image-300x300.png',
        offline_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/18ca67c6863077dd-channel_offline_image-1920x1080.jpeg',
        view_count: 7438,
        email: '',
        created_at: '2014-09-10T13:42:43Z',
      },
    },
    {
      id: 'e3d312df-a8e2-464d-95d4-594b20b3be7b',
      channelId: '156632065',
      userId: '128644134',
      twitchUser: {
        id: '156632065',
        login: 'vs_code',
        display_name: 'VS_Code',
        type: '',
        broadcaster_type: 'affiliate',
        description: 'всщоде, высщоде, вчкоде, вщкод, вшкод, вскод, противокод',
        profile_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/bb4e7f62-e948-45f4-b501-b684768b5e9b-profile_image-300x300.png',
        offline_image_url:
          'https://static-cdn.jtvnw.net/jtv_user_pictures/693b0e93-7e27-46b0-8d0b-81ef4779dd79-channel_offline_image-1920x1080.png',
        view_count: 2229,
        email: '',
        created_at: '2017-05-12T15:04:38Z',
      },
    },
  ];

  const [searchState, setSearch] = useState('');

  return (
    <Menu transition="skew-down" shadow="md" width={200}>
      <Menu.Target>
        <Avatar
          src="https://avatars.githubusercontent.com/u/10353856?s=460&u=88394dfd67727327c1f7670a1764dc38a8a24831&v=4"
          alt="it's me"
        />
      </Menu.Target>
      <Menu.Dropdown>
        <Menu.Label>S4tont</Menu.Label>
        <TextInput
          placeholder="Search user..."
          value={searchState}
          onChange={(event) => setSearch(event.currentTarget.value)}
          style={{ marginBottom: 5 }}
        />
        <ScrollArea type="auto" style={{ height: 250 }}>
          {dashboards
            .filter((d) => (searchState !== '' ? d.twitchUser.login.includes(searchState) : true))
            .map((d) => (
              <Menu.Item
                icon={<Image src={d.twitchUser.profile_image_url} height={20} />}
                key={d.userId + d.id}
              >
                {d.twitchUser.login}
              </Menu.Item>
            ))}
        </ScrollArea>

        <Menu.Divider />
        <Menu.Item color="red" icon={<IconLogout size={14} />}>
          Logout
        </Menu.Item>
      </Menu.Dropdown>
    </Menu>
  );
}
