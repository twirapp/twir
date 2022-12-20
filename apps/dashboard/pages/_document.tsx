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
        <Head />
        <body>
          <Main />
          <NextScript />
        </body>
      </Html>
    );
  }
}
