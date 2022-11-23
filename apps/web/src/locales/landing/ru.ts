import { NavMenuTabs } from '@/data/landing/navMenu.js';
import { BasicPlanFeatures, PlanId, ProPlanFeatures } from '@/data/landing/pricingPlans.js';
import { TeamMemberId } from '@/data/landing/team.js';
import type ILandingLocale from '@/locales/landing/interface.js';

const messages: ILandingLocale = {
  navMenu: [
    { id: NavMenuTabs.features, name: 'Функции' },
    { id: NavMenuTabs.reviews, name: 'Отзывы' },
    { id: NavMenuTabs.team, name: 'Команда' },
    { id: NavMenuTabs.pricing, name: 'Прайсинг' },
  ],
  buttons: {
    buyPlan: 'Купить план',
    getStarted: 'Выбрать план',
    learnMore: 'Узнать больше',
    login: 'Войти',
    startForFree: 'Начать бесплатно',
    tryFeature: 'Попробовать функцию',
  },
  tagline:
    'Создан стримерами. Сделано для стримеров. Используется стримерами. Для стримеров с любовью.',
  sections: {
    features: {
      title: 'Функции бота',
      featuresInDev: 'Функции в разработке',
      content: [
        {
          name: 'Команды',
          description:
            'Мощная система команд с алиасами, счётчиками, всевозможными переменными на все случаи жизни.',
        },
        {
          name: 'Модерация',
          description:
            'Не хватает помощников, чтобы модерировать чат? Не беда. В нашей системе вы найдёте всё что вам нужно, включая быстрое удаление сообщений через nuke',
        },
        {
          name: 'Таймеры',
          description:
            'Простая система, но с воодушевлением стала популярной системой анонсов от Twitch',
        },
        {
          name: 'Приветствия',
          description: 'Хотите как-то выделить ваших любимчиков? Добавьте приветствие!',
        },
      ],
    },
    firstScreen: {
      title: '<span>Идеальный бот</span> для вашего стрима',
    },
    footer: {
      rights: '© Tsuwari {year}. Все права защищены.',
    },
    integrations: {
      preTitle: 'Интеграции',
      title: 'Множество встроенных интеграций с внешними приложения',
      description:
        'Praesent dolor quis aliquam nulla id in orci. Mi sit pulvinar nunc blandit egestas cras. Sed habitant amet ultrices vitae. At volutpat enim vel quam dignissim ut justo.',
    },
    pricing: {
      title: 'У нас есть план, который идеально подходит для вас',
      features: 'Функции',
      perMonth: 'в месяц',
      plans: {
        [PlanId.basic]: {
          name: 'Базовый план',
          price: 0,
          features: {
            [BasicPlanFeatures.first]: {
              name: 'Unlimited commands',
              status: 'accessible',
            },
            [BasicPlanFeatures.second]: {
              name: 'Unlimited commands',
              status: 'accessible',
            },
            [BasicPlanFeatures.last]: {
              name: 'Unlimited commands',
              status: 'limited',
            },
          },
        },
        [PlanId.pro]: {
          name: 'Профессиональный план',
          price: 10,
          features: {
            [ProPlanFeatures.first]: {
              name: 'Unlimited commands',
              status: 'accessible',
            },
            [ProPlanFeatures.second]: {
              name: 'Unlimited commands',
              status: 'accessible',
            },
            [ProPlanFeatures.last]: {
              name: 'Unlimited commands',
              status: 'accessible',
            },
          },
        },
      },
    },
    reviews: {
      title: 'Отзывы стримеров и других зрителей',
    },
    statLine: {
      statPlaceholder: 'Aliquam nulla',
    },
    subscribeForUpdates: {
      title: 'Подписаться на обновления',
      description:
        'Non rhoncus, neque arcu, commodo malesuada sed porttitor dictumst integer. Suscipit dictum quam ut blandit amet.',
      inputPlaceholder: 'Введите свой email',
    },
    team: {
      title: 'Наша команда',
      description:
        'Backend часть была полностью написана Satont, ранние версии сайта тоже. Позже к нам присоединился Melkam и нарисовал новый, великолепный дизайн, и воплотил наши задумки в жизнь.',
      founder: 'Основатель',
      members: {
        [TeamMemberId.Satont]: 'Backend разработчик',
        [TeamMemberId.Melkam]: 'UI-UX дизайнер Frontend разработчик',
      },
    },
  },
};

export default messages;
