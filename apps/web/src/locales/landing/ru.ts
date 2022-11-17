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
            'Facilisi eget laoreet quam fringilla pulvinar diam. Risus massa ut pellentesque mi. Facilisi lobortis erat nibh diam nunc nunc. Sit natoque risus, ut malesuada',
        },
        {
          name: 'Модерация',
          description:
            'Facilisi eget laoreet quam fringilla pulvinar diam. Risus massa ut pellentesque mi. Facilisi lobortis erat nibh diam nunc nunc. Sit natoque risus, ut malesuada',
        },
        {
          name: 'Таймеры',
          description:
            'Facilisi eget laoreet quam fringilla pulvinar diam. Risus massa ut pellentesque mi. Facilisi lobortis erat nibh diam nunc nunc. Sit natoque risus, ut malesuada',
        },
        {
          name: 'Приветствия',
          description:
            'Facilisi eget laoreet quam fringilla pulvinar diam. Risus massa ut pellentesque mi. Facilisi lobortis erat nibh diam nunc nunc. Sit natoque risus, ut malesuada',
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
        'Sed eget leo adipiscing lectus nunc laoreet. Scelerisque est justo, pellentesque ut eu sit in. Suspendisse venenatis, odio dui a. Vivamus in fames augue blandit ut non sagittis, sagittis, pretium. Mollis rhoncus, pretium, morbi',
      founder: 'Основатель',
      members: {
        [TeamMemberId.Satont]: 'Backend разработчик',
        [TeamMemberId.LwGerry]: 'Backend разработчик',
        [TeamMemberId.Melkam]: 'UI-UX дизайнер Frontend разработчик',
      },
    },
  },
};

export default messages;
