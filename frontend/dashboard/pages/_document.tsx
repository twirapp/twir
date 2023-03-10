import { createGetInitialProps } from '@mantine/next';
import nextI18nextConfig from 'next-i18next.config.js';
import Document, { Head, Html, Main, NextScript } from 'next/document';

const getInitialProps = createGetInitialProps();

export default class _Document extends Document {
  static getInitialProps = getInitialProps;

  render() {
    const currentLocale = this.props.__NEXT_DATA__.locale || nextI18nextConfig.i18n.defaultLocale;

    return (
      <Html lang={currentLocale}>
        <Head>
          <link rel="shortcut icon" href="/dashboard/TsuwariInCircle.svg" sizes="any" type="image/svg+xml" />

          <meta name="application-name" content="Twir app" />
          <meta name="apple-mobile-web-app-capable" content="yes" />
          <meta name="apple-mobile-web-app-status-bar-style" content="default" />
          <meta name="apple-mobile-web-app-title" content="Twir app" />
          <meta name="description" content="Twir app" />
          <meta name="format-detection" content="telephone=no" />
          <meta name="mobile-web-app-capable" content="yes" />
          <meta name="msapplication-config" content="/dashboard/icons/browserconfig.xml" />
          <meta name="msapplication-TileColor" content="#2B5797" />
          <meta name="msapplication-tap-highlight" content="no" />
          <meta name="theme-color" content="#000000" />
          <meta name="darkreader-lock" />


          <link rel="manifest" href="/dashboard/manifest.json" />
          <link rel="stylesheet" href="https://fonts.googleapis.com/css?family=Roboto:300,400,500" />

          <meta property="og:type" content="website" />
          <meta property="og:title" content="Twir app" />
          <meta property="og:description" content="Twir bot dashboard" />
          <meta property="og:site_name" content="Twir app" />
          <meta property="og:url" content="https://twir.app" />
          <meta property="og:image" content="https://twir.app/dashboard/icons/apple-touch-icon.png" />
        </Head>
        <body>
          <Main />
          <NextScript />
        </body>
      </Html>
    );
  }
}
