import { Navbar } from "./Navbar";
import { Container } from "reactstrap";
import React from "react";
import { Sidebar, SidebarItem } from "./Sidebar";
import { Footer } from "./Footer";
import { AdminAnalyticsModal } from "../components/AdminAnalyticsModal";
type SidebarLayoutProps = {
  children: React.ReactNode;
  sidebarItems: SidebarItem[];
};

export function SidebarLayout({ children, sidebarItems }: SidebarLayoutProps) {
  return (
    <div className="admin-layout">
      <Sidebar sidebarItems={sidebarItems} />
      <div className="superuser-navbar">
        <Navbar showLogo={false} />
      </div>
      <div className="page-wrapper wrapper-sidebar">
        <Container className="mt-4 mb-4 admin-page-content">
          {children}
        </Container>
        <AdminAnalyticsModal />
        <Footer />
      </div>
    </div>
  );
}
