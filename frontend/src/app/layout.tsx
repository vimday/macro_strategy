import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import { AntdRegistry } from '@ant-design/nextjs-registry';
import { Providers } from './providers';
import '@ant-design/v5-patch-for-react-19';

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "Macro Strategy | 巨策略",
  description: "A comprehensive platform for testing and comparing macro trading strategies",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body className={inter.className}>
        <AntdRegistry>
          <Providers>
            {children}
          </Providers>
        </AntdRegistry>
      </body>
    </html>
  );
}
