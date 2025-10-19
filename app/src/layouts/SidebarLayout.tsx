import { Navbar } from "./Navbar";
import { Container } from "reactstrap";
import React from "react";
import { Sidebar, SidebarItem } from "./Sidebar";
import { Footer } from "./Footer";

type SidebarLayoutProps = {
  children: React.ReactNode;
  sidebarItems: SidebarItem[];
};

export function SidebarLayout({ children, sidebarItems }: SidebarLayoutProps) {
  return (
    <React.Fragment>
      <Sidebar sidebarItems={sidebarItems} />
      <div className="page-wrapper">
        <div className="superuser-navbar">
          <Navbar showLogo={false} />
        </div>
        <Container className="mt-4 mb-4" style={{ minHeight: "calc(100vh - 190px)" }}>
          {children}
        </Container>
        <Footer />
      </div>
    </React.Fragment>
  );
}
