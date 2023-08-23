import React from 'react';
import ReactDOM from 'react-dom/client';
import { createBrowserRouter, RouterProvider } from 'react-router-dom';

import { Alerts } from './pages/alerts.js';
import { OBS } from './pages/obs';
import { TTS } from './pages/tts';

const router = createBrowserRouter([
	{
		path: '/:apiKey/tts',
		element: <TTS/>,
	},
	{
		path: '/:apiKey/obs',
		element: <OBS/>,
	},
	{
		path: '/:apiKey/alerts',
		element: <Alerts/>,
	},
], {
	basename: '/overlays',
});

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
	<RouterProvider router={router}/>,
);
