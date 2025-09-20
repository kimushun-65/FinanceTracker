import type { Metadata } from 'next';
import { Geist, Geist_Mono } from 'next/font/google';
import './globals.css';

const geistSans = Geist({
  variable: '--font-geist-sans',
  subsets: ['latin'],
});

const geistMono = Geist_Mono({
  variable: '--font-geist-mono',
  subsets: ['latin'],
});

export const metadata: Metadata = {
  title: 'FinSight',
  description: 'Personal finance management application',
  manifest: '/manifest.json',
  icons: {
    icon: [
      { url: '/image/logo_1/favicon-16x16.png', sizes: '16x16', type: 'image/png' },
      { url: '/image/logo_1/favicon-32x32.png', sizes: '32x32', type: 'image/png' },
      { url: '/image/logo_1/favicon.ico', sizes: '48x48' },
    ],
    apple: [
      { url: '/image/logo_1/apple-touch-icon.png', sizes: '180x180', type: 'image/png' },
    ],
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang='en'>
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased`}
      >
        {children}
      </body>
    </html>
  );
}
