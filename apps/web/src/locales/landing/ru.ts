import type ILandingLocale from '@/types/landingLocaleInterface.js';
import { NavMenuTabs } from '@/types/navMenu.js';

const messages: ILandingLocale = {
  navMenu: [
    { id: NavMenuTabs.features, name: 'Функции' },
    { id: NavMenuTabs.pricing, name: 'Прайсинг' },
    { id: NavMenuTabs.reviews, name: 'Отзывы' },
    { id: NavMenuTabs.team, name: 'Команда' },
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
      title: 'Возможности бота',
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
      title: 'Бот имеет встроенный API для самых необходимых приложений',
      description:
        'Praesent dolor quis aliquam nulla id in orci. Mi sit pulvinar nunc blandit egestas cras. Sed habitant amet ultrices vitae. At volutpat enim vel quam dignissim ut justo.',
    },
    pricing: {
      title: 'У нас есть план, который идеально подходит для вас',
      features: 'Функции',
      perMonth: 'в месяц',
      plans: [
        {
          id: 1,
          name: 'Базовый план',
          features: [
            { id: 1, name: 'Unlimited commands' },
            { id: 2, name: '24 hours access' },
            { id: 3, name: '5 integrations' },
            { id: 4, name: 'Unlimited commands' },
            { id: 5, name: 'Maximum 3 users' },
            { id: 6, name: 'Maximum 3 users' },
          ],
        },
        {
          id: 2,
          name: 'Профессиональный план',
          features: [
            { id: 1, name: 'Unlimited commands' },
            { id: 2, name: '24 hours access' },
            { id: 3, name: '5 integrations' },
            { id: 4, name: 'Unlimited commands' },
            { id: 5, name: 'Maximum 3 users' },
            { id: 6, name: 'Maximum 3 users' },
          ],
        },
      ],
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
      members: [
        {
          id: 1,
          role: 'Backend разработчик',
        },
        {
          id: 2,
          role: 'Backend разработчик',
        },
        {
          id: 3,
          role: 'UI-UX дизайнер Frontend developer',
        },
      ],
    },
  },
};

export default messages;
