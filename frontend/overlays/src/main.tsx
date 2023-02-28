import React from 'react';
import ReactDOM from 'react-dom/client';
import { createBrowserRouter, RouterProvider } from 'react-router-dom';

import { TTS } from './pages/tts';

const router = createBrowserRouter([
  {
    path: '/:apiKey/tts',
    element: <TTS />,
  },
], {
  basename: '/overlays',
});

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <RouterProvider router={router} />,
);
